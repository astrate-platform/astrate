# ⚠️ Local development keys — NEVER use outside localhost

This keypair is **committed to the repository** so `docker compose --profile
full up` gives a working realm (`test`) and dashboard login with zero
ceremony. Anyone with this repo can mint valid JWTs for any deployment that
uses these keys.

For anything that is not a throwaway local stack, generate your own pair and
point `realm.jwt_public_key_file` (or `ASTRATE_REALM_JWT_PUBLIC_KEY_FILE`) at
it:

```sh
openssl genrsa -out realm_private.pem 2048
openssl rsa -in realm_private.pem -pubout -out realm_public.pem
```

Mint a dashboard login token with astartectl:

```sh
astartectl utils gen-jwt all-realm-apis -k deploy/devrealm/realm_private.pem
```
