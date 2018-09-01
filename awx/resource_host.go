package awx

import (
	"fmt"
	"strconv"

	awxgo "github.com/Colstuwjx/awx-go"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:     schema.TypeInt,
				Optional: true,
				Default:  "",
			},
			"inventory": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
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

	_, res, _ := awxService.ListHosts(map[string]string{"name": d.Get("name").(string)})
	if len(res.Results) >= 1 {
		return fmt.Errorf("Host %s with id %d already exists", res.Results[0].Name, res.Results[0].ID)
	}

	result, err := awxService.CreateHost(map[string]interface{}{
		"name":         d.Get("name").(string),
		"organization": d.Get("organization").(string),
		"description":  d.Get("description").(string),
		"kind":         d.Get("kind").(string),
		"host_filter":  d.Get("host_filter").(string),
		"variables":    d.Get("variables").(string),
	}, map[string]string{})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceHostRead(d, m)

}

func resourceHostUpdate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.HostService
	_, res, _ := awxService.ListInventories(map[string]string{"name": d.Get("name").(string)})
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	if len(res.Results) >= 1 {

		_, err = awxService.UpdateHost(id, map[string]interface{}{
			"name":         d.Get("name").(string),
			"organization": d.Get("organization").(string),
			"description":  d.Get("description").(string),
			"kind":         d.Get("kind").(string),
			"host_filter":  d.Get("host_filter").(string),
			"variables":    d.Get("variables").(string),
		}, nil)
		if err != nil {
			return err
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
	r, err := awxService.GetHost(id, map[string]string{})
	if err != nil {
		return err
	}
	d = setHostResourceData(d, r)
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
	d.Set("organization", strconv.Itoa(r.Organization))
	d.Set("description", r.Description)
	d.Set("kind", r.Kind)
	d.Set("host_filter", r.HostFilter)
	d.Set("variables", normalizeJsonYaml(r.Variables))
	return d
}
