# Sample POC

hydra + pay_app + client

## Native App OAuth0 Authentication Flow

1. create client

```sh
# Create client
docker compose exec hydra hydra create oauth2-client --endpoint http://localhost:4445/ --format json --name poc_client01 --response-type code --grant-type authorization_code,refresh_token --redirect-uri http://localhost:3001/callback | jq

```

{
  "client_id": "f4d402cd-a3b9-446b-857f-4411ac2bde71",
  "client_name": "poc_client01",
  "client_secret": "mba2vumlP3G82.wVVC683sQt39",
  "client_secret_expires_at": 0,
  "client_uri": "",
  "created_at": "2024-08-23T02:55:02Z",
  "grant_types": [
    "authorization_code",
    "refresh_token"
  ],
  "jwks": {},
  "logo_uri": "",
  "metadata": {},
  "owner": "",
  "policy_uri": "",
  "redirect_uris": [
    "http://localhost:3001/callback"
  ],
  "registration_access_token": "ory_at_K-ij59y9x4yx7YyMW1u7hYKFbg1jtW93g0tOCH5dXng.eA8xEzuQAprIJt3kHsP5aFJ46oQMk9fq-2AFTLS-4gQ",
  "registration_client_uri": "http://localhost:4444/oauth2/register/",
  "request_object_signing_alg": "RS256",
  "response_types": [
    "code"
  ],
  "scope": "offline_access offline openid",
  "skip_consent": false,
  "skip_logout_consent": false,
  "subject_type": "public",
  "token_endpoint_auth_method": "client_secret_basic",
  "tos_uri": "",
  "updated_at": "2024-08-23T02:55:02.0892Z",
  "userinfo_signed_response_alg": "none"
}


1. open blow page
localhost:4444/oauth2/auth?client_id=${client_id}&redirect_uri=http%3A%2F%2Flocalhost%3A3001%2Fcallback&response_type=code&state=aaaaaaaaaaaaaaaaaa
localhost:4444/oauth2/auth?client_id=f4d402cd-a3b9-446b-857f-4411ac2bde71&redirect_uri=http%3A%2F%2Flocalhost%3A3002%2Fcallback&response_type=code&state=aaaaaaaaaaaaaaaaaa
