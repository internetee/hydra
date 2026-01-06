# Changes in branch: internetee:hydra:develop-eeid (compare with ory/hydra:master)

Date: 2026-01-06  
Made by [Sergei Tsõganov](https://github.com/maricavor)

Branch: internetee:hydra:develop-eeid  
Target base: ory/hydra master

This document summarizes the notable changes introduced on this branch, their rationale, and migration/compatibility guidance.

---

## Summary (high level)

- OIDC discovery document and OIDC-related behavior have been significantly simplified / hardened. Many optional discovery properties were removed or commented out; the discovery now advertises a minimal set of metadata.
- Default client scope was reduced to `openid` (removed `offline_access` / `offline` from defaults).
- Well-known JWKS / JWT keys exposure was aligned with the configured access token strategy.
- The OIDC userinfo behavior was changed: claim normalization, removal of many standard JWT claims from the userinfo response, pulling `profile_attributes` into top-level claims, and setting `auth_time` from `iat`.
- Persistence of consent/login flows now consults the DB for authoritative `login_was_used` and enforces `login_was_used` when persisted for some states.
- Tests and snapshots were updated to reflect these behavior changes.
- `Version` was bumped to `2.2`.
- A development Dockerfile was added for local builds and a distroless runner image is provided.

---

## Notable changes (by area)

### OIDC discovery and oauth2 handler
Files:
- `oauth2/handler.go`
- `oauth2/.snapshots/TestHandlerWellKnown-hsm_enabled=false.json`

What changed:
- The discovery structure (`.well-known/openid-configuration`) was simplified. The following fields were removed / commented-out or significantly reduced:
  - `registration_endpoint`
  - `revocation_endpoint`
  - `response_modes_supported`
  - `code_challenge_methods_supported` (PKCE methods)
  - back/front-channel logout support fields (`frontchannel_logout_*`, `backchannel_logout_*`)
  - request object/request URI registration/claims parameter flags
  - verifiable credentials / credentials endpoint metadata (commented out)
  - many response types (now only `code`) and grant types (now only `authorization_code`)
  - token endpoint auth methods reduced to `client_secret_basic`
- Two new discovery fields were added:
  - `claim_types_supported` (value example: `["normal"]`)
  - `ui_locales_supported` (value example: `["et","en","ru"]`)
- The HTTP handler that produces discovery now returns a compact/minimal set of metadata. Many previously advertised features are intentionally not exposed any more by default.

Impact:
- Relying parties (clients) that read the discovery document to detect support for features such as:
  - dynamic client registration,
  - revocation endpoint,
  - PKCE methods,
  - multiple grant/response types,
  - logout endpoints,
  - request parameter/URI features,
  may no longer find them in the discovery and must be configured explicitly or adjusted.
- Clients should not assume `offline_access` is part of the default scope set anymore.

Migration guidance:
- If you need discovery to advertise any of the commented/removed features, configure your deployment to add those strings back in the provider configuration (or adjust your customization of the discover handler back to the previous shape).
- Update clients that are expecting `response_types_supported` or `grant_types_supported` arrays with values other than `code` / `authorization_code`.
- If your clients rely on `token_endpoint_auth_methods_supported` containing `client_secret_post`, `private_key_jwt` or `none`, update or reconfigure them to work with `client_secret_basic` or expose those via configuration.

### Default client scope and provider defaults
Files:
- `driver/config/provider.go`
- `driver/config/provider_test.go`

What changed:
- Default client scope changed from `["offline_access", "offline", "openid"]` to `["openid"]`.
- OIDC discovery supported scope list likewise removed `offline_access`/`offline` from the defaults.
- Well-known keys now include the OAuth2 JWT key only if `AccessTokenStrategy` is `AccessTokenJWTStrategy`.

Impact:
- Any client or automation assuming offline-access/refresh token scopes are present by default must explicitly request them.
- Tests and snapshots updated accordingly.

Migration guidance:
- Explicitly request `offline_access` if you need refresh token/consent behavior.
- Review code that reads `WellKnownKeys` or expects `hydra.jwt.access-token` to be present unconditionally.

### Userinfo changes
Files:
- `oauth2/handler.go`
- `oauth2/handler_test.go`

What changed:
- The userinfo response construction deletes/filters a lot of JWT claims previously returned in userinfo:
  - deleted: `nonce`, `at_hash`, `c_hash`, `exp`, `sid`, `jti`, `state`, `iat`, `iss`, `nbf`, `aud`, `rat`, `amr`
- `auth_time` is set to the original `iat` claim.
- If a `profile_attributes` map exists in the session claims it is merged into the top-level userinfo claims (except subject).
- The logic that attempted to ensure `aud` contained the client ID has been commented out.
- Tests previously asserting on `aud` values were commented-out to align with the new behavior.

Impact:
- Userinfo responses will be leaner and will not include various JWT internals that previously might have been present.
- If systems depended on an `aud` claim included in userinfo, they must not rely on it any more.

Migration guidance:
- Consumers of userinfo should be reviewed and adjusted to rely on the stable claims that remain (and/or explicitly include required claims in session attributes).
- If `aud` is required in userinfo for your integrations, adjust the code to add it explicitly before returning or restore the original logic guarded by appropriate configuration.

### Internal models (httpclient)
Files:
- `internal/httpclient/model_oidc_configuration.go`

What changed:
- New JSON fields added to the OIDC configuration model used by the internal HTTP client:
  - `claim_types_supported`
  - `ui_locales_supported`

Impact:
- Clients using the internal generated client will be able to read the new properties.

### Persistence / Consent flow changes
Files:
- `persistence/sql/persister_consent.go`

What changed:
- On GetLoginRequest: the DB is consulted for the canonical flow using `GetFlow` and if present the `login_was_used` from DB is used overriding encoded flow values.
- On VerifyAndInvalidateConsentRequest: when a flow is first persisted, `login_was_used` is set to true for certain states (various FlowState* values). This ensures the persisted flow reflects the login-used status correctly.
- Some debug logs (commented) were added.

Impact:
- Flows encoded in client-side tokens/challenges that may have stale `login_was_used` values are now reconciled with the stored DB flow.
- This prevents subtle re-use or stale behavior around login usage flags.

Migration guidance:
- No configuration change required, but operators should be aware that the DB is now considered authoritative for `login_was_used`.
- Review any custom persistence or flow-manipulating code to ensure compatibility with the stronger DB-driven behavior.

### Tests & snapshots
Files:
- Numerous `.snapshots` under `client/`, `cmd/`, `oauth2/`, `cmd/`, `client/`
- Unit tests adjusted to reflect removed fields and behaviors (aud assertions commented out; verifiable credentials checks gated on presence).

What changed:
- Snapshots updated to expect `scope: "openid"` rather than `offline_access offline openid`.
- Tests that previously required certain discovery elements or userinfo aud presence were either removed/commented or guarded.

Impact:
- Test behavior is aligned with the new defaults. If you run tests against a customized deployment that still advertises the earlier features, you may need to adapt snapshots or tests accordingly.

### Version & build
Files:
- `driver/config/buildinfo.go`
- `.docker/Dockerfile-dev` (new)

What changed:
- Build info version string changed from `"master"` to `"2.2"`.
- A `.docker/Dockerfile-dev` was added which:
  - builds Hydra using golang:1.22 (builder stage),
  - copies the binary into a distroless Debian 12 nonroot runner,
  - exposes ports 4444 (public) and 4445 (admin),
  - uses `hydra serve all` as default command.

Impact:
- CI / developer environments can use `.docker/Dockerfile-dev` to produce a small runtime image for local tests.
- Update any automation that depends on the `Version` string if necessary.

---

## Breaking changes and compatibility notes

- Discovery now advertises a minimal subset by default. This can break automatic feature-detection by clients. Review all relying parties and custom integrations that read `.well-known/openid-configuration`.
- Default client scopes have changed; automated client creation which assumed refresh tokens from default scopes will no longer get them. Explicitly request `offline_access` where required.
- `userinfo` no longer provides `aud` or many JWT internals by default — integrations that relied on these must be updated.
- Tests and snapshots were changed to match the new defaults; if you maintain forks that expected previous defaults, align your snapshots/tests.
- If you rely on `client_secret_post`, `private_key_jwt`, or `none` as token auth methods, they will not be advertised by default; adjust configuration or client logic.

---

## Migration checklist (recommended steps if upgrading)
- Review and update clients to:
  - Explicitly request `offline_access` when refresh tokens are needed.
  - Avoid depending on `aud` being present in userinfo.
  - Avoid reliance on discovery for features that are no longer advertised by default (PKCE, registration, revocation, additional response types).
- If you need older discovery behavior, restore/extend the discover handler in your deployment or set provider config values that re-expose the necessary discovery entries.
- Rebuild or re-deploy with the new `Version` if you want to reflect the `2.2` label.
- Run the updated test-suite and update snapshot expectations if you maintain custom tests or CI that compare snapshots.

---

## Files touched (non-exhaustive focused list)
- Added:
  - `.docker/Dockerfile-dev`
- Modified:
  - `driver/config/buildinfo.go` (Version -> 2.2)
  - `driver/config/provider.go` (WellKnownKeys, DefaultClientScope, OIDCDiscoverySupportedScope)
  - `driver/config/provider_test.go` (tests updated)
  - `internal/httpclient/model_oidc_configuration.go` (new fields)
  - `oauth2/handler.go` (discovery and userinfo logic simplified/changed)
  - `persistence/sql/persister_consent.go` (login_was_used DB reconciliation and persistence logic)
  - Many `.snapshots` and tests across `client/`, `cmd/`, `oauth2/` to match new behavior

---
