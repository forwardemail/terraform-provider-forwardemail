package forwardemail

import (
	"context"
	"fmt"
	"strings"

	"github.com/forwardemail/forwardemail-api-go/forwardemail"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlias() *schema.Resource {
	return &schema.Resource{
		Description: "A resource to create Forward Email domain aliases.",
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Fully qualified domain name (FQDN).",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Alias name.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable to disable this alias.",
			},
			"recipient_verification": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable to require recipients to click an email verification link for emails to flow through.",
			},
			"recipients": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of recipients as valid email addresses, fully-qualified domain names (FQDN), IP addresses, or webhook URL's.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alias description.",
			},
			"labels": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of labels.",
			},
		},
		CreateContext: resourceAliasCreate,
		ReadContext:   resourceAliasRead,
		UpdateContext: resourceAliasUpdate,
		DeleteContext: resourceAliasDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAliasImport,
		},
	}
}
func resourceAliasCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*forwardemail.Client)
	if !ok {
		return diag.Errorf("meta is not of type *forwardemail.Client")
	}

	domain := d.Get("domain").(string)
	name := d.Get("name").(string)

	params := forwardemail.AliasParameters{
		HasRecipientVerification: toBool(d.Get("recipient_verification")),
		IsEnabled:                toBool(d.Get("enabled")),
		Recipients:               toSliceOfStrings(toChanges(nil, d.Get("recipients"))),
		Labels:                   toSliceOfStrings(toChanges(nil, d.Get("labels"))),
		Description:              d.Get("description").(string),
	}

	alias, err := client.CreateAlias(domain, name, params)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range map[string]interface{}{
		"domain":                 alias.Domain.Name,
		"recipient_verification": alias.HasRecipientVerification,
		"enabled":                alias.IsEnabled,
		"recipients":             alias.Recipients,
		"labels":                 alias.Labels,
		"description":            alias.Description,
	} {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(name)

	return nil
}

func resourceAliasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*forwardemail.Client)
	if !ok {
		return diag.Errorf("meta is not of type *forwardemail.Client")
	}

	domain := d.Get("domain").(string)
	name := d.Get("name").(string)

	alias, err := client.GetAlias(domain, name)
	if err != nil {
		return diag.FromErr(err)
	}

	for k, v := range map[string]interface{}{
		"domain":                 alias.Domain.Name,
		"recipient_verification": alias.HasRecipientVerification,
		"enabled":                alias.IsEnabled,
		"recipients":             alias.Recipients,
		"labels":                 alias.Labels,
		"description":            alias.Description,
	} {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceAliasUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*forwardemail.Client)
	if !ok {
		return diag.Errorf("meta is not of type *forwardemail.Client")
	}

	domain := d.Get("domain").(string)
	name := d.Id()

	// N.B.: we can't use d.GetChange because Forward Email API is working not as expected.
	// So instead of passing only changed parameters we need to pass all of them.
	params := forwardemail.AliasParameters{}
	params.HasRecipientVerification = toBool(toChange(nil, d.Get("recipient_verification")))
	params.IsEnabled = toBool(toChange(nil, d.Get("enabled")))
	params.Recipients = toSliceOfStrings(toChanges(nil, d.Get("recipients")))
	params.Labels = toSliceOfStrings(toChanges(nil, d.Get("labels")))

	_, err := client.UpdateAlias(domain, name, params)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAliasDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, ok := meta.(*forwardemail.Client)
	if !ok {
		return diag.Errorf("meta is not of type *forwardemail.Client")
	}

	domain := d.Get("domain").(string)
	name := d.Id()

	err := client.DeleteAlias(domain, name)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func updateResourceDataFromAlias(d *schema.ResourceData, alias *forwardemail.Alias) error {
	d.SetId(alias.Id)
	if err := d.Set("name", alias.Name); err != nil {
		return err
	}
	if err := d.Set("domain", alias.Domain.Name); err != nil {
		return err
	}
	return nil
}

func resourceAliasImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client, ok := meta.(*forwardemail.Client)
	if !ok {
		return nil, fmt.Errorf("meta is not of type *forwardemail.Client")
	}

	importID := d.Id()

	// If the import ID contains "@", treat it as ALIAS@DOMAIN.
	if strings.Contains(importID, "@") {
		parts := strings.SplitN(importID, "@", 2)
		aliasName := parts[0]
		domain := parts[1]

		if aliasName == "" || domain == "" {
			return nil, fmt.Errorf("invalid import format %q, expected ALIAS@DOMAIN", importID)
		}

		alias, err := client.GetAlias(domain, aliasName)
		if err != nil {
			return nil, fmt.Errorf("error fetching forwardemail_alias %s@%s: %w", aliasName, domain, err)
		}

		err = updateResourceDataFromAlias(d, alias)
		if err != nil {
			return nil, fmt.Errorf("updateResourceDataFromAlias failed: %w", err)
		}

		return []*schema.ResourceData{d}, nil
	}

	// Otherwise, treat the import ID as DOMAIN_ID:ALIAS_ID.
	parts := strings.SplitN(importID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("invalid import format %q, expected ALIAS@DOMAIN or DOMAIN_ID:ALIAS_ID", importID)
	}

	domainID := parts[0]
	aliasID := parts[1]

	alias, err := client.GetAlias(domainID, aliasID)
	if err != nil {
		return nil, fmt.Errorf("error fetching forwardemail_alias %s:%s: %w", domainID, aliasID, err)
	}

	err = updateResourceDataFromAlias(d, alias)
	if err != nil {
		return nil, fmt.Errorf("updateResourceDataFromAlias failed: %w", err)
	}

	return []*schema.ResourceData{d}, nil
}

// toSliceOfStrings converts slice of interfaces into pointer to slice of strings.
func toSliceOfStrings(vs []interface{}) *[]string {
	var stringSlice []string
	for _, v := range vs {
		if str, ok := v.(string); ok {
			stringSlice = append(stringSlice, str)
		}
	}

	return &stringSlice
}

// toChanges converts interface into slice of interfaces.
func toChanges(p, c interface{}) []interface{} {
	if cmp.Equal(p, c) {
		return nil
	}

	switch v := c.(type) {
	case []interface{}:
		return v
	}

	return nil
}
