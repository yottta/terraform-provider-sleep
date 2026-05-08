package provider

import (
	"context"
	"terraform-provider-sleep/internal/provider/datasources"
	"terraform-provider-sleep/internal/provider/ephemerals"
	"terraform-provider-sleep/internal/provider/resources"
	"terraform-provider-sleep/internal/sleep"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider                       = &sleepProvider{}
	_ provider.ProviderWithEphemeralResources = &sleepProvider{}
	_ provider.ProviderWithFunctions          = &sleepProvider{}
)

type sleepProviderModel struct {
	MaxDefault types.Number `tfsdk:"max_default"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &sleepProvider{
			version: version,
		}
	}
}

type sleepProvider struct {
	version string
}

func (p *sleepProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "sleep"
	resp.Version = p.version
}

func (p *sleepProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"max_default": schema.NumberAttribute{
				Optional:    true,
				Description: "The number of seconds to be used as default maximum number of seconds a resource action can sleep. If no configured, it will use the provider default value of 60 seconds",
			},
		},
	}
}

func (p *sleepProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config sleepProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	sleeper := sleep.NewSleeper(config.MaxDefault)
	resp.DataSourceData = sleeper
	resp.ResourceData = sleeper
	resp.EphemeralResourceData = sleeper
}

func (p *sleepProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewSleeper,
	}
}

func (p *sleepProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewSleeper,
	}
}

func (p *sleepProvider) Functions(context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *sleepProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		ephemerals.NewSleeper,
	}
}
