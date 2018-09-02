package awx

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// awx_host test case
func TestAccAWXHost(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHostConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStateHost("name", "testacc-host_1"),
					testAccCheckStateHost("description", "AWX Acc test host"),
				),
			},
		},
	})
}

func testAccCheckStateHost(skey, svalue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["awx_host.testacc-host_1"]
		if !ok {
			return fmt.Errorf("awx_host.testacc-host_1 not found")
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		cr := rs.Primary

		if cr.Attributes[skey] != svalue {
			return fmt.Errorf("%s != %s (actual: %s)", skey, svalue, cr.Attributes[skey])
		}

		return nil
	}
}

const testAccHostConfig = `
resource "awx_host" "testacc-host_1" {
	name         = "testacc-host_1"
	description  = "AWX Acc test host"
	inventory_id = "1"
	variables = <<VARIABLES
---
api_server_enabled: false
VARIABLES

  }
`
