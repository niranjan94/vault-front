package api

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/GeertJohan/go.rice"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/niranjan94/vault-front/src/api/auth"
	"github.com/niranjan94/vault-front/src/api/routes"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	e        *echo.Echo
	cert     *tls.Certificate
	certLock = new(sync.RWMutex)
)

func StartApiServer(box *rice.HTTPBox, withGracefulExit bool) *echo.Echo {
	// already configured, restarting listener at runtime is not currently supported
	if e != nil {
		return e
	}

	listenAddress := viper.GetString("api.listenAddress")
	useSsl := viper.GetBool("api.ssl.enabled")

	e = echo.New()
	e.HideBanner = true
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 2 * time.Minute

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.BodyLimit("100K"))
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: viper.GetStringSlice("cors.allowedOrigins"),
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			return next(c)
		}
	})

	if useSsl {
		e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
			XSSProtection:         "1; mode=block",
			ContentTypeNosniff:    "nosniff",
			XFrameOptions:         "SAMEORIGIN",
			ContentSecurityPolicy: "default-src 'self' blob: https://api.github.com;",
		}))
	}

	assetHandler := http.FileServer(box)
	e.GET("/*", echo.WrapHandler(assetHandler))

	guest := e.Group("/api/v1")

	guest.POST("/auth/session", routes.Login())

	authenticated := e.Group("/api/v1")
	authenticated.Use(auth.RequireTokenAuthentication())

	authenticated.GET("/databases", routes.GetAllowedDatabases())
	authenticated.POST("/databases/credentials", routes.GetCredential())

	authenticated.DELETE("/auth/session", routes.Logout())

	// Start server
	go func() {
		if !useSsl {
			if err := e.Start(listenAddress); err != nil {
				e.Logger.Info("shutting down the server")
				return
			}
		} else {
			e.TLSServer.TLSConfig = new(tls.Config)
			e.TLSServer.TLSConfig.MinVersion = tls.VersionTLS12
			e.TLSServer.TLSConfig.PreferServerCipherSuites = true
			e.TLSServer.TLSConfig.CipherSuites = []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			}

			c, err := tls.LoadX509KeyPair(viper.GetString("api.ssl.crt"), viper.GetString("api.ssl.key"))
			if err != nil {
				log.Fatalln(err.Error())
				return
			}
			certLock.Lock()
			cert = &c
			certLock.Unlock()

			e.TLSServer.TLSConfig.GetCertificate = GetCertificate
			e.TLSServer.Addr = listenAddress

			if err := e.StartServer(e.TLSServer); err != nil {
				e.Logger.Info("shutting down the server")
				return
			}
		}
	}()

	if withGracefulExit {
		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	}

	return e
}

func GetCertificate(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
	certLock.RLock()
	defer certLock.RUnlock()

	if cert == nil {
		return nil, errors.New("no certificate configured")
	}
	return cert, nil
}
