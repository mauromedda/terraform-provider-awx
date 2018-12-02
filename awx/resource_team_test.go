package awx

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// awx_team test case
func TestAccAWXTeam(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeamConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStateTeam("name", "automation_team"),
					testAccCheckStateTeam("description", "Automation Team"),
					testAccCheckStateTeam("organization_id", "1"),
				),
			},
		},
	})
}

func testAccCheckStateTeam(skey, svalue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["awx_team.testacc-team_1"]
		if !ok {
			return fmt.Errorf("awx_team.testacc-team_1 not found")
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

const testAccTeamConfig = `
resource "awx_team" "testacc-team_1" {
	name = "automation_team"
	description = "Automation Team"
	organization_id = "1"
  }
`
