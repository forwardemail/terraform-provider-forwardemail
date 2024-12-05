package forwardemail

import (
	"context"

	"github.com/abagayev/go-forwardemail/forwardemail"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Description: "The API key for Forward Email management",
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("FORWARDEMAIL_API_KEY", nil),
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"forwardemail_domain": resourceDomain(),
			"forwardemail_alias":  resourceAlias(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"forwardemail_account": dataSourceAccount(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(c context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client := forwardemail.NewClient(forwardemail.ClientOptions{
		ApiKey: d.Get("api_key").(string),
	})

	return client, diag.Diagnostics{}
}
