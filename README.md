# Codius Auth

![](https://github.com/codius/codius-auth/workflows/Docker%20CI/badge.svg)

Auth server that authenticates based on Interledger payment.

### Environment Variables

#### AUTH_PRICE
* Type: Number
* Description: The amount required to have been paid to grant authorization. Denominated in receiver's asset (code and scale).

#### PORT
* Type: String
* Description: The port that webhook server will listen on.
* Default: 8080

#### RBAC_USER
* Type: String
* Description: Kubernetes user as whom valid tokens grant authorization.

#### RECEIPT_VERIFIER_URL
* Type: String
* Description: URL of the [receipt verifier](https://github.com/coilhq/receipt-verifier/) with which to deduct paid balances.

### Routes

#### `/token`
Kubernetes [webhook token authentication](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#webhook-token-authentication) endpoint for requests to the Kubernetes API Server.

This is intended for requests to create Codius services. The token is expected to be the Web Monetization `requestId` from visiting the Codius host's home page.

Grants authentication as [`RBAC_USER`](#rbac_user) after successfully debiting [`AUTH_PRICE`](#auth_price) from the token's balance at the [receipt verifier](https://github.com/coilhq/receipt-verifier).
