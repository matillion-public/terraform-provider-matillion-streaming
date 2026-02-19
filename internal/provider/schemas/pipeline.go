package schemas

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// PipelineSchema returns the complete pipeline resource schema
func PipelineSchema() schema.Schema {
	return schema.Schema{
		MarkdownDescription: "Manages a streaming pipeline for data transfer between sources and targets",
		Attributes: map[string]schema.Attribute{
			"pipeline_id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the streaming pipeline",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the project containing this streaming pipeline",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the streaming pipeline",
				Required:            true,
			},
			"agent_id": schema.StringAttribute{
				MarkdownDescription: "ID of the agent to use",
				Required:            true,
			},
			"postgres_source":   PostgresSourceSchema(),
			"sql_server_source": SQLServerSourceSchema(),
			"mysql_source":      MySQLSourceSchema(),
			"oracle_source":     OracleSourceSchema(),
			"db2_ibm_i_source":  DB2IbmISourceSchema(),
			"snowflake_target":  SnowflakeTargetSchema(),
			"s3_target":         S3TargetSchema(),
			"abs_target":        ABSTargetSchema(),
			"advanced_properties": schema.MapAttribute{
				MarkdownDescription: "Advanced configuration properties for the streaming pipeline",
				Optional:            true,
				ElementType:         types.StringType,
			},
		},
	}
}
