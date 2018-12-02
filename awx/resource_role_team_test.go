package awx

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// awx_team test case
func TestAccAWXTeamRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamRoleConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStateTeamRole("role", "admin"),
					testAccCheckStateTeamRole("resource_type", "inventory"),
					testAccCheckStateTeamRole("resource_name", "Demo Inventory"),
					testAccCheckStateTeamRole("organization_id", "1"),
				),
			},
		},
	})
}

func testAccCheckStateTeamRole(skey, svalue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["awx_team_role.testacc-team_role_1"]
		if !ok {
			return fmt.Errorf("awx_team_role.testacc-team_role_1 not found")
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

const testAccTeamRoleConfig = `
resource "awx_team_role" "testacc-team_role_1" {
	team_id = 4
	organization_id = 1
	resource_type = "inventory"
	resource_name = "Demo Inventory"
	role = "admin"
  }
`
