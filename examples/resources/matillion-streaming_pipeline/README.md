# Streaming Pipeline Examples

This directory contains examples for creating streaming pipelines with the `matillion-streaming_pipeline` resource.

## Examples

### [resource.tf](./resource.tf)
Complete PostgreSQL to Snowflake pipeline with all configuration options demonstrated.

### [postgres_to_s3.tf](./postgres_to_s3.tf)
PostgreSQL to S3 pipeline - simple example showing streaming data to S3.

### [mysql_to_snowflake.tf](./mysql_to_snowflake.tf)
MySQL to Snowflake pipeline - demonstrates MySQL-specific configuration (no database field).

### [sqlserver_to_snowflake.tf](./sqlserver_to_snowflake.tf)
SQL Server to Snowflake pipeline - shows SQL Server encryption settings

### [oracle_to_abs.tf](./oracle_to_abs.tf)
Oracle to Azure Blob Storage pipeline - shows Oracle PDB configuration

### [db2_to_s3.tf](./db2_to_s3.tf)
DB2 for IBM i to S3 pipeline - demonstrates DB2 configuration (no database field) and advanced properties.

## Prerequisites

Before using these examples, ensure you have:

1. **Created an agent** - Use the `matillion-streaming_agent` resource or create one in the Matillion console
2. **Project ID** - Your Matillion Data Productivity Cloud project ID
3. **Secrets configured** - Database passwords and authentication credentials stored in your cloud secret manager (AWS Secrets Manager, Azure Key Vault, or Google Secret Manager)

## Supported Sources

- PostgreSQL
- MySQL
- SQL Server
- Oracle
- DB2 for IBM i

## Supported Targets

- Snowflake
- Amazon S3
- Azure Blob Storage (ABS)

## Key Configuration Notes

### Database Sources

- **PostgreSQL, SQL Server, Oracle**: Include a `database` field
- **MySQL, DB2**: Do NOT include a `database` field
- **Oracle only**: Can optionally specify a `pdb` (pluggable database)

### Secret References

All password and authentication credentials must reference secrets from a cloud secret manager:

```hcl
password = {
  type = "aws_secrets_manager"  # or "azure_key_vault" or "google_secret_manager"
  name = "my-secret-name"
  key  = "password"              # optional: specific key within the secret
}
```

### Snowflake Authentication

Snowflake requires private key authentication:

```hcl
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
```

### Table Selection

Specify which tables to stream using the `tables` list:

```hcl
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
```

## More Information

For complete schema documentation, see the [Terraform Registry documentation](https://registry.terraform.io/providers/matillion/matillion-streaming/latest/docs/resources/pipeline).
