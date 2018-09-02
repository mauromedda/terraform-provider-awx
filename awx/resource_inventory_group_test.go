package awx

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// awx_example test case
func TestAccAWXInventoryGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStateGroup("name", "testacc-grp_1"),
					testAccCheckStateGroup("description", "AWX Acc test group"),
				),
			},
		},
	})
}

func testAccCheckStateGroup(skey, svalue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["awx_inventory_group.testacc-grp"]
		if !ok {
			return fmt.Errorf("awx_inventory_group.testacc-grp not found")
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

const testAccInventoryGroupConfig = `
resource "awx_inventory" "testacc" {
	name = "testacc-grp"
	organization_id = 1
	description = "AWX Acc test"
}

resource "awx_inventory_group" "testacc-grp" {
	name = "testacc-grp_1"
	inventory_id = "${awx_inventory.testacc.id}"
	description = "AWX Acc test group"
}
`
