package idp_oauth

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddGenericOAuthProvider(ctx, &admin.AddGenericOAuthProviderRequest{
		Name:                  idp_utils.StringValue(d, idp_utils.NameVar),
		ClientId:              idp_utils.StringValue(d, idp_utils.ClientIDVar),
		ClientSecret:          idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		AuthorizationEndpoint: idp_utils.StringValue(d, AuthorizationEndpointVar),
		TokenEndpoint:         idp_utils.StringValue(d, TokenEndpointVar),
		UserEndpoint:          idp_utils.StringValue(d, UserEndpointVar),
		IdAttribute:           idp_utils.StringValue(d, IdAttributeVar),
		Scopes:                idp_utils.ScopesValue(d),
		ProviderOptions:       idp_utils.ProviderOptionsValue(d),
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
	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.UpdateGenericOAuthProvider(ctx, &admin.UpdateGenericOAuthProviderRequest{
		Id:                    d.Id(),
		Name:                  idp_utils.StringValue(d, idp_utils.NameVar),
		ClientId:              idp_utils.StringValue(d, idp_utils.ClientIDVar),
		ClientSecret:          idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		AuthorizationEndpoint: idp_utils.StringValue(d, AuthorizationEndpointVar),
		TokenEndpoint:         idp_utils.StringValue(d, TokenEndpointVar),
		UserEndpoint:          idp_utils.StringValue(d, UserEndpointVar),
		IdAttribute:           idp_utils.StringValue(d, IdAttributeVar),
		Scopes:                idp_utils.ScopesValue(d),
		ProviderOptions:       idp_utils.ProviderOptionsValue(d),
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
	client, err := helper.GetAdminClient(ctx, clientinfo)
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
	specificCfg := cfg.GetOauth()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		idp_utils.NameVar:              idp.GetName(),
		idp_utils.ClientIDVar:          specificCfg.GetClientId(),
		idp_utils.ClientSecretVar:      idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		idp_utils.ScopesVar:            specificCfg.GetScopes(),
		AuthorizationEndpointVar:       specificCfg.GetAuthorizationEndpoint(),
		TokenEndpointVar:               specificCfg.GetTokenEndpoint(),
		UserEndpointVar:                specificCfg.GetUserEndpoint(),
		IdAttributeVar:                 specificCfg.GetIdAttribute(),
		idp_utils.IsLinkingAllowedVar:  generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar: generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:    generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:      generalCfg.GetIsAutoUpdate(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)
	return nil
}
