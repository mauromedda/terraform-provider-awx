package awx

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// awx_project test case
func TestAccAWXProject(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccProjectConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStateProject("name", "testacc-prj_1"),
					testAccCheckStateProject("description", "AWX Acc test project"),
					testAccCheckStateProject("organization_id", "1"),
					testAccCheckStateProject("scm_type", "git"),
					testAccCheckStateProject("scm_update_on_launch", "true"),
					testAccCheckStateProject("scm_url", "https://github.com/ansible/ansible-tower-samples"),
				),
			},
		},
	})
}

func testAccCheckStateProject(skey, svalue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["awx_project.testacc-prj_1"]
		if !ok {
			return fmt.Errorf("awx_project.testacc-prj_1 not found")
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

const testAccProjectConfig = `
resource "awx_project" "testacc-prj_1" {
	name = "testacc-prj_1"
	description = "AWX Acc test project"
	scm_type = "git"
	scm_url = "https://github.com/ansible/ansible-tower-samples"
	scm_update_on_launch = true
	organization_id = "1"
  }
`
