package provider

import (
	"context"
	"fmt"
	pb "github.com/Abubakarr99/taskManager/proto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/genproto/googleapis/type/date"
)

func resourceTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaskCreate,
		ReadContext:   resourceTaskRead,
		UpdateContext: resourceTaskUpdate,
		DeleteContext: resourceTaskDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique ID of the task",
			},
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The title of the task",
			},
			"urgency": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The level of urgency of the task",
			},
		},
	}
}

// resourceTaskCreate creates the task
func resourceTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clientFromMeta(meta)
	if err != nil {
		diag.FromErr(err)
	}
	task := &pb.Task{
		Title:   d.Get("title").(string),
		Urgency: parseUrgency(d.Get("urgency").(string)),
		DueDate: parseDate(d.Get("due_date").(string)),
	}
	resp, err := client.AddTasks(ctx, []*pb.Task{task})
	if err != nil {
		return append(diag.FromErr(err))
	}
	d.SetId(resp[0])
	return nil
}

func resourceTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := clientFromMeta(meta)
	if err != nil {
		diag.FromErr(err)
	}
	Tasks := []*pb.Task{
		{
			Id:      d.Id(),
			Title:   d.Get("title").(string),
			Urgency: parseUrgency(d.Get("urgency").(string)),
			DueDate: parseDate(d.Get("due_date").(string)),
		},
	}
	err = client.UpdateTasks(ctx, Tasks)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	// TO DO: implement a search here to find a task before calling the update.
	return diags
}

func resourceTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client, err := clientFromMeta(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	Ids := []string{d.Id()}
	err = client.DeleteTasks(ctx, Ids)
	if err != nil {
		return append(diags, diag.FromErr(err)...)
	}
	d.SetId("")
	return nil
}

func resourceTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clientFromMeta(meta)
	if err != nil {
		return diag.FromErr(err)
	}
	tasks, err := findTaskInTaskManager(ctx, client, findTaskRequest{ID: d.Id()})
	if err != nil {
		return diag.FromErr(err)
	}
	if len(tasks) == 0 {
		return nil
	}
	d.Set("title", tasks[0].Title)
	d.Set("urgency", tasks[0].Urgency.String())
	d.Set("due_date", tasks[0].DueDate.String())
	return nil
}

func parseUrgency(urgency string) pb.Urgency {
	switch urgency {
	case "VERY_HIGH":
		return pb.Urgency_VERY_HIGH
	case "MODERATE":
		return pb.Urgency_MODERATE
	case "LOW":
		return pb.Urgency_LOW
	default:
		return pb.Urgency_LOW
	}
}

func parseDate(dateStr string) *date.Date {
	var year, month, day int
	fmt.Sscanf(dateStr, "%d-%d-%d", &year, &month, &day)
	return &date.Date{Year: int32(year), Month: int32(month), Day: int32(day)}
}
