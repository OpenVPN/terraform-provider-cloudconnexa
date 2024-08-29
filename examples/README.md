# Important
As for now currently examples inside this folder are just to show how some specific resources could look like, or how they could be implemented.
Examples inside this folder are not intended to provide complete picture or full working example.
Please be careful to understand what resources you copy from this folder and give it a double check if it will suite your needs (or not) and adjust this code to your specific needs.

## Some generic consideration
If you have small amount of users it is okay to use local users (created inside CloudConnexa), examples in this folder will help you to achieve this.
If you have quite big amount of users - at some point it would be beneficial to look into managing users via LDAP or SAML. It is not possible (yet) configure LDAP and SAML settings via Terraform provider, this functionality will be added in the future releases.
If you have configured LDAP or SAML via CloudConnexa Admin UI - then using Terraform resource "cloudconnexa_user" will be not needed.
