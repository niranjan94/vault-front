
api = {
  listenAddress = "0.0.0.0:8000"
  ssl = {
    crt = "/path/to/ssl.crt"
    enabled = false
    key = "/path/to/ssl.key"
  }
}

disableMLock = false

cors = {
  allowedOrigins = ["*"]
}

vault = {
  authMode = "aws"
  address = "http://127.0.0.1:8200"
  token = "xyzzy"
}
