path "auth/userpass/users/*" {
  capabilities = ["read", "update", "create"]
}

path "auth/userpass/login/*" {
  capabilities = ["create"]
}

path "kv/metadata/auth/userpass/*" {
  capabilities = ["read", "update", "delete", "create"]
}

path "totp/code/*" {
  capabilities = ["create", "read", "update"]
}

path "totp/keys/*" {
  capabilities = ["create", "delete", "read", "update"]
}

path "auth/token/lookup-accessor" {
  capabilities = ["update"]
}

path "database/roles" {
  capabilities = ["list"]
}