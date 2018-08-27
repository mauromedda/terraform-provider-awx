provider "awx" {
}


resource "awx_inventory" "test" {
	name = "test1002"
	organization = 1
}

resource "awx_inventory" "test1" {
	name = "test1003"
	organization = 1
	variables = <<EOF
---
test: 1
other: 2
EOF
}

resource "awx_inventory_group" "test1_gr1" {
	name = "test_group"
	inventory = "${awx_inventory.test1.id}"
	description = "test group 1"
}
