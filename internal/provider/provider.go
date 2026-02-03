package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-matillion-streaming/internal/client"
)

// MatillionStreamingProvider defines the provider implementation.
type MatillionStreamingProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type matillionStreamingProviderModel struct {
	AccountId types.String `tfsdk:"account_id"`
	Region    types.String `tfsdk:"region"`
}

func (p MatillionStreamingProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "matillion-streaming"
	resp.Version = p.version
}

func (p MatillionStreamingProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"region": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Region to use (eu or us).",
				Validators: []validator.String{
					stringvalidator.OneOf("eu", "us"),
				},
			},
		},
	}
}

func (p MatillionStreamingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Provider")
	var config matillionStreamingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	accountId := config.AccountId.ValueString()
	ctx = tflog.SetField(ctx, "account_id", accountId)

	// Determine Region
	var region client.Region
	regionStr := config.Region.ValueString()

	if regionStr == "eu" {
		region = client.RegionEU
	} else { // regionStr == "us" (validated by OneOf)
		region = client.RegionUS
	}

	ctx = tflog.SetField(ctx, "region", region)
	c, err := client.NewClient(accountId, region)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Matillion Streaming API Client",
			"Failed to create authenticated Matillion Streaming API Client: "+err.Error()+
				"\n\nEnsure MATILLION_CLIENT_ID and MATILLION_CLIENT_SECRET environment variables are set",
		)
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c

	tflog.Info(ctx, "Configured client with region: "+string(region))
}

func (p MatillionStreamingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p MatillionStreamingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAgentDefinitionResource,
		NewPipelineResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MatillionStreamingProvider{
			version: version,
		}
	}
}
