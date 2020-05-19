# Codius Token Auth Webhook

Kubernetes [webhook token authentication](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#webhook-token-authentication) server that grants authorization based on Interledger payment.

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
* Description: URL of the [receipt verifier](https://github.com/wilsonianb/receipt-verifier/) with which to deduct paid balances.
