package test

import (
	"testing"

	"github.com/tfsec/tfsec/internal/app/tfsec/checks"
	"github.com/tfsec/tfsec/internal/app/tfsec/scanner"
)

func Test_AZURDPAccessNotAllowedFromInternet(t *testing.T) {

	var tests = []struct {
		name                  string
		source                string
		mustIncludeResultCode scanner.RuleCode
		mustExcludeResultCode scanner.RuleCode
	}{
		{
			name: "check RDP access from * causes a failure",
			source: `
resource "azurerm_network_security_rule" "bad_example" {
     name                        = "bad_example_security_rule"
     direction                   = "Inbound"
     access                      = "Allow"
     protocol                    = "TCP"
     source_port_range           = "*"
     destination_port_range      = ["3389"]
     source_address_prefix       = "*"
     destination_address_prefix  = "*"
}
`,
			mustIncludeResultCode: checks.AZURDPAccessNotAllowedFromInternet,
		},
		{
			name: "check RDP access from * is ok when mode is deny",
			source: `
resource "azurerm_network_security_rule" "example_deny" {
     name                        = "example_deny_security_rule"
     direction                   = "Inbound"
     access                      = "Deny"
     protocol                    = "TCP"
     source_port_range           = "*"
     destination_port_range      = ["3389"]
     source_address_prefix       = "*"
     destination_address_prefix  = "*"
}
`,
			mustExcludeResultCode: checks.AZURDPAccessNotAllowedFromInternet,
		},
		{
			name: "check RDP access from 0.0.0.0 causes a failure",
			source: `
resource "azurerm_network_security_rule" "bad_example" {
     name                        = "bad_example_security_rule"
     direction                   = "Inbound"
     access                      = "Allow"
     protocol                    = "TCP"
     source_port_range           = "0.0.0.0"
     destination_port_range      = ["3389"]
     source_address_prefix       = "*"
     destination_address_prefix  = "*"
}
`,
			mustIncludeResultCode: checks.AZURDPAccessNotAllowedFromInternet,
		},
		{
			name: "check RDP access from /0 causes a failure",
			source: `
resource "azurerm_network_security_rule" "bad_example" {
     name                        = "bad_example_security_rule"
     direction                   = "Inbound"
     access                      = "Allow"
     protocol                    = "TCP"
     source_port_range           = "/0"
     destination_port_range      = ["3389"]
     source_address_prefix       = "*"
     destination_address_prefix  = "*"
}
`,
			mustIncludeResultCode: checks.AZURDPAccessNotAllowedFromInternet,
		},
		{
			name: "check RDP access from internet causes a failure",
			source: `
resource "azurerm_network_security_rule" "bad_example" {
     name                        = "bad_example_security_rule"
     direction                   = "Inbound"
     access                      = "Allow"
     protocol                    = "TCP"
     source_port_range           = "internet"
     destination_port_range      = ["3389"]
     source_address_prefix       = "*"
     destination_address_prefix  = "*"
}
`,
			mustIncludeResultCode: checks.AZURDPAccessNotAllowedFromInternet,
		},
		{
			name: "check RDP access from internet causes a failure",
			source: `
resource "azurerm_network_security_rule" "bad_example" {
     name                        = "bad_example_security_rule"
     direction                   = "Inbound"
     access                      = "Allow"
     protocol                    = "TCP"
     source_port_range           = "any"
     destination_port_range      = ["3389"]
     source_address_prefix       = "*"
     destination_address_prefix  = "*"
}
`,
			mustIncludeResultCode: checks.AZURDPAccessNotAllowedFromInternet,
		},
		{
			name: "check RDP access from * causes a failure on security group",
			source: `
resource "azurerm_network_security_group" "example" {
  name                = "tf-appsecuritygroup"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  
  security_rule {
	 source_port_range           = "any"
     destination_port_range      = ["3389", "80", "443"]
     source_address_prefix       = "*"
     destination_address_prefix  = "*"
  }
}
`,
			mustIncludeResultCode: checks.AZURDPAccessNotAllowedFromInternet,
		},
		{
			name: "check RDP access from * is ok when access mode is deny",
			source: `
resource "azurerm_network_security_group" "example_deny" {
  name                = "tf-appsecuritygroup"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  
  security_rule {
     access                      = "Deny"
     source_port_range           = "any"
     destination_port_range      = ["3389", "80", "443"]
     source_address_prefix       = "*"
     destination_address_prefix  = "*"
  }
}
`,
			mustExcludeResultCode: checks.AZURDPAccessNotAllowedFromInternet,
		},
		{
			name: "check RDP access from multiple security rules causes a failure on security group",
			source: `
resource "azurerm_network_security_group" "example" {
  name                = "tf-appsecuritygroup"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  
  security_rule {
	 source_port_range           = "any"
     destination_port_range      = ["3389", "80", "443"]
     source_address_prefix       = "4.53.160.75"
     destination_address_prefix  = "*"
  }

  security_rule {
	 source_port_range           = "any"
     destination_port_range      = ["3389", "80", "443"]
     source_address_prefix       = "internet"
     destination_address_prefix  = "*"
  }
}
`,
			mustIncludeResultCode: checks.AZURDPAccessNotAllowedFromInternet,
		},
		{
			name: "check RDP is acceptable from a specific source",
			source: `
resource "azurerm_network_security_rule" "good_example" {
     name                        = "good_example_security_rule"
     direction                   = "Inbound"
     access                      = "Allow"
     protocol                    = "TCP"
     source_port_range           = "*"
     destination_port_range      = ["3389"]
     source_address_prefix       = "4.53.160.75"
     destination_address_prefix  = "*"
}
`,
			mustExcludeResultCode: checks.AZURDPAccessNotAllowedFromInternet,
		},
		{
			name: "check RDP is acceptable from a specific source on security group",
			source: `
resource "azurerm_network_security_group" "example" {
  name                = "tf-appsecuritygroup"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  
  security_rule {
	 source_port_range           = "any"
     destination_port_range      = ["3389"]
     source_address_prefix       = "4.53.160.75"
     destination_address_prefix  = "*"
  }
}
`,
			mustExcludeResultCode: checks.AZURDPAccessNotAllowedFromInternet,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			results := scanSource(test.source)
			assertCheckCode(t, test.mustIncludeResultCode, test.mustExcludeResultCode, results)
		})
	}

}
