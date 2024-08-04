package user_grant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "represents role grants",
		Schema: map[string]*schema.Schema{
			grantIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the usergrant",
			},
			OrgNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				Computed:    true,
			},
			roleNamesVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A set of all roles for a user.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			projectNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the org.",
				Computed:    true,
			},
			roleStatusVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status of role",
				Computed:    true,
			},
			userNameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A domain of the org.",
				Computed:    true,
			},
		},
		ReadContext: read,
	}
}

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "represents role grants",
		Schema: map[string]*schema.Schema{
			OrgNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
			},
			userGrantDataVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all grantid and userids.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						grantIDVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "grantID",
						},
						UserIDVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "userid",
						},
					},
				},
			},
		},
		ReadContext: list,
	}
}
