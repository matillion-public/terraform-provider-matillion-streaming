package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SnowflakeTargetSchema returns the complete Snowflake target schema
func SnowflakeTargetSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Snowflake target configuration",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"connection": schema.SingleNestedAttribute{
				MarkdownDescription: "Snowflake connection configuration",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"account_name": schema.StringAttribute{
						MarkdownDescription: "Snowflake account name",
						Required:            true,
					},
					"username": schema.StringAttribute{
						MarkdownDescription: "Snowflake username",
						Required:            true,
					},
					"authentication": schema.SingleNestedAttribute{
						MarkdownDescription: "Snowflake authentication configuration",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"private_key": func() schema.SingleNestedAttribute {
								ref := SnowflakeSecretReferenceSchema()
								ref.MarkdownDescription = "Private key secret reference"
								return ref
							}(),
							"passphrase": func() schema.SingleNestedAttribute {
								ref := SnowflakeSecretReferenceSchema()
								ref.MarkdownDescription = "Passphrase secret reference"
								ref.Required = false
								ref.Optional = true
								return ref
							}(),
						},
					},
					"jdbc_properties": schema.MapAttribute{
						MarkdownDescription: "JDBC properties",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "Snowflake role",
				Required:            true,
			},
			"warehouse": schema.StringAttribute{
				MarkdownDescription: "Snowflake warehouse",
				Required:            true,
			},
			"database": schema.StringAttribute{
				MarkdownDescription: "Snowflake database",
				Required:            true,
			},
			"stage_schema": schema.StringAttribute{
				MarkdownDescription: "Snowflake stage schema",
				Required:            true,
			},
			"stage_name": schema.StringAttribute{
				MarkdownDescription: "Snowflake stage name",
				Required:            true,
			},
			"stage_prefix": schema.StringAttribute{
				MarkdownDescription: "Snowflake stage prefix",
				Optional:            true,
			},
			"table_schema": schema.StringAttribute{
				MarkdownDescription: "Snowflake table schema",
				Required:            true,
			},
			"table_prefix_type": schema.StringAttribute{
				MarkdownDescription: "Table prefix type. Valid values: `prefix`, `source_database_and_schema`, `none`",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("prefix", "source_database_and_schema", "none"),
				},
			},
			"transformation_type": schema.StringAttribute{
				MarkdownDescription: "Transformation type. Valid values: `copy_table`, `copy_table_soft_delete`, `change_log`",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("copy_table", "copy_table_soft_delete", "change_log"),
				},
			},
			"temporal_mapping": schema.StringAttribute{
				MarkdownDescription: "Temporal mapping configuration. Valid values: `native`, `epoch`",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("native", "epoch"),
				},
			},
		},
	}
}

// S3TargetSchema returns the complete S3 target schema
func S3TargetSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "S3 target configuration",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"bucket": schema.StringAttribute{
				MarkdownDescription: "S3 bucket name",
				Required:            true,
			},
			"prefix": schema.StringAttribute{
				MarkdownDescription: "S3 prefix",
				Optional:            true,
			},
			"decimal_mapping": schema.StringAttribute{
				MarkdownDescription: "Decimal mapping configuration. Valid values: `logical`, `legacy`",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("logical", "legacy"),
				},
			},
		},
	}
}

// ABSTargetSchema returns the complete Azure Blob Storage target schema
func ABSTargetSchema() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		MarkdownDescription: "Azure Blob Storage target configuration",
		Optional:            true,
		Attributes: map[string]schema.Attribute{
			"container": schema.StringAttribute{
				MarkdownDescription: "Azure Blob Storage container name",
				Required:            true,
			},
			"prefix": schema.StringAttribute{
				MarkdownDescription: "ABS prefix",
				Optional:            true,
			},
			"account_name": schema.StringAttribute{
				MarkdownDescription: "Azure storage account name",
				Required:            true,
			},
			"account_key": func() schema.SingleNestedAttribute {
				ref := SecretReferenceSchema()
				ref.MarkdownDescription = "Secret reference for Azure storage account key"
				return ref
			}(),
			"decimal_mapping": schema.StringAttribute{
				MarkdownDescription: "Decimal mapping configuration. Valid values: `logical`, `legacy`",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("logical", "legacy"),
				},
			},
		},
	}
}
