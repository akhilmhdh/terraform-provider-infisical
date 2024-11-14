---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "infisical_integration_gcp_secret_manager Resource - terraform-provider-infisical"
subcategory: ""
description: |-
  Create GCP Secret Manager integration & save to Infisical. Only Machine Identity authentication is supported for this data source
---

# infisical_integration_gcp_secret_manager (Resource)

Create GCP Secret Manager integration & save to Infisical. Only Machine Identity authentication is supported for this data source

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
  client_id     = "<machine-identity-client-id>"
  client_secret = "<machine-identity-client-secret>"
}

variable "service_account_json" {
  type        = string
  description = "Google Cloud service account JSON key"
}



resource "infisical_integration_gcp_secret_manager" "gcp-integration" {
  project_id           = "your-project-id"
  service_account_json = var.service_account_json
  environment          = "dev"
  secret_path          = "/"
  gcp_project_id       = "gcp-project-id"

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment` (String) The slug of the environment to sync to GCP Secret Manager (prod, dev, staging, etc).
- `gcp_project_id` (String) The ID of the GCP project.
- `project_id` (String) The ID of your Infisical project.
- `secret_path` (String) The secret path in Infisical to sync secrets from.
- `service_account_json` (String, Sensitive) Service account json for the GCP project.

### Optional

- `options` (Attributes) Integration options (see [below for nested schema](#nestedatt--options))

### Read-Only

- `integration_auth_id` (String) The ID of the integration auth, used internally by Infisical.
- `integration_id` (String) The ID of the integration, used internally by Infisical.

<a id="nestedatt--options"></a>
### Nested Schema for `options`

Optional:

- `secret_prefix` (String) The prefix to add to the secret name in GCP Secret Manager.
- `secret_suffix` (String) The suffix to add to the secret name in GCP Secret Manager.