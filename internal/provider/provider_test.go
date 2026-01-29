package provider

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"log"
	"os"
	"strings"
	"sync"
	"terraform-provider-matillion-streaming/internal/client"
	"testing"
)

// TestMain triggers the sweeper functions if the sweep flag is present
func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"matillion-streaming": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// Verify required environment variables are set
	requiredEnvVars := []string{
		"MATILLION_ACCOUNT_ID",
		"MATILLION_PROJECT_ID",
		"MATILLION_CLIENT_ID",
		"MATILLION_CLIENT_SECRET",
		"MATILLION_REGION",
	}

	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			t.Fatalf("%s must be set for acceptance tests", env)
		}
	}
}

func testAccGetRegion() string {
	if region := os.Getenv("MATILLION_REGION"); region != "" {
		return region
	}
	return "us"
}

// sharedClientForRegion returns a client configured for the sweeper tests
// The region parameter is ignored, but kept for compatibility with the sweeper interface
func sharedClientForRegion(_ string) (*client.Client, error) {
	accountID := os.Getenv("MATILLION_ACCOUNT_ID")
	if accountID == "" {
		return nil, fmt.Errorf("MATILLION_ACCOUNT_ID must be set for sweeper tests")
	}

	clientID := os.Getenv("MATILLION_CLIENT_ID")
	if clientID == "" {
		return nil, fmt.Errorf("MATILLION_CLIENT_ID must be set for sweeper tests")
	}

	clientSecret := os.Getenv("MATILLION_CLIENT_SECRET")
	if clientSecret == "" {
		return nil, fmt.Errorf("MATILLION_CLIENT_SECRET must be set for sweeper tests")
	}

	regionStr := testAccGetRegion()
	var region client.Region
	switch regionStr {
	case "eu":
		region = client.RegionEU
	case "us":
		region = client.RegionUS
	default:
		return nil, fmt.Errorf("invalid region %s, must be eu or us", regionStr)
	}

	c, err := client.NewClient(accountID, region)

	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	return c, nil
}

// testSweepMatchesTestRunID checks if a resource name belongs to the current test run.
// It matches resources that contain the current test run ID in their name.
func testSweepMatchesTestRunID(name string) bool {
	runID := testAccGetTestRunID()
	return strings.Contains(name, fmt.Sprintf("-%s", runID))
}

// testRunID stores the unique identifier for this test run
var testRunID string
var testRunIDOnce sync.Once

// testAccGetTestRunID returns a unique identifier for this test run.
// It checks the TEST_RUN_ID environment variable first, then generates
// an 8-character UUID if not set. The ID is cached for the test suite duration.
func testAccGetTestRunID() string {
	testRunIDOnce.Do(func() {
		if id := os.Getenv("TEST_RUN_ID"); id != "" {
			testRunID = id
		} else {
			// Generate 8-character UUID
			testRunID = strings.ReplaceAll(uuid.New().String(), "-", "")[:8]
		}
		log.Printf("[INFO] Test Run ID: %s", testRunID)
	})
	return testRunID
}

// testAccResourceName generates a unique resource name for acceptance tests.
// Format: {prefix}-{test-run-id}
// Example: tf-acc-pipeline-a1b2c3d4
func testAccResourceName(prefix string) string {
	runID := testAccGetTestRunID()
	return fmt.Sprintf("%s-%s", prefix, runID)
}
