package ionoscloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
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
				Type:        schema.TypeBool,
				Description: "Create backup unit privilege.",
				Optional:    true,
			},
			"create_internet_access": {
				Type:        schema.TypeBool,
				Description: "Create internet access privilege.",
				Optional:    true,
			},
			"create_k8s_cluster": {
				Type:        schema.TypeBool,
				Description: "Create Kubernetes cluster privilege.",
				Optional:    true,
			},
			"create_flow_log": {
				Type:        schema.TypeBool,
				Description: "Create Flow Logs privilege.",
				Optional:    true,
			},
			"access_and_manage_monitoring": {
				Type: schema.TypeBool,
				Description: "Privilege for a group to access and manage monitoring related functionality " +
					"(access metrics, CRUD on alarms, alarm-actions etc) using Monotoring-as-a-Service (MaaS).",
				Optional: true,
			},
			"access_and_manage_certificates": {
				Type:        schema.TypeBool,
				Description: "Privilege for a group to access and manage certificates.",
				Optional:    true,
			},
			"manage_dbaas": {
				Type:        schema.TypeBool,
				Description: "Privilege for a group to manage DBaaS related functionality",
				Optional:    true,
			},
			"user_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"user_ids"},
				Deprecated:    "Please use user_ids for adding users to the group, since user_id will be removed in the future",
			},
			"user_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ConflictsWith: []string{"user_id"},
			},
			"users": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
		Timeouts:      &resourceDefaultTimeouts,
		SchemaVersion: 1,
		//StateUpgraders: []schema.StateUpgrader{
		//	{
		//		Type:    resourceGroup0().CoreConfigSchema().ImpliedType(),
		//		Upgrade: resourceGroupUpgradeV0,
		//		Version: 0,
		//	},
		//},
	}
}

