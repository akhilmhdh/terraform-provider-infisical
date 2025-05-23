package resource

import (
	"context"
	"fmt"
	infisical "terraform-provider-infisical/internal/client"
	pkg "terraform-provider-infisical/internal/pkg/input"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource = &IntegrationAWSSecretsManagerResource{}
)

// NewIntegrationAwsSecretsManagerResource is a helper function to simplify the provider implementation.
func NewIntegrationAwsSecretsManagerResource() resource.Resource {
	return &IntegrationAWSSecretsManagerResource{}
}

// IntegrationAwsSecretsManager is the resource implementation.
type IntegrationAWSSecretsManagerResource struct {
	client *infisical.Client
}

type AwsSecretsManagerMetadataStruct struct {
	SecretAWSTag []infisical.AwsTag `json:"secretAWSTag,omitempty"`
	SecretPrefix string             `json:"secretPrefix,omitempty"`
}

type AwsSecretsManagerOptions struct {
	AwsTags          []infisical.AwsTag `tfsdk:"aws_tags" json:"secretAWSTag,omitempty"`
	SecretPrefix     *string            `tfsdk:"secret_prefix"`
	MetadataSyncMode *string            `tfsdk:"metadata_sync_mode"`
}

// projectResourceSourceModel describes the data source data model.
type IntegrationAWSSecretsManagerResourceModel struct {
	AccessKeyID     types.String `tfsdk:"access_key_id"`
	SecretAccessKey types.String `tfsdk:"secret_access_key"`
	AssumeRoleArn   types.String `tfsdk:"assume_role_arn"`

	ProjectID types.String `tfsdk:"project_id"`

	IntegrationAuthID types.String `tfsdk:"integration_auth_id"`
	IntegrationID     types.String `tfsdk:"integration_id"`

	Environment types.String `tfsdk:"environment"`
	SecretPath  types.String `tfsdk:"secret_path"`
	AWSRegion   types.String `tfsdk:"aws_region"`

	MappingBehavior types.String `tfsdk:"mapping_behavior"`
	AWSPath         types.String `tfsdk:"secrets_manager_path"`

	Options types.Object `tfsdk:"options"`
}

const METADATA_SYNC_MODE_SECRET_METADATA = "secret-metadata"

// Metadata returns the resource type name.
func (r *IntegrationAWSSecretsManagerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_integration_aws_secrets_manager"
}

