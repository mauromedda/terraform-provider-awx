provider "awx" {
}


resource "awx_inventory" "test" {
	name = "test1002"
	organization = 1
}
resource "awx_inventory" "test1" {
	name = "test1003"
	organization = 1
}