//
//func resourceGroup0() *schema.Resource {
//	return &schema.Resource{
//		Schema: map[string]*schema.Schema{
//			"user_id": {
//				Type:     schema.TypeString,
//				Optional: true,
//			},
//		},
//		Timeouts: &resourceDefaultTimeouts,
//	}
//}
//
//func resourceGroupUpgradeV0(_ context.Context, state map[string]interface{}, _ interface{}) (map[string]interface{}, error) {
//	oldState := state
//	var oldData string
//	if d, ok := oldState["user_id"].(string); ok {
//		oldData = d
//		var users []string
//		users = append(users, oldData)
//		state["user_ids"] = users
//	}
//
//	return state, nil
//}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	request := ionoscloud.Group{
		Properties: &ionoscloud.GroupProperties{},
	}

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
	manageDbaas := d.Get("manage_dbaas").(bool)
	request.Properties.AccessAndManageCertificates = &tempAccessAndManageCertificates
	request.Properties.ManageDBaaS = &manageDbaas
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
	if userVal, userOK := d.GetOk("user_id"); userOK {
		userID := userVal.(string)
		log.Printf("[INFO] Adding user %+v to group...", userID)
		if err := addUserToGroup(userID, d.Id(), ctx, d, meta); err != nil {
			return diag.FromErr(err)
		}
	}

	if usersVal, usersOK := d.GetOk("user_ids"); usersOK {
		usersList := usersVal.(*schema.Set)
		if usersList.List() != nil {
			for _, userItem := range usersList.List() {
				userID := userItem.(string)
				log.Printf("[INFO] Adding user %+v to group...", userID)
				if err := addUserToGroup(userID, d.Id(), ctx, d, meta); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}
	return resourceGroupRead(ctx, d, meta)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(SdkBundle).CloudApiClient

	group, apiResponse, err := client.UserManagementApi.UmGroupsFindById(ctx, d.Id()).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
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
	tempManageDBaaS := d.Get("manage_dbaas").(bool)

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
			ManageDBaaS:                 &tempManageDBaaS,
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

	if d.HasChange("user_id") {
		oldValue, newValue := d.GetChange("user_id")

		userIdToAdd := newValue.(string)
		userIdToRemove := oldValue.(string)

		log.Printf("[INFO] User to add: %+v", userIdToAdd)
		log.Printf("[INFO] User to remove: %+v", userIdToRemove)

		if userIdToAdd != "" {
			if err := addUserToGroup(userIdToAdd, d.Id(), ctx, d, meta); err != nil {
				return diag.FromErr(err)
			}
		}

		if userIdToRemove != "" {
			if err := deleteUserFromGroup(userIdToRemove, d.Id(), ctx, d, meta); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("user_ids") {
		oldValues, newValues := d.GetChange("user_ids")
		oldUsersList := convertSlice(oldValues.(*schema.Set).List())
		newUsersList := convertSlice(newValues.(*schema.Set).List())

		newUsers := utils.DiffSliceOneWay(newUsersList, oldUsersList)
		deletedUsers := utils.DiffSliceOneWay(oldUsersList, newUsersList)

		if newUsers != nil && len(newUsers) > 0 {
			log.Printf("[INFO] New users to add: %+v", newUsers)
			for _, userID := range newUsers {
				if err := addUserToGroup(userID, d.Id(), ctx, d, meta); err != nil {
					return diag.FromErr(err)
				}
			}
		}

		if deletedUsers != nil && len(deletedUsers) > 0 {
			log.Printf("[INFO] Users to delete: %+v", deletedUsers)
			for _, userID := range deletedUsers {
				if err := deleteUserFromGroup(userID, d.Id(), ctx, d, meta); err != nil {
					return diag.FromErr(err)
				}
			}
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
		if httpNotFound(apiResponse) {
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
				return fmt.Errorf("error while setting name property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.CreateDataCenter != nil {
			err := d.Set("create_datacenter", *group.Properties.CreateDataCenter)
			if err != nil {
				return fmt.Errorf("error while setting create_datacenter property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.CreateSnapshot != nil {
			err := d.Set("create_snapshot", *group.Properties.CreateSnapshot)
			if err != nil {
				return fmt.Errorf("error while setting create_snapshot property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.ReserveIp != nil {
			err := d.Set("reserve_ip", *group.Properties.ReserveIp)
			if err != nil {
				return fmt.Errorf("error while setting reserve_ip property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.AccessActivityLog != nil {
			err := d.Set("access_activity_log", *group.Properties.AccessActivityLog)
			if err != nil {
				return fmt.Errorf("error while setting access_activity_log property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.CreatePcc != nil {
			err := d.Set("create_pcc", *group.Properties.CreatePcc)
			if err != nil {
				return fmt.Errorf("error while setting create_pcc property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.S3Privilege != nil {
			err := d.Set("s3_privilege", *group.Properties.S3Privilege)
			if err != nil {
				return fmt.Errorf("error while setting s3_privilege property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.CreateBackupUnit != nil {
			err := d.Set("create_backup_unit", *group.Properties.CreateBackupUnit)
			if err != nil {
				return fmt.Errorf("error while setting create_backup_unit property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.CreateInternetAccess != nil {
			err := d.Set("create_internet_access", *group.Properties.CreateInternetAccess)
			if err != nil {
				return fmt.Errorf("error while setting create_internet_access property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.CreateK8sCluster != nil {
			err := d.Set("create_k8s_cluster", *group.Properties.CreateK8sCluster)
			if err != nil {
				return fmt.Errorf("error while setting create_k8s_cluster property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.CreateFlowLog != nil {
			err := d.Set("create_flow_log", *group.Properties.CreateFlowLog)
			if err != nil {
				return fmt.Errorf("error while setting create_flow_log property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.AccessAndManageMonitoring != nil {
			err := d.Set("access_and_manage_monitoring", *group.Properties.AccessAndManageMonitoring)
			if err != nil {
				return fmt.Errorf("error while setting access_and_manage_monitoring property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.AccessAndManageCertificates != nil {
			err := d.Set("access_and_manage_certificates", *group.Properties.AccessAndManageCertificates)
			if err != nil {
				return fmt.Errorf("error while setting access_and_manage_certificates property for group %s: %w", d.Id(), err)
			}
		}

		if group.Properties.ManageDBaaS != nil {
			err := d.Set("manage_dbaas", *group.Properties.ManageDBaaS)
			if err != nil {
				return fmt.Errorf("error while setting manage_dbaas property for group %s: %w", d.Id(), err)
			}
		}

		users, apiResponse, err := client.UserManagementApi.UmGroupsUsersGet(ctx, d.Id()).Depth(1).Execute()
		logApiRequestTime(apiResponse)
		if err != nil {
			return fmt.Errorf("an error occured while UmGroupsUsersGet %s %w", d.Id(), err)
		}

		usersEntries := make([]interface{}, 0)
		if users.Items != nil && len(*users.Items) > 0 {
			usersEntries = make([]interface{}, len(*users.Items))
			for userIndex, user := range *users.Items {
				userEntry := make(map[string]interface{})

				if user.Id != nil {
					userEntry["id"] = *user.Id
				}

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

func addUserToGroup(userId, groupId string, ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).CloudApiClient
	userToAdd := ionoscloud.User{
		Id: &userId,
	}

	_, apiResponse, err := client.UserManagementApi.UmGroupsUsersPost(ctx, groupId).User(userToAdd).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return fmt.Errorf("an error occured while adding %s user to group ID %s %w", userId, groupId, err)
	}

	log.Printf("[INFO] Added user %s to group %s", userId, groupId)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		return errState
	}

	return nil
}

func deleteUserFromGroup(userId, groupId string, ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).CloudApiClient

	apiResponse, err := client.UserManagementApi.UmGroupsUsersDelete(ctx, groupId, userId).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		return fmt.Errorf("an error occured while deleting %s user from group ID %s %w", userId, groupId, err)
	}

	log.Printf("[INFO] Deleted user %s from group %s", userId, groupId)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForStateContext(ctx)
	if errState != nil {
		return errState
	}

	return nil
}
