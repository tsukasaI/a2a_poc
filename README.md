# Sample POC

hydra + pay_app + client

## Native App OAuth0 Authentication Flow

1. create client

```sh
# Create client
docker compose exec hydra hydra create oauth2-client --token-endpoint-auth-method none --endpoint http://localhost:4445/ --format json --name poc_client01 --response-type code --grant-type authorization_code,refresh_token --redirect-uri http://localhost:3002/callback --redirect-uri testDeepLink://mobile | jq

docker compose exec hydra hydra create oauth2-client --token-endpoint-auth-method none --endpoint http://localhost:4445/ --format json --name poc_client02 --response-type code --grant-type authorization_code,refresh_token --redirect-uri testDeepLink://mobile | jq

```

```json
{
  "client_id": "3a121f8a-9802-4efc-b23d-214a90cda035",
  "client_name": "poc_client02",
  "client_secret_expires_at": 0,
  "client_uri": "",
  "created_at": "2024-08-29T06:36:33Z",
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
    "testDeepLink://mobile"
  ],
  "registration_access_token": "ory_at_3rOnPu7Eraiknq4phuU8oDWKWfSWE2HNcyiZWOtHdr8.qzcaxFlUwhrPAXH0sLiTWjw983OqzPeIM1Rc-YsrTxk",
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
  "updated_at": "2024-08-29T06:36:32.52215Z",
  "userinfo_signed_response_alg": "none"
}

{
  "client_id": "d03a0fa9-3059-4d16-bf51-1decc9bae9b0",
  "client_name": "poc_client01",
  "client_secret_expires_at": 0,
  "client_uri": "",
  "created_at": "2024-08-30T08:22:25Z",
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
  "registration_access_token": "ory_at_3w4hweL_MH8iaYaB6ydvBQ0Hi-zFqwnzSkFU8Tj4tlU.NLfWe7IwosAbW81GlsF9n0N9XbsMDIJIe71zzZgjK2U",
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
  "updated_at": "2024-08-30T08:22:24.965254Z",
  "userinfo_signed_response_alg": "none"
}
```

1. open blow page
localhost:4444/oauth2/auth?client_id=${client_id}&redirect_uri=http%3A%2F%2Flocalhost%3A3002%2Fcallback&response_type=code&state=aaaaaaaaaaaaaaaaaa

localhost:4444/oauth2/auth?client_id=3a121f8a-9802-4efc-b23d-214a90cda035&redirect_uri=testDeepLink%3A%2F%2Fmobile&response_type=code&state=aaaaaaaaaaaaaaaaaa&code_challenge_method=S256&code_challenge=h3HrXnwQFqrpz-hXmJgmFCV6TTCTlGB3415JorTSexc

できたこと
GoでPKCE作成してcode_challenge=に貼る
iOSのsafariに貼り付けて進める

TODO
クライアントからのフロー作成
- Goのoauthのリポジトリをcloneしてデバッグする（POSTデータを見る）
- Clientでcode exchange実施
