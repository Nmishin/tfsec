package database

import (
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/rules/azure/database"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
	"github.com/aquasecurity/tfsec/pkg/rule"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		BadExample: []string{`
 resource "azurerm_resource_group" "example" {
   name     = "example-resources"
   location = "West Europe"
 }
 
 resource "azurerm_postgresql_server" "example" {
   name                = "example-psqlserver"
   location            = azurerm_resource_group.example.location
   resource_group_name = azurerm_resource_group.example.name
 
   administrator_login          = "psqladminun"
   administrator_login_password = "H@Sh1CoR3!"
 
   sku_name   = "GP_Gen5_4"
   version    = "9.6"
   storage_mb = 640000
 }
 `},
		GoodExample: []string{`
 resource "azurerm_resource_group" "example" {
   name     = "example-resources"
   location = "West Europe"
 }
 
 resource "azurerm_postgresql_server" "example" {
   name                = "example-psqlserver"
   location            = azurerm_resource_group.example.location
   resource_group_name = azurerm_resource_group.example.name
 
   administrator_login          = "psqladminun"
   administrator_login_password = "H@Sh1CoR3!"
 
   sku_name   = "GP_Gen5_4"
   version    = "9.6"
   storage_mb = 640000
 }
 
 resource "azurerm_postgresql_configuration" "example" {
 	name                = "log_connections"
 	resource_group_name = azurerm_resource_group.example.name
 	server_name         = azurerm_postgresql_server.example.name
 	value               = "on"
   }
   
   `},
		Links: []string{
			"https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/postgresql_configuration",
			"https://docs.microsoft.com/en-us/azure/postgresql/concepts-server-logs#configure-logging"},
		RequiredTypes: []string{
			"resource",
		},
		RequiredLabels: []string{
			"azurerm_postgresql_server",
		},
		Base: database.CheckPostgresConfigurationLogConnections,
		CheckTerraform: func(resourceBlock block.Block, module block.Module) (results rules.Results) {
			referencingBlocks := module.GetReferencingResources(resourceBlock, "azurerm_postgresql_configuration", "server_name")
			for _, refBlock := range referencingBlocks {
				if nameAttr := refBlock.GetAttribute("name"); nameAttr.IsNotNil() && nameAttr.Equals("log_connections") {
					if valAttr := refBlock.GetAttribute("value"); valAttr.IsNotNil() && valAttr.Equals("on", block.IgnoreCase) {
						return
					}
				}
			}
			results.Add("Resource does not have a corresponding log configuration enabling 'log_connections'", resourceBlock)
			return results
		},
	})
}
