package awx

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	awxgo "github.com/mauromedda/awx-go"
)

func resourceOrganizationObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceOrganizationCreate,
		Read:   resourceOrganizationRead,
		Delete: resourceOrganizationDelete,
		Update: resourceOrganizationUpdate,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the organization.",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional organization decription.",
			},

			"custom_virtualenv": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The path of the custom virtualenv.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceOrganizationCreate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.OrganizationService
	_, res, err := awxService.ListOrganizations(map[string]string{
		"name": d.Get("name").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) >= 1 {
		return fmt.Errorf("Organization with name %s already exists",
			d.Get("name").(string))
	}

	result, err := awxService.CreateOrganization(map[string]interface{}{
		"name":              d.Get("name").(string),
		"description":       d.Get("description").(string),
		"custom_virtualenv": d.Get("custom_virtualenv").(string),
	}, map[string]string{})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceOrganizationRead(d, m)
}

func resourceOrganizationUpdate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.OrganizationService
	_, res, err := awxService.ListOrganizations(map[string]string{
		"id": d.Id()},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return fmt.Errorf("Organization with name %s doesn't exists",
			d.Get("name").(string))
	}
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	_, err = awxService.UpdateOrganization(id, map[string]interface{}{
		"name":              d.Get("name").(string),
		"description":       d.Get("description").(string),
		"custom_virtualenv": d.Get("custom_virtualenv").(string),
	}, map[string]string{})
	if err != nil {
		return err
	}

	return resourceOrganizationRead(d, m)
}

func resourceOrganizationRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.OrganizationService
	_, res, err := awxService.ListOrganizations(map[string]string{
		"name": d.Get("name").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}
	d = setOrganizationResourceData(d, res.Results[0])
	return nil
}

func resourceOrganizationDelete(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.OrganizationService
	id, err := strconv.Atoi(d.Id())
	_, res, err := awxService.ListOrganizations(map[string]string{
		"name": d.Get("name").(string)})
	if len(res.Results) == 0 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	if _, err = awxService.DeleteOrganization(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func setOrganizationResourceData(d *schema.ResourceData, r *awxgo.Organization) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("custom_virtualenv", r.CustomVirtualEnv)
	return d
}
