package ionoscloud

import (
	"context"
	"fmt"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/profitbricks/profitbricks-sdk-go/v5"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupCreate,
		Read:   resourceGroupRead,
		Update: resourceGroupUpdate,
		Delete: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client
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

	usertoAdd := d.Get("user_id").(string)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)

	if cancel != nil {
		defer cancel()
	}

	group, apiRsponse, err := client.UserManagementApi.UmGroupsPost(ctx).Group(request).Execute()

	log.Printf("[DEBUG] GROUP ID: %s", *group.Id)

	if err != nil {
		return fmt.Errorf("An error occured while creating a group: %s", err)
	}

	d.SetId(*group.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiRsponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	//add users to group if any is provided
	if usertoAdd != "" {
		user := ionoscloud.User{
			Id: &usertoAdd,
		}
		_, apiResponse, err := client.UserManagementApi.UmGroupsUsersPost(ctx, d.Id()).User(user).Execute()
		if err != nil {
			return fmt.Errorf("An error occured while adding %s user to group ID %s %s", usertoAdd, d.Id(), err)
		}
		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
		if errState != nil {
			return errState
		}
	}
	return resourceGroupRead(d, meta)
}

func resourceGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	group, _, err := client.UserManagementApi.UmGroupsFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiError, ok := err.(profitbricks.ApiError); ok {
			if apiError.HttpStatusCode() == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("An error occured while fetching a Group ID %s %s", d.Id(), err)
	}

	if group.Properties.Name != nil {
		err := d.Set("name", *group.Properties.Name)
		if err != nil {
			return fmt.Errorf("Error while setting name property for group %s: %s", d.Id(), err)
		}
	}

	if group.Properties.CreateDataCenter != nil {
		err := d.Set("create_datacenter", *group.Properties.CreateDataCenter)
		if err != nil {
			return fmt.Errorf("Error while setting create_datacenter property for group %s: %s", d.Id(), err)
		}
	}

	if group.Properties.CreateSnapshot != nil {
		err := d.Set("create_snapshot", *group.Properties.CreateSnapshot)
		if err != nil {
			return fmt.Errorf("Error while setting create_snapshot property for group %s: %s", d.Id(), err)
		}
	}

	if group.Properties.ReserveIp != nil {
		err := d.Set("reserve_ip", *group.Properties.ReserveIp)
		if err != nil {
			return fmt.Errorf("Error while setting reserve_ip property for group %s: %s", d.Id(), err)
		}
	}

	if group.Properties.ReserveIp != nil {
		err := d.Set("access_activity_log", *group.Properties.AccessActivityLog)
		if err != nil {
			return fmt.Errorf("Error while setting access_activity_log property for group %s: %s", d.Id(), err)
		}
	}

	users, _, err := client.UserManagementApi.UmGroupsUsersGet(ctx, d.Id()).Execute()
	if err != nil {
		return fmt.Errorf("An error occured while ListGroupUsers %s %s", d.Id(), err)
	}

	var usersArray = []ionoscloud.UserProperties{}
	if len(*users.Items) > 0 {
		for _, usr := range *users.Items {
			usersArray = append(usersArray, *usr.Properties)
		}
		d.Set("users", usersArray)
	}

	return nil
}

func resourceGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	tempCreateDataCenter := d.Get("create_datacenter").(bool)
	tempCreateSnapshot := d.Get("create_snapshot").(bool)
	tempReserveIp := d.Get("reserve_ip").(bool)
	tempAccessActivityLog := d.Get("access_activity_log").(bool)
	usertoAdd := d.Get("user_id").(string)

	groupReq := ionoscloud.Group{
		Properties: &ionoscloud.GroupProperties{
			CreateDataCenter:  &tempCreateDataCenter,
			CreateSnapshot:    &tempCreateSnapshot,
			ReserveIp:         &tempReserveIp,
			AccessActivityLog: &tempAccessActivityLog,
		},
	}

	_, newValue := d.GetChange("name")
	newValueStr := newValue.(string)
	groupReq.Properties.Name = &newValueStr

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.UserManagementApi.UmGroupsPut(ctx, d.Id()).Group(groupReq).Execute()
	if err != nil {
		return fmt.Errorf("An error occured while patching a group ID %s %s", d.Id(), err)
	}
	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	//add users to group if any is provided
	if usertoAdd != "" {

		user := ionoscloud.User{
			Id: &usertoAdd,
		}

		_, apiResponse, err := client.UserManagementApi.UmGroupsUsersPost(ctx, d.Id()).User(user).Execute()
		if err != nil {
			return fmt.Errorf("An error occured while adding %s user to group ID %s %s", usertoAdd, d.Id(), err)
		}

		// Wait, catching any errors
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
		if errState != nil {
			return errState
		}
	}
	return resourceGroupRead(d, meta)
}

func resourceGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(SdkBundle).Client

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)

	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.UserManagementApi.UmGroupsDelete(ctx, d.Id()).Execute()
	if err != nil {
		//try again in 20 seconds
		time.Sleep(20 * time.Second)
		_, apiResponse, err = client.UserManagementApi.UmGroupsDelete(ctx, d.Id()).Execute()

		if err != nil {
			if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
				if apiResponse.Response.StatusCode != 404 {
					return fmt.Errorf("An error occured while deleting a group %s %s", d.Id(), err)
				}
			}
		}
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
	if errState != nil {
		return errState
	}

	d.SetId("")
	return nil
}
