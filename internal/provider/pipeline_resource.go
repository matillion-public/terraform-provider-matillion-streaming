package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-matillion-streaming/internal/client"
	"terraform-provider-matillion-streaming/internal/provider/models"
	"terraform-provider-matillion-streaming/internal/provider/schemas"
)

var (
	_ resource.Resource                     = &pipelineResource{}
	_ resource.ResourceWithConfigure        = &pipelineResource{}
	_ resource.ResourceWithConfigValidators = &pipelineResource{}
	_ resource.ResourceWithImportState      = &pipelineResource{}
)

type pipelineResource struct {
	client *client.Client
}

func NewPipelineResource() resource.Resource {
	return &pipelineResource{}
}

func (r *pipelineResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pipeline"
}

func (r *pipelineResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.PipelineSchema()
}

func (r *pipelineResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}

	r.client = apiClient
}

func (r *pipelineResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("snowflake_target"),
			path.MatchRoot("s3_target"),
			path.MatchRoot("abs_target"),
		),
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("postgres_source"),
			path.MatchRoot("sql_server_source"),
			path.MatchRoot("mysql_source"),
			path.MatchRoot("oracle_source"),
			path.MatchRoot("db2_ibm_i_source"),
		),
	}
}

func (r *pipelineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.PipelineResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "creating plan "+plan.Name.ValueString())

	projectId := plan.ProjectId
	name := plan.Name
	agent := plan.AgentId

	clientSource, err := mapSourceToClientSource(&plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error mapping source",
			"Failed to map source: "+err.Error(),
		)
		return
	}

	// Map the resource model to the client's target struct
	clientTarget, err := mapTargetToClientTarget(&plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error mapping target",
			"Failed to map target: "+err.Error(),
		)
		return
	}

	advancedProperties := make(map[string]string)
	for k, v := range plan.AdvancedProperties {
		advancedProperties[k] = v.ValueString()
	}

	opts := client.PipelineOptions{
		Name:               name.ValueString(),
		AgentId:            agent.ValueString(),
		StreamingSource:    clientSource,
		StreamingTarget:    clientTarget,
		AdvancedProperties: advancedProperties,
	}

	id, err := r.client.Pipelines.Create(projectId.ValueString(), opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating streaming pipeline",
			"Failed to create streaming pipeline: "+err.Error(),
		)
		return
	}

	tflog.Info(ctx, "created plan "+plan.Name.ValueString()+" id="+id)

	plan.PipelineId = types.StringValue(id)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *pipelineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.PipelineResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pipeline, err := r.client.Pipelines.Get(state.ProjectId.ValueString(), state.PipelineId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Streaming Pipeline",
			"Could not read streaming pipeline id "+state.PipelineId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map basic fields
	state.PipelineId = types.StringValue(pipeline.StreamingPipelineId)
	state.Name = types.StringValue(pipeline.Name)
	state.AgentId = types.StringValue(pipeline.AgentId)

	// Map advanced properties
	if pipeline.AdvancedProperties != nil {
		state.AdvancedProperties = make(map[string]types.String)
		for k, v := range pipeline.AdvancedProperties {
			state.AdvancedProperties[k] = types.StringValue(v)
		}
	}

	// Map source from API response
	if err = mapSourceFromAPI(pipeline.StreamingSource, &state); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Streaming Pipeline Source",
			"Could not map source from streaming pipeline "+state.PipelineId.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map target from API response
	if err = mapTargetFromAPI(pipeline.StreamingTarget, &state); err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Streaming Pipeline Target",
			"Could not map target from streaming pipeline "+state.PipelineId.ValueString()+": "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *pipelineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.PipelineResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clientSource, err := mapSourceToClientSource(&plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error mapping source",
			"Failed to map source: "+err.Error(),
		)
		return
	}

	clientTarget, err := mapTargetToClientTarget(&plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error mapping target",
			"Failed to map target: "+err.Error(),
		)
		return
	}

	advancedProperties := make(map[string]string)
	for k, v := range plan.AdvancedProperties {
		advancedProperties[k] = v.ValueString()
	}

	opts := client.PipelineOptions{
		Name:               plan.Name.ValueString(),
		AgentId:            plan.AgentId.ValueString(),
		StreamingSource:    clientSource,
		StreamingTarget:    clientTarget,
		AdvancedProperties: advancedProperties,
	}

	_, err = r.client.Pipelines.Replace(plan.ProjectId.ValueString(), plan.PipelineId.ValueString(), opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating streaming pipeline",
			"Could not update streaming pipeline id "+plan.PipelineId.ValueString()+": "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *pipelineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.PipelineResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Pipelines.Delete(state.ProjectId.ValueString(), state.PipelineId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting streaming pipeline",
			"Could not delete streaming pipeline id "+state.PipelineId.ValueString()+": "+err.Error(),
		)
	}
}

func (r *pipelineResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Parse the import ID format: project_id:pipeline_id
	parts := strings.Split(req.ID, ":")
	if len(parts) != 2 {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Import ID must be in the format: project_id:pipeline_id",
		)
		return
	}

	projectId := parts[0]
	pipelineId := parts[1]

	// Set the attributes in the state
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("pipeline_id"), pipelineId)...)
}
