package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"os"
	"testing"
)

func TestAccAgentCredentialsDataSource(t *testing.T) {
	accountId := os.Getenv("MATILLION_ACCOUNT_ID")
	resourceName := testAccResourceName("test-agent-creds")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAgentCredentialsDataSourceConfig(accountId, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the data source attributes are set
					resource.TestCheckResourceAttrSet("data.matillion-streaming_agent_credentials.test", "agent_id"),
					resource.TestCheckResourceAttrSet("data.matillion-streaming_agent_credentials.test", "client_id"),
					resource.TestCheckResourceAttrSet("data.matillion-streaming_agent_credentials.test", "client_secret"),
					// Verify the agent_id matches the resource
					resource.TestCheckResourceAttrPair(
						"data.matillion-streaming_agent_credentials.test", "agent_id",
						"matillion-streaming_agent.test", "agent_id",
					),
				),
			},
		},
	})
}

func testAccAgentCredentialsDataSourceConfig(accountId, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id = "%s"
  region     = "%s"
}

resource "matillion-streaming_agent" "test" {
  name           = "%s"
  description    = "Test agent for credentials data source"
  deployment     = "fargate"
  cloud_provider = "aws"
}

data "matillion-streaming_agent_credentials" "test" {
  agent_id = matillion-streaming_agent.test.agent_id
}
`, accountId, testAccGetRegion(), resourceName)
}
