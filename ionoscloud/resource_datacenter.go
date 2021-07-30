package ionoscloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"log"
	"regexp"
	"strings"
)

func resourceDatacenter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatacenterCreate,
		ReadContext:   resourceDatacenterRead,
		UpdateContext: resourceDatacenterUpdate,
		DeleteContext: resourceDatacenterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{

			//Datacenter parameters
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"location": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.All(validation.StringIsNotWhiteSpace),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sec_auth_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		Timeouts: &resourceDefaultTimeouts,
	}
}

func resourceDatacenterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	datacenterName := d.Get("name").(string)
	datacenterLocation := d.Get("location").(string)

	datacenter := ionoscloud.Datacenter{
		Properties: &ionoscloud.DatacenterProperties{
			Name:     &datacenterName,
			Location: &datacenterLocation,
		},
	}

	if attr, ok := d.GetOk("description"); ok {
		attrStr := attr.(string)
		datacenter.Properties.Description = &attrStr
	}

	if attr, ok := d.GetOk("sec_auth_protection"); ok {
		attrStr := attr.(bool)
		datacenter.Properties.SecAuthProtection = &attrStr
	}

	createdDatacenter, apiResponse, err := client.DataCentersApi.DatacentersPost(ctx).Datacenter(datacenter).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("error creating data center (%s) (%s)", d.Id(), err))
		return diags
	}
	d.SetId(*createdDatacenter.Id)

	log.Printf("[INFO] DataCenter Id: %s", d.Id())

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

	return resourceDatacenterRead(ctx, d, meta)
}

func resourceDatacenterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	datacenter, apiResponse, err := client.DataCentersApi.DatacentersFindById(ctx, d.Id()).Execute()

	if err != nil {
		if apiResponse != nil && apiResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		diags := diag.FromErr(fmt.Errorf("error while fetching a data center ID %s %s", d.Id(), err))
		return diags
	}

	if datacenter.Properties.Name != nil {
		err := d.Set("name", *datacenter.Properties.Name)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting name property for datacenter %s: %s", d.Id(), err))
			return diags
		}
	}

	if datacenter.Properties.Location != nil {
		err := d.Set("location", *datacenter.Properties.Location)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting location property for datacenter %s: %s", d.Id(), err))
			return diags
		}
	}

	if datacenter.Properties.Description != nil {
		err := d.Set("description", *datacenter.Properties.Description)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting description property for datacenter %s: %s", d.Id(), err))
			return diags
		}
	}

	if datacenter.Properties.SecAuthProtection != nil {
		err := d.Set("sec_auth_protection", *datacenter.Properties.SecAuthProtection)
		if err != nil {
			diags := diag.FromErr(fmt.Errorf("error while setting sec_auth_protection property for datacenter %s: %s", d.Id(), err))
			return diags
		}
	}

	return nil
}

func resourceDatacenterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient
	obj := ionoscloud.DatacenterProperties{}

	if d.HasChange("name") {
		_, newName := d.GetChange("name")
		newNameStr := newName.(string)
		obj.Name = &newNameStr
	}

	if d.HasChange("description") {
		_, newDescription := d.GetChange("description")
		newDescriptionStr := newDescription.(string)
		obj.Description = &newDescriptionStr
	}

	if d.HasChange("location") {
		oldLocation, newLocation := d.GetChange("location")
		diags := diag.FromErr(fmt.Errorf("data center is created in %s location. You can not change location of the data center to %s; it requires recreation of the data center", oldLocation, newLocation))
		return diags
	}

	if d.HasChange("sec_auth_protection") {
		_, newSecAuthProtection := d.GetChange("sec_auth_protection")
		newSecAuthProtectionStr := newSecAuthProtection.(bool)
		obj.SecAuthProtection = &newSecAuthProtectionStr
	}
	_, apiResponse, err := client.DataCentersApi.DatacentersPatch(ctx, d.Id()).Datacenter(obj).Execute()

	if err != nil {
		diags := diag.FromErr(fmt.Errorf("an error occured while update the data center ID %s %s", d.Id(), err))
		return diags
	}

	// Wait, catching any errors
	_, errState := getStateChangeConf(meta, d, apiResponse.Header.Get("Location"), schema.TimeoutUpdate).WaitForStateContext(ctx)
	if errState != nil {
		diags := diag.FromErr(errState)
		return diags
	}

	return resourceDatacenterRead(ctx, d, meta)
}

func resourceDatacenterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := meta.(SdkBundle).CloudApiClient

	apiResponse, err := client.DataCentersApi.DatacentersDelete(ctx, d.Id()).Execute()

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

func getImage(client *ionoscloud.APIClient, dcId string, imageName string, imageType string) (*ionoscloud.Image, error) {

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if imageName == "" {
		return nil, fmt.Errorf("imageName not suplied")
	}

	dc, _, err := client.DataCentersApi.DatacentersFindById(ctx, dcId).Execute()

	if err != nil {
		log.Print(fmt.Errorf("error while fetching a data center ID %s %s", dcId, err))
		return nil, err
	}

	images, _, err := client.ImagesApi.ImagesGet(ctx).Execute()

	if err != nil {
		log.Print(fmt.Errorf("error while fetching the list of images %s", err))
		return nil, err
	}

	if len(*images.Items) > 0 {
		for _, i := range *images.Items {
			imgName := ""
			if i.Properties.Name != nil && *i.Properties.Name != "" {
				imgName = *i.Properties.Name
			}

			if imageType == "SSD" {
				imageType = "HDD"
			}

			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(imageName)) && *i.Properties.ImageType == imageType && *i.Properties.Location == *dc.Properties.Location {
				return &i, err
			}

			if imgName != "" && strings.ToLower(imageName) == strings.ToLower(*i.Id) && *i.Properties.ImageType == imageType && *i.Properties.Location == *dc.Properties.Location {
				return &i, err
			}

		}
	}
	return nil, err
}

func getSnapshotId(client *ionoscloud.APIClient, snapshotName string) string {

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if snapshotName == "" {
		return ""
	}

	snapshots, _, err := client.SnapshotsApi.SnapshotsGet(ctx).Execute()

	if err != nil {
		log.Print(fmt.Errorf("error while fetching the list of snapshots %s", err))
	}

	if len(*snapshots.Items) > 0 {
		for _, i := range *snapshots.Items {
			imgName := ""
			if *i.Properties.Name != "" {
				imgName = *i.Properties.Name
			}

			if imgName != "" && strings.Contains(strings.ToLower(imgName), strings.ToLower(snapshotName)) {
				return *i.Id
			}
		}
	}
	return ""
}

func getImageAlias(client *ionoscloud.APIClient, imageAlias string, location string) string {

	ctx, cancel := context.WithTimeout(context.Background(), *resourceDefaultTimeouts.Default)

	if cancel != nil {
		defer cancel()
	}

	if imageAlias == "" {
		return ""
	}
	parts := strings.SplitN(location, "/", 2)
	if len(parts) != 2 {
		log.Print(fmt.Errorf("invalid location id %s", location))
	}

	locations, _, err := client.LocationsApi.LocationsFindByRegionIdAndId(ctx, parts[0], parts[1]).Execute()

	if err != nil {
		log.Print(fmt.Errorf("error while fetching the list of locations %s", err))
	}

	if len(*locations.Properties.ImageAliases) > 0 {
		for _, i := range *locations.Properties.ImageAliases {
			alias := ""
			if i != "" {
				alias = i
			}

			if alias != "" && strings.ToLower(alias) == strings.ToLower(imageAlias) {
				return i
			}
		}
	}
	return ""
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
