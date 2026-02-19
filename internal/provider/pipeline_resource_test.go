package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"log"
	"os"
	"regexp"
	"testing"
)

func init() {
	resource.AddTestSweepers("matillion-streaming_pipeline", &resource.Sweeper{
		Name: "matillion-streaming_pipeline",
		F:    testSweepStreamingPipelines,
	})
}

func testSweepStreamingPipelines(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	// Test project ID used in acceptance tests
	projectId := os.Getenv("MATILLION_PROJECT_ID")

	pipelines, err := client.Pipelines.List(projectId)
	if err != nil {
		return fmt.Errorf("error listing streaming pipelines: %s", err)
	}

	for _, pipeline := range pipelines {
		// Only delete test pipelines (those with names starting with "tf-acc-pipeline-")
		if testSweepMatchesTestRunID(pipeline.Name) {
			log.Printf("[INFO] Deleting streaming pipeline: %s (ID: %s)", pipeline.Name, pipeline.StreamingPipelineId)
			if err = client.Pipelines.Delete(projectId, pipeline.StreamingPipelineId); err != nil {
				log.Printf("[ERROR] Failed to delete streaming pipeline %s (ID: %s): %s", pipeline.Name, pipeline.StreamingPipelineId, err)
			}
		}
	}

	return nil
}

// TestAccStreamingPipelineResource tests PostgreSQL -> Snowflake streaming pipeline
func TestAccStreamingPipelineResource_postgres_snowflake(t *testing.T) {
	resourceName := testAccResourceName("tf-acc-pipeline")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingPipelineConfig_postgres_snowflake(t, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "name", resourceName),
					resource.TestCheckResourceAttrSet("matillion-streaming_pipeline.test", "project_id"),
					resource.TestCheckResourceAttrSet("matillion-streaming_pipeline.test", "agent_id"),
					resource.TestCheckResourceAttrSet("matillion-streaming_pipeline.test", "pipeline_id"),

					// PostgreSQL source checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "postgres_source.connection.host", "localhost"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "postgres_source.connection.port", "5432"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "postgres_source.connection.database", "testdb"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "postgres_source.connection.username", "postgres"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "postgres_source.connection.password.type", "aws_secrets_manager"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "postgres_source.connection.password.name", "test-postgres-password"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "postgres_source.tables.0.schema", "public"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "postgres_source.tables.0.table", "users"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "postgres_source.connection.jdbc_properties.ssl", "require"),

					// Snowflake target checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.connection.account_name", "test-account.eu-central-1"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.connection.username", "SNOWFLAKE_USERNAME"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.connection.authentication.private_key.secret_type", "aws_secrets_manager"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.connection.authentication.private_key.secret_name", "snowflake-private-key"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.connection.authentication.passphrase.secret_type", "aws_secrets_manager"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.connection.authentication.passphrase.secret_name", "snowflake-passphrase"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.connection.jdbc_properties.networkTimeout", "300000"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.role", "STREAMING_ROLE"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.warehouse", "STREAMING_WH"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.database", "STREAMING_DB"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.stage_schema", "STREAMING"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.stage_name", "STREAMING_STAGE"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.stage_prefix", "test/streaming"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.table_schema", "STREAMING"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.table_prefix_type", "prefix"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.transformation_type", "copy_table"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.temporal_mapping", "native"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_pipeline.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingPipelineImportStateIdFunc("matillion-streaming_pipeline.test"),
				ImportStateVerifyIdentifierAttribute: "pipeline_id",
			},
			{
				Config: testAccStreamingPipelineConfig_postgres_snowflake_updated(t, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "name", resourceName+"-updated"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.stage_prefix", "test/streaming/updated"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "snowflake_target.transformation_type", "copy_table_soft_delete"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test", "advanced_properties.buffer.size", "2000"),
				),
			},
		},
	})
}

func testAccStreamingPipelineConfig_postgres_snowflake(t *testing.T, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test" {
  name       = "%s"
  project_id = "%s"
  agent_id   = "%s"

  postgres_source = {
    connection = {
      host     = "localhost"
      port     = 5432
      database = "testdb"
      username     = "postgres"
      password = {
        type = "aws_secrets_manager"
        name = "test-postgres-password"
      }
      jdbc_properties = {
        "ssl" = "require"
      }
    }
    tables = [
      {
        schema = "public"
        table  = "users"
      }
    ]
  }

  snowflake_target = {
    connection = {
      account_name = "test-account.eu-central-1"
      username     = "SNOWFLAKE_USERNAME"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-private-key"
        }
        passphrase = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-passphrase"
        }
      }
      jdbc_properties = {
      	"networkTimeout": "300000"
      }
    }
    role                = "STREAMING_ROLE"
    warehouse          = "STREAMING_WH"
    database           = "STREAMING_DB"
    stage_schema       = "STREAMING"
    stage_name         = "STREAMING_STAGE"
    stage_prefix       = "test/streaming"
    table_schema       = "STREAMING"
    table_prefix_type  = "prefix"
    transformation_type = "copy_table"
    temporal_mapping   = "native"
  }

  advanced_properties = {
    "buffer.size" = "1000"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), resourceName, testAccGetProjectID(t), testAccGetAgentID(t))
}

