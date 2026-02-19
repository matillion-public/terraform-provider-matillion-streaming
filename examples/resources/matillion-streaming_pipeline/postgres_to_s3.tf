# Example: PostgreSQL to S3 streaming pipeline
resource "matillion-streaming_pipeline" "postgres_to_s3" {
  name       = "postgres-to-s3-pipeline"
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
    }
    tables = [
      {
        schema = "public"
        table  = "events"
      }
    ]
  }

  # S3 target configuration
  s3_target = {
    bucket          = "my-streaming-bucket"
    prefix          = "streaming/postgres/events"
    decimal_mapping = "logical"
  }
}
