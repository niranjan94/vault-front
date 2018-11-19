#!/usr/bin/env bash

if [[ -z "$1" ]]; then
    trap "exit" INT TERM
    trap "kill 0 && exit" EXIT
fi

export PATH=$PATH:$(pwd)/.dev

echo "Starting vault server in dev mode..."
vault server -dev > /dev/null 2>&1 &

export VAULT_ADDR=http://127.0.0.1:8200
export PII_VAULT_ADDRESS=http://127.0.0.1:8200;
export TEST_WORKING_DIR=$(pwd);

sleep 1

echo "Setting up vault..."
vault secrets enable totp;
vault secrets enable -version=1 kv;
vault secrets enable -version=1 -path=database kv; # Mount a kv backend at /database for testing purposes
vault auth enable userpass;

# Write the manager policy
vault policy write manager vault/policies/vault-front.hcl;

# Write user policies with varying database access for testing purposes
vault policy write test-user-policy scripts/testing/test-user-policy.hcl
vault policy write test-user-policy-restricted scripts/testing/test-user-policy-restricted.hcl

# Write dummy database configs to the database mount clone
vault kv put database/config/development-master @scripts/testing/db-configs/development-master.json
vault kv put database/config/testing-master @scripts/testing/db-configs/testing-master.json
vault kv put database/config/staging-master @scripts/testing/db-configs/staging-master.json

# Write dummy database credentials to the database mount clone
vault kv put database/creds/database-role-one @scripts/testing/db-creds/database-role-one.json
vault kv put database/creds/database-role-two @scripts/testing/db-creds/database-role-two.json
vault kv put database/creds/database-role-three @scripts/testing/db-creds/database-role-three.json
vault kv put database/creds/database-role-three @scripts/testing/db-creds/database-role-four.json

# Write dummy database credentials to the database mount clone for creds listing purposes
vault kv put database/roles/database-role-one @scripts/testing/db-roles/database-role-one.json
vault kv put database/roles/database-role-two @scripts/testing/db-roles/database-role-two.json
vault kv put database/roles/database-role-three @scripts/testing/db-roles/database-role-three.json
vault kv put database/roles/database-role-three @scripts/testing/db-roles/database-role-four.json

export PII_VAULT_TOKEN=$(vault token create -policy=manager -field=token);
export PII_TEST_VAULT_TOKEN=$(vault token create -policy=manager -field=token);

echo "Running tests..."

go test -coverprofile=coverage.out -v ./...