func testAccStreamingPipelineConfig_postgres_snowflake_updated(t *testing.T, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test" {
  name       = "%s-updated"
  project_id = "%s"
  agent_id   = "%s"

  postgres_source = {
    connection = {
      host     = "localhost"
      port     = 5432
      database = "testdb"
      username     = "postgres"
      password = {
        type = "aws_secrets_manager"
        name = "test-postgres-password"
      }
      jdbc_properties = {
        "ssl" = "require"
      }
    }
    tables = [
      {
        schema = "public"
        table  = "users"
      }
    ]
  }

  snowflake_target = {
    connection = {
      account_name = "test-account.eu-central-1"
      username     = "SNOWFLAKE_USERNAME"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-private-key"
        }
        passphrase = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-passphrase"
        }
      }
      jdbc_properties = {
      	"networkTimeout": "300000"
      }
    }
    role                = "STREAMING_ROLE"
    warehouse          = "STREAMING_WH"
    database           = "STREAMING_DB"
    stage_schema       = "STREAMING"
    stage_name         = "STREAMING_STAGE"
    stage_prefix       = "test/streaming/updated"
    table_schema       = "STREAMING"
    table_prefix_type  = "prefix"
    transformation_type = "copy_table_soft_delete"
    temporal_mapping   = "native"
  }

  advanced_properties = {
    "buffer.size" = "2000"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), resourceName, testAccGetProjectID(t), testAccGetAgentID(t))
}

// TestAccStreamingPipelineResource_postgres_snowflake_optional_values_not_present tests Snowflake target without optional values
func TestAccStreamingPipelineResource_postgres_snowflake_optional_values_not_present(t *testing.T) {
	resourceName := testAccResourceName("tf-acc-pipeline")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingPipelineConfig_postgres_snowflake_optional_values_not_present(t, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "name", resourceName),
					resource.TestCheckResourceAttrSet("matillion-streaming_pipeline.test_no_optional", "project_id"),
					resource.TestCheckResourceAttrSet("matillion-streaming_pipeline.test_no_optional", "agent_id"),
					resource.TestCheckResourceAttrSet("matillion-streaming_pipeline.test_no_optional", "pipeline_id"),
					resource.TestCheckNoResourceAttr("matillion-streaming_pipeline.test_no_optional", "postgres_source.connection.jdbc_properties.%"),

					// Snowflake target checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.connection.account_name", "test-account.eu-central-1"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.connection.username", "SNOWFLAKE_USERNAME"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.connection.authentication.private_key.secret_type", "aws_secrets_manager"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.connection.authentication.private_key.secret_name", "snowflake-private-key"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.role", "STREAMING_ROLE"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.warehouse", "STREAMING_WH"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.database", "STREAMING_DB"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.stage_schema", "STREAMING"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.stage_name", "STREAMING_STAGE"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.table_schema", "STREAMING"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.table_prefix_type", "source_database_and_schema"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.transformation_type", "copy_table_soft_delete"),
					resource.TestCheckNoResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.stage_prefix"),
					resource.TestCheckNoResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.temporal_mapping"),
					resource.TestCheckNoResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.connection.authentication.passphrase.secret_type"),
					resource.TestCheckNoResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.connection.authentication.passphrase.secret_name"),
					resource.TestCheckNoResourceAttr("matillion-streaming_pipeline.test_no_optional", "snowflake_target.connection.jdbc_properties.%"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_pipeline.test_no_optional",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingPipelineImportStateIdFunc("matillion-streaming_pipeline.test_no_optional"),
				ImportStateVerifyIdentifierAttribute: "pipeline_id",
			},
		},
	})
}

