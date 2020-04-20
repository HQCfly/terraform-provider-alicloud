package alicloud

import (
	slsPop "github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func resourceAlicloudLogAudit() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudLogAuditCreate,
		Read:   resourceAlicloudLogAuditRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"variable_map_string": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceAlicloudLogAuditCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := slsPop.CreateAnalyzeAppLogRequest()
	request.Domain = "sls.aliyuncs.com"
	request.AppType = "audit"
	request.DisplayName = d.Get("display_name").(string)
	request.VariableMap = d.Get("variable_map_string").(string)
	response := slsPop.CreateCreateAppResponse()

	if err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		_, err := client.WithLogPopClient(func(client *slsPop.Client) (interface{}, error) {
			return nil, client.DoAction(request, response)
		})
		if err != nil {
			if IsExpectedErrors(err, []string{LogClientTimeout}) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug("CreateLogAudit", response, request)
		return nil
	}); err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_log_audit", "CreateLogAudit", AliyunLogGoSdkERROR)
	}

	return resourceAlicloudLogAuditRead(d, meta)
}

func resourceAlicloudLogAuditRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	request := slsPop.CreateDescribeAppRequest()
	_, err := client.WithLogPopClient(func(client *slsPop.Client) (interface{}, error) {
		return nil, client.DoAction(request, slsPop.CreateCreateAppResponse())
	})
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	return nil
}