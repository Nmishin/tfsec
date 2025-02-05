package apigateway

import (
	"github.com/aquasecurity/defsec/rules/aws/apigateway"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
	"github.com/aquasecurity/tfsec/pkg/rule"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		BadExample: []string{`
 resource "aws_api_gateway_rest_api" "test" {
	
 }

 resource "aws_api_gateway_stage" "bad_example" {
   stage_name    = "prod"
   rest_api_id   = aws_api_gateway_rest_api.test.id
   deployment_id = aws_api_gateway_deployment.test.id
   xray_tracing_enabled = false
 }
 `},
		GoodExample: []string{`
 resource "aws_api_gateway_rest_api" "test" {
	
 }

 resource "aws_api_gateway_stage" "good_example" {
   stage_name    = "prod"
   rest_api_id   = aws_api_gateway_rest_api.test.id
   deployment_id = aws_api_gateway_deployment.test.id
   xray_tracing_enabled = true
 }
 `},
		Links: []string{
			"https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/api_gateway_stage#xray_tracing_enabled",
		},
		Base: apigateway.CheckEnableTracing,
	})
}
