package forwardemail

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceForwardemailAlias(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccAliasResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("forwardemail_domain.test", "name", "example.com"),
					resource.TestCheckResourceAttr("forwardemail_alias.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("forwardemail_alias.test", "name", "postmaster"),
				),
			},
			{
				ResourceName:      "forwardemail_alias.test",
				ImportState:       true,
				ImportStateId:     "postmaster@example.com",
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAliasResourceConfig() string {
	return `
resource "forwardemail_domain" "test" {
  name = "example.com"
}

resource "forwardemail_alias" "test" {
  name       = "postmaster"
  domain     = forwardemail_domain.test.name
  recipients = ["user@example.com"]
}
`
}
