package compute

import (
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/rules/google/compute"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
	"github.com/aquasecurity/tfsec/pkg/rule"
	"github.com/zclconf/go-cty/cty"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		BadExample: []string{`
 resource "google_compute_project_metadata" "default" {
   metadata = {
 	enable-oslogin = false
   }
 }
 `},
		GoodExample: []string{`
 resource "google_compute_project_metadata" "default" {
   metadata = {
     enable-oslogin = true
   }
 }
 `},
		Links: []string{
			"https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_project_metadata#",
		},
		RequiredTypes: []string{
			"resource",
		},
		RequiredLabels: []string{
			"google_compute_project_metadata",
		},
		Base: compute.CheckProjectLevelOslogin,
		CheckTerraform: func(resourceBlock block.Block, _ block.Module) (results rules.Results) {
			metadataAttr := resourceBlock.GetAttribute("metadata")
			val := metadataAttr.MapValue("enable-oslogin")
			if val.Type() == cty.NilType {
				results.Add("Resource'%s' has OS Login disabled by default", resourceBlock)
				return
			}
			if val.Type() == cty.Bool && val.False() {
				results.Add("Resource'%s' has OS Login explicitly disabled", resourceBlock)
				return
			}
			return results
		},
	})
}
