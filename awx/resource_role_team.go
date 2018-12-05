package awx

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	awxgo "github.com/mauromedda/awx-go"
)

func resourceTeamRoleObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamRoleGrant,
		Read:   resourceTeamRoleRead,
		Delete: resourceTeamRoleRevoke,

		Schema: map[string]*schema.Schema{
			"team_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					validResourceTypes := map[string]bool{"admin": true, "read": true, "use": true,
						"member": true, "execute": true, "adhoc": true, "update": true, "auditor": true,
						"project admin": true, "workflow admin": true, "inventory admin": true, "job template admin": true}
					value := v.(string)
					if !validResourceTypes[value] {
						errors = append(errors, fmt.Errorf("%q must match one of the valid vaules", k))
					}
					return
				},
			},
			"resource_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					validResourceTypes := map[string]bool{"inventory": true, "team": true, "organization": true,
						"job_template": true, "credential": true, "project": true}
					value := v.(string)
					if !validResourceTypes[value] {
						errors = append(errors, fmt.Errorf("%q must match one of inventory, team, organization, job_template, credential or project", k))
					}
					return
				},
			},
			"resource_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func resourceTeamRoleGrant(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.TeamService
	_, res, err := awxService.ListTeams(map[string]string{
		"id": d.Get("team_id").(string)},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return fmt.Errorf("Team with Id %s doesn't exists",
			d.Get("team_id").(string))
	}
	id, _ := strconv.Atoi(d.Get("team_id").(string))
	roleID, err := getRoleID(d, m)
	if err == nil {
		err = awxService.GrantRole(id, roleID)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	d.SetId(d.Get("team_id").(string))
	return resourceTeamRoleRead(d, m)

}

func resourceTeamRoleRevoke(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.TeamService

	_, res, err := awxService.ListTeams(map[string]string{
		"id": d.Get("team_id").(string)},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return fmt.Errorf("Team with Id %s doesn't exists",
			d.Get("team_id").(string))
	}
	roleID, err := getRoleID(d, m)
	if err == nil {
		id, _ := strconv.Atoi(d.Get("team_id").(string))
		err = awxService.RevokeRole(id, roleID)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	d.SetId("")
	return resourceTeamRoleRead(d, m)
}

func resourceTeamRoleRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.TeamService
	_, res, err := awxService.ListTeams(map[string]string{
		"id": d.Get("team_id").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}
	d = setTeamRoleResourceData(d, res.Results[0])
	return nil
}

func setTeamRoleResourceData(d *schema.ResourceData, r *awxgo.Team) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("team_id", r.ID)
	d.Set("resource_name", d.Get("resource_name").(string))
	d.Set("resource_type", d.Get("resource_type").(string))
	d.Set("role", d.Get("role").(string))
	return d
}
