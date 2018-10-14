package awx

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	awxgo "github.com/mauromedda/awx-go"
)

func resourceJobTemplateObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobTemplateCreate,
		Read:   resourceJobTemplateRead,
		Delete: resourceJobTemplateDelete,
		Update: resourceJobTemplateUpdate,

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
			// Run, Check, Scan
			"job_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "One of: run, check, scan",
			},
			"inventory_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"playbook": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"credential_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vault_credential_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"forks": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"limit": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			//0,1,2,3,4,5
			"verbosity": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "One of 0,1,2,3,4,5",
			},
			"extra_vars": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"job_tags": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"force_handlers": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"skip_tags": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"start_at_task": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"use_fact_cache": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"host_config_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ask_diff_mode_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_limit_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_tags_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_verbosity_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_inventory_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_variables_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_credential_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"survey_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"become_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"diff_mode": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"ask_skip_tags_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_simultaneous": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"custom_virtualenv": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ask_job_type_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceJobTemplateCreate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService
	var jobID int
	var finished time.Time
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"name":    d.Get("name").(string),
		"project": d.Get("project_id").(string)},
	)

	if len(res.Results) >= 1 {
		return fmt.Errorf("JobTemplate with name %s already exists",
			d.Get("name").(string))
	}
	_, prj, err := awx.ProjectService.ListProjects(map[string]string{
		"id": d.Get("project_id").(string)},
	)
	if err != nil {
		return err
	}
	if prj.Results[0].SummaryFields.CurrentJob["id"] != nil {
		jobID = int(prj.Results[0].SummaryFields.CurrentJob["id"].(float64))
	} else if prj.Results[0].SummaryFields.LastJob["id"] != nil {
		jobID = int(prj.Results[0].SummaryFields.LastJob["id"].(float64))
	}

	if jobID != 0 {
		// check if finished is 0
		for finished.IsZero() {
			prj, _ := awx.ProjectUpdatesService.ProjectUpdateGet(jobID)
			if prj != nil {
				finished = prj.Finished
				time.Sleep(1 * time.Second)
			}
		}
	}
	result, err := awxService.CreateJobTemplate(map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"job_type":                 d.Get("job_type").(string),
		"inventory":                AtoipOr(d.Get("inventory_id").(string), nil),
		"project":                  AtoipOr(d.Get("project_id").(string), nil),
		"playbook":                 d.Get("playbook").(string),
		"forks":                    d.Get("forks").(int),
		"limit":                    d.Get("limit").(string),
		"verbosity":                d.Get("verbosity").(int),
		"extra_vars":               d.Get("extra_vars").(string),
		"job_tags":                 d.Get("job_tags").(string),
		"force_handlers":           d.Get("force_handlers").(bool),
		"skip_tags":                d.Get("skip_tags").(string),
		"start_at_task":            d.Get("start_at_task").(string),
		"timeout":                  d.Get("timeout").(int),
		"use_fact_cache":           d.Get("use_fact_cache").(bool),
		"host_config_key":          d.Get("host_config_key").(string),
		"ask_diff_mode_on_launch":  d.Get("ask_diff_mode_on_launch").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"ask_tags_on_launch":       d.Get("ask_tags_on_launch").(bool),
		"ask_skip_tags_on_launch":  d.Get("ask_skip_tags_on_launch").(bool),
		"ask_job_type_on_launch":   d.Get("ask_job_type_on_launch").(bool),
		"ask_verbosity_on_launch":  d.Get("ask_verbosity_on_launch").(bool),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_credential_on_launch": d.Get("ask_credential_on_launch").(bool),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"become_enabled":           d.Get("become_enabled").(bool),
		"diff_mode":                d.Get("diff_mode").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"custom_virtualenv":        AtoipOr(d.Get("custom_virtualenv").(string), nil),
		"credential":               AtoipOr(d.Get("credential_id").(string), nil),
		"vault_credential":         AtoipOr(d.Get("vault_credential_id").(string), nil),
	}, map[string]string{})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceJobTemplateRead(d, m)
}

func resourceJobTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"id":      d.Id(),
		"project": d.Get("project_id").(string)},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return fmt.Errorf("JobTemplate with name %s doesn't exists",
			d.Get("name").(string))
	}
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	_, err = awxService.UpdateJobTemplate(id, map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"job_type":                 d.Get("job_type").(string),
		"inventory":                AtoipOr(d.Get("inventory_id").(string), nil),
		"project":                  AtoipOr(d.Get("project_id").(string), nil),
		"playbook":                 d.Get("playbook").(string),
		"forks":                    d.Get("forks").(int),
		"limit":                    d.Get("limit").(string),
		"verbosity":                d.Get("verbosity").(int),
		"extra_vars":               d.Get("extra_vars").(string),
		"job_tags":                 d.Get("job_tags").(string),
		"force_handlers":           d.Get("force_handlers").(bool),
		"skip_tags":                d.Get("skip_tags").(string),
		"start_at_task":            d.Get("start_at_task").(string),
		"timeout":                  d.Get("timeout").(int),
		"use_fact_cache":           d.Get("use_fact_cache").(bool),
		"host_config_key":          d.Get("host_config_key").(string),
		"ask_diff_mode_on_launch":  d.Get("ask_diff_mode_on_launch").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"ask_tags_on_launch":       d.Get("ask_tags_on_launch").(bool),
		"ask_skip_tags_on_launch":  d.Get("ask_skip_tags_on_launch").(bool),
		"ask_job_type_on_launch":   d.Get("ask_job_type_on_launch").(bool),
		"ask_verbosity_on_launch":  d.Get("ask_verbosity_on_launch").(bool),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_credential_on_launch": d.Get("ask_credential_on_launch").(bool),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"become_enabled":           d.Get("become_enabled").(bool),
		"diff_mode":                d.Get("diff_mode").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"custom_virtualenv":        AtoipOr(d.Get("custom_virtualenv").(string), nil),
		"credential":               AtoipOr(d.Get("credential_id").(string), nil),
		"vault_credential":         AtoipOr(d.Get("vault_credential_id").(string), nil),
	}, map[string]string{})
	if err != nil {
		return err
	}

	return resourceJobTemplateRead(d, m)
}

func resourceJobTemplateRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"name":    d.Get("name").(string),
		"project": d.Get("project_id").(string)},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}
	d = setJobTemplateResourceData(d, res.Results[0])
	return nil
}

func resourceJobTemplateDelete(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"id":      d.Id(),
		"project": d.Get("project_id").(string)},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		d.SetId("")
		return nil
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	_, err = awxService.DeleteJobTemplate(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func setJobTemplateResourceData(d *schema.ResourceData, r *awxgo.JobTemplate) *schema.ResourceData {
	d.Set("allow_simultaneous", r.AllowSimultaneous)
	d.Set("ask_credential_on_launch", r.AskCredentialOnLaunch)
	d.Set("ask_job_type_on_launch", r.AskJobTypeOnLaunch)
	d.Set("ask_limit_on_launch", r.AskLimitOnLaunch)
	d.Set("ask_skip_tags_on_launch", r.AskSkipTagsOnLaunch)
	d.Set("ask_tags_on_launch", r.AskTagsOnLaunch)
	d.Set("ask_variables_on_launch", r.AskVariablesOnLaunch)
	d.Set("credential_id", r.Credential)
	d.Set("description", r.Description)
	d.Set("extra_vars", r.ExtraVars)
	d.Set("force_handlers", r.ForceHandlers)
	d.Set("forks", r.Forks)
	d.Set("host_config_key", r.HostConfigKey)
	d.Set("inventory_id", r.Inventory)
	d.Set("job_tags", r.JobTags)
	d.Set("job_type", r.JobType)
	d.Set("diff_mode", r.DiffMode)
	d.Set("custom_virtualenv", r.CustomVirtualenv)
	d.Set("vault_credential_id", r.VaultCredential)
	d.Set("limit", r.Limit)
	d.Set("name", r.Name)
	d.Set("become_enabled", r.BecomeEnabled)
	d.Set("use_fact_cache", r.UseFactCache)
	d.Set("playbook", r.Playbook)
	d.Set("project_id", r.Project)
	d.Set("skip_tags", r.SkipTags)
	d.Set("start_at_task", r.StartAtTask)
	d.Set("survey_enabled", r.SurveyEnabled)
	d.Set("verbosity", r.Verbosity)
	return d
}
