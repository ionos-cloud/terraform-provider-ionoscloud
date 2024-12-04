package vpn

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpnSdk "github.com/ionos-cloud/sdk-go-bundle/products/vpn/v2"
)

// GetMaintenanceWindowData returns the MaintenanceWindow data from the schema
func GetMaintenanceWindowData(d *schema.ResourceData) *vpnSdk.MaintenanceWindow {
	var maintenanceWindow vpnSdk.MaintenanceWindow

	if timeV, ok := d.GetOk("maintenance_window.0.time"); ok {
		maintenanceWindow.Time = timeV.(string)
	}
	if dayOfTheWeek, ok := d.GetOk("maintenance_window.0.day_of_the_week"); ok {
		maintenanceWindow.DayOfTheWeek = vpnSdk.DayOfTheWeek(dayOfTheWeek.(string))
	}

	return &maintenanceWindow
}
