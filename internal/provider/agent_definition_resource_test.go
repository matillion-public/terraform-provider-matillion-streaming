package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"log"
	"os"
	"testing"
)

func init() {
	resource.AddTestSweepers("matillion-streaming_agent", &resource.Sweeper{
		Name: "matillion-streaming_agent",
		F:    testSweepStreamingAgents,
	})
}

func testSweepStreamingAgents(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	agents, err := client.Agents.List()
	if err != nil {
		return fmt.Errorf("error listing agents: %s", err)
	}

	for _, agent := range agents {
		// Only delete streaming test agents from this test run
		if agent.AgentType == "streaming" && testSweepMatchesTestRunID(agent.Name) {
			log.Printf("[INFO] Deleting streaming agent: %s (ID: %s)", agent.Name, agent.AgentId)
			if err = client.Agents.Delete(agent.AgentId); err != nil {
				log.Printf("[ERROR] Failed to delete streaming agent %s (ID: %s): %s", agent.Name, agent.AgentId, err)
			}
		}
	}

	return nil
}

func TestAccAgentDefinitionResource(t *testing.T) {
	accountId := os.Getenv("MATILLION_ACCOUNT_ID")

	resourceName := testAccResourceName("test-strm-agent-aws")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingAgentConfig_basic(accountId, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_agent.test", "name", resourceName),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test", "description", "Test streaming agent for AWS"),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test", "deployment", "fargate"),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test", "cloud_provider", "aws"),
					resource.TestCheckResourceAttrSet("matillion-streaming_agent.test", "agent_id"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_agent.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingAgentImportStateIdFunc("matillion-streaming_agent.test"),
				ImportStateVerifyIdentifierAttribute: "agent_id",
			},
			{
				Config: testAccStreamingAgentConfig_updated(accountId, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_agent.test", "name", resourceName),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test", "description", "Updated test streaming agent for AWS"),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test", "deployment", "fargate"),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test", "cloud_provider", "aws"),
				),
			},
		},
	})

}

func TestAccAgentDefinitionResource_minimal(t *testing.T) {
	accountId := os.Getenv("MATILLION_ACCOUNT_ID")

	resourceName := testAccResourceName("test-strm-agent-min")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingAgentConfig_minimal(accountId, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_agent.test_minimal", "name", resourceName),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test_minimal", "deployment", "fargate"),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test_minimal", "cloud_provider", "aws"),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test_minimal", "description", ""),
					resource.TestCheckResourceAttrSet("matillion-streaming_agent.test_minimal", "agent_id"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_agent.test_minimal",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingAgentImportStateIdFunc("matillion-streaming_agent.test_minimal"),
				ImportStateVerifyIdentifierAttribute: "agent_id",
			},
		},
	})
}

func TestAccAgentDefinitionResource_azure(t *testing.T) {
	accountID := os.Getenv("MATILLION_ACCOUNT_ID")

	resourceName := testAccResourceName("test-strm-agent-azure")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingAgentConfig_azure(accountID, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_agent.test_azure", "name", resourceName),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test_azure", "description", "Test streaming agent for Azure"),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test_azure", "deployment", "aci"),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test_azure", "cloud_provider", "azure"),
					resource.TestCheckResourceAttrSet("matillion-streaming_agent.test_azure", "agent_id"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_agent.test_azure",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingAgentImportStateIdFunc("matillion-streaming_agent.test_azure"),
				ImportStateVerifyIdentifierAttribute: "agent_id",
			},
		},
	})
}

func TestAccAgentDefinitionResource_gcp(t *testing.T) {
	accountID := os.Getenv("MATILLION_ACCOUNT_ID")

	resourceName := testAccResourceName("test-strm-agent-gcp")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingAgentConfig_gcp(accountID, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_agent.test_gcp", "name", resourceName),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test_gcp", "description", "Test streaming agent for GCP"),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test_gcp", "deployment", "gke"),
					resource.TestCheckResourceAttr("matillion-streaming_agent.test_gcp", "cloud_provider", "gcp"),
					resource.TestCheckResourceAttrSet("matillion-streaming_agent.test_gcp", "agent_id"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_agent.test_gcp",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingAgentImportStateIdFunc("matillion-streaming_agent.test_gcp"),
				ImportStateVerifyIdentifierAttribute: "agent_id",
			},
		},
	})
}

func testAccStreamingAgentConfig_basic(accountId, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id = "%s"
  region     = "%s"
}

resource "matillion-streaming_agent" "test" {
  name           = "%s"
  description    = "Test streaming agent for AWS"
  deployment     = "fargate"
  cloud_provider = "aws"
}
`, accountId, testAccGetRegion(), resourceName)
}

func testAccStreamingAgentConfig_updated(accountId, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id = "%s"
  region     = "%s"
}

resource "matillion-streaming_agent" "test" {
  name           = "%s"
  description    = "Updated test streaming agent for AWS"
  deployment     = "fargate"
  cloud_provider = "aws"
}
`, accountId, testAccGetRegion(), resourceName)
}

func testAccStreamingAgentConfig_minimal(accountId, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_agent" "test_minimal" {
  name           = "%s"
  deployment     = "fargate"
  cloud_provider = "aws"
  description    = ""
}
`, accountId, testAccGetRegion(), resourceName)
}

func testAccStreamingAgentConfig_azure(accountId, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_agent" "test_azure" {
  name           = "%s"
  description    = "Test streaming agent for Azure"
  deployment     = "aci"
  cloud_provider = "azure"
}
`, accountId, testAccGetRegion(), resourceName)
}

func testAccStreamingAgentConfig_gcp(accountId, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_agent" "test_gcp" {
  name           = "%s"
  description    = "Test streaming agent for GCP"
  deployment     = "gke"
  cloud_provider = "gcp"
}
`, accountId, testAccGetRegion(), resourceName)
}

func testAccStreamingAgentImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		agentId := rs.Primary.Attributes["agent_id"]
		if agentId == "" {
			return "", fmt.Errorf("no agent_id found in state for %s", resourceName)
		}

		return agentId, nil
	}
}
