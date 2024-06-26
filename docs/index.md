---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "infisical Provider"
subcategory: ""
description: |-
  This provider allows you to interact with Infisical
---

# infisical Provider

This provider allows you to interact with Infisical

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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `client_id` (String, Sensitive) Machine identity client ID. Used to fetch/modify secrets for a given project
- `client_secret` (String, Sensitive) Machine identity client secret. Used to fetch/modify secrets for a given project
- `host` (String) Used to point the client to fetch secrets from your self hosted instance of Infisical. If not host is provided, https://app.infisical.com is the default host.
- `service_token` (String, Sensitive) (DEPRECATED, USE MACHINE IDENTITY), Used to fetch/modify secrets for a given project
