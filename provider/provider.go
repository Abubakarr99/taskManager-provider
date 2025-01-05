package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"taskmanager_task": resourceTask(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"taskmanager_task": datasourceTask(),
		},
		Schema: map[string]*schema.Schema{
			"taskmanager_host": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TASKMANAGER_HOST", nil),
			},
		},
	}
}

// configureProvider configures the terraform provider for taskManager
func configureProvider(d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	host := d.Get("taskmanager_host").(string)
	client := NewClient(host)
	return client, nil
}
