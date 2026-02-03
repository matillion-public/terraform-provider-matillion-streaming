# Example: DB2 for IBM i to S3 streaming pipeline
resource "matillion-streaming_pipeline" "db2_to_s3" {
  name       = "db2-to-s3-pipeline"
  project_id = "your-project-id"
  agent_id   = "your-agent-id"

  # DB2 for IBM i source configuration (note: no database field for DB2)
  db2_ibm_i_source = {
    connection = {
      host     = "as400.example.com"
      port     = 446
      username = "STREAMING"
      password = {
        type = "aws_secrets_manager"
        name = "db2-streaming-password"
      }
      jdbc_properties = {
        "secure" = "true"
      }
    }
    tables = [
      {
        schema = "MYLIB"
        table  = "CUSTOMERS"
      },
      {
        schema = "MYLIB"
        table  = "INVENTORY"
      }
    ]
  }

  # S3 target configuration
  s3_target = {
    bucket          = "my-streaming-bucket"
    prefix          = "streaming/db2/data"
    decimal_mapping = "logical"
  }

  advanced_properties = {
    "batch.size" = "500"
  }
}
