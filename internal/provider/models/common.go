package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// TableResourceModel represents a database table reference in the provider schema
type TableResourceModel struct {
	Schema types.String `tfsdk:"schema"`
	Table  types.String `tfsdk:"table"`
}

// SecretReferenceModel represents a reference to a secret stored in an external system
type SecretReferenceModel struct {
	Type types.String `tfsdk:"type"`
	Name types.String `tfsdk:"name"`
	Key  types.String `tfsdk:"key"`
}

// SnowflakeSecretReferenceModel represents a secret reference for Snowflake
type SnowflakeSecretReferenceModel struct {
	SecretType types.String `tfsdk:"secret_type"`
	SecretName types.String `tfsdk:"secret_name"`
}
