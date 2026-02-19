package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SecretReferenceSchema returns a reusable schema for secret references
func SecretReferenceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Secret reference for password",
		Required:            true,
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				MarkdownDescription: "Secret reference type. Valid values: `aws_secrets_manager`, `azure_key_vault`, `google_secret_manager`",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("aws_secrets_manager", "azure_key_vault", "google_secret_manager"),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Secret reference name",
				Required:            true,
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "Secret reference key",
				Optional:            true,
			},
		},
	}
}

// SnowflakeSecretReferenceSchema returns a reusable schema for Snowflake secret references (no key field)
func SnowflakeSecretReferenceSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Required: true,
		Attributes: map[string]schema.Attribute{
			"secret_type": schema.StringAttribute{
				MarkdownDescription: "Secret type. Valid values: `aws_secrets_manager`, `azure_key_vault`, `google_secret_manager`",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("aws_secrets_manager", "azure_key_vault", "google_secret_manager"),
				},
			},
			"secret_name": schema.StringAttribute{
				MarkdownDescription: "Secret name",
				Required:            true,
			},
		},
	}
}

// BasicDatabaseConnectionSchema returns a reusable schema for basic database connections
func BasicDatabaseConnectionSchema(dbType string) schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Required: true,
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				MarkdownDescription: dbType + " host",
				Required:            true,
			},
			"port": schema.Int32Attribute{
				MarkdownDescription: dbType + " port",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: dbType + " username",
				Required:            true,
			},
			"password": SecretReferenceSchema(),
			"jdbc_properties": schema.MapAttribute{
				MarkdownDescription: "JDBC properties",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}
