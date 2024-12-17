<a href="https://terraform.io">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset=".github/terraform_logo_dark.svg">
    <source media="(prefers-color-scheme: light)" srcset=".github/terraform_logo_light.svg">
    <img src=".github/terraform_logo_light.svg" alt="Terraform logo" title="Terraform" align="right" height="50">
  </picture>
</a>

# Terraform Forward Email Provider

The [Forward Email Provider]([https://registry.terraform.io/providers/](https://registry.terraform.io/providers/The-Infra-Company/forwardemail/latest)) enables [Terraform](https://terraform.io) to manage [Forward Email](https://forwardemail.net/) resources.

## Usage

```terraform
terraform {
  required_version = ">= 1.0.0"

  required_providers {
    spacelift = {
      source  = "forwardemail/forwardemail"
      version = ">= 1.0"
    }
  }
}

provider "forwardemail" {
  api_key = "XXXXXXX"
}

resource "forwardemail_domain" "default" {
  name = "infracompany.com"
}
```

