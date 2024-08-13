package idp_github_es

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_utils"
)

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddGitHubEnterpriseServerProvider(ctx, &admin.AddGitHubEnterpriseServerProviderRequest{
		Name:                  idp_utils.StringValue(d, idp_utils.NameVar),
		ClientId:              idp_utils.StringValue(d, idp_utils.ClientIDVar),
		ClientSecret:          idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		Scopes:                idp_utils.ScopesValue(d),
		ProviderOptions:       idp_utils.ProviderOptionsValue(d),
		AuthorizationEndpoint: idp_utils.StringValue(d, AuthorizationEndpointVar),
		TokenEndpoint:         idp_utils.StringValue(d, TokenEndpointVar),
		UserEndpoint:          idp_utils.StringValue(d, UserEndpointVar),
	})
	if err != nil {
		return diag.Errorf("failed to create idp: %v", err)
	}
	d.SetId(resp.GetId())
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.UpdateGitHubEnterpriseServerProvider(ctx, &admin.UpdateGitHubEnterpriseServerProviderRequest{
		Id:                    d.Id(),
		Name:                  idp_utils.StringValue(d, idp_utils.NameVar),
		ClientId:              idp_utils.StringValue(d, idp_utils.ClientIDVar),
		ClientSecret:          idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		Scopes:                idp_utils.ScopesValue(d),
		ProviderOptions:       idp_utils.ProviderOptionsValue(d),
		AuthorizationEndpoint: idp_utils.StringValue(d, AuthorizationEndpointVar),
		TokenEndpoint:         idp_utils.StringValue(d, TokenEndpointVar),
		UserEndpoint:          idp_utils.StringValue(d, UserEndpointVar),
	})
	if err != nil {
		return diag.Errorf("failed to update idp: %v", err)
	}
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.GetProviderByID(ctx, &admin.GetProviderByIDRequest{Id: helper.GetID(d, idp_utils.IdpIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get idp")
	}
	idp := resp.GetIdp()
	cfg := idp.GetConfig()
	specificCfg := cfg.GetGithubEs()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		idp_utils.NameVar:              idp.GetName(),
		idp_utils.ClientIDVar:          specificCfg.GetClientId(),
		idp_utils.ClientSecretVar:      idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		idp_utils.ScopesVar:            specificCfg.GetScopes(),
		idp_utils.IsLinkingAllowedVar:  generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar: generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:    generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:      generalCfg.GetIsAutoUpdate(),
		AuthorizationEndpointVar:       specificCfg.GetAuthorizationEndpoint(),
		TokenEndpointVar:               specificCfg.GetTokenEndpoint(),
		UserEndpointVar:                specificCfg.GetUserEndpoint(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)
	return nil
}