func testAccStreamingPipelineConfig_postgres_snowflake_optional_values_not_present(t *testing.T, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test_no_optional" {
  name       = "%s"
  project_id = "%s"
  agent_id   = "%s"

  postgres_source = {
    connection = {
      host     = "localhost"
      port     = 5432
      database = "testdb"
      username     = "postgres"
      password = {
        type = "aws_secrets_manager"
        name = "test-postgres-password"
      }
    }
    tables = [
      {
        schema = "public"
        table  = "users"
      }
    ]
  }

  snowflake_target = {
    connection = {
      account_name = "test-account.eu-central-1"
      username     = "SNOWFLAKE_USERNAME"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-private-key"
        }
      }
    }
    role                = "STREAMING_ROLE"
    warehouse          = "STREAMING_WH"
    database           = "STREAMING_DB"
    stage_schema       = "STREAMING"
    stage_name         = "STREAMING_STAGE"
    table_schema       = "STREAMING"
    table_prefix_type  = "source_database_and_schema"
    transformation_type = "copy_table_soft_delete"
  }

  advanced_properties = {
    "buffer.size" = "1000"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), resourceName, testAccGetProjectID(t), testAccGetAgentID(t))
}

// TestAccStreamingPipelineResource_sqlserver_snowflake tests SQL Server -> Snowflake streaming pipeline
func TestAccStreamingPipelineResource_sqlserver_snowflake(t *testing.T) {
	resourceName := testAccResourceName("tf-acc-pipeline-sql")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingPipelineConfig_sqlserver_snowflake(t, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "name", resourceName),
					resource.TestCheckResourceAttrSet("matillion-streaming_pipeline.test_sql", "pipeline_id"),

					// SQL Server source checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "sql_server_source.connection.host", "localhost"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "sql_server_source.connection.port", "1433"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "sql_server_source.connection.database", "testdb"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "sql_server_source.connection.username", "SA"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "sql_server_source.connection.password.type", "aws_secrets_manager"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "sql_server_source.connection.password.name", "test-sqlserver-password"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "sql_server_source.connection.jdbc_properties.%", "0"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "sql_server_source.tables.0.schema", "dbo"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "sql_server_source.tables.0.table", "customers"),

					// Snowflake Target checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "snowflake_target.table_prefix_type", "none"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "snowflake_target.transformation_type", "change_log"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "snowflake_target.temporal_mapping", "epoch"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "snowflake_target.connection.jdbc_properties.%", "0"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_pipeline.test_sql",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingPipelineImportStateIdFunc("matillion-streaming_pipeline.test_sql"),
				ImportStateVerifyIdentifierAttribute: "pipeline_id",
			},
			{
				// Second apply to verify no drift with empty jdbc_properties
				Config: testAccStreamingPipelineConfig_sqlserver_snowflake(t, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_sql", "snowflake_target.connection.jdbc_properties.%", "0"),
				),
			},
		},
	})
}

func testAccStreamingPipelineConfig_sqlserver_snowflake(t *testing.T, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test_sql" {
  name       = "%s"
  project_id = "%s"
  agent_id   = "%s"

  sql_server_source = {
    connection = {
      host     = "localhost"
      port     = 1433
      database = "testdb"
      username     = "SA"
      password = {
        type = "aws_secrets_manager"
        name = "test-sqlserver-password"
      }
      jdbc_properties = {}
    }
    tables = [
      {
        schema = "dbo"
        table  = "customers"
      }
    ]
  }

  snowflake_target = {
    connection = {
      account_name = "test-account.eu-central-1"
      username     = "SNOWFLAKE_USERNAME"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-private-key"
        }
        passphrase = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-passphrase"
        }
      }
      jdbc_properties = {}
    }
    role                = "STREAMING_ROLE"
    warehouse          = "STREAMING_WH"
    database           = "STREAMING_DB"
    stage_schema       = "STREAMING"
    stage_name         = "STREAMING_STAGE"
    table_schema       = "STREAMING"
    table_prefix_type  = "none"
    transformation_type = "change_log"
    temporal_mapping   = "epoch"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), resourceName, testAccGetProjectID(t), testAccGetAgentID(t))
}

// TestAccStreamingPipelineResource_validation tests validation errors
func TestAccStreamingPipelineResource_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccStreamingPipelineConfig_missing_name(t),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccStreamingPipelineConfig_missing_projectId(t),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccStreamingPipelineConfig_missing_agentId(t),
				ExpectError: regexp.MustCompile("Missing required argument"),
			},
			{
				Config:      testAccStreamingPipelineConfig_invalid_transformation_type(t),
				ExpectError: regexp.MustCompile("Invalid Attribute Value Match"),
			},
			{
				Config:      testAccStreamingPipelineConfig_missing_source(t),
				ExpectError: regexp.MustCompile("Exactly one of these attributes must be configured"),
			},
			{
				Config:      testAccStreamingPipelineConfig_missing_target(t),
				ExpectError: regexp.MustCompile("Exactly one of these attributes must be configured"),
			},
		},
	})
}

