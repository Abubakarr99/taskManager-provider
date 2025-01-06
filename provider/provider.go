package provider

import (
	"context"
	"github.com/Abubakarr99/taskManager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"taskmanager_task": resourceTask(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"taskmanager_task": dataSourceTask(),
		},
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TASKMANAGER_HOST", nil),
			},
		},
		ConfigureContextFunc: configureProvider,
	}
}

// configureProvider configures the terraform provider for taskManager
func configureProvider(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	host, ok := d.Get("host").(string)
	if !ok {
		return nil, diag.Errorf("the host (127.0.0.1:6742) must be provided or TASKMANAGER_HOST env")
	}
	c, err := client.New(host)
	if err != nil {
		return nil, append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "unable to create the task manager client",
			Detail:   "Unable to connect to the task manager client",
		})
	}
	return c, diags
}
