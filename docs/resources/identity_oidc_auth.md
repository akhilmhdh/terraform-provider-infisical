---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "infisical_identity_oidc_auth Resource - terraform-provider-infisical"
subcategory: ""
description: |-
  Create and manage identity oidc auth in Infisical.
---

# infisical_identity_oidc_auth (Resource)

Create and manage identity oidc auth in Infisical.

## Example Usage

```terraform
terraform {
  required_providers {
    infisical = {
      # version = <latest version>
      source = "infisical/infisical"
    }
  }
}

provider "infisical" {
  host          = "https://app.infisical.com" # Only required if using self hosted instance of Infisical, default is https://app.infisical.com
  client_id     = "<>"
  client_secret = "<>"
}

resource "infisical_project" "example" {
  name = "example"
  slug = "example"
}

resource "infisical_identity" "machine-identity-1" {
  name   = "machine-identity-1"
  role   = "admin"
  org_id = "<>"
}

resource "infisical_identity_oidc_auth" "oidc-auth" {
  identity_id        = infisical_identity.machine-identity-1.id
  oidc_discovery_url = "<>"
  bound_issuer       = "<>"
  bound_audiences    = ["sample-audience"]
  bound_subject      = "<>"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `bound_issuer` (String) The unique identifier of the identity provider issuing the OIDC tokens.
- `identity_id` (String) The ID of the identity to attach the configuration onto.
- `oidc_discovery_url` (String) The URL used to retrieve the OpenID Connect configuration from the identity provider.

### Optional

- `access_token_max_ttl` (Number) The maximum lifetime for an access token in seconds. This value will be referenced at renewal time. Default: 2592000
- `access_token_num_uses_limit` (Number) The maximum number of times that an access token can be used; a value of 0 implies infinite number of uses. Default:0
- `access_token_trusted_ips` (Attributes List) A list of IPs or CIDR ranges that access tokens can be used from. You can use 0.0.0.0/0, to allow usage from any network address... (see [below for nested schema](#nestedatt--access_token_trusted_ips))
- `access_token_ttl` (Number) The lifetime for an access token in seconds. This value will be referenced at renewal time. Default: 2592000
- `bound_audiences` (List of String) The comma-separated list of intended recipients.
- `bound_claims` (Map of String) The attributes that should be present in the JWT for it to be valid. The provided values can be a glob pattern.
- `bound_subject` (String) The expected principal that is the subject of the JWT.
- `oidc_ca_certificate` (String) The PEM-encoded CA cert for establishing secure communication with the Identity Provider endpoints

### Read-Only

- `id` (String) The ID of the oidc auth.

<a id="nestedatt--access_token_trusted_ips"></a>
### Nested Schema for `access_token_trusted_ips`

Optional:

- `ip_address` (String)