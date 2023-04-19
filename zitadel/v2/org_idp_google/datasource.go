package org_idp_google

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a Google IdP of the organization.",
		Schema: map[string]*schema.Schema{
			org_idp_utils.OrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
			},
			idp_utils.IdpIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			idp_utils.NameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the IDP",
			},
			idp_utils.ClientIDVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "client id generated by the identity provider",
			},
			idp_utils.ClientSecretVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "client secret generated by the identity provider",
				Sensitive:   true,
			},
			idp_utils.ScopesVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "the scopes requested by ZITADEL during the request on the identity provider",
			},
			idp_utils.IsLinkingAllowedVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "enable if users should be able to link an existing ZITADEL user with an external account",
			},
			idp_utils.IsCreationAllowedVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "enable if users should be able to create a new account in ZITADEL when using an external account",
			},
			idp_utils.IsAutoCreationVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "enable if a new account in ZITADEL should be created automatically when login with an external account",
			},
			idp_utils.IsAutoUpdateVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "enable if a the ZITADEL account fields should be updated automatically on each login",
			},
		},
		ReadContext: read,
		Importer:    &schema.ResourceImporter{StateContext: org_idp_utils.ImportIDPWithOrgAndClientSecret},
	}
}
