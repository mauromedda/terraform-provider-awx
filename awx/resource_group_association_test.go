package awx

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// awx_host test case
func TestAccAWXGroupAssociation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { TestAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGroupAssociationConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStateGroupAssociation("name", "k8s-node-1_k8s-nodes"),
					testAccCheckStateGroupAssociation("inventory_id", "1"),
				),
			},
		},
	})
}

func testAccCheckStateGroupAssociation(skey, svalue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["awx_group_association.k8s-node-1_k8s-nodes"]
		if !ok {
			return fmt.Errorf("awx_group_association.k8s-node-1_k8s-nodes not found")
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

const testAccGroupAssociationConfig = `
resource "awx_host" "testacc-host_1" {
	name         = "testacc-host_1"
	description  = "AWX Acc test host"
	inventory_id = "1"
	variables = <<VARIABLES
---
api_server_enabled: false
VARIABLES
  }

resource "awx_inventory_group" "k8s-nodes" {
	name         = "k8s-nodes"
	inventory_id = "1"
  }

resource "awx_group_association" "k8s-node-1_k8s-nodes" {
	name         = "k8s-node-1_k8s-nodes"
	inventory_id = "1"
	group_id     = "${awx_inventory_group.k8s-nodes.id}"
	host_id      = "${awx_host.testacc-host_1.id}"
  }
`