// Schema defines the schema for the resource.
func (r *IntegrationAWSSecretsManagerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Create AWS Secrets Manager integration & save to Infisical. Only Machine Identity authentication is supported for this data source",
		Attributes: map[string]schema.Attribute{
			"options": schema.SingleNestedAttribute{
				Description: "Integration options",
				Optional:    true,
				Computed:    true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"secret_prefix":      types.StringType,
							"metadata_sync_mode": types.StringType,
							"aws_tags":           types.SetType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"key": types.StringType, "value": types.StringType}}},
						},
						map[string]attr.Value{
							"secret_prefix":      types.StringValue(""),
							"metadata_sync_mode": types.StringValue(""),
							"aws_tags":           types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{"key": types.StringType, "value": types.StringType}}),
						},
					),
				),
				Attributes: map[string]schema.Attribute{
					"secret_prefix": schema.StringAttribute{
						Optional:    true,
						Description: "The prefix to add to the secret name in AWS Secrets Manager.",
					},
					"aws_tags": schema.SetNestedAttribute{
						Description: "Tags to attach to the AWS Secrets Manager secrets.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Description: "The key of the tag.",
									Optional:    true,
								},
								"value": schema.StringAttribute{
									Description: "The value of the tag.",
									Optional:    true,
								},
							},
						},
					},
					"metadata_sync_mode": schema.StringAttribute{
						Optional:    true,
						Description: "The sync mode for AWS tags. The supported options are `secret-metadata` and `custom`. If `secret-metadata` is selected, the metadata of the Infisical secrets are used as tags in AWS (only supported for one-to-one integrations). If `custom` is selected, then the key/value pairs in the `aws_tags` field is used.",
					},
				},
			},

			"integration_auth_id": schema.StringAttribute{
				Computed:      true,
				Description:   "The ID of the integration auth, used internally by Infisical.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

			"integration_id": schema.StringAttribute{
				Computed:      true,
				Description:   "The ID of the integration, used internally by Infisical.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},

			"aws_region": schema.StringAttribute{
				Required:    true,
				Description: "The AWS region to sync secrets to. (us-east-1, us-east-2, etc)",
			},

			"access_key_id": schema.StringAttribute{
				Sensitive:   true,
				Optional:    true,
				Description: "The AWS access key ID. Used to authenticate with AWS Secrets Manager. You must either set secret_access_key and access_key_id, or set assume_role_arn to assume a role.",
			},

			"secret_access_key": schema.StringAttribute{
				Sensitive:   true,
				Optional:    true,
				Description: "The AWS secret access key. Used to authenticate with AWS Secrets Manager. You must either set secret_access_key and access_key_id, or set assume_role_arn to assume a role.",
			},

			"assume_role_arn": schema.StringAttribute{
				Optional:    true,
				Description: "The ARN of the role to assume when syncing secrets to AWS Secrets Manager. You must either set secret_access_key and access_key_id, or set assume_role_arn to assume a role.",
			},

			"project_id": schema.StringAttribute{
				Required:      true,
				Description:   "The ID of your Infisical project.",
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},

			"environment": schema.StringAttribute{
				Required:    true,
				Description: "The slug of the environment to sync to AWS Secrets Manager (prod, dev, staging, etc).",
			},

			"mapping_behavior": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(infisical.AWS_MAPPING_BEHAVIOR_MANY_TO_ONE),
				Description: "The behavior of the mapping. Can be 'many-to-one' or 'one-to-one'. Many to One: All Infisical secrets will be mapped to a single AWS secret. One to One: Each Infisical secret will be mapped to its own AWS secret.",
			},

			"secrets_manager_path": schema.StringAttribute{
				Optional:    true,
				Description: "The path in AWS Secrets Manager to sync secrets to. This is required if mapping_behavior is 'many-to-one'.",
			},

			"secret_path": schema.StringAttribute{
				Required:    true,
				Description: "The secret path in Infisical to sync secrets from.",
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *IntegrationAWSSecretsManagerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*infisical.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *IntegrationAWSSecretsManagerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if !r.client.Config.IsMachineIdentityAuth {
		resp.Diagnostics.AddError(
			"Unable to create integration",
			"Only Machine Identity authentication is supported for this operation",
		)
		return
	}

	// Retrieve values from plan
	var plan IntegrationAWSSecretsManagerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.MappingBehavior.ValueString() == infisical.AWS_MAPPING_BEHAVIOR_MANY_TO_ONE && (plan.AWSPath.IsNull() || plan.AWSPath.ValueString() == "") {
		resp.Diagnostics.AddError(
			"Invalid plan",
			"secrets_manager_path is required when mapping_behavior is 'many-to-one'",
		)
		return
	}

	if plan.MappingBehavior.ValueString() == infisical.AWS_MAPPING_BEHAVIOR_ONE_TO_ONE && (!plan.AWSPath.IsNull() && plan.AWSPath.ValueString() != "") {
		resp.Diagnostics.AddError(
			"Invalid plan",
			"secrets_manager_path should not be used when mapping_behavior is 'one-to-one'",
		)
		return
	}

	authMethod, err := pkg.ValidateAwsInputCredentials(plan.AccessKeyID, plan.SecretAccessKey, plan.AssumeRoleArn)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating AWS credentials",
			err.Error(),
		)
		return
	}

	createIntegrationAuthRequest := infisical.CreateIntegrationAuthRequest{
		ProjectID:   plan.ProjectID.ValueString(),
		Integration: infisical.IntegrationAuthTypeAwsSecretsManager,
	}

	if authMethod == pkg.AwsAuthMethodAccessKey {
		createIntegrationAuthRequest.AccessId = plan.AccessKeyID.ValueString()
		createIntegrationAuthRequest.AccessToken = plan.SecretAccessKey.ValueString()
	} else if authMethod == pkg.AwsAuthMethodAssumeRole {
		createIntegrationAuthRequest.AWSAssumeIamRoleArn = plan.AssumeRoleArn.ValueString()
	}

	// Create integration auth first
	auth, err := r.client.CreateIntegrationAuth(createIntegrationAuthRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create integration auth",
			err.Error(),
		)
		return
	}

	var planOptions AwsSecretsManagerOptions

	if !plan.Options.IsNull() {
		diags := plan.Options.As(ctx, &planOptions, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if plan.MappingBehavior.ValueString() != infisical.AWS_MAPPING_BEHAVIOR_ONE_TO_ONE && planOptions.MetadataSyncMode != nil && *planOptions.MetadataSyncMode == METADATA_SYNC_MODE_SECRET_METADATA {
		resp.Diagnostics.AddError(
			"Invalid plan",
			"cannot use secret metadata sync mode when mapping_behavior is 'many-to-one'",
		)
		return
	}

	// Convert metadata to map[string]interface{} if needed
	metadataMap := map[string]interface{}{}

	metadataMap["mappingBehavior"] = plan.MappingBehavior.ValueString()

	if planOptions.SecretPrefix != nil {
		metadataMap["secretPrefix"] = *planOptions.SecretPrefix
	} else {
		metadataMap["secretPrefix"] = ""
	}

	if planOptions.AwsTags != nil {
		metadataMap["secretAWSTag"] = planOptions.AwsTags
	} else {
		metadataMap["secretAWSTag"] = []infisical.AwsTag{}
	}

	if planOptions.MetadataSyncMode != nil && *planOptions.MetadataSyncMode != "" {
		metadataMap["metadataSyncMode"] = planOptions.MetadataSyncMode
	}

	request := infisical.CreateIntegrationRequest{
		IntegrationAuthID: auth.IntegrationAuth.ID,
		Region:            plan.AWSRegion.ValueString(),
		Metadata:          metadataMap,
		SecretPath:        plan.SecretPath.ValueString(),
		SourceEnvironment: plan.Environment.ValueString(),
	}

	if plan.MappingBehavior.ValueString() == infisical.AWS_MAPPING_BEHAVIOR_MANY_TO_ONE {
		request.App = plan.AWSPath.ValueString()
	}

	// Create the integration
	integration, err := r.client.CreateIntegration(request)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create integration",
			err.Error(),
		)
		return
	}

	plan.IntegrationAuthID = types.StringValue(auth.IntegrationAuth.ID)
	plan.IntegrationID = types.StringValue(integration.Integration.ID)
	plan.Environment = types.StringValue(integration.Integration.Environment.Slug)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *IntegrationAWSSecretsManagerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	if !r.client.Config.IsMachineIdentityAuth {
		resp.Diagnostics.AddError(
			"Unable to read integration",
			"Only Machine Identity authentication is supported for this operation",
		)
		return
	}

	// Get current state
	var state IntegrationAWSSecretsManagerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.client.GetIntegration(infisical.GetIntegrationRequest{
		ID: state.IntegrationID.ValueString(),
	})

	if err != nil {
		if err == infisical.ErrNotFound {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError(
				"Unable to read integration",
				err.Error(),
			)
		}
		return
	}

	var planOptions AwsSecretsManagerOptions

	if !state.Options.IsNull() {
		diags := state.Options.As(ctx, &planOptions, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if (planOptions.SecretPrefix != nil && integration.Integration.Metadata.SecretPrefix != *planOptions.SecretPrefix) || integration.Integration.Metadata.SecretPrefix != "" {
		planOptions.SecretPrefix = &integration.Integration.Metadata.SecretPrefix
	}

	if (planOptions.MetadataSyncMode != nil && integration.Integration.Metadata.MetadataSyncMode != *planOptions.MetadataSyncMode) || integration.Integration.Metadata.MetadataSyncMode != "" {
		planOptions.MetadataSyncMode = &integration.Integration.Metadata.MetadataSyncMode
	}

	if len(integration.Integration.Metadata.SecretAWSTag) > 0 {
		planOptions.AwsTags = integration.Integration.Metadata.SecretAWSTag
	}

	// Create a new types.Object from the modified planOptions
	optionsObj, diags := types.ObjectValueFrom(ctx, state.Options.AttributeTypes(ctx), planOptions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state.Options.
	state.Options = optionsObj
	state.SecretPath = types.StringValue(integration.Integration.SecretPath)
	state.IntegrationAuthID = types.StringValue(integration.Integration.IntegrationAuthID)
	state.Environment = types.StringValue(integration.Integration.Environment.Slug)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *IntegrationAWSSecretsManagerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if !r.client.Config.IsMachineIdentityAuth {
		resp.Diagnostics.AddError(
			"Unable to update integration",
			"Only Machine Identity authentication is supported for this operation",
		)
		return
	}

	// Retrieve values from plan
	var plan IntegrationAWSSecretsManagerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IntegrationAWSSecretsManagerResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var planOptions AwsSecretsManagerOptions

	if !plan.Options.IsNull() {
		diags := plan.Options.As(ctx, &planOptions, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if plan.MappingBehavior.ValueString() != infisical.AWS_MAPPING_BEHAVIOR_ONE_TO_ONE && planOptions.MetadataSyncMode != nil && *planOptions.MetadataSyncMode == METADATA_SYNC_MODE_SECRET_METADATA {
		resp.Diagnostics.AddError(
			"Invalid plan",
			"cannot use secret metadata sync mode when mapping_behavior is 'many-to-one'",
		)
		return
	}

	if plan.MappingBehavior.ValueString() == infisical.AWS_MAPPING_BEHAVIOR_MANY_TO_ONE && (plan.AWSPath.IsNull() || plan.AWSPath.ValueString() == "") {
		resp.Diagnostics.AddError(
			"Invalid plan",
			"secrets_manager_path is required when mapping_behavior is 'many-to-one'",
		)
		return
	}

	if plan.MappingBehavior.ValueString() == infisical.AWS_MAPPING_BEHAVIOR_ONE_TO_ONE && (!plan.AWSPath.IsNull() && plan.AWSPath.ValueString() != "") {
		resp.Diagnostics.AddError(
			"Invalid plan",
			"secrets_manager_path should not be used when mapping_behavior is 'one-to-one'",
		)
		return
	}

	updateIntegrationAuthRequest := infisical.UpdateIntegrationAuthRequest{
		Integration:       infisical.IntegrationAuthTypeAwsSecretsManager,
		IntegrationAuthId: plan.IntegrationAuthID.ValueString(),
	}

	authMethod, err := pkg.ValidateAwsInputCredentials(plan.AccessKeyID, plan.SecretAccessKey, plan.AssumeRoleArn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error validating AWS credentials",
			err.Error(),
		)
		return
	}

	if authMethod == pkg.AwsAuthMethodAccessKey {
		updateIntegrationAuthRequest.AccessId = plan.AccessKeyID.ValueString()
		updateIntegrationAuthRequest.AccessToken = plan.SecretAccessKey.ValueString()
	} else if authMethod == pkg.AwsAuthMethodAssumeRole {
		updateIntegrationAuthRequest.AWSAssumeIamRoleArn = plan.AssumeRoleArn.ValueString()
	}

	_, err = r.client.UpdateIntegrationAuth(updateIntegrationAuthRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating integration auth",
			err.Error(),
		)
		return
	}

	// Convert metadata to map[string]interface{} if needed
	metadataMap := map[string]interface{}{}

	if planOptions.SecretPrefix != nil {
		metadataMap["secretPrefix"] = planOptions.SecretPrefix
	} else {
		metadataMap["secretPrefix"] = ""
	}

	if planOptions.AwsTags != nil {
		metadataMap["secretAWSTag"] = planOptions.AwsTags
	} else {
		metadataMap["secretAWSTag"] = []infisical.AwsTag{}
	}

	if planOptions.MetadataSyncMode != nil && *planOptions.MetadataSyncMode != "" {
		metadataMap["metadataSyncMode"] = planOptions.MetadataSyncMode
	}

	updateIntegrationRequest := infisical.UpdateIntegrationRequest{
		ID:          state.IntegrationID.ValueString(),
		Metadata:    metadataMap,
		Environment: plan.Environment.ValueString(),
		SecretPath:  plan.SecretPath.ValueString(),
		Region:      plan.AWSRegion.ValueString(),
		IsActive:    true,
	}

	if plan.MappingBehavior.ValueString() == infisical.AWS_MAPPING_BEHAVIOR_MANY_TO_ONE {
		updateIntegrationRequest.App = plan.AWSPath.ValueString()
	}

	// Update the integration
	updatedIntegration, err := r.client.UpdateIntegration(updateIntegrationRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating integration",
			err.Error(),
		)
		return
	}

	plan.SecretPath = types.StringValue(updatedIntegration.Integration.SecretPath)
	plan.IntegrationAuthID = types.StringValue(updatedIntegration.Integration.IntegrationAuthID)
	plan.Environment = types.StringValue(updatedIntegration.Integration.Environment.Slug)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *IntegrationAWSSecretsManagerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if !r.client.Config.IsMachineIdentityAuth {
		resp.Diagnostics.AddError(
			"Unable to delete AWS Secrets Manager integration",
			"Only Machine Identity authentication is supported for this operation",
		)
		return
	}

	var state IntegrationAWSSecretsManagerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteIntegrationAuth(infisical.DeleteIntegrationAuthRequest{
		ID: state.IntegrationAuthID.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting AWS Secrets Manager Integration",
			"Couldn't delete AWS Secrets Manager integration from your Infiscial project, unexpected error: "+err.Error(),
		)
		return
	}
}
