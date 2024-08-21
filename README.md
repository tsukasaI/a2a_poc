# Sample POC

hydra + pay_app + client

## Native App OAuth0 Authentication Flow

1. create client

```
docker compose exec hydra hydra create oauth2-client --endpoint http://localhost:4445/ --format json --name poc_client01 --response-type code --grant-type authorization_code,refresh_token --redirect-uri http://127.0.0.1:3001/callback | jq

```

localhost:4444/oauth2/auth?client_id=${client_id}&redirect_uri=http%3A%2F%2F127.0.0.1%3A3001%2Fcallback&response_type=code&state=aaa
