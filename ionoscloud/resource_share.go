package ionoscloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
)

func resourceShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceShareCreate,
		Read:   resourceShareRead,
		Update: resourceShareUpdate,
		Delete: resourceShareDelete,
		Schema: map[string]*schema.Schema{
			"edit_privilege": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"share_privilege": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceShareCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	request := ionoscloud.GroupShare{
		Properties: &ionoscloud.GroupShareProperties{},
	}

	tempSharePrivilege := d.Get("edit_privilege").(bool)
	request.Properties.SharePrivilege = &tempSharePrivilege
	tempEditPrivilege := d.Get("share_privilege").(bool)
	request.Properties.EditPrivilege = &tempEditPrivilege

	rsp, apiResponse, err := client.UserManagementApi.UmGroupsSharesPost(context.TODO(),
		d.Get("group_id").(string), d.Get("resource_id").(string)).Resource(request).Execute()

	log.Printf("[DEBUG] SHARE ID: %s", *rsp.Id)

	if err != nil {
		return fmt.Errorf("An error occured while creating a share: %s", err)
	}
	d.SetId(*rsp.Id)

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutCreate).WaitForState()
	if errState != nil {
		if IsRequestFailed(err) {
			// Request failed, so resource was not created, delete resource from state file
			d.SetId("")
		}
		return errState
	}

	return resourceShareRead(d, meta)
}

func resourceShareRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	rsp, apiResponse, err := client.UserManagementApi.UmGroupsSharesFindByResourceId(context.TODO(),
		d.Get("group_id").(string), d.Get("resource_id").(string)).Execute()
	if err != nil {
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("An error occured while fetching a Share ID %s %s", d.Id(), err)
	}

	d.Set("edit_privilege", *rsp.Properties.EditPrivilege)
	d.Set("share_privilege", *rsp.Properties.SharePrivilege)
	return nil
}

func resourceShareUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	tempSharePrivilege := d.Get("share_privilege").(bool)
	tempEditPrivilege := d.Get("edit_privilege").(bool)

	shareReq := ionoscloud.GroupShare{
		Properties: &ionoscloud.GroupShareProperties{
			EditPrivilege:  &tempEditPrivilege,
			SharePrivilege: &tempSharePrivilege,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)
	if cancel != nil {
		defer cancel()
	}

	_, apiResponse, err := client.UserManagementApi.UmGroupsSharesPut(ctx,
		d.Get("group_id").(string), d.Get("resource_id").(string)).Resource(shareReq).Execute()
	if err != nil {
		return fmt.Errorf("An error occured while patching a share ID %s %s", d.Id(), err)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceShareRead(d, meta)
}

func resourceShareDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}

	groupId := d.Get("group_id").(string)
	resourceId := d.Get("resource_id").(string)

	_, apiResponse, err := client.UserManagementApi.UmGroupsSharesDelete(ctx, groupId, resourceId).Execute()
	if err != nil {
		if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
			return err
		}
		//try again in 20 seconds
		time.Sleep(20 * time.Second)
		_, apiResponse, err := client.UserManagementApi.UmGroupsSharesDelete(ctx, groupId, resourceId).Execute()
		if err != nil {
			if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
				if apiResponse == nil || apiResponse.Response.StatusCode != 404 {
					return fmt.Errorf("an error occured while deleting a share %s %s", d.Id(), err)
				}
			}
		}
	}

	// Wait, catching any errors
	if apiResponse.Header.Get("Location") != "" {
		_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutDelete).WaitForState()
		if errState != nil {
			return errState
		}
	}

	d.SetId("")
	return nil
}
