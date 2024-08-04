package user_grant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/object"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/user"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveUserGrant(helper.CtxWithOrgID(ctx, d), &management.RemoveUserGrantRequest{
		GrantId: d.Id(),
		UserId:  d.Get(UserIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete usergrant: %v", err)
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateUserGrant(helper.CtxWithOrgID(ctx, d), &management.UpdateUserGrantRequest{
		GrantId:  d.Id(),
		UserId:   d.Get(UserIDVar).(string),
		RoleKeys: helper.GetOkSetToStringSlice(d, RoleKeysVar),
	})
	if err != nil {
		return diag.Errorf("failed to update usergrant: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddUserGrant(helper.CtxWithOrgID(ctx, d), &management.AddUserGrantRequest{
		UserId:         d.Get(UserIDVar).(string),
		ProjectGrantId: d.Get(projectGrantIDVar).(string),
		ProjectId:      d.Get(projectIDVar).(string),
		RoleKeys:       helper.GetOkSetToStringSlice(d, RoleKeysVar),
	})
	if err != nil {
		return diag.Errorf("failed to create usergrant: %v", err)
	}
	d.SetId(resp.GetUserGrantId())
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.GetUserGrantByID(helper.CtxWithOrgID(ctx, d), &management.GetUserGrantByIDRequest{
		GrantId: helper.GetID(d, grantIDVar),
		UserId:  d.Get(UserIDVar).(string),
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get user grant")
	}
	grant := resp.GetUserGrant()
	set := map[string]interface{}{
		UserIDVar:       grant.GetUserId(),
		RoleKeysVar:     grant.GetRoleKeys(),
		helper.OrgIDVar: grant.GetDetails().GetResourceOwner(),
	}
	if grant.GetProjectId() != "" {
		set[projectIDVar] = grant.GetProjectId()
	}
	if grant.GetProjectGrantId() != "" {
		set[projectGrantIDVar] = grant.GetProjectGrantId()
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of usergrant: %v", k, err)
		}
	}
	d.SetId(grant.GetId())
	return nil
}

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	orgName := d.Get(OrgNameVar).(string)
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	req := &management.ListUserGrantRequest{}

	req.Queries = append(req.Queries, &user.UserGrantQuery{
		Query: &user.UserGrantQuery_OrgNameQuery{
			OrgNameQuery: &user.UserGrantOrgNameQuery{
				OrgName: orgName,
				Method:  object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE,
			},
		},
	})
	resp, err := client.ListUserGrants(ctx, req)

	if err != nil {
		return diag.Errorf("error while getting roles by orgName %s: %v", orgName, err)
	}
	results := []map[string]interface{}{}
	for _, roleGrant := range resp.Result {
		results = append(results, map[string]interface{}{
			UserIDVar:  roleGrant.UserId,
			grantIDVar: roleGrant.Id,
		})
	}
	// If the ID is blank, the datasource is deleted and not usable.
	d.SetId("-")
	return diag.FromErr(d.Set(userGrantDataVar, results))

}
