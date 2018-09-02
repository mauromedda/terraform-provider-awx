provider "awx" {
}


resource "awx_inventory" "test" {
	name = "test1001"
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

resource "awx_inventory" "test2" {
        name = "test1002"
        organization = 1
        variables = <<EOF
---
test: 1
other: 2
EOF
}

resource "awx_host" "host1" {
	name = "host1"
	inventory = "${awx_inventory.test1.id}"
	description = "prova host1"
	variables = "ansible_host: localhost"
}

resource "awx_inventory_group" "test1_gr1" {
	name = "test_group"
	inventory = "${awx_inventory.test1.id}"
	description = "test group 1 ciao"
}

resource "awx_inventory_group" "test1_gr2" {
        name = "test_group_2"
        inventory = "${awx_inventory.test2.id}"
        description = "test group 2 ciao"
}

resource "awx_group_association" "host_group_1" {
	name = "association-route-table"
	inventory = "${awx_inventory.test1.id}"
	host_id = "${awx_host.host1.id}"
	group_id = "${awx_inventory_group.test1_gr1.id}"
}

