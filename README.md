# Terraform Provider for Matillion Streaming

[![Terraform Registry](https://img.shields.io/badge/terraform-registry-623CE4?logo=terraform)](https://registry.terraform.io/providers/matillion-public/matillion-streaming/latest)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-1.25-00ADD8?logo=go)](go.mod)

The Matillion Streaming Terraform provider lets you manage [Matillion Streaming](https://docs.matillion.com/data-productivity-cloud/streaming/docs/streaming-pipelines/) resources as infrastructure-as-code.
You can declaratively create and manage streaming agents and pipelines that capture data from a source database and write it to either cloud storage or a cloud data warehouse.

Full documentation is available on the [Terraform Registry](https://registry.terraform.io/providers/matillion-public/matillion-streaming/latest/docs).



## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- A Matillion Data Productivity Cloud account with [API credentials](https://docs.matillion.com/data-productivity-cloud/api/docs/authentication/)
- [Go](https://golang.org/doc/install) >= 1.25 (for building the provider from source)



## Using the Provider

Add the provider to your Terraform configuration:

```hcl
terraform {
  required_providers {
    matillion-streaming = {
      source  = "matillion-public/matillion-streaming"
      version = "~> 1.0"
    }
  }
}

provider "matillion-streaming" {
  account_id = "your-account-id"
  region     = "eu" # or "us"
}
```



## Authentication

The provider authenticates using OAuth2 client credentials. Set the following environment variables before running Terraform:

```bash
export MATILLION_CLIENT_ID="your-client-id"
export MATILLION_CLIENT_SECRET="your-client-secret"
```

You can see how to obtain these credentials from the [Matillion Docs](https://docs.matillion.com/data-productivity-cloud/api/docs/authentication/).



## Resources & Data Sources

| Type | Name | Description |
|------|------|-------------|
| Resource | [`matillion-streaming_agent`](https://registry.terraform.io/providers/matillion-public/matillion-streaming/latest/docs/resources/agent) | Manages a streaming agent in DPC |
| Resource | [`matillion-streaming_pipeline`](https://registry.terraform.io/providers/matillion-public/matillion-streaming/latest/docs/resources/pipeline) | Manages a streaming pipeline between a source and target |
| Data Source | [`matillion-streaming_agent_credentials`](https://registry.terraform.io/providers/matillion-public/matillion-streaming/latest/docs/data-sources/agent_credentials) | Retrieves the client credentials for a streaming agent |

### Supported Sources

- PostgreSQL
- MySQL
- Oracle
- SQL Server
- DB2 for IBM i

### Supported Targets

- Snowflake
- Amazon S3
- Azure Blob Storage



## Example

More examples covering additional source and target combinations can be found in the [`examples/`](./examples) folder.

The following example creates a streaming agent and a PostgreSQL-to-Snowflake pipeline:

```hcl
resource "matillion-streaming_agent" "example" {
  name           = "my-streaming-agent"
  description    = "Production streaming agent"
  deployment     = "fargate"
  cloud_provider = "aws"
}

data "matillion-streaming_agent_credentials" "example" {
  agent_id = matillion-streaming_agent.example.agent_id
}

resource "matillion-streaming_pipeline" "example" {
  name       = "postgres-to-snowflake"
  project_id = "your-project-id"
  agent_id   = matillion-streaming_agent.example.agent_id

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
      { schema = "public", table = "users" },
      { schema = "public", table = "orders" }
    ]
  }

  snowflake_target = {
    connection = {
      account_name = "myorg.us-east-1"
      username     = "STREAMING_USER"
      authentication = {
        private_key = {
          secret_type = "aws_secrets_manager"
          secret_name = "snowflake-private-key"
        }
      }
    }
    role                = "STREAMING_ROLE"
    warehouse           = "STREAMING_WH"
    database            = "ANALYTICS_DB"
    stage_schema        = "STREAMING_STAGE"
    stage_name          = "ANALYTICS_STAGE"
    table_schema        = "PUBLIC"
    table_prefix_type   = "none"
    transformation_type = "copy_table"
  }
}
```

## Development

### Building the Provider

```bash
make build
```

To install the provider into your local Terraform plugin cache:

```bash
make install
```

### Generating Documentation

Documentation under `docs/` is generated from provider schema and template files. To regenerate:

```bash
make generate
```

### Running Tests

Acceptance tests (requires real Matillion DPC credentials and sets `TF_ACC=1`):

```bash
export MATILLION_CLIENT_ID="..."
export MATILLION_CLIENT_SECRET="..."
export MATILLION_ACCOUNT_ID="..."
export MATILLION_REGION="..."
export MATILLION_PROJECT_ID="..."
export MATILLION_AGENT_ID="..."  # any valid UUID
make testacc
```

> **Note:** Acceptance tests create real resources in your Matillion DPC account.

