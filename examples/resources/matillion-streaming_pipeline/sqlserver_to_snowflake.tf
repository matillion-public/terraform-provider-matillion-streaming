# Example: SQL Server to Snowflake streaming pipeline
resource "matillion-streaming_pipeline" "sqlserver_to_snowflake" {
  name       = "sqlserver-to-snowflake-pipeline"
  project_id = "your-project-id"
  agent_id   = "your-agent-id"

  # SQL Server source configuration
  sql_server_source = {
    connection = {
      host     = "sqlserver.example.com"
      port     = 1433
      database = "ProductionDB"
      username = "streaming_user"
      password = {
        type = "aws_secrets_manager"
        name = "sqlserver-streaming-password"
      }
      jdbc_properties = {
        "encrypt"                = "true"
        "trustServerCertificate" = "false"
      }
    }
    tables = [
      {
        schema = "dbo"
        table  = "Customers"
      },
      {
        schema = "sales"
        table  = "Orders"
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
    stage_prefix        = "streaming/sqlserver"
    table_schema        = "PUBLIC"
    table_prefix_type   = "prefix"
    transformation_type = "copy_table_soft_delete"
    temporal_mapping    = "native"
  }
}
