#
# Build the frontend UI
#

FROM node as ui

WORKDIR ui

COPY ui/yarn.lock ui/package.json ./

RUN yarn

COPY ui .

RUN yarn build

#
# Build the back-end server with the embedded-frontend
#

FROM golang:alpine as server

RUN apk update && \
    apk add curl git file && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/niranjan94/vault-front/

RUN go get -v github.com/GeertJohan/go.rice/rice

COPY Gopkg.lock Gopkg.toml ./

RUN dep ensure -v -vendor-only

COPY src src
COPY main.go .

COPY --from=ui /ui/dist ui/dist

RUN rice embed-go -v && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -v && \
    file vault-front

#
# Build the final scratch image with just the server binary
#

FROM scratch

WORKDIR /app

COPY --from=server /go/src/github.com/niranjan94/vault-front/vault-front /app

EXPOSE 8000

CMD ["/app/vault-front"]

