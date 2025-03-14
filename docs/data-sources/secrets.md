---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "infisical_secrets Data Source - terraform-provider-infisical"
subcategory: ""
description: |-
  Interact with Infisical secrets
---

# infisical_secrets (Data Source)

Interact with Infisical secrets

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
  host = "https://app.infisical.com" # Only required if using self hosted instance of Infisical, default is https://app.infisical.com
  auth = {
    universal = {
      client_id     = "<machine-identity-client-id>"
      client_secret = "<machine-identity-client-secret>"
    }
  }
}

data "infisical_secrets" "common_secrets" {
  env_slug     = "dev"
  workspace_id = "<project id>" // project ID
  folder_path  = "/"
}

output "all-project-secrets" {
  value = nonsensitive(data.infisical_secrets.common_secrets.secrets["SECRET-NAME"].value)
}

output "all-project-secrets" {
  value = nonsensitive(data.infisical_secrets.common_secrets.secrets["SECRET-NAME"].comment)
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `env_slug` (String) The environment from where secrets should be fetched from
- `folder_path` (String) The path to the folder from where secrets should be fetched from

### Optional

- `workspace_id` (String) The Infisical project ID (Required for Machine Identity auth, and service tokens with multiple scopes)

### Read-Only

- `secrets` (Attributes Map) (see [below for nested schema](#nestedatt--secrets))

<a id="nestedatt--secrets"></a>
### Nested Schema for `secrets`

Read-Only:

- `comment` (String) The secret comment
- `secret_type` (String) The secret type (shared or personal)
- `value` (String, Sensitive) The secret value
