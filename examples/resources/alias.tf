resource "forwardemail_alias" "infracompany" {
  name   = "Infra Company"
  domain = "infracompany.com"

  recipients = ["james@rhodes.com"]
}
