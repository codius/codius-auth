# Codius Auth

Auth server that authenticates based on Interledger payment.

### Environment Variables

#### CODIUS_HOST_URL
* Type: String
* Description: Root URL of the Codius host.

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

#### `/forward`
[Forward authentication](https://docs.traefik.io/v1.7/configuration/entrypoints/#forward-authentication) endpoint for requests to a Codius service.

This has been tested with [Traefik v1.7](https://docs.traefik.io/v1.7/configuration/backends/kubernetes/#authentication)

>If the response code is 2XX, access is granted and the original request is performed. Otherwise, the response from the authentication server is returned.

Returns 200 after successfully debiting [`AUTH_PRICE`](#auth_price) from the Codius service's balance at the [receipt verifier](https://github.com/coilhq/receipt-verifier).

Otherwise, returns 303 redirect to `/{ID}/402` if the service's balance is insufficient.

#### `/token`
Kubernetes [webhook token authentication](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#webhook-token-authentication) endpoint for requests to the Kubernetes API Server.

This is intended for requests to create Codius services. The token is expected to be the Web Monetization `requestId` from visiting the Codius host's home page.

Grants authentication as [`RBAC_USER`](#rbac_user) after successfully debiting [`AUTH_PRICE`](#auth_price) from the token's balance at the [receipt verifier](https://github.com/coilhq/receipt-verifier).
