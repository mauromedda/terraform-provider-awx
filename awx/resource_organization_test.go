package awx

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// awx_organization test case
func TestAccAWXOrganization(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStateOrganization("name", "automation_organization"),
					testAccCheckStateOrganization("description", "Automation Organization"),
				),
			},
		},
	})
}

func testAccCheckStateOrganization(skey, svalue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["awx_organization.testacc-organization_1"]
		if !ok {
			return fmt.Errorf("awx_organization.testacc-organization_1 not found")
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

const testAccOrganizationConfig = `
resource "awx_organization" "testacc-organization_1" {
	name = "automation_organization"
	description = "Automation Organization"
  }
`
