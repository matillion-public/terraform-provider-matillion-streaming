package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// SnowflakeTargetModel represents a Snowflake target configuration in the provider schema
type SnowflakeTargetModel struct {
	Connection         SnowflakeConnectionModel `tfsdk:"connection"`
	Role               types.String             `tfsdk:"role"`
	Warehouse          types.String             `tfsdk:"warehouse"`
	Database           types.String             `tfsdk:"database"`
	StageSchema        types.String             `tfsdk:"stage_schema"`
	StageName          types.String             `tfsdk:"stage_name"`
	StagePrefix        types.String             `tfsdk:"stage_prefix"`
	TableSchema        types.String             `tfsdk:"table_schema"`
	TablePrefixType    types.String             `tfsdk:"table_prefix_type"`
	TransformationType types.String             `tfsdk:"transformation_type"`
	TemporalMapping    types.String             `tfsdk:"temporal_mapping"`
}

// S3TargetModel represents an Amazon S3 target configuration in the provider schema
type S3TargetModel struct {
	Bucket         types.String `tfsdk:"bucket"`
	Prefix         types.String `tfsdk:"prefix"`
	DecimalMapping types.String `tfsdk:"decimal_mapping"`
}

// AbsTargetModel represents an Azure Blob Storage target configuration in the provider schema
type AbsTargetModel struct {
	Container      types.String         `tfsdk:"container"`
	Prefix         types.String         `tfsdk:"prefix"`
	AccountName    types.String         `tfsdk:"account_name"`
	AccountKey     SecretReferenceModel `tfsdk:"account_key"`
	DecimalMapping types.String         `tfsdk:"decimal_mapping"`
}
