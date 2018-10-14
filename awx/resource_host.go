package awx

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	awxgo "github.com/mauromedda/awx-go"
)

func resourceHostObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceHostCreate,
		Read:   resourceHostRead,
		Delete: resourceHostDelete,
		Update: resourceHostUpdate,

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
			"inventory_id": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"group_ids": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
			"enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  "",
			},
			"instance_id": &schema.Schema{
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

func resourceHostCreate(d *schema.ResourceData, m interface{}) error {

	awx := m.(*awxgo.AWX)
	awxService := awx.HostService

	inv := d.Get("inventory_id").(int)
	_, res, _ := awxService.ListHosts(map[string]string{
		"name":      d.Get("name").(string),
		"inventory": strconv.Itoa(inv)},
	)
	if len(res.Results) >= 1 {
		return fmt.Errorf("Host %s with id %d already exists", res.Results[0].Name, res.Results[0].ID)
	}

	result, err := awxService.CreateHost(map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory_id").(int),
		"enabled":     d.Get("enabled").(bool),
		"instance_id": d.Get("instance_id").(string),
		"variables":   d.Get("variables").(string),
	}, map[string]string{})
	if err != nil {
		return err
	}

	hostID := result.ID
	if d.IsNewResource() {
		rawGroups := d.Get("group_ids").([]interface{})
		for _, v := range rawGroups {

			_, err := awxService.AssociateGroup(hostID, map[string]interface{}{
				"id": v.(int),
			}, map[string]string{})
			if err != nil {
				return err
			}
		}
	}
	d.SetId(strconv.Itoa(result.ID))
	return resourceHostRead(d, m)
}

func resourceHostUpdate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.HostService
	_, res, _ := awxService.ListHosts(map[string]string{"id": d.Id()})
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	if len(res.Results) >= 1 {

		_, err = awxService.UpdateHost(id, map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"inventory":   d.Get("inventory_id").(int),
			"enabled":     d.Get("enabled").(bool),
			"instance_id": d.Get("instance_id").(string),
			"variables":   d.Get("variables").(string),
		}, nil)
		if err != nil {
			return err
		}

		if d.HasChange("group_ids") {
			rawGroups := d.Get("group_ids").([]interface{})
			for _, v := range rawGroups {

				_, err := awxService.AssociateGroup(id, map[string]interface{}{
					"id": v.(int),
				}, map[string]string{})
				if err != nil {
					return err
				}
			}
		}
		return resourceHostRead(d, m)
	}

	return fmt.Errorf("Host %s with id %d doesn't exist", d.Get("name").(string), id)

}

func resourceHostRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.HostService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Host %d not found", id)
	}
	_, res, err := awxService.ListHosts(map[string]string{"id": d.Id()})
	if err != nil {
		return err
	}
	d = setHostResourceData(d, res.Results[0])
	return nil
}

func resourceHostDelete(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.HostService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	if _, err := awxService.DeleteHost(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func setHostResourceData(d *schema.ResourceData, r *awxgo.Host) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("inventory_id", r.Inventory)
	d.Set("enabled", r.Enabled)
	d.Set("instance_id", r.InstanceID)
	d.Set("variables", normalizeJsonYaml(r.Variables))
	d.Set("group_ids", d.Get("group_ids").([]interface{}))
	return d
}
