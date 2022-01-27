package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGroupImporter,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"create_datacenter": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"reserve_ip": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"access_activity_log": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_pcc": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"s3_privilege": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_backup_unit": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_internet_access": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_k8s_cluster": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"create_flow_log": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"access_and_manage_monitoring": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"access_and_manage_certificates": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"first_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"administrator": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"force_sec_auth": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	request := ionoscloud.Group{
		Properties: &ionoscloud.GroupProperties{},
	}

	log.Printf("[DEBUG] NAME %s", d.Get("name"))
	groupName := d.Get("name").(string)
	if d.Get("name") != nil {
		request.Properties.Name = &groupName
	}

	tempCreateDataCenter := d.Get("create_datacenter").(bool)
	request.Properties.CreateDataCenter = &tempCreateDataCenter
	tempCreateSnapshot := d.Get("create_snapshot").(bool)
	request.Properties.CreateSnapshot = &tempCreateSnapshot
	tempReserveIp := d.Get("reserve_ip").(bool)
	request.Properties.ReserveIp = &tempReserveIp
	tempAccessActivityLog := d.Get("access_activity_log").(bool)
	request.Properties.AccessActivityLog = &tempAccessActivityLog
	tempCreatePcc := d.Get("create_pcc").(bool)
	request.Properties.CreatePcc = &tempCreatePcc
	tempS3Privilege := d.Get("s3_privilege").(bool)
	request.Properties.S3Privilege = &tempS3Privilege
	tempCreateBackupUnit := d.Get("create_backup_unit").(bool)
	request.Properties.CreateBackupUnit = &tempCreateBackupUnit
	tempCreateInternetAccess := d.Get("create_internet_access").(bool)
	request.Properties.CreateInternetAccess = &tempCreateInternetAccess
	tempCreateK8sCluster := d.Get("create_k8s_cluster").(bool)
	request.Properties.CreateK8sCluster = &tempCreateK8sCluster
	tempCreateFlowLog := d.Get("create_flow_log").(bool)
	request.Properties.CreateFlowLog = &tempCreateFlowLog
	tempAccessAndManageMonitoring := d.Get("access_and_manage_monitoring").(bool)
	request.Properties.AccessAndManageMonitoring = &tempAccessAndManageMonitoring
	tempAccessAndManageCertificates := d.Get("access_and_manage_certificates").(bool)
	request.Properties.AccessAndManageCertificates = &tempAccessAndManageCertificates

	usertoAdd := d.Get("user_id").(string)

	group, apiResponse, err := client.UserManagementApi.UmGroupsPost(ctx).Group(request).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while creating a group: %w", err))
		return diags
	}

	log.Printf("[DEBUG] GROUP ID: %s", *group.Id)

	d.SetId(*group.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		diags := diag.FromErr(errState)
		return diags
	}

	//add users to group if any is provided
	if usertoAdd != "" {
		user := ionoscloud.User{
			Id: &usertoAdd,
		}
		_, apiResponse, err := client.UserManagementApi.UmGroupsUsersPost(ctx, d.Id()).User(user).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while adding %s user to group ID %s %w", usertoAdd, d.Id(), err))
			return diags
		}
		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(errState)
			return diags
		}
	}
	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	group, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("an error occured while fetching a Group ID %s %w", d.Id(), err))
		return diags
	}

	if err := setGroupData(ctx, client, d, &group); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	tempCreateDataCenter := d.Get("create_datacenter").(bool)
	tempCreateSnapshot := d.Get("create_snapshot").(bool)
	tempReserveIp := d.Get("reserve_ip").(bool)
	tempAccessActivityLog := d.Get("access_activity_log").(bool)
	tempCreatePcc := d.Get("create_pcc").(bool)
	tempS3Privilege := d.Get("s3_privilege").(bool)
	tempCreateBackupUnit := d.Get("create_backup_unit").(bool)
	tempCreateInternetAccess := d.Get("create_internet_access").(bool)
	tempCreateK8sCluster := d.Get("create_k8s_cluster").(bool)
	tempCreateFlowLog := d.Get("create_flow_log").(bool)
	tempAccessAndManageMonitoring := d.Get("access_and_manage_monitoring").(bool)
	tempAccessAndManageCertificates := d.Get("access_and_manage_certificates").(bool)

	userToAdd := d.Get("user_id").(string)

	groupReq := ionoscloud.Group{
		Properties: &ionoscloud.GroupProperties{
			CreateDataCenter:            &tempCreateDataCenter,
			CreateSnapshot:              &tempCreateSnapshot,
			ReserveIp:                   &tempReserveIp,
			AccessActivityLog:           &tempAccessActivityLog,
			CreatePcc:                   &tempCreatePcc,
			S3Privilege:                 &tempS3Privilege,
			CreateBackupUnit:            &tempCreateBackupUnit,
			CreateInternetAccess:        &tempCreateInternetAccess,
			CreateK8sCluster:            &tempCreateK8sCluster,
			CreateFlowLog:               &tempCreateFlowLog,
			AccessAndManageMonitoring:   &tempAccessAndManageMonitoring,
			AccessAndManageCertificates: &tempAccessAndManageCertificates,
		},
	}

	_, newValue := d.GetChange("name")
	newValueStr := newValue.(string)
	groupReq.Properties.Name = &newValueStr

	_, apiResponse, err := client.UserManagementApi.UmGroupsPut(ctx, d.Id()).Group(groupReq).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while patching a group ID %s %w", d.Id(), err))
		return diags
	}
	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	//add users to group if any is provided
	if userToAdd != "" {

		user := ionoscloud.User{
			Id: &userToAdd,
		}

		_, apiResponse, err := client.UserManagementApi.UmGroupsUsersPost(ctx, d.Id()).User(user).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("an error occured while adding %s user to group ID %s %w", userToAdd, d.Id(), err))
			return diags
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
		if errState != nil {
			diags := diag.FromErr(errState)
			return diags
		}
	}
	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	apiResponse, err := client.UserManagementApi.UmGroupsDelete(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		diags := diag.FromErr(err)
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	d.SetId("")
	return nil
}

func resourceGroupImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(SdkBundle).CloudApiClient

	grpId := d.Id()

	group, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, grpId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if apiResponse != nil && apiResponse.Response != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil, fmt.Errorf("an error occured while trying to fetch the group %q", grpId)
		}
		return nil, fmt.Errorf("group does not exist%q", grpId)
	}

	log.Printf("[INFO] group found: %+v", group)

	if err := setGroupData(ctx, client, d, &group); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

func setGroupData(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData, group *ionoscloud.Group) error {

	if group.Id != nil {
		d.SetId(*group.Id)
	}

	if group.Properties != nil {
		if group.Properties.Name != nil {
			err := d.Set("name", *group.Properties.Name)
			if err != nil {
				return fmt.Errorf("error while setting name property for group %s: %s", d.Id(), err)
			}
		}

		if group.Properties.CreateDataCenter != nil {
			err := d.Set("create_datacenter", *group.Properties.CreateDataCenter)
			if err != nil {
				return fmt.Errorf("error while setting create_datacenter property for group %s: %s", d.Id(), err)
			}
		}

		if group.Properties.CreateSnapshot != nil {
			err := d.Set("create_snapshot", *group.Properties.CreateSnapshot)
			if err != nil {
				return fmt.Errorf("error while setting create_snapshot property for group %s: %s", d.Id(), err)
			}
		}

		if group.Properties.ReserveIp != nil {
			err := d.Set("reserve_ip", *group.Properties.ReserveIp)
			if err != nil {
				return fmt.Errorf("error while setting reserve_ip property for group %s: %s", d.Id(), err)
			}
		}

		if group.Properties.AccessActivityLog != nil {
			err := d.Set("access_activity_log", *group.Properties.AccessActivityLog)
			if err != nil {
				return fmt.Errorf("error while setting access_activity_log property for group %s: %s", d.Id(), err)
			}
		}

		if group.Properties.CreatePcc != nil {
			err := d.Set("create_pcc", *group.Properties.CreatePcc)
			if err != nil {
				return fmt.Errorf("error while setting create_pcc property for group %s: %s", d.Id(), err)
			}
		}

		if group.Properties.S3Privilege != nil {
			err := d.Set("s3_privilege", *group.Properties.S3Privilege)
			if err != nil {
				return fmt.Errorf("error while setting s3_privilege property for group %s: %s", d.Id(), err)
			}
		}

		if group.Properties.CreateBackupUnit != nil {
			err := d.Set("create_backup_unit", *group.Properties.CreateBackupUnit)
			if err != nil {
				return fmt.Errorf("error while setting create_backup_unit property for group %s: %s", d.Id(), err)
			}
		}

		if group.Properties.CreateInternetAccess != nil {
			err := d.Set("create_internet_access", *group.Properties.CreateInternetAccess)
			if err != nil {
				return fmt.Errorf("error while setting create_internet_access property for group %s: %s", d.Id(), err)
			}
		}

		if group.Properties.CreateK8sCluster != nil {
			err := d.Set("create_k8s_cluster", *group.Properties.CreateK8sCluster)
			if err != nil {
				return fmt.Errorf("error while setting create_k8s_cluster property for group %s: %s", d.Id(), err)
			}
		}

		users, apiResponse, err := client.UserManagementApi.UmGroupsUsersGet(ctx, d.Id()).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occured while ListGroupUsers %s %w", d.Id(), err)
		}

		usersEntries := make([]interface{}, 0)
		if users.Items != nil && len(*users.Items) > 0 {
			usersEntries = make([]interface{}, len(*users.Items))
			for userIndex, user := range *users.Items {
				userEntry := make(map[string]interface{})
				if user.Properties != nil {
					if user.Properties.Firstname != nil {
						userEntry["first_name"] = *user.Properties.Firstname
					}

					if user.Properties.Lastname != nil {
						userEntry["last_name"] = *user.Properties.Lastname
					}

					if user.Properties.Email != nil {
						userEntry["email"] = *user.Properties.Email
					}

					if user.Properties.Administrator != nil {
						userEntry["administrator"] = *user.Properties.Administrator
					}

					if user.Properties.ForceSecAuth != nil {
						userEntry["force_sec_auth"] = *user.Properties.ForceSecAuth
					}
				}
				usersEntries[userIndex] = userEntry
			}

			if len(usersEntries) > 0 {
				if err := d.Set("users", usersEntries); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
