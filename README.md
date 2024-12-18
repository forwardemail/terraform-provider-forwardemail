<a href="https://terraform.io">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset=".github/terraform_logo_dark.svg">
    <source media="(prefers-color-scheme: light)" srcset=".github/terraform_logo_light.svg">
    <img src=".github/terraform_logo_light.svg" alt="Terraform logo" title="Terraform" align="right" height="50">
  </picture>
</a>

# Terraform Forward Email Provider

The [Forward Email Provider]([https://registry.terraform.io/providers/](https://registry.terraform.io/providers/forwardemail/forwardemail/latest)) enables [Terraform](https://terraform.io) to manage [Forward Email](https://forwardemail.net/) resources.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.12
- [Go](https://golang.org/doc/install) >= 1.13 (to build the provider plugin)

## Building the Provider

Clone repository:

```sh
git clone https://github.com/forwardemail/terraform-provider-forwardemail.git
```

Enter the provider directory and build the provider:

```sh
cd terraform-provider-forwardemail
make build
```

To use a released provider in your Terraform environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider. To specify a particular provider version when installing released providers, see the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions above), follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) After placing the custom-built provider into your plugins directory, run `terraform init` to initialize it.

## Examples

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
  name = "example_domain.com"
}
```

## Testing the Provider

In order to test the provider, you can run `go test`.

```sh
go test -v ./...
```

In order to run the full suite of Acceptance tests, you'll need a paid Forward Email account.

You'll also need to set the API key environment variable:

```sh
export FORWARDEMAIL_API_KEY=<API_KEY>
make testacc
```

## Contributing

For bug reports & feature requests, please use the [issue tracker](https://github.com/forwardemail/terraform-provider-forwardemail/issues).

PRs are welcome! We follow the typical "fork-and-pull" Git workflow.
 1. **Fork** the repo on GitHub
 2. **Clone** the project to your own machine
 3. **Commit** changes to your own branch
 4. **Push** your work back up to your fork
 5. Submit a **Pull Request** so that we can review your changes

> [!TIP]
> Be sure to merge the latest changes from "upstream" before making a pull request!

### Many Thanks to Our Contributors

<a href="https://github.com/forwardemail/terraform-provider-forwardemail/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=forwardemail/terraform-provider-forwardemail&max=24" />
</a>