func testAccStreamingPipelineConfig_missing_name(t *testing.T) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test" {
  project_id = "%s"
  agent_id   = "%s"

  postgres_source = {
    connection = {
      host     = "localhost"
      port     = 5432
      database = "testdb"
      username     = "postgres"
      password = {
        type = "aws_secrets_manager"
        name = "test-postgres-password"
      }
      jdbc_properties = {
        "ssl" = "require"
      }
    }
    tables = [
      {
        schema = "public"
        table  = "users"
      }
    ]
  }

  snowflake_target = {
    connection = {
      account_name = "test-account.eu-central-1"
      username     = "SNOWFLAKE_USERNAME"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-private-key"
        }
        passphrase = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-passphrase"
        }
      }
    }
    role                = "STREAMING_ROLE"
    warehouse          = "STREAMING_WH"
    database           = "STREAMING_DB"
    stage_schema       = "STREAMING"
    stage_name         = "STREAMING_STAGE"
    stage_prefix       = "test/streaming"
    table_schema       = "STREAMING"
    table_prefix_type  = "prefix"
    transformation_type = "copy_table"
    temporal_mapping   = "native"
  }

  advanced_properties = {
    "buffer.size" = "1000"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), testAccGetProjectID(t), testAccGetAgentID(t))
}

func testAccStreamingPipelineConfig_missing_projectId(t *testing.T) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test" {
  name = "missing project id"
  agent_id   = "%s"

  postgres_source = {
    connection = {
      host     = "localhost"
      port     = 5432
      database = "testdb"
      username     = "postgres"
      password = {
        type = "aws_secrets_manager"
        name = "test-postgres-password"
      }
      jdbc_properties = {
        "ssl" = "require"
      }
    }
    tables = [
      {
        schema = "public"
        table  = "users"
      }
    ]
  }

  snowflake_target = {
    connection = {
      account_name = "test-account.eu-central-1"
      username     = "SNOWFLAKE_USERNAME"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-private-key"
        }
        passphrase = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-passphrase"
        }
      }
    }
    role                = "STREAMING_ROLE"
    warehouse          = "STREAMING_WH"
    database           = "STREAMING_DB"
    stage_schema       = "STREAMING"
    stage_name         = "STREAMING_STAGE"
    stage_prefix       = "test/streaming"
    table_schema       = "STREAMING"
    table_prefix_type  = "prefix"
    transformation_type = "copy_table"
    temporal_mapping   = "native"
  }

  advanced_properties = {
    "buffer.size" = "1000"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), testAccGetAgentID(t))
}

func testAccStreamingPipelineConfig_missing_agentId(t *testing.T) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test" {
  name = "missing agentId"
  project_id   = "%s"

  postgres_source = {
    connection = {
      host     = "localhost"
      port     = 5432
      database = "testdb"
      username     = "postgres"
      password = {
        type = "aws_secrets_manager"
        name = "test-postgres-password"
      }
      jdbc_properties = {
        "ssl" = "require"
      }
    }
    tables = [
      {
        schema = "public"
        table  = "users"
      }
    ]
  }

  snowflake_target = {
    connection = {
      account_name = "test-account.eu-central-1"
      username     = "SNOWFLAKE_USERNAME"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-private-key"
        }
        passphrase = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-passphrase"
        }
      }
    }
    role                = "STREAMING_ROLE"
    warehouse          = "STREAMING_WH"
    database           = "STREAMING_DB"
    stage_schema       = "STREAMING"
    stage_name         = "STREAMING_STAGE"
    stage_prefix       = "test/streaming"
    table_schema       = "STREAMING"
    table_prefix_type  = "prefix"
    transformation_type = "copy_table"
    temporal_mapping   = "native"
  }

  advanced_properties = {
    "buffer.size" = "1000"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), testAccGetProjectID(t))
}

func testAccStreamingPipelineConfig_invalid_transformation_type(t *testing.T) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test_invalid" {
  name       = "test-invalid"
  project_id = "%s"
  agent_id   = "%s"

  postgres_source = {
    connection = {
      host     = "localhost"
      port     = 5432
      database = "testdb"
      username     = "postgres"
      password = {
        type = "aws_secrets_manager"
        name = "test-password"
      }
    }
    tables = [
      {
        schema = "dbo"
        table  = "customers"
      }
    ]
  }

  snowflake_target = {
    connection = {
      account_name = "test-account"
      username     = "user"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "key"
        }
        passphrase = {
          secret_type = "aws_secrets_manager"
          secret_name = "pass"
        }
      }
    }
    role                = "ROLE"
    warehouse          = "WH"
    database           = "DB"
    stage_schema       = "SCHEMA"
    stage_name         = "STAGE"
    table_schema       = "SCHEMA"
    table_prefix_type  = "prefix"
    transformation_type = "INVALID_TYPE"
    temporal_mapping   = "native"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), testAccGetProjectID(t), testAccGetAgentID(t))
}

