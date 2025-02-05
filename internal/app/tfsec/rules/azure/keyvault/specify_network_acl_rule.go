package keyvault

import (
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/rules/azure/keyvault"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
	"github.com/aquasecurity/tfsec/pkg/rule"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		LegacyID: "AZU020",
		BadExample: []string{`
 resource "azurerm_key_vault" "bad_example" {
     name                        = "examplekeyvault"
     location                    = azurerm_resource_group.bad_example.location
     enabled_for_disk_encryption = true
     soft_delete_retention_days  = 7
     purge_protection_enabled    = false
 }
 `},
		GoodExample: []string{`
 resource "azurerm_key_vault" "good_example" {
     name                        = "examplekeyvault"
     location                    = azurerm_resource_group.good_example.location
     enabled_for_disk_encryption = true
     soft_delete_retention_days  = 7
     purge_protection_enabled    = false
 
     network_acls {
         bypass = "AzureServices"
         default_action = "Deny"
     }
 }
 `},
		Links: []string{
			"https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/key_vault#network_acls",
		},
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"azurerm_key_vault"},
		Base:           keyvault.CheckSpecifyNetworkAcl,
		CheckTerraform: func(resourceBlock block.Block, _ block.Module) (results rules.Results) {

			defaultActionAttr := resourceBlock.GetNestedAttribute("network_acls.default_action")
			if defaultActionAttr.IsNil() {
				results.Add("Resource specifies does not specify a network acl block with default action.", resourceBlock)
				return
			}

			if !defaultActionAttr.Equals("Deny") {
				results.Add("Resource specifies does not specify a network acl block.", defaultActionAttr)
			}

			return results
		},
	})
}
