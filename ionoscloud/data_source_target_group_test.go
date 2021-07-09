package ionoscloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTargetGroup_matchId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTargetGroupCreateResources,
			},
			{
				Config: testAccDataSourceTargetGroupMatchId,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_target_group.test_target_group", "name", "test_datasource_target_group"),
				),
			},
		},
	})
}

func TestAccDataSourceTargetGroup_matchName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceTargetGroupCreateResources,
			},
			{
				Config: testAccDataSourceTargetGroupMatchName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ionoscloud_target_group.test_target_group", "name", "test_datasource_target_group"),
				),
			},
		},
	})
}

const testAccDataSourceTargetGroupCreateResources = `
resource "ionoscloud_target_group" "target_group" {
 name = "test_datasource_target_group"
 algorithm = "SOURCE_IP"
 protocol = "HTTP"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "123"
   health_check {
     check = true
     check_interval = 1000
   }
 }
 health_check {
     connect_timeout = 1000
 }
 http_health_check {
     path = "/monitoring"
     match_type = "RESPONSE_BODY"
     response = "200"
   }
}
`

const testAccDataSourceTargetGroupMatchId = `
resource "ionoscloud_target_group" "target_group" {
 name = "test_datasource_target_group"
 algorithm = "SOURCE_IP"
 protocol = "HTTP"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "123"
   health_check {
     check = true
     check_interval = 1000
   }
 }
 health_check {
     connect_timeout = 1000
 }
 http_health_check {
     path = "/monitoring"
     match_type = "RESPONSE_BODY"
     response = "200"
   }
}

data "ionoscloud_target_group" "test_target_group" {
  id			= ionoscloud_target_group.target_group.id
}
`

const testAccDataSourceTargetGroupMatchName = `
resource "ionoscloud_target_group" "target_group" {
 name = "test_datasource_target_group"
 algorithm = "SOURCE_IP"
 protocol = "HTTP"
 targets {
   ip = "22.231.2.2"
   port = "8080"
   weight = "123"
   health_check {
     check = true
     check_interval = 1000
   }
 }
 health_check {
     connect_timeout = 1000
 }
 http_health_check {
     path = "/monitoring"
     match_type = "RESPONSE_BODY"
     response = "200"
   }
}

data "ionoscloud_target_group" "test_target_group" {
  name			= "test_datasource_"
}
`
