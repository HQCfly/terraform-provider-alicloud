package alicloud

import (
	"fmt"
	"github.com/aliyun/aliyun-log-go-sdk"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
	"testing"
)

func TestAccAlicloudLogAudit_basic(t *testing.T) {
	var v *sls.LogProject
	resourceId := "alicloud_log_audit.default"
	ra := resourceAttrInit(resourceId, logAuditMap)
	serviceFunc := func() interface{} {
		return &LogService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacclogaudit-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceLogAuditConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"display_name":        name,
					"variable_map_string": "${alicloud_log_audit.foo.variable_map_string}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"display_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func resourceLogAuditConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	    default = "%s"
	}
	resource "alicloud_log_audit" "foo" {
	    display_name = "${var.name}"
	    variable_map_string = "{}"
	}
	`, name)
}

var logAuditMap = map[string]string{
	"display_name": CHECKSET,
}
