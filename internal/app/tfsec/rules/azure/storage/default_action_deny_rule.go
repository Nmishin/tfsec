package storage

import (
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/rules/azure/storage"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
	"github.com/aquasecurity/tfsec/pkg/rule"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		LegacyID: "AZU012",
		BadExample: []string{`
 resource "azurerm_storage_account_network_rules" "bad_example" {
   
   default_action             = "Allow"
   ip_rules                   = ["127.0.0.1"]
   virtual_network_subnet_ids = [azurerm_subnet.test.id]
   bypass                     = ["Metrics"]
 }
 `},
		GoodExample: []string{`
 resource "azurerm_storage_account_network_rules" "good_example" {
   
   default_action             = "Deny"
   ip_rules                   = ["127.0.0.1"]
   virtual_network_subnet_ids = [azurerm_subnet.test.id]
   bypass                     = ["Metrics"]
 }
 `},
		Links: []string{
			"https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/storage_account_network_rules#default_action",
		},
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"azurerm_storage_account", "azurerm_storage_account_network_rules"},
		Base:           storage.CheckDefaultActionDeny,
		CheckTerraform: func(resourceBlock block.Block, _ block.Module) (results rules.Results) {

			if resourceBlock.IsResourceType("azurerm_storage_account") {
				if resourceBlock.MissingChild("network_rules") {
					return
				}
				resourceBlock = resourceBlock.GetBlock("network_rules")
			}

			defaultAction := resourceBlock.GetAttribute("default_action")
			if defaultAction.IsNotNil() && defaultAction.Equals("Allow", block.IgnoreCase) {
				results.Add("Resource defines a default_action of Allow. It should be Deny.", defaultAction)
			}

			return results
		},
	})
}
