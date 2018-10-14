package awx

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// awx_user test case
func TestAccAWXUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStateUser("username", "mauromedda"),
					testAccCheckStateUser("password", "password"),
					testAccCheckStateUser("first_name", "Mauro"),
					testAccCheckStateUser("last_name", "Medda"),
					testAccCheckStateUser("is_superuser", "true"),
					testAccCheckStateUser("email", "medda.mauro@test.td"),
				),
			},
		},
	})
}

func testAccCheckStateUser(skey, svalue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["awx_user.testacc-user_1"]
		if !ok {
			return fmt.Errorf("awx_user.testacc-user_1 not found")
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

const testAccUserConfig = `
resource "awx_user" "testacc-user_1" {
	username = "mauromedda"
	password = "password"
	first_name = "Mauro"
	last_name = "Medda"
	is_superuser = true
	email = "medda.mauro@test.td"
  }
`
