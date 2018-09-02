package awx

import (
	"fmt"
	"strconv"
	"time"

	awxgo "github.com/Colstuwjx/awx-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGroupAssociationObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceGroupAssociationCreate,
		Read:   resourceGroupAssociationRead,
		Delete: resourceGroupAssociationDelete,
		Update: resourceGroupAssociationUpdate,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"inventory_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"host_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"group_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func resourceGroupAssociationCreate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxServiceHost := awx.HostService
	awxServiceGroup := awx.GroupService
	var id, inv int
	id = d.Get("host_id").(int)
	inv = d.Get("inventory_id").(int)
	_, resHost, _ := awxServiceHost.ListHosts(map[string]string{
		"id":        strconv.Itoa(id),
		"inventory": strconv.Itoa(inv)},
	)
	if len(resHost.Results) == 0 {
		return fmt.Errorf("Host %d not found in inventory %d", d.Get("host_id").(int), inv)
	}
	id = d.Get("group_id").(int)
	_, resGroup, _ := awxServiceGroup.ListGroups(map[string]string{
		"id":        strconv.Itoa(id),
		"inventory": strconv.Itoa(inv)},
	)
	if len(resGroup.Results) == 0 {
		return fmt.Errorf("Group %d not found in inventory %d", d.Get("group_id").(int), inv)
	}

	result, err := awxServiceHost.AssociateGroup(d.Get("host_id").(int), map[string]interface{}{
		"id": d.Get("group_id").(int),
	}, map[string]string{})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceGroupAssociationRead(d, m)
}

func resourceGroupAssociationUpdate(d *schema.ResourceData, m interface{}) error {
	d = setGroupAssociationResourceData(d)
	return nil
}

func resourceGroupAssociationRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceGroupAssociationDelete(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxServiceHost := awx.HostService
	awxServiceGroup := awx.GroupService
	var id, inv int
	id = d.Get("host_id").(int)
	inv = d.Get("inventory_id").(int)
	_, resHost, _ := awxServiceHost.ListHosts(map[string]string{
		"id":        strconv.Itoa(id),
		"inventory": strconv.Itoa(inv)},
	)
	if len(resHost.Results) == 0 {
		return fmt.Errorf("Host %d not found in inventory %d", d.Get("host_id").(int), inv)
	}
	id = d.Get("group_id").(int)
	_, resGroup, _ := awxServiceGroup.ListGroups(map[string]string{
		"id":        strconv.Itoa(id),
		"inventory": strconv.Itoa(inv)},
	)
	if len(resGroup.Results) == 0 {
		return fmt.Errorf("Group %d not found in inventory %d", d.Get("group_id").(int), inv)
	}

	_, err := awxServiceHost.DisAssociateGroup(d.Get("host_id").(int), map[string]interface{}{
		"id": d.Get("group_id").(int),
	}, map[string]string{})
	if err != nil {
		return err
	}

	d.SetId("")
	return resourceGroupAssociationRead(d, m)
}

func setGroupAssociationResourceData(d *schema.ResourceData) *schema.ResourceData {
	d.Set("name", d.Get("name").(string))
	d.Set("inventory_id", d.Get("inventory_id").(int))
	d.Set("host_id", d.Get("host_id").(int))
	d.Set("group_id", d.Get("group_id").(int))
	return d

}
