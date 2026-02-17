package provider

import (
	"context"
	"fmt"
	"terraform-provider-matillion-streaming/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &agentCredentialsDataSource{}
	_ datasource.DataSourceWithConfigure = &agentCredentialsDataSource{}
)

func NewAgentCredentialsDataSource() datasource.DataSource {
	return &agentCredentialsDataSource{}
}

type agentCredentialsDataSource struct {
	client *client.Client
}

type agentCredentialsDataSourceModel struct {
	AgentId      types.String `tfsdk:"agent_id"`
	ClientId     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (d *agentCredentialsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_agent_credentials"
}

func (d *agentCredentialsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the credentials for a streaming agent",
		Attributes: map[string]schema.Attribute{
			"agent_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the agent for which to retrieve credentials",
				Required:            true,
			},
			"client_id": schema.StringAttribute{
				MarkdownDescription: "The client ID for authenticating the agent",
				Computed:            true,
			},
			"client_secret": schema.StringAttribute{
				MarkdownDescription: "The client secret for authenticating the agent",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (d *agentCredentialsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config agentCredentialsDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	agentCredentials, err := d.client.Agents.GetCredentials(config.AgentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Agent Credentials",
			"Could not read credentials for agent id "+config.AgentId.ValueString()+": "+err.Error(),
		)
		return
	}

	config.ClientId = types.StringValue(agentCredentials.ClientId)
	config.ClientSecret = types.StringValue(agentCredentials.ClientSecret)

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *agentCredentialsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)

		return
	}

	d.client = c
}
