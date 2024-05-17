---
page_title: "zitadel_idp_saml Data Source - terraform-provider-zitadel"
subcategory: ""
description: |-
  Datasource representing a SAML IDP on the instance.
---

# zitadel_idp_saml (Data Source)

Datasource representing a SAML IDP on the instance.

## Example Usage

```terraform
data "zitadel_idp_saml" "default" {
  id = "123456789012345678"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The ID of this resource.

### Read-Only

- `binding` (String) The binding
- `is_auto_creation` (Boolean) enabled if a new account in ZITADEL are created automatically on login with an external account
- `is_auto_update` (Boolean) enabled if a the ZITADEL account fields are updated automatically on each login
- `is_creation_allowed` (Boolean) enabled if users are able to create a new account in ZITADEL when using an external account
- `is_linking_allowed` (Boolean) enabled if users are able to link an existing ZITADEL user with an external account
- `metadata_xml` (String) The metadata XML as plain string
- `name` (String) Name of the IDP
- `with_signed_request` (String) Whether the SAML IDP requires signed requests