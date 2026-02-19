# Example: PostgreSQL to Snowflake streaming pipeline
resource "matillion-streaming_pipeline" "postgres_to_snowflake" {
  name       = "postgres-to-snowflake-pipeline"
  project_id = "your-project-id"
  agent_id   = "your-agent-id"

  # PostgreSQL source configuration
  postgres_source = {
    connection = {
      host     = "postgres.example.com"
      port     = 5432
      database = "production_db"
      username = "streaming_user"
      password = {
        type = "aws_secrets_manager"
        name = "postgres-streaming-password"
      }
      jdbc_properties = {
        "ssl" = "require"
      }
    }
    tables = [
      {
        schema = "public"
        table  = "users"
      },
      {
        schema = "public"
        table  = "orders"
      }
    ]
  }

  # Snowflake target configuration
  snowflake_target = {
    connection = {
      account_name = "myorg.us-east-1"
      username     = "STREAMING_USER"
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
        "networkTimeout" : "300000"
      }
    }
    role                = "STREAMING_ROLE"
    warehouse           = "STREAMING_WH"
    database            = "ANALYTICS_DB"
    stage_schema        = "STREAMING_STAGE"
    stage_name          = "ANALYTICS_STAGE"
    stage_prefix        = "streaming/postgres"
    table_schema        = "PUBLIC"
    table_prefix_type   = "prefix"
    transformation_type = "copy_table"
    temporal_mapping    = "native"
  }

  # Optional advanced properties
  advanced_properties = {
    "buffer.size" = "1000"
  }
}

# Output the pipeline ID
output "pipeline_id" {
  value       = matillion-streaming_pipeline.postgres_to_snowflake.pipeline_id
  description = "The unique identifier of the created pipeline"
}
