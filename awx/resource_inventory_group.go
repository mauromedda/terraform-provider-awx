package awx

import (
	"fmt"
	"strconv"

	awxgo "github.com/Colstuwjx/awx-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceInventoryGroupObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceInventoryGroupCreate,
		Read:   resourceInventoryGroupRead,
		Update: resourceInventoryGroupUpdate,
		Delete: resourceInventoryGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"inventory": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"variables": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "",
				StateFunc: normalizeJsonYaml,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceInventoryGroupCreate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.GroupService

	_, res, _ := awxService.ListGroups(map[string]string{
		"name":      d.Get("name").(string),
		"inventory": d.Get("inventory").(string),
	})
	if len(res.Results) >= 1 {
		return fmt.Errorf("InventoryGroup %s with id %d already exists", res.Results[0].Name, res.Results[0].ID)
	}

	result, err := awxService.CreateGroup(map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory").(string),
		"variables":   d.Get("variables").(string),
	}, map[string]string{})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventoryGroupRead(d, m)

}

func resourceInventoryGroupUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceInventoryGroupDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceInventoryGroupRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.GroupService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("InventoryGroup %d not found", id)
	}
	_, res, err := awxService.ListGroups(map[string]string{
		"name":      d.Get("name").(string),
		"inventory": d.Get("inventory").(string),
	})
	if err != nil {
		return err
	}
	d = setInventoryGroupResourceData(d, res.Results[0])
	return nil
}

func setInventoryGroupResourceData(d *schema.ResourceData, r *awxgo.Group) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("inventory", r.Inventory)
	d.Set("variables", normalizeJsonYaml(r.Variables))
	return d
}
