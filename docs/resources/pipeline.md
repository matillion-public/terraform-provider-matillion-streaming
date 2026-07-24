---
page_title: "matillion-streaming_pipeline Resource - matillion-streaming"
subcategory: ""
description: |-
    Manages a streaming pipeline for data transfer between sources and targets
---

# matillion-streaming_pipeline (Resource)

Manages a streaming pipeline for data transfer between sources and targets

## Example Usage

### PostgreSQL to Snowflake

```terraform
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
```

### PostgreSQL to S3

```terraform
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
```

### MySQL to Snowflake

```terraform
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
```

### SQL Server to Snowflake

```terraform
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
```

### Oracle to Azure Blob Storage

```terraform
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
```

### DB2 for IBM i to S3

```terraform
# Example: DB2 for IBM i to S3 streaming pipeline
resource "matillion-streaming_pipeline" "db2_to_s3" {
  name       = "db2-to-s3-pipeline"
  project_id = "your-project-id"
  agent_id   = "your-agent-id"

  # DB2 for IBM i source configuration (note: no database field for DB2)
  db2_ibm_i_source = {
    connection = {
      host     = "as400.example.com"
      port     = 8471
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
```

### MySQL to Google Cloud Storage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `agent_id` (String) ID of the agent to use
- `name` (String) Name of the streaming pipeline
- `project_id` (String) Identifier of the project containing this streaming pipeline

### Optional

