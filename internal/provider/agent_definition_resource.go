package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-matillion-streaming/internal/client"
)

var (
	// Ensure provider defined types fully satisfy framework interfaces
	_ resource.Resource                = &agentDefinitionResource{}
	_ resource.ResourceWithConfigure   = &agentDefinitionResource{}
	_ resource.ResourceWithImportState = &agentDefinitionResource{}
)

// agentDefinitionResource defines the streaming agent resource implementation.
type agentDefinitionResource struct {
	client *client.Client
}

// agentDefinitionResourceModel maps streaming agent resource schema data.
type agentDefinitionResourceModel struct {
	AgentId       types.String `tfsdk:"agent_id"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	Deployment    types.String `tfsdk:"deployment"`
	CloudProvider types.String `tfsdk:"cloud_provider"`
}

// NewAgentDefinitionResource returns a new instance of the streaming agent resource.
func NewAgentDefinitionResource() resource.Resource {
	return &agentDefinitionResource{}
}

func (r *agentDefinitionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent"
}

func (r *agentDefinitionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a streaming agent definition in DPC",
		Attributes: map[string]schema.Attribute{
			"agent_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The unique identifier of the agent.",
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the agent (1 - 30 characters).",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.LengthAtMost(30),
				},
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A description for the agent.",
				// When the description is not provided, the API returns an empty string rather than null.
				// Setting a default empty string with Computed:true ensures Terraform won't detect a drift
				// between the planned state (null) and the applied state (empty string).
				Default: stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(500),
				},
			},
			"cloud_provider": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The cloud provider for the agent. Supported values: aws, azure, gcp. This value cannot be changed after creation.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"aws",
						"azure",
						"gcp",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"deployment": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The deployment type for the agent. Supported values: fargate, eks, aci, aks, gke, gce. This value cannot be changed after creation.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"fargate",
						"eks",
						"aci",
						"aks",
						"gke",
						"gce",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *agentDefinitionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData))
		return
	}

	r.client = apiClient
}

func (r *agentDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan agentDefinitionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := client.CreateAgentOptions{
		Name:          plan.Name.ValueString(),
		Deployment:    plan.Deployment.ValueString(),
		CloudProvider: plan.CloudProvider.ValueString(),
		AgentType:     "streaming",
	}

	if !plan.Description.IsNull() {
		opts.Description = plan.Description.ValueString()
	}

	agentId, err := r.client.Agents.Create(opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating streaming agent definition",
			"Failed to create streaming agent definition: "+err.Error(),
		)
		return
	}

	agentDefinition, err := r.client.Agents.Get(agentId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading newly created streaming agent definition",
			"Failed to read newly created streaming agent definition: "+err.Error(),
		)
		return
	}

	state := agentDefinitionResourceModel{
		AgentId:       types.StringValue(agentDefinition.AgentId),
		Name:          types.StringValue(agentDefinition.Name),
		Description:   types.StringValue(agentDefinition.Description),
		Deployment:    types.StringValue(agentDefinition.Deployment),
		CloudProvider: types.StringValue(agentDefinition.CloudProvider),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *agentDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state agentDefinitionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	agentDefinition, err := r.client.Agents.Get(state.AgentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Streaming Agent Definition",
			"Could not read streaming agent definition id "+state.AgentId.ValueString()+": "+err.Error(),
		)
		return
	}

	if agentDefinition.AgentType != "streaming" {
		resp.Diagnostics.AddError(
			"Error Reading Streaming Agent Definition",
			"Resource with ID "+state.AgentId.ValueString()+" is not a streaming agent",
		)
		return
	}

	state.Name = types.StringValue(agentDefinition.Name)
	state.Description = types.StringValue(agentDefinition.Description)
	state.Deployment = types.StringValue(agentDefinition.Deployment)
	state.CloudProvider = types.StringValue(agentDefinition.CloudProvider)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *agentDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan agentDefinitionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	opts := client.UpdateAgentOptions{
		Name: plan.Name.ValueString(),
	}

	if !plan.Description.IsNull() {
		opts.Description = plan.Description.ValueString()
	}

	err := r.client.Agents.Update(plan.AgentId.ValueString(), opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating streaming agent definition",
			"Could not update streaming agent definition id "+plan.AgentId.ValueString()+": "+err.Error(),
		)
		return
	}

	agentDefinition, err := r.client.Agents.Get(plan.AgentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading updated streaming agent definition",
			"Failed to read updated streaming agent definition: "+err.Error(),
		)
		return
	}

	state := agentDefinitionResourceModel{
		AgentId:       types.StringValue(agentDefinition.AgentId),
		Name:          types.StringValue(agentDefinition.Name),
		Description:   types.StringValue(agentDefinition.Description),
		Deployment:    types.StringValue(agentDefinition.Deployment),
		CloudProvider: types.StringValue(agentDefinition.CloudProvider),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *agentDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state agentDefinitionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.Agents.Delete(state.AgentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Streaming Agent Definition",
			"Could not delete streaming agent definition, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *agentDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	agentId := req.ID

	agentDefinition, err := r.client.Agents.Get(agentId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Streaming Agent Definition",
			"Could not import streaming agent definition id "+agentId+": "+err.Error(),
		)
		return
	}

	if agentDefinition.AgentType != "streaming" {
		resp.Diagnostics.AddError(
			"Error Importing Streaming Agent Definition",
			"Resource with ID "+agentId+" is not a streaming agent",
		)
		return
	}

	state := agentDefinitionResourceModel{
		AgentId:       types.StringValue(agentDefinition.AgentId),
		Name:          types.StringValue(agentDefinition.Name),
		Description:   types.StringValue(agentDefinition.Description),
		Deployment:    types.StringValue(agentDefinition.Deployment),
		CloudProvider: types.StringValue(agentDefinition.CloudProvider),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
