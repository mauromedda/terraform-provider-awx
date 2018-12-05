#--- Terraform AWX Provider example file ---#

variable "name" {
  type    = "string"
  default = "foo"
}

provider "awx" {
  username = "admin"
  password = "password"
  endpoint = "http://192.168.99.100:30967"
}

data "template_file" "default_yaml" {
  template = "${file("${path.module}/example.yaml")}"

  vars {
    name = "${var.name}"
  }
}

resource "awx_inventory" "default" {
  name            = "alpha"
  organization_id = 1
  variables       = "${data.template_file.default_yaml.rendered}"
}

resource "awx_inventory_group" "k8s-nodes" {
  name         = "k8s-nodes"
  inventory_id = "${awx_inventory.default.id}"
}

resource "awx_inventory_group" "etcd" {
  name         = "etcd"
  inventory_id = "${awx_inventory.default.id}"
}

resource "awx_host" "k8s-nodes" {
  count        = 3
  name         = "k8s-node-${count.index}.awx.local"
  description  = "Kubernetes minion ${count.index}"
  inventory_id = "${awx_inventory.default.id}"
  group_ids    = ["${awx_inventory_group.etcd.id}", "${awx_inventory_group.k8s-nodes.id}"]

  variables = <<VARIABLES
---
api_server_enabled: false
VARIABLES
}

resource "awx_host" "k8s-node" {
  name         = "k8s-node-4.awx.local"
  description  = "Kubernetes minion ${count.index}"
  inventory_id = "${awx_inventory.default.id}"
  group_ids    = ["${awx_inventory_group.etcd.id}"]

  variables = <<VARIABLES
---
api_server_enabled: false
VARIABLES
}

resource "awx_group_association" "k8s-node-1_k8s-nodes" {
  name         = "k8s-node-1_k8s-nodes"
  inventory_id = "${awx_inventory.default.id}"
  group_id     = "${awx_inventory_group.k8s-nodes.id}"
  host_id      = "${awx_host.k8s-node.id}"
}

resource "awx_project" "alpha" {
  name                 = "alpha"
  scm_type             = "git"
  scm_url              = "https://github.com/ansible/ansible-tower-samples"
  scm_update_on_launch = true
  organization_id      = 1
}

resource "awx_job_template" "alpha" {
  name         = "alpha"
  description  = "Alpha job template example"
  project_id   = "${awx_project.alpha.id}"
  job_type     = "run"
  inventory_id = "${awx_inventory.default.id}"
  playbook     = "hello_world.yml"
}

resource "awx_user" "test" {
  username     = "mauromedda"
  password     = "password"
  email        = "medda.mauro@gmail.com"
  is_superuser = true
}

resource "awx_user_role" "inventory_test_admin" {
	user_id = "${awx_user.test.id}"
	resource_type = "inventory"
	resource_name = "${awx_inventory.default.name}"
	role = "admin"
	organization_id = 1
}

resource "awx_user_role" "project_alpha_use" {
        user_id = "${awx_user.test.id}"
        resource_type = "project"
        resource_name = "${awx_project.alpha.name}"
        role = "use"
        organization_id = 1
}

resource "awx_user_role" "jobtemplate_alpha_execute" {
        user_id = "${awx_user.test.id}"
        resource_type = "job_template"
        resource_name = "${awx_job_template.alpha.name}"
        role = "execute"
        organization_id = 1
}


resource "awx_team" "automation_team" {
        name = "automation_team"
        description = "Automation team"
        organization_id = 1
}

resource "awx_user_role" "automation_team_member" {
        user_id = "${awx_user.test.id}"
        resource_type = "team"
        resource_name = "${awx_team.automation_team.name}"
        role = "member"
        organization_id = 1
}

resource "awx_team_role" "automation_team_inventory_admin" {
        team_id = "${awx_team.automation_team.id}"
        resource_type = "inventory"
        resource_name = "${awx_inventory.default.name}"
        role = "admin"
        organization_id = 1
}

resource "awx_organization" "deftunix_org" {
	name = "deftunix-org"
	description = "deftunix organization"
}
