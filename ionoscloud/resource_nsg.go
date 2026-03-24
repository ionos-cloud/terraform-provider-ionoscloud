package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

func resourceNSG() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNSGCreate,
		ReadContext:   resourceNSGRead,
		UpdateContext: resourceNSGUpdate,
		DeleteContext: resourceNSGDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceNSGImport,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"datacenter_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
			},
			"location": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The location of the resource. This field should be used only if you are also using a file configuration and should not be configured otherwise.",
				ForceNew:    true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceNSGCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	datacenterID := d.Get("datacenter_id").(string)
	sgName := d.Get("name").(string)
	sgDescription := d.Get("description").(string)
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	sg := ionoscloud.SecurityGroupRequest{
		Properties: &ionoscloud.SecurityGroupProperties{
			Name:        &sgName,
			Description: &sgDescription,
		},
	}

	securityGroup, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsPost(ctx, datacenterID).SecurityGroup(sg).Execute()
	apiResponse.LogInfo()
	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while creating a Network Security Group for datacenter dcID: %s, %w", datacenterID, err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.StatusCode})
	}
	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutCreate); errState != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while waiting for Network Security Group to be created for datacenter dcID: %s,  %w", datacenterID, errState), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutCreate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}
	d.SetId(*securityGroup.Id)

	return diagutil.ToDiags(d, setNSGData(d, &securityGroup), nil)
}

func resourceNSGRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	datacenterID := d.Get("datacenter_id").(string)
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	securityGroup, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFindById(ctx, datacenterID, d.Id()).Depth(2).Execute()
	apiResponse.LogInfo()
	if err != nil {
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while retrieving a network security group: %w", err), nil)
	}

	if err := setNSGData(d, &securityGroup); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}
	return nil
}

func resourceNSGUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	datacenterID := d.Get("datacenter_id").(string)
	sgName := d.Get("name").(string)
	sgDescription := d.Get("description").(string)
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	sg := ionoscloud.SecurityGroupRequest{
		Properties: &ionoscloud.SecurityGroupProperties{
			Name:        &sgName,
			Description: &sgDescription,
		},
	}

	_, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsPut(ctx, datacenterID, d.Id()).SecurityGroup(sg).Execute()
	apiResponse.LogInfo()
	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while updating network security group: dcID: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.StatusCode})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutUpdate); errState != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while waiting for Network Security Group to be updated for datacenter dcID: %s,  %w", datacenterID, errState), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutUpdate).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return resourceNSGRead(ctx, d, meta)
}

func resourceNSGDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	datacenterID := d.Get("datacenter_id").(string)
	location := d.Get("location").(string)

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return diag.FromErr(err)
	}

	apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsDelete(ctx, datacenterID, d.Id()).Execute()
	apiResponse.LogInfo()
	if err != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while deleting a network security group: %w", err), &diagutil.ErrorContext{RequestID: diagutil.ExtractRequestID(requestLocation), StatusCode: apiResponse.StatusCode})
	}

	if errState := bundleclient.WaitForStateChange(ctx, meta, d, apiResponse, schema.TimeoutDelete); errState != nil {
		requestLocation, _ := apiResponse.Location()
		return diagutil.ToDiags(d, fmt.Errorf("an error occurred while waiting for Network Security Group to be deleted for datacenter dcID: %s,  %w", datacenterID, errState), &diagutil.ErrorContext{Timeout: d.Timeout(schema.TimeoutDelete).String(), RequestID: diagutil.ExtractRequestID(requestLocation)})
	}

	return nil
}

func resourceNSGImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	importID := d.Id()

	location, parts := splitImportID(importID, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid import identifier: expected one of <location>:<datacenter-id>/<nsg-id> or <datacenter-id>/<nsg-id>, got: %s", importID)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, diagutil.ToError(d, fmt.Errorf("failed validating import identifier %q: %w", importID, err), nil)
	}

	datacenterID := parts[0]
	nsgID := parts[1]

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(location)
	if err != nil {
		return nil, err
	}

	nsg, apiResponse, err := client.SecurityGroupsApi.DatacentersSecuritygroupsFindById(ctx, datacenterID, nsgID).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("unable to find Network Security Group %q", nsgID), nil)
		}
		return nil, diagutil.ToError(d, fmt.Errorf("an error occurred while retrieving the Network Security Group, %w", err), nil)
	}

	log.Printf("[INFO] Datacenter found: %+v", nsg)
	if err = d.Set("datacenter_id", datacenterID); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}
	if err = d.Set("location", location); err != nil {
		return nil, err
	}
	if err = setNSGData(d, &nsg); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	return []*schema.ResourceData{d}, nil
}

func setNSGData(d *schema.ResourceData, securityGroup *ionoscloud.SecurityGroup) error {

	if securityGroup.Id != nil {
		d.SetId(*securityGroup.Id)
	}

	if securityGroup.Properties != nil {
		if securityGroup.Properties.Name != nil {
			err := d.Set("name", *securityGroup.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting Network Security Group name  %s: %w", d.Id(), err)
			}
		}
		if securityGroup.Properties.Description != nil {
			err := d.Set("description", *securityGroup.Properties.Description)
			if err != nil {
				return fmt.Errorf("error while setting Network Security Group description  %s: %w", d.Id(), err)
			}
		}
	}
	var ruleIDs []string
	if securityGroup.Entities != nil {
		if securityGroup.Entities.Rules != nil && securityGroup.Entities.Rules.Items != nil {
			ruleIDs = make([]string, 0, len(*securityGroup.Entities.Rules.Items))
			for _, rule := range *securityGroup.Entities.Rules.Items {
				ruleIDs = append(ruleIDs, *rule.Id)
			}
		}
	}
	if err := d.Set("rule_ids", ruleIDs); err != nil {
		return fmt.Errorf("error while setting rule_ids property for NSG  %s: %w", d.Id(), err)
	}
	return nil
}
