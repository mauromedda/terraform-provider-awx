#--- Terraform AWX Provider example file ---#

variable "name" {
  type    = "string"
  default = "foo"
}

provider "awx" {
  username = "admin"
  password = "password"
  endpoint = "http://192.168.99.100:31838"
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

resource "awx_host" "k8s-node-1" {
  name         = "k8s-node-1.awx.local"
  description  = "Kubernetes minion 1"
  inventory_id = "${awx_inventory.default.id}"

  variables = <<VARIABLES
---
api_server_enabled: false
VARIABLES
}

resource "awx_group_association" "k8s-node-1_k8s-nodes" {
  name         = "k8s-node-1_k8s-nodes"
  inventory_id = "${awx_inventory.default.id}"
  group_id     = "${awx_inventory_group.k8s-nodes.id}"
  host_id      = "${awx_host.k8s-node-1.id}"
}

resource "awx_group_association" "k8s-node-1_etcd" {
  name         = "k8s-node-1_etcd"
  inventory_id = "${awx_inventory.default.id}"
  group_id     = "${awx_inventory_group.etcd.id}"
  host_id      = "${awx_host.k8s-node-1.id}"
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
