package models

import "github.com/hashicorp/terraform-plugin-framework/types"

// PipelineResourceModel represents the Terraform resource schema for a streaming pipeline.
// It maps the provider's resource configuration to the internal representation used by the provider.
type PipelineResourceModel struct {
	PipelineId         types.String            `tfsdk:"pipeline_id"`
	ProjectId          types.String            `tfsdk:"project_id"`
	Name               types.String            `tfsdk:"name"`
	AgentId            types.String            `tfsdk:"agent_id"`
	PostgresSource     *PostgresSourceModel    `tfsdk:"postgres_source"`
	SqlServerSource    *SqlServerSourceModel   `tfsdk:"sql_server_source"`
	MySqlSource        *MysqlSourceModel       `tfsdk:"mysql_source"`
	OracleSource       *OracleSourceModel      `tfsdk:"oracle_source"`
	Db2IbmISource      *Db2IbmISourceModel     `tfsdk:"db2_ibm_i_source"`
	SnowflakeTarget    *SnowflakeTargetModel   `tfsdk:"snowflake_target"`
	S3Target           *S3TargetModel          `tfsdk:"s3_target"`
	ABSTarget          *AbsTargetModel         `tfsdk:"abs_target"`
	AdvancedProperties map[string]types.String `tfsdk:"advanced_properties"`
}
