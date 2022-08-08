
resource zitadel_org_member org_member {
  depends_on = [zitadel_org.org, zitadel_human_user.human_user]

  org_id  = zitadel_org.org.id
  user_id = zitadel_human_user.human_user.id
  roles   = ["ORG_OWNER"]
}