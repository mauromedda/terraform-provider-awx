package awx

import (
	"fmt"
	"strconv"

	awxgo "github.com/Colstuwjx/awx-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceInventoryObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceInventoryCreate,
		Read:   resourceInventoryRead,
		Delete: resourceInventoryDelete,
		Update: resourceInventoryUpdate,

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
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"host_filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
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

func resourceInventoryCreate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.InventoriesService

	_, res, _ := awxService.ListInventories(map[string]string{
		"name":         d.Get("name").(string),
		"organization": d.Get("organization_id").(string),
	})
	if len(res.Results) >= 1 {
		return fmt.Errorf("Inventory %s with id %d already exists", res.Results[0].Name, res.Results[0].ID)
	}

	result, err := awxService.CreateInventory(map[string]interface{}{
		"name":         d.Get("name").(string),
		"organization": d.Get("organization_id").(string),
		"description":  d.Get("description").(string),
		"kind":         d.Get("kind").(string),
		"host_filter":  d.Get("host_filter").(string),
		"variables":    d.Get("variables").(string),
	}, map[string]string{})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventoryRead(d, m)

}

func resourceInventoryUpdate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.InventoriesService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	_, res, _ := awxService.ListInventories(map[string]string{"id": d.Id()})
	if len(res.Results) >= 1 {

		_, err = awxService.UpdateInventory(id, map[string]interface{}{
			"name":         d.Get("name").(string),
			"organization": d.Get("organization_id").(string),
			"description":  d.Get("description").(string),
			"kind":         d.Get("kind").(string),
			"host_filter":  d.Get("host_filter").(string),
			"variables":    d.Get("variables").(string),
		}, nil)
		if err != nil {
			return err
		}

		return resourceInventoryRead(d, m)
	}

	return fmt.Errorf("Inventory %s with id %d doesn't exist", d.Get("name").(string), id)

}

func resourceInventoryRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.InventoriesService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Inventory %d not found", id)
	}
	r, err := awxService.GetInventory(id, map[string]string{})
	if err != nil {
		return err
	}
	d = setInventoryResourceData(d, r)
	return nil
}

func resourceInventoryDelete(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.InventoriesService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	if _, err := awxService.DeleteInventory(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func setInventoryResourceData(d *schema.ResourceData, r *awxgo.Inventory) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("organization_id", strconv.Itoa(r.Organization))
	d.Set("description", r.Description)
	d.Set("kind", r.Kind)
	d.Set("host_filter", r.HostFilter)
	d.Set("variables", normalizeJsonYaml(r.Variables))
	return d
}
