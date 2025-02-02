---
page_title: "glesys_emailaccount Resource - terraform-provider-glesys"
subcategory: ""
description: |-
  Create a GleSYS Email Account.
---
# glesys_emailaccount (Resource)
Create a GleSYS Email Account.
## Example Usage
```terraform
# Setup an email account

resource "glesys_emailaccount" "bob" {
  emailaccount       = "bob@example.com"
  password           = "SecretPassword123"
  autorespond        = "yes"
  autorespondmessage = "I'm away."
  quotaingib         = 2
}
resource "glesys_emailaccount" "alice" {
  emailaccount  = "alice@example.com"
  password      = "PasswordSecret321"
  antispamlevel = 5
}
```
<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `emailaccount` (String) Email account name

### Optional

- `antispamlevel` (Number) Email Account antispam level. `0-5`
- `antivirus` (String) Email Account enable Antivirus. `yes/no`
- `autorespond` (String) Email Account Autoresponse. `yes/no`
- `autorespondmessage` (String) Email Account Autoresponse message.
- `password` (String) Email Account password
- `quotaingib` (Number) Email Account Quota (GiB)
- `rejectspam` (String) Email Account Reject spam setting. `yes/no`

### Read-Only

- `autorespondsaveemail` (String) Email Account Save emails on autorespond.
- `created` (String) Email Account created date
- `displayname` (String) Email Account displayname
- `id` (String) The ID of this resource.
- `modified` (String) Email Account modification date
## Import
Import is supported using the following syntax:
```shell
# Email account import.
$ terraform import glesys_emailaccount.alice alice@example.com
```
