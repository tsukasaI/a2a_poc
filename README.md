# Sample POC

hydra + mobile_pay_app + client

## Native App OAuth0 Authentication Flow

1. create client

```sh
docker compose exec hydra hydra create oauth2-client --token-endpoint-auth-method none --endpoint http://localhost:4445/ --format json --name poc_client02 --response-type code --grant-type authorization_code,refresh_token --redirect-uri http://localhost:3002/callback  --redirect-uri testDeepLink://mobile | jq

```

sample response

```json
{
  "client_id": "4d9b630b-abfa-4aaf-82c7-9a50518cb68b",
  "client_name": "poc_client01",
  "client_secret_expires_at": 0,
  "client_uri": "",
  "created_at": "2024-09-02T03:27:18Z",
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
    "http://localhost:3002/callback",
    "testDeepLink://mobile"
  ],
  "registration_access_token": "ory_at_xMpDkLy_VRwRollXVy7PW5yx17yy-krsz5qw3Nlw-5c.Rww4xtkzyuw5S2JZP7H2sIpGx-W3jh6XgxHtIUVPNYE",
  "registration_client_uri": "http://localhost:4444/oauth2/register/",
  "request_object_signing_alg": "RS256",
  "response_types": [
    "code"
  ],
  "scope": "offline_access offline openid",
  "skip_consent": false,
  "skip_logout_consent": false,
  "subject_type": "public",
  "token_endpoint_auth_method": "none",
  "tos_uri": "",
  "updated_at": "2024-09-02T03:27:17.880303Z",
  "userinfo_signed_response_alg": "none"
}

{
  "client_id": "579bdbf0-42d0-4576-bbee-3e9846356643",
  "client_name": "poc_client02",
  "client_secret_expires_at": 0,
  "client_uri": "",
  "created_at": "2024-09-02T04:17:46Z",
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
    "http://localhost:3002/callback",
    "http://localhost:3001/consent",
    "testDeepLink://mobile"
  ],
  "registration_access_token": "ory_at_NV0pvaV2WDOOZe0XjbfBDgUdzBdcIE95r_4bAPIC1RM.M4sM04EUOSjnMf3W8ezDqjxpDlptRFmvsSl7NigMTCg",
  "registration_client_uri": "http://localhost:4444/oauth2/register/",
  "request_object_signing_alg": "RS256",
  "response_types": [
    "code"
  ],
  "scope": "offline_access offline openid",
  "skip_consent": false,
  "skip_logout_consent": false,
  "subject_type": "public",
  "token_endpoint_auth_method": "none",
  "tos_uri": "",
  "updated_at": "2024-09-02T04:17:45.640365Z",
  "userinfo_signed_response_alg": "none"
}
```

1. Set client id on `clinet_app/src/App.tsx` > `clientId`

1. Run Client and mobile app
