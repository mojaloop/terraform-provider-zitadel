package idp_github_es_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils/test_utils"
)

func TestAccZITADELInstanceIdPGitHubES(t *testing.T) {
	resourceName := "zitadel_idp_github_es"
	frame, err := test_utils.NewInstanceTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunBasicLifecyleTest(t, frame, func(name, secret string) string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  name                   = "%s"
  client_id              = "aclientid"
  client_secret          = "%s"
  scopes                 = ["two", "scopes"]
  authorization_endpoint = "https://auth.endpoint"
  token_endpoint         = "https://token.endpoint"
  user_endpoint          = "https://user.endpoint"
  is_linking_allowed     = false
  is_creation_allowed    = true
  is_auto_creation       = false
  is_auto_update         = true
}`, resourceName, frame.UniqueResourcesID, name, secret)
	}, "client_secret")
}
