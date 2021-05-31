package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"first_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"last_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"email": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"password": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"administrator": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"force_sec_auth": {
				Type:     schema.TypeBool,
				Required: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceUserCreate(d *schema.ResourceData, meta interface{}) error {
<<<<<<< HEAD
	client := meta.(*ionoscloud.APIClient)
=======
	client := meta.(*ionoscloud.APIClient)

>>>>>>> master
	request := ionoscloud.UserPost{
		Properties: &ionoscloud.UserPropertiesPost{},
	}

	log.Printf("[DEBUG] NAME %s", d.Get("first_name"))

	if d.Get("first_name") != nil {
		firstName := d.Get("first_name").(string)
		request.Properties.Firstname = &firstName
	}
	if d.Get("last_name") != nil {
		lastName := d.Get("last_name").(string)
		request.Properties.Lastname = &lastName
	}
	if d.Get("email") != nil {
		email := d.Get("email").(string)
		request.Properties.Email = &email
	}
	if d.Get("password") != nil {
		password := d.Get("password").(string)
		request.Properties.Password = &password
	}

	administrator := d.Get("administrator").(bool)
	forceSecAuth := d.Get("force_sec_auth").(bool)
	request.Properties.Administrator = &administrator
	request.Properties.ForceSecAuth = &forceSecAuth

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Create)
	if cancel != nil {
		defer cancel()
	}

	rsp, apiResponse, err := client.UserManagementApi.UmUsersPost(ctx).User(request).Execute()
	if rsp.Id != nil {
		log.Printf("[DEBUG] USER ID: %s", *rsp.Id)
	}

	if err != nil {
		payload := "<nil>"
		if apiResponse != nil {
			payload = string(apiResponse.Payload)
		}
		return fmt.Errorf("an error occured while creating a user: %s; payload: %s", err, payload)
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
	return resourceUserRead(d, meta)
}

func resourceUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)
	if cancel != nil {
		defer cancel()
	}

	rsp, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
			if apiResponse != nil && apiResponse.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
		return fmt.Errorf("an error occured while fetching a User ID %s %s", d.Id(), err)
	}

	d.Set("first_name", *rsp.Properties.Firstname)
	d.Set("last_name", *rsp.Properties.Lastname)
	d.Set("email", *rsp.Properties.Email)
	d.Set("administrator", *rsp.Properties.Administrator)
	d.Set("force_sec_auth", *rsp.Properties.ForceSecAuth)
	return nil
}

func resourceUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Update)
	if cancel != nil {
		defer cancel()
	}

	rsp, apiResponse, err := client.UserManagementApi.UmUsersFindById(ctx, d.Id()).Execute()

	if err != nil {
		return fmt.Errorf("an error occured while fetching a User ID %s %s", d.Id(), err)
	}

	administrator := d.Get("administrator").(bool)
	forceSecAuth := d.Get("force_sec_auth").(bool)
	userReq := ionoscloud.UserPut{
		Properties: &ionoscloud.UserPropertiesPut{
			Administrator: &administrator,
			ForceSecAuth:  &forceSecAuth,
		},
	}

	if d.HasChange("first_name") {
		_, newValue := d.GetChange("first_name")
		firstName := newValue.(string)
		userReq.Properties.Firstname = &firstName

	} else {
		userReq.Properties.Firstname = rsp.Properties.Firstname
	}

	if d.HasChange("last_name") {
		_, newValue := d.GetChange("last_name")
		lastName := newValue.(string)
		userReq.Properties.Lastname = &lastName
	} else {
		userReq.Properties.Lastname = rsp.Properties.Lastname
	}

	if d.HasChange("email") {
		_, newValue := d.GetChange("email")
		email := newValue.(string)
		userReq.Properties.Email = &email
	} else {
		userReq.Properties.Email = rsp.Properties.Email
	}

	rsp, apiResponse, err = client.UserManagementApi.UmUsersPut(ctx, d.Id()).User(userReq).Execute()
	if err != nil {
		return fmt.Errorf("an error occured while patching a user ID %s %s payload: %s", d.Id(), err, apiResponse.Payload)
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForState()
	if errState != nil {
		return errState
	}

	return resourceUserRead(d, meta)
}

func resourceUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ionoscloud.APIClient)

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Delete)
	if cancel != nil {
		defer cancel()
	}
	_, apiResponse, err := client.UserManagementApi.UmUsersDelete(ctx, d.Id()).Execute()
	if err != nil {
		//try again in 20 seconds
		time.Sleep(20 * time.Second)
		_, _, err := client.UserManagementApi.UmUsersDelete(ctx, d.Id()).Execute()
		if err != nil {
			if _, ok := err.(ionoscloud.GenericOpenAPIError); ok {
				if apiResponse == nil || apiResponse.Response.StatusCode != 404 {
					return fmt.Errorf("an error occured while deleting a user %s %s", d.Id(), err)
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