func testAccStreamingPipelineConfig_missing_source(t *testing.T) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test_invalid" {
  name       = "test-invalid"
  project_id = "%s"
  agent_id   = "%s"

  snowflake_target = {
    connection = {
      account_name = "test-account"
      username     = "user"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "key"
        }
        passphrase = {
          secret_type = "aws_secrets_manager"
          secret_name = "pass"
        }
      }
    }
    role                = "ROLE"
    warehouse          = "WH"
    database           = "DB"
    stage_schema       = "SCHEMA"
    stage_name         = "STAGE"
    table_schema       = "SCHEMA"
    table_prefix_type  = "prefix"
    transformation_type = "copy_table"
    temporal_mapping   = "native"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), testAccGetProjectID(t), testAccGetAgentID(t))
}

func testAccStreamingPipelineConfig_missing_target(t *testing.T) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test_invalid" {
  name       = "test-invalid"
  project_id = "%s"
  agent_id   = "%s"

  postgres_source = {
    connection = {
      host     = "localhost"
      port     = 5432
      database = "testdb"
      username     = "postgres"
      password = {
        type = "aws_secrets_manager"
        name = "test-password"
      }
    }
    tables = [
      {
        schema = "dbo"
        table  = "customers"
      }
    ]
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), testAccGetProjectID(t), testAccGetAgentID(t))
}

// TestAccStreamingPipelineResource_postgres_s3 tests PostgreSQL -> S3 streaming pipeline
func TestAccStreamingPipelineResource_postgres_s3(t *testing.T) {
	resourceName := testAccResourceName("tf-acc-s3-pipeline")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingPipelineConfig_postgres_s3(t, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_s3", "name", resourceName),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_s3", "s3_target.bucket", "test-streaming-bucket"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_s3", "s3_target.prefix", "streaming/data"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_s3", "s3_target.decimal_mapping", "legacy"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_pipeline.test_s3",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingPipelineImportStateIdFunc("matillion-streaming_pipeline.test_s3"),
				ImportStateVerifyIdentifierAttribute: "pipeline_id",
			},
		},
	})
}

func testAccStreamingPipelineConfig_postgres_s3(t *testing.T, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test_s3" {
  name       = "%s"
  project_id = "%s"
  agent_id   = "%s"

  postgres_source = {
    connection = {
      host     = "localhost"
      port     = 5432
      database = "testdb"
      username     = "postgres"
      password = {
        type = "aws_secrets_manager"
        name = "test-postgres-password"
      }
    }
    tables = [
      {
        schema = "public"
        table  = "events"
      }
    ]
  }

  s3_target = {
    bucket         = "test-streaming-bucket"
    prefix         = "streaming/data"
    decimal_mapping = "legacy"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), resourceName, testAccGetProjectID(t), testAccGetAgentID(t))
}

// TestAccStreamingPipelineResource_postgres_s3_optional_values_not_present tests S3 Target without optional values
func TestAccStreamingPipelineResource_postgres_s3_optional_values_not_present(t *testing.T) {
	resourceName := testAccResourceName("tf-acc-s3-pipeline")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingPipelineConfig_postgres_s3_optional_values_not_present(t, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_s3_no_optional", "name", resourceName),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_s3_no_optional", "s3_target.bucket", "test-streaming-bucket"),
					resource.TestCheckNoResourceAttr("matillion-streaming_pipeline.test_s3_no_optional", "s3_target.prefix"),
					resource.TestCheckNoResourceAttr("matillion-streaming_pipeline.test_s3_no_optional", "s3_target.decimal_mapping"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_pipeline.test_s3_no_optional",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingPipelineImportStateIdFunc("matillion-streaming_pipeline.test_s3_no_optional"),
				ImportStateVerifyIdentifierAttribute: "pipeline_id",
			},
		},
	})

}

