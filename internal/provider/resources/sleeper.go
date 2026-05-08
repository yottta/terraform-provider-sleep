package resources

import (
	"context"
	"fmt"
	"terraform-provider-sleep/internal/sleep"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource               = &sleeper{}
	_ resource.ResourceWithConfigure  = &sleeper{}
	_ resource.ResourceWithModifyPlan = &sleeper{}
)

type sleeperModel struct {
	Id              types.String `tfsdk:"id"`
	PlanSleep       types.Number `tfsdk:"plan_sleep"`
	ApplySleep      types.Number `tfsdk:"apply_sleep"`
	DestroySleep    types.Number `tfsdk:"destroy_sleep"`
	ReadSleep       types.Number `tfsdk:"read_sleep"`
	UpdateSleep     types.Number `tfsdk:"update_sleep"`
	String          types.String `tfsdk:"string"`
	StringWo        types.String `tfsdk:"string_wo"`
	StringWoVersion types.Number `tfsdk:"string_wo_version"`
}

func NewSleeper() resource.Resource {
	return &sleeper{}
}

type sleeper struct {
	sleeper sleep.Sleeper
}

// Metadata returns the resource type name.
func (r *sleeper) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sleeper"
}

// Schema defines the schema for the resource.
func (r *sleeper) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The id of the resource. This will be generated on creation",
			},
			"plan_sleep": schema.NumberAttribute{
				Optional:    true,
				Description: "The number of seconds the provider should sleep before returning from ModifyPlan. If not provided, it will use a random value between [0,N] where N is the maximum number of seconds configured at the provider level.",
			},
			"apply_sleep": schema.NumberAttribute{
				Optional:    true,
				Description: "The number of seconds the provider should sleep before returning from Create. If not provided, it will use a random value between [0,N] where N is the maximum number of seconds configured at the provider level.",
			},
			"destroy_sleep": schema.NumberAttribute{
				Optional:    true,
				Description: "The number of seconds the provider should sleep before returning from Delete. If not provided, it will use a random value between [0,N] where N is the maximum number of seconds configured at the provider level.",
			},
			"read_sleep": schema.NumberAttribute{
				Optional:    true,
				Description: "The number of seconds the provider should sleep before returning from Read. If not provided, it will use a random value between [0,N] where N is the maximum number of seconds configured at the provider level.",
			},
			"update_sleep": schema.NumberAttribute{
				Optional:    true,
				Description: "The number of seconds the provider should sleep before returning from Update. If not provided, it will use a random value between [0,N] where N is the maximum number of seconds configured at the provider level.",
			},
			"string": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.PreferWriteOnlyAttribute(path.MatchRoot("string_wo")),
					stringvalidator.ConflictsWith(path.MatchRoot("string_wo")),
				},
				Description: "This is an attribute that can be configured to store the information in the state",
			},
			"string_wo": schema.StringAttribute{
				Optional:  true,
				WriteOnly: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("string")),
					stringvalidator.AlsoRequires(path.MatchRoot("string_wo_version")),
				},
				Description: "This is an attribute that can be configured with ephemeral values. This value will be null in the returned state",
			},
			"string_wo_version": schema.NumberAttribute{
				Optional:    true,
				WriteOnly:   true,
				Description: "Value associated with `string_wo`",
			},
		},
	}
}

func (r *sleeper) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var config *sleeperModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config != nil {
		r.sleeper.Sleep(ctx, config.PlanSleep)
	}
}

func (r *sleeper) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan, config sleeperModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.sleeper.Sleep(ctx, config.ApplySleep)

	if plan.Id.IsNull() || plan.Id.IsUnknown() {
		plan.Id = types.StringValue(uuid.NewString())
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sleeper) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sleeperModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.sleeper.Sleep(ctx, state.ReadSleep)
}

func (r *sleeper) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sleeperModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Id.IsNull() || plan.Id.IsUnknown() {
		plan.Id = types.StringValue(uuid.NewString())
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.sleeper.Sleep(ctx, plan.UpdateSleep)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sleeper) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sleeperModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.sleeper.Sleep(ctx, state.DestroySleep)
}

func (r *sleeper) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(sleep.Sleeper)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Sleeper type",
			fmt.Sprintf("Expected sleep.Sleeper, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.sleeper = c
}
