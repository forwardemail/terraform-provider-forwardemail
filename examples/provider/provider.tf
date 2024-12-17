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
