# Example: MySQL to Snowflake streaming pipeline
resource "matillion-streaming_pipeline" "mysql_to_snowflake" {
  name       = "mysql-to-snowflake-pipeline"
  project_id = "your-project-id"
  agent_id   = "your-agent-id"

  # MySQL source configuration (note: no database field for MySQL)
  mysql_source = {
    connection = {
      host     = "mysql.example.com"
      port     = 3306
      username = "streaming_user"
      password = {
        type = "aws_secrets_manager"
        name = "mysql-streaming-password"
      }
      jdbc_properties = {
        "useSSL" = "true"
      }
    }
    tables = [
      {
        schema = "myapp"
        table  = "customers"
      }
    ]
  }

  # Snowflake target configuration
  snowflake_target = {
    connection = {
      account_name = "myorg.us-west-2"
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
    table_schema        = "PUBLIC"
    table_prefix_type   = "none"
    transformation_type = "change_log"
    temporal_mapping    = "native"
  }
}
