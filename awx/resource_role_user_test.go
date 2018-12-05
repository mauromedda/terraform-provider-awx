package awx

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// awx_user test case
func TestAccAWXUserRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStateUserRole("role", "admin"),
					testAccCheckStateUserRole("resource_type", "inventory"),
					testAccCheckStateUserRole("resource_name", "Demo Inventory"),
				),
			},
		},
	})
}

func testAccCheckStateUserRole(skey, svalue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["awx_user_role.testacc-user_role_1"]
		if !ok {
			return fmt.Errorf("awx_user_role.testacc-user_role_1 not found")
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

const testAccUserRoleConfig = `
resource "awx_user_role" "testacc-user_role_1" {
	user_id = 4
	organization_id = 1
	resource_type = "inventory"
	resource_name = "Demo Inventory"
	role = "admin"
  }

  resource "awx_user_role" "testacc-user_role_2" {
	user_id = 4
	organization_id = 1
	resource_type = "organization"
	resource_name = "organization"
	role = "inventory admin"
  }
`
