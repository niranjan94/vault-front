
api = {
  listenAddress = "0.0.0.0:8000"
  ssl = {
    crt = "/path/to/test.ssl.crt"
    enabled = false
    key = "/path/to/test.ssl.key"
  }
}

disableMLock = false

cors = {
  allowedOrigins = ["*"]
}

vault = {
  address = "http://127.0.0.1:8200"
}
