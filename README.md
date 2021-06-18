# hsdp-function-vault-read

Read Vault keys from a HSDP Vault instance

# usage

```terraform
resource "hsdp_function" "vault_read" {
  name         = "vault-read"
  docker_image = "philipslabs/hsdp-function-vault-read:v0.0.4"

  environment = {
    VREAD_ENDPOINT            = "https://vproxy.us-east.philips-healthsuite.com/"
    VREAD_ROLE_ID             = "XXX"
    VREAD_SECRET_ID           = "YYY"
    VREAD_SPACE_SECRET_PATH   = "/v1/cf/8cb5a2ea-d20a-4ea0-815b-742075dc92ba/secret"
    VREAD_SERVICE_SECRET_PATH = "/v1/cf/51536c9b-f91c-402a-87f5-406258c792df/secret"
  }

  backend {
    credentials = module.siderite_backend.credentials
  }

}

output "sync_endpoint" {
  value = hsdp_function.vault_read.endpoint
}
```

Append `/vault/read/service/:key` to the `sync_endpoint.value` with the right token in `Authorization` header (fixed or IAM) and the function will
read all the values from the cofngiured Vault under `:key` 

## curl example

Write something to the vault:

```shell
vault write cf/51536c9b-f91c-402a-87f5-406258c792df/secret/andy \
         secret1="value stored in vault" \
         config="or a config item"
```

Then retrieve the data via the `hsdp_function`:

```shell
curl -H "Authorization: Token ZZZ" \
  https://hsdp-func-gateway-www.eu-west.philips-healthsuite.com/function/aaa/vault/read/service/andy
```
output:
```json
{
   "secret1": "value stored in vault",
   "config": "or a config item"
}
```

# license

License is MIT
