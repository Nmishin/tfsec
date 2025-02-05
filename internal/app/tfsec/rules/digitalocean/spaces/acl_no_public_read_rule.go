package spaces

import (
	"github.com/aquasecurity/defsec/rules"
	"github.com/aquasecurity/defsec/rules/digitalocean/spaces"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
	"github.com/aquasecurity/tfsec/pkg/rule"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		LegacyID: "DIG005",
		BadExample: []string{`
 resource "digitalocean_spaces_bucket" "bad_example" {
   name   = "public_space"
   region = "nyc3"
   acl    = "public-read"
 }
 
 resource "digitalocean_spaces_bucket_object" "index" {
   region       = digitalocean_spaces_bucket.bad_example.region
   bucket       = digitalocean_spaces_bucket.bad_example.name
   key          = "index.html"
   content      = "<html><body><p>This page is empty.</p></body></html>"
   content_type = "text/html"
   acl          = "public-read"
 }
 `},
		GoodExample: []string{`
 resource "digitalocean_spaces_bucket" "good_example" {
   name   = "private_space"
   region = "nyc3"
   acl    = "private"
 }
   
 resource "digitalocean_spaces_bucket_object" "index" {
   region       = digitalocean_spaces_bucket.good_example.region
   bucket       = digitalocean_spaces_bucket.good_example.name
   key          = "index.html"
   content      = "<html><body><p>This page is empty.</p></body></html>"
   content_type = "text/html"
 }
 `},
		Links: []string{
			"https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/spaces_bucket#acl",
			"https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/spaces_bucket_object#acl",
		},
		RequiredTypes:  []string{"resource"},
		RequiredLabels: []string{"digitalocean_spaces_bucket", "digitalocean_spaces_bucket_object"},
		Base:           spaces.CheckAclNoPublicRead,
		CheckTerraform: func(resourceBlock block.Block, _ block.Module) (results rules.Results) {

			if resourceBlock.HasChild("acl") {
				aclAttr := resourceBlock.GetAttribute("acl")
				if aclAttr.Equals("public-read", block.IgnoreCase) {
					results.Add("Resource has a publicly readable acl.", aclAttr)
				}
			}
			return results
		},
	})
}
