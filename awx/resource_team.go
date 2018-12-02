package awx

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	awxgo "github.com/mauromedda/awx-go"
)

func resourceTeamObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamCreate,
		Read:   resourceTeamRead,
		Delete: resourceTeamDelete,
		Update: resourceTeamUpdate,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the Team.",
			},

			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional team decription.",
			},

			"organization_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the organization.",
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

func resourceTeamCreate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.TeamService
	_, res, err := awxService.ListTeams(map[string]string{
		"name": d.Get("name").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) >= 1 {
		return fmt.Errorf("Team with name %s already exists",
			d.Get("name").(string))
	}

	result, err := awxService.CreateTeam(map[string]interface{}{
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"organization": d.Get("organization_id").(string),
	}, map[string]string{})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceTeamRead(d, m)
}

func resourceTeamUpdate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.TeamService
	_, res, err := awxService.ListTeams(map[string]string{
		"id": d.Id()},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return fmt.Errorf("Team with name %s doesn't exists",
			d.Get("name").(string))
	}
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	_, err = awxService.UpdateTeam(id, map[string]interface{}{
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"organization": d.Get("organization_id").(string),
	}, map[string]string{})
	if err != nil {
		return err
	}

	return resourceTeamRead(d, m)
}

func resourceTeamRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.TeamService
	_, res, err := awxService.ListTeams(map[string]string{
		"name": d.Get("name").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}
	d = setTeamResourceData(d, res.Results[0])
	return nil
}

func resourceTeamDelete(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.TeamService
	id, err := strconv.Atoi(d.Id())
	_, res, err := awxService.ListTeams(map[string]string{
		"name": d.Get("name").(string)})
	if len(res.Results) == 0 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	if _, err = awxService.DeleteTeam(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func setTeamResourceData(d *schema.ResourceData, r *awxgo.Team) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("organization", r.Organization)
	return d
}
