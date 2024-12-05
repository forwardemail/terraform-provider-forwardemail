package forwardemail

import (
	"context"

	"github.com/abagayev/go-forwardemail/forwardemail"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDomain() *schema.Resource {
	return &schema.Resource{
		Description: "A resource to create Forward Email domains.",
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Fully qualified domain name (FQDN) or IP address.",
			},
			"adult_content_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable Spam Scanner adult content protection on this domain.",
			},
			"phishing_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable Spam Scanner phishing protection on this domain.",
			},
			"executable_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable Spam Scanner executable protection on this domain.",
			},
			"virus_protection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable Spam Scanner virus protection on this domain.",
			},
			"recipient_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Global domain default for whether to require alias recipients to click an email verification link for emails to flow through.",
			},
		},
		CreateContext: resourceDomainCreate,
		ReadContext:   resourceDomainRead,
		UpdateContext: resourceDomainUpdate,
		DeleteContext: resourceDomainDelete,
	}
}

func resourceDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*forwardemail.Client)
	if !ok {
		return diag.Errorf("could not get forwardemail client")
	}

	name := d.Get("name").(string)

	params := forwardemail.DomainParameters{
		HasAdultContentProtection: toBool(d.Get("adult_content_protection")),
		HasPhishingProtection:     toBool(d.Get("phishing_protection")),
		HasExecutableProtection:   toBool(d.Get("executable_protection")),
		HasVirusProtection:        toBool(d.Get("virus_protection")),
		HasRecipientVerification:  toBool(d.Get("recipient_verification")),
	}

	domain, err := client.CreateDomain(name, params)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range map[string]interface{}{
		"adult_content_protection": domain.HasAdultContentProtection,
		"phishing_protection":      domain.HasPhishingProtection,
		"executable_protection":    domain.HasExecutableProtection,
		"virus_protection":         domain.HasVirusProtection,
		"recipient_verification":   domain.HasRecipientVerification,
	} {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(name)

	return nil
}

func resourceDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*forwardemail.Client)
	if !ok {
		return diag.Errorf("could not get forwardemail client")
	}

	name := d.Id()

	domain, err := client.GetDomain(name)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range map[string]interface{}{
		"adult_content_protection": domain.HasAdultContentProtection,
		"phishing_protection":      domain.HasPhishingProtection,
		"executable_protection":    domain.HasExecutableProtection,
		"virus_protection":         domain.HasVirusProtection,
		"recipient_verification":   domain.HasRecipientVerification,
	} {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceDomainUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*forwardemail.Client)
	if !ok {
		return diag.Errorf("could not get forwardemail client")
	}

	name := d.Id()

	params := forwardemail.DomainParameters{}
	params.HasAdultContentProtection = toBool(toChange(d.GetChange("adult_content_protection")))
	params.HasPhishingProtection = toBool(toChange(d.GetChange("phishing_protection")))
	params.HasExecutableProtection = toBool(toChange(d.GetChange("executable_protection")))
	params.HasVirusProtection = toBool(toChange(d.GetChange("virus_protection")))
	params.HasRecipientVerification = toBool(toChange(d.GetChange("recipient_verification")))

	_, err := client.UpdateDomain(name, params)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*forwardemail.Client)
	if !ok {
		return diag.Errorf("could not get forwardemail client")
	}

	name := d.Id()

	err := client.DeleteDomain(name)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// toBool returns a pointer to the bool value passed in.
func toBool(v interface{}) *bool {
	if b, ok := v.(bool); ok {
		return &b
	}

	return nil
}

func toChange(p, c interface{}) interface{} {
	if cmp.Equal(p, c) {
		return nil
	}

	return c
}