func testAccStreamingPipelineConfig_postgres_s3_optional_values_not_present(t *testing.T, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test_s3_no_optional" {
  name       = "%s"
  project_id = "%s"
  agent_id   = "%s"

  postgres_source = {
    connection = {
      host     = "localhost"
      port     = 5432
      database = "testdb"
      username     = "postgres"
      password = {
        type = "aws_secrets_manager"
        name = "test-postgres-password"
      }
    }
    tables = [
      {
        schema = "public"
        table  = "events"
      }
    ]
  }

  s3_target = {
    bucket         = "test-streaming-bucket"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), resourceName, testAccGetProjectID(t), testAccGetAgentID(t))
}

// TestAccStreamingPipelineResource_mysql_snowflake tests MySQL -> Snowflake streaming pipeline
func TestAccStreamingPipelineResource_mysql_snowflake(t *testing.T) {
	resourceName := testAccResourceName("tf-acc-mysql-pipeline")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingPipelineConfig_mysql_snowflake(t, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_mysql", "name", resourceName),
					resource.TestCheckResourceAttrSet("matillion-streaming_pipeline.test_mysql", "pipeline_id"),

					// MySQL source checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_mysql", "mysql_source.connection.host", "localhost"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_mysql", "mysql_source.connection.port", "3306"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_mysql", "mysql_source.connection.username", "mysql_user"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_mysql", "mysql_source.connection.password.type", "aws_secrets_manager"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_mysql", "mysql_source.connection.password.name", "test-mysql-password"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_mysql", "mysql_source.tables.0.schema", "testdb"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_mysql", "mysql_source.tables.0.table", "users"),

					// Snowflake target checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_mysql", "snowflake_target.transformation_type", "copy_table"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_pipeline.test_mysql",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingPipelineImportStateIdFunc("matillion-streaming_pipeline.test_mysql"),
				ImportStateVerifyIdentifierAttribute: "pipeline_id",
			},
		},
	})
}

func testAccStreamingPipelineConfig_mysql_snowflake(t *testing.T, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test_mysql" {
  name       = "%s"
  project_id = "%s"
  agent_id   = "%s"

  mysql_source = {
    connection = {
      host     = "localhost"
      port     = 3306
      username     = "mysql_user"
      password = {
        type = "aws_secrets_manager"
        name = "test-mysql-password"
      }
      jdbc_properties = {
        "useSSL" = "false"
      }
    }
    tables = [
      {
        schema = "testdb"
        table  = "users"
      }
    ]
  }

  snowflake_target = {
    connection = {
      account_name = "test-account.eu-central-1"
      username     = "SNOWFLAKE_USERNAME"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-private-key"
        }
        passphrase = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-passphrase"
        }
      }
    }
    role                = "STREAMING_ROLE"
    warehouse          = "STREAMING_WH"
    database           = "STREAMING_DB"
    stage_schema       = "STREAMING"
    stage_name         = "STREAMING_STAGE"
    table_schema       = "STREAMING"
    table_prefix_type  = "prefix"
    transformation_type = "copy_table"
    temporal_mapping   = "native"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), resourceName, testAccGetProjectID(t), testAccGetAgentID(t))
}

// TestAccStreamingPipelineResource_oracle_snowflake tests Oracle -> Snowflake streaming pipeline
func TestAccStreamingPipelineResource_oracle_snowflake(t *testing.T) {
	resourceName := testAccResourceName("tf-acc-oracle-pipeline")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingPipelineConfig_oracle_snowflake(t, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_oracle", "name", resourceName),
					resource.TestCheckResourceAttrSet("matillion-streaming_pipeline.test_oracle", "pipeline_id"),

					// Oracle source checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_oracle", "oracle_source.connection.host", "oracle-host.example.com"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_oracle", "oracle_source.connection.port", "1521"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_oracle", "oracle_source.connection.database", "ORCL"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_oracle", "oracle_source.connection.pdb", "XEPDB1"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_oracle", "oracle_source.connection.username", "oracle_user"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_oracle", "oracle_source.connection.password.type", "aws_secrets_manager"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_oracle", "oracle_source.connection.password.name", "test-oracle-password"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_oracle", "oracle_source.tables.0.schema", "HR"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_oracle", "oracle_source.tables.0.table", "EMPLOYEES"),

					// Snowflake target checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_oracle", "snowflake_target.transformation_type", "change_log"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_pipeline.test_oracle",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingPipelineImportStateIdFunc("matillion-streaming_pipeline.test_oracle"),
				ImportStateVerifyIdentifierAttribute: "pipeline_id",
			},
		},
	})
}