- `abs_target` (Attributes) Azure Blob Storage target configuration (see [below for nested schema](#nestedatt--abs_target))
- `advanced_properties` (Map of String) Advanced configuration properties for the streaming pipeline
- `db2_ibm_i_source` (Attributes) DB2 for IBM i source configuration (see [below for nested schema](#nestedatt--db2_ibm_i_source))
- `gcs_target` (Attributes) Google Cloud Storage target configuration (see [below for nested schema](#nestedatt--gcs_target))
- `mysql_source` (Attributes) MySQL source configuration (see [below for nested schema](#nestedatt--mysql_source))
- `oracle_source` (Attributes) Oracle source configuration (see [below for nested schema](#nestedatt--oracle_source))
- `postgres_source` (Attributes) Postgres source configuration (see [below for nested schema](#nestedatt--postgres_source))
- `s3_target` (Attributes) S3 target configuration (see [below for nested schema](#nestedatt--s3_target))
- `snowflake_target` (Attributes) Snowflake target configuration (see [below for nested schema](#nestedatt--snowflake_target))
- `sql_server_source` (Attributes) SQL Server source configuration (see [below for nested schema](#nestedatt--sql_server_source))

### Read-Only

- `pipeline_id` (String) Unique identifier for the streaming pipeline

<a id="nestedatt--abs_target"></a>
### Nested Schema for `abs_target`

Required:

- `account_key` (Attributes) Secret reference for Azure storage account key (see [below for nested schema](#nestedatt--abs_target--account_key))
- `account_name` (String) Azure storage account name
- `container` (String) Azure Blob Storage container name

Optional:

- `decimal_mapping` (String) Decimal mapping configuration. Valid values: `logical`, `legacy`
- `prefix` (String) ABS prefix

<a id="nestedatt--abs_target--account_key"></a>
### Nested Schema for `abs_target.account_key`

Required:

- `name` (String) Secret reference name
- `type` (String) Secret reference type. Valid values: `aws_secrets_manager`, `azure_key_vault`, `google_secret_manager`

Optional:

- `key` (String) Secret reference key



<a id="nestedatt--db2_ibm_i_source"></a>
### Nested Schema for `db2_ibm_i_source`

Required:

- `connection` (Attributes) (see [below for nested schema](#nestedatt--db2_ibm_i_source--connection))
- `tables` (Attributes List) (see [below for nested schema](#nestedatt--db2_ibm_i_source--tables))

<a id="nestedatt--db2_ibm_i_source--connection"></a>
### Nested Schema for `db2_ibm_i_source.connection`

Required:

- `host` (String) DB2 for IBM i host
- `password` (Attributes) Secret reference for password (see [below for nested schema](#nestedatt--db2_ibm_i_source--connection--password))
- `port` (Number) DB2 for IBM i port
- `username` (String) DB2 for IBM i username

Optional:

- `jdbc_properties` (Map of String) JDBC properties

<a id="nestedatt--db2_ibm_i_source--connection--password"></a>
### Nested Schema for `db2_ibm_i_source.connection.password`

Required:

- `name` (String) Secret reference name
- `type` (String) Secret reference type. Valid values: `aws_secrets_manager`, `azure_key_vault`, `google_secret_manager`

Optional:

- `key` (String) Secret reference key



<a id="nestedatt--db2_ibm_i_source--tables"></a>
### Nested Schema for `db2_ibm_i_source.tables`

Required:

- `schema` (String)
- `table` (String)



<a id="nestedatt--gcs_target"></a>
### Nested Schema for `gcs_target`

Required:

- `bucket` (String) Google Cloud Storage bucket name

Optional:

- `decimal_mapping` (String) Decimal mapping configuration. Valid values: `logical`, `legacy`
- `prefix` (String) GCS prefix


<a id="nestedatt--mysql_source"></a>
### Nested Schema for `mysql_source`

Required:

- `connection` (Attributes) (see [below for nested schema](#nestedatt--mysql_source--connection))
- `tables` (Attributes List) (see [below for nested schema](#nestedatt--mysql_source--tables))

<a id="nestedatt--mysql_source--connection"></a>
### Nested Schema for `mysql_source.connection`

Required:

- `host` (String) MySQL host
- `password` (Attributes) Secret reference for password (see [below for nested schema](#nestedatt--mysql_source--connection--password))
- `port` (Number) MySQL port
- `username` (String) MySQL username

Optional:

- `jdbc_properties` (Map of String) JDBC properties

<a id="nestedatt--mysql_source--connection--password"></a>
### Nested Schema for `mysql_source.connection.password`

Required:

- `name` (String) Secret reference name
- `type` (String) Secret reference type. Valid values: `aws_secrets_manager`, `azure_key_vault`, `google_secret_manager`

Optional:

- `key` (String) Secret reference key



<a id="nestedatt--mysql_source--tables"></a>
### Nested Schema for `mysql_source.tables`

Required:

- `schema` (String)
- `table` (String)



<a id="nestedatt--oracle_source"></a>
### Nested Schema for `oracle_source`

Required:

- `connection` (Attributes) (see [below for nested schema](#nestedatt--oracle_source--connection))
- `tables` (Attributes List) (see [below for nested schema](#nestedatt--oracle_source--tables))

<a id="nestedatt--oracle_source--connection"></a>
### Nested Schema for `oracle_source.connection`

Required:

- `database` (String) Oracle database name
- `host` (String) Oracle host
- `password` (Attributes) Secret reference for password (see [below for nested schema](#nestedatt--oracle_source--connection--password))
- `port` (Number) Oracle port
- `username` (String) Oracle username

Optional:

- `jdbc_properties` (Map of String) JDBC properties
- `pdb` (String) Oracle pluggable database (PDB)

<a id="nestedatt--oracle_source--connection--password"></a>
### Nested Schema for `oracle_source.connection.password`

Required:

- `name` (String) Secret reference name
- `type` (String) Secret reference type. Valid values: `aws_secrets_manager`, `azure_key_vault`, `google_secret_manager`

Optional:

- `key` (String) Secret reference key



<a id="nestedatt--oracle_source--tables"></a>
### Nested Schema for `oracle_source.tables`

Required:

- `schema` (String)
- `table` (String)



<a id="nestedatt--postgres_source"></a>
### Nested Schema for `postgres_source`

Required:

- `connection` (Attributes) (see [below for nested schema](#nestedatt--postgres_source--connection))
- `tables` (Attributes List) (see [below for nested schema](#nestedatt--postgres_source--tables))

<a id="nestedatt--postgres_source--connection"></a>
### Nested Schema for `postgres_source.connection`

Required:

- `database` (String) Postgres database name
- `host` (String) Postgres host
- `password` (Attributes) Secret reference for password (see [below for nested schema](#nestedatt--postgres_source--connection--password))
- `port` (Number) Postgres port
- `username` (String) Postgres username

Optional:

- `jdbc_properties` (Map of String) JDBC properties

<a id="nestedatt--postgres_source--connection--password"></a>
### Nested Schema for `postgres_source.connection.password`

Required:

- `name` (String) Secret reference name
- `type` (String) Secret reference type. Valid values: `aws_secrets_manager`, `azure_key_vault`, `google_secret_manager`

Optional:

- `key` (String) Secret reference key



<a id="nestedatt--postgres_source--tables"></a>
### Nested Schema for `postgres_source.tables`

Required:

- `schema` (String)
- `table` (String)



<a id="nestedatt--s3_target"></a>
### Nested Schema for `s3_target`

Required:

- `bucket` (String) S3 bucket name

Optional:

- `decimal_mapping` (String) Decimal mapping configuration. Valid values: `logical`, `legacy`
- `prefix` (String) S3 prefix


<a id="nestedatt--snowflake_target"></a>
### Nested Schema for `snowflake_target`

Required:

- `connection` (Attributes) Snowflake connection configuration (see [below for nested schema](#nestedatt--snowflake_target--connection))
- `database` (String) Snowflake database
- `role` (String) Snowflake role
- `stage_name` (String) Snowflake stage name
- `stage_schema` (String) Snowflake stage schema
- `table_prefix_type` (String) Table prefix type. Valid values: `prefix`, `source_database_and_schema`, `none`
- `table_schema` (String) Snowflake table schema
- `transformation_type` (String) Transformation type. Valid values: `copy_table`, `copy_table_soft_delete`, `change_log`
- `warehouse` (String) Snowflake warehouse

Optional:

- `stage_prefix` (String) Snowflake stage prefix
- `temporal_mapping` (String) Temporal mapping configuration. Valid values: `native`, `epoch`

<a id="nestedatt--snowflake_target--connection"></a>
### Nested Schema for `snowflake_target.connection`

Required:

- `account_name` (String) Snowflake account name
- `authentication` (Attributes) Snowflake authentication configuration (see [below for nested schema](#nestedatt--snowflake_target--connection--authentication))
- `username` (String) Snowflake username

Optional:

- `jdbc_properties` (Map of String) JDBC properties

<a id="nestedatt--snowflake_target--connection--authentication"></a>
### Nested Schema for `snowflake_target.connection.authentication`

Required:

- `private_key` (Attributes) Private key secret reference (see [below for nested schema](#nestedatt--snowflake_target--connection--authentication--private_key))

Optional:

- `passphrase` (Attributes) Passphrase secret reference (see [below for nested schema](#nestedatt--snowflake_target--connection--authentication--passphrase))

<a id="nestedatt--snowflake_target--connection--authentication--private_key"></a>
### Nested Schema for `snowflake_target.connection.authentication.private_key`

Required:

- `secret_name` (String) Secret name
- `secret_type` (String) Secret type. Valid values: `aws_secrets_manager`, `azure_key_vault`, `google_secret_manager`


<a id="nestedatt--snowflake_target--connection--authentication--passphrase"></a>
### Nested Schema for `snowflake_target.connection.authentication.passphrase`

Required:

- `secret_name` (String) Secret name
- `secret_type` (String) Secret type. Valid values: `aws_secrets_manager`, `azure_key_vault`, `google_secret_manager`





<a id="nestedatt--sql_server_source"></a>
### Nested Schema for `sql_server_source`

Required:

- `connection` (Attributes) (see [below for nested schema](#nestedatt--sql_server_source--connection))
- `tables` (Attributes List) (see [below for nested schema](#nestedatt--sql_server_source--tables))

<a id="nestedatt--sql_server_source--connection"></a>
### Nested Schema for `sql_server_source.connection`

Required:

- `database` (String) SQL Server database name
- `host` (String) SQL Server host
- `password` (Attributes) Secret reference for password (see [below for nested schema](#nestedatt--sql_server_source--connection--password))
- `port` (Number) SQL Server port
- `username` (String) SQL Server username

Optional:

- `jdbc_properties` (Map of String) JDBC properties

<a id="nestedatt--sql_server_source--connection--password"></a>
### Nested Schema for `sql_server_source.connection.password`

Required:

- `name` (String) Secret reference name
- `type` (String) Secret reference type. Valid values: `aws_secrets_manager`, `azure_key_vault`, `google_secret_manager`

Optional:

- `key` (String) Secret reference key



<a id="nestedatt--sql_server_source--tables"></a>
### Nested Schema for `sql_server_source.tables`

Required:

- `schema` (String)
- `table` (String)

## Import

Import is supported using the following syntax:

The [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used, for example:

```shell
#!/bin/bash
# Example: Import an existing streaming pipeline into Terraform state

# Before importing, ensure the resource block is defined in your .tf file, for example:
# resource "matillion-streaming_pipeline" "example" {
#   name       = "production-pipeline"
#   project_id = "proj-abc123"
#   agent_id   = "agent-123"
#   # ... rest of configuration
# }

# The import ID format is: project_id:pipeline_id
terraform import matillion-streaming_pipeline.example "your-project-id:your-pipeline-id"
```
