# Example: Oracle to Azure Blob Storage streaming pipeline
resource "matillion-streaming_pipeline" "oracle_to_abs" {
  name       = "oracle-to-abs-pipeline"
  project_id = "your-project-id"
  agent_id   = "your-agent-id"

  # Oracle source configuration
  oracle_source = {
    connection = {
      host     = "oracle.example.com"
      port     = 1521
      database = "ORCL"
      pdb      = "ORCLPDB"
      username = "streaming_user"
      password = {
        type = "azure_key_vault"
        name = "oracle-streaming-password"
      }
    }
    tables = [
      {
        schema = "APP_SCHEMA"
        table  = "TRANSACTIONS"
      }
    ]
  }

  # Azure Blob Storage target configuration
  abs_target = {
    container       = "streaming-data"
    prefix          = "streaming/oracle/transactions"
    account_name    = "mystorageaccount"
    decimal_mapping = "logical"
    account_key = {
      type = "azure_key_vault"
      name = "storage-account-key"
      key  = "key1"
    }
  }
}
