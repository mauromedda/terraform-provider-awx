package awx

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// awx_example test case
func TestAccAWXInventory(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInventoryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckState("name", "testacc"),
					testAccCheckState("organization_id", "1"),
					testAccCheckState("description", "AWX Acc test"),
				),
			},
			{
				ResourceName:      "awx_inventory.testacc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckState(skey, svalue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["awx_inventory.testacc"]
		if !ok {
			return fmt.Errorf("awx_inventory.testacc not found")
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

const testAccInventoryConfig = `
resource "awx_inventory" "testacc" {
	name = "testacc"
	organization_id = 1
	description = "AWX Acc test"
}
`