func testAccStreamingPipelineConfig_oracle_snowflake(t *testing.T, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test_oracle" {
  name       = "%s"
  project_id = "%s"
  agent_id   = "%s"

  oracle_source = {
    connection = {
      host     = "oracle-host.example.com"
      port     = 1521
      database = "ORCL"
      pdb      = "XEPDB1"
      username     = "oracle_user"
      password = {
        type = "aws_secrets_manager"
        name = "test-oracle-password"
      }
      jdbc_properties = {
        "oracle.jdbc.timezoneAsRegion" = "false"
      }
    }
    tables = [
      {
        schema = "HR"
        table  = "EMPLOYEES"
      },
      {
        schema = "HR"
        table  = "DEPARTMENTS"
      }
    ]
  }

  snowflake_target = {
    connection = {
      account_name = "test-account.eu-central-1"
      username     = "SNOWFLAKE_USERNAME"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-private-key"
        }
        passphrase = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-passphrase"
        }
      }
    }
    role                = "STREAMING_ROLE"
    warehouse          = "STREAMING_WH"
    database           = "STREAMING_DB"
    stage_schema       = "STREAMING"
    stage_name         = "STREAMING_STAGE"
    table_schema       = "STREAMING"
    table_prefix_type  = "prefix"
    transformation_type = "change_log"
    temporal_mapping   = "native"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), resourceName, testAccGetProjectID(t), testAccGetAgentID(t))
}

// TestAccStreamingPipelineResource_db2_ibm_i_snowflake tests DB2 for IBM i -> Snowflake streaming pipeline
func TestAccStreamingPipelineResource_db2_ibm_i_snowflake(t *testing.T) {
	resourceName := testAccResourceName("tf-acc-db2-pipeline")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingPipelineConfig_db2_ibm_i_snowflake(t, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_db2", "name", resourceName),
					resource.TestCheckResourceAttrSet("matillion-streaming_pipeline.test_db2", "pipeline_id"),

					// DB2 for IBM i source checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_db2", "db2_ibm_i_source.connection.host", "as400-host.example.com"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_db2", "db2_ibm_i_source.connection.port", "446"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_db2", "db2_ibm_i_source.connection.username", "db2_user"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_db2", "db2_ibm_i_source.connection.password.type", "aws_secrets_manager"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_db2", "db2_ibm_i_source.connection.password.name", "test-db2-password"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_db2", "db2_ibm_i_source.tables.0.schema", "MYLIB"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_db2", "db2_ibm_i_source.tables.0.table", "CUSTOMER"),

					// Snowflake target checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_db2", "snowflake_target.transformation_type", "copy_table_soft_delete"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_pipeline.test_db2",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingPipelineImportStateIdFunc("matillion-streaming_pipeline.test_db2"),
				ImportStateVerifyIdentifierAttribute: "pipeline_id",
			},
		},
	})
}

func testAccStreamingPipelineConfig_db2_ibm_i_snowflake(t *testing.T, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test_db2" {
  name       = "%s"
  project_id = "%s"
  agent_id   = "%s"

  db2_ibm_i_source = {
    connection = {
      host     = "as400-host.example.com"
      port     = 446
      username     = "db2_user"
      password = {
        type = "aws_secrets_manager"
        name = "test-db2-password"
      }
      jdbc_properties = {
        "translate binary" = "true"
      }
    }
    tables = [
      {
        schema = "MYLIB"
        table  = "CUSTOMER"
      },
      {
        schema = "MYLIB"
        table  = "ORDERS"
      }
    ]
  }

  snowflake_target = {
    connection = {
      account_name = "test-account.eu-central-1"
      username     = "SNOWFLAKE_USERNAME"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-private-key"
        }
        passphrase = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-passphrase"
        }
      }
    }
    role                = "STREAMING_ROLE"
    warehouse          = "STREAMING_WH"
    database           = "STREAMING_DB"
    stage_schema       = "STREAMING"
    stage_name         = "STREAMING_STAGE"
    table_schema       = "STREAMING"
    table_prefix_type  = "prefix"
    transformation_type = "copy_table_soft_delete"
    temporal_mapping   = "native"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), resourceName, testAccGetProjectID(t), testAccGetAgentID(t))
}

// TestAccStreamingPipelineResource_postgres_abs tests PostgreSQL -> Azure Blob Storage streaming pipeline
func TestAccStreamingPipelineResource_postgres_abs(t *testing.T) {
	resourceName := testAccResourceName("tf-acc-abs-pipeline")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingPipelineConfig_postgres_abs(t, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_abs", "name", resourceName),
					resource.TestCheckResourceAttrSet("matillion-streaming_pipeline.test_abs", "pipeline_id"),

					// PostgreSQL source checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_abs", "postgres_source.connection.host", "localhost"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_abs", "postgres_source.connection.port", "5432"),

					// Azure Blob Storage target checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_abs", "abs_target.container", "test-container"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_abs", "abs_target.account_name", "teststorageaccount"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_abs", "abs_target.prefix", "streaming/data"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_abs", "abs_target.decimal_mapping", "logical"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_pipeline.test_abs",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingPipelineImportStateIdFunc("matillion-streaming_pipeline.test_abs"),
				ImportStateVerifyIdentifierAttribute: "pipeline_id",
			},
		},
	})
}

