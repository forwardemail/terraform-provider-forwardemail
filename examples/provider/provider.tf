terraform {
  required_version = ">= 1.0.0"

  required_providers {
    spacelift = {
      source  = "the-infra-company/forwardemail"
      version = ">= 1.0"
    }
  }
}

provider "forwardemail" {
  api_key = "XXXXXXX"
}
