# Example: MySQL to Google Cloud Storage streaming pipeline
resource "matillion-streaming_pipeline" "mysql_to_gcs" {
  name       = "mysql-to-gcs-pipeline"
  project_id = "your-project-id"
  agent_id   = "your-agent-id"

  # MySQL source configuration
  mysql_source = {
    connection = {
      host     = "mysql.example.com"
      port     = 3306
      username = "streaming_user"
      password = {
        type = "aws_secrets_manager"
        name = "mysql-streaming-password"
      }
    }
    tables = [
      {
        schema = "app"
        table  = "orders"
      }
    ]
  }

  # Google Cloud Storage target configuration
  gcs_target = {
    bucket          = "streaming-data"
    prefix          = "streaming/mysql/orders"
    decimal_mapping = "logical"
  }
}