func testAccStreamingPipelineConfig_postgres_abs(t *testing.T, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test_abs" {
  name       = "%s"
  project_id = "%s"
  agent_id   = "%s"

  postgres_source = {
    connection = {
      host     = "localhost"
      port     = 5432
      database = "testdb"
      username     = "postgres"
      password = {
        type = "aws_secrets_manager"
        name = "test-postgres-password"
      }
      jdbc_properties = {
        "ssl" = "require"
      }
    }
    tables = [
      {
        schema = "public"
        table  = "events"
      }
    ]
  }

  abs_target = {
    container       = "test-container"
    account_name    = "teststorageaccount"
    account_key     = {
      type = "aws_secrets_manager"
      name = "test-storage-key"
    }
    prefix          = "streaming/data"
    decimal_mapping = "logical"
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), resourceName, testAccGetProjectID(t), testAccGetAgentID(t))
}

// TestAccStreamingPipelineResource_postgres_abs tests Azure Blob Storage without optional values
func TestAccStreamingPipelineResource_postgres_abs_optional_values_not_present(t *testing.T) {
	resourceName := testAccResourceName("tf-acc-abs-pipeline")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccStreamingPipelineConfig_postgres_abs_optional_values_not_present(t, resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_abs_no_optional", "name", resourceName),
					resource.TestCheckResourceAttrSet("matillion-streaming_pipeline.test_abs_no_optional", "pipeline_id"),

					// PostgreSQL source checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_abs_no_optional", "postgres_source.connection.host", "localhost"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_abs_no_optional", "postgres_source.connection.port", "5432"),

					// Azure Blob Storage target checks
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_abs_no_optional", "abs_target.container", "test-container"),
					resource.TestCheckResourceAttr("matillion-streaming_pipeline.test_abs_no_optional", "abs_target.account_name", "teststorageaccount"),
					resource.TestCheckNoResourceAttr("matillion-streaming_pipeline.test_abs_no_optional", "abs_target.prefix"),
					resource.TestCheckNoResourceAttr("matillion-streaming_pipeline.test_abs_no_optional", "abs_target.decimal_mapping"),
				),
			},
			{
				ResourceName:                         "matillion-streaming_pipeline.test_abs_no_optional",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateIdFunc:                    testAccStreamingPipelineImportStateIdFunc("matillion-streaming_pipeline.test_abs_no_optional"),
				ImportStateVerifyIdentifierAttribute: "pipeline_id",
			},
		},
	})
}

func testAccStreamingPipelineConfig_postgres_abs_optional_values_not_present(t *testing.T, resourceName string) string {
	return fmt.Sprintf(`
provider "matillion-streaming" {
  account_id     = "%s"
  region  = "%s"
}

resource "matillion-streaming_pipeline" "test_abs_no_optional" {
  name       = "%s"
  project_id = "%s"
  agent_id   = "%s"

  postgres_source = {
    connection = {
      host     = "localhost"
      port     = 5432
      database = "testdb"
      username     = "postgres"
      password = {
        type = "aws_secrets_manager"
        name = "test-postgres-password"
      }
      jdbc_properties = {
        "ssl" = "require"
      }
    }
    tables = [
      {
        schema = "public"
        table  = "events"
      }
    ]
  }

  abs_target = {
    container       = "test-container"
    account_name    = "teststorageaccount"
    account_key     = {
      type = "aws_secrets_manager"
      name = "test-storage-key"
    }
  }
}
`, testAccGetAccountID(t), testAccGetRegion(), resourceName, testAccGetProjectID(t), testAccGetAgentID(t))
}

func testAccStreamingPipelineImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}

		projectId := rs.Primary.Attributes["project_id"]
		pipelineId := rs.Primary.Attributes["pipeline_id"]

		if projectId == "" {
			return "", fmt.Errorf("no project_id found in state for %s", resourceName)
		}
		if pipelineId == "" {
			return "", fmt.Errorf("no pipeline_id found in state for %s", resourceName)
		}

		return fmt.Sprintf("%s:%s", projectId, pipelineId), nil
	}
}
