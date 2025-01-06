package provider

import (
	"context"
	"errors"
	"github.com/Abubakarr99/taskManager/client"
	pb "github.com/Abubakarr99/taskManager/proto"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
)

func dataSourceTask() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaskRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique ID of the task",
			},
			"start_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The start date to filter task",
			},
			"end_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The end date to filter task",
			},
			"urgency": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceTaskRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := clientFromMeta(meta)
	if err != nil {
		diag.FromErr(err)
	}
	tasks, err := findTaskInTaskManager(ctx, client, findTaskRequest{
		ID:        d.Get("id").(string),
		urgency:   d.Get("urgency").(string),
		startDate: d.Get("start_date").(string),
		endDate:   d.Get("end_date").(string),
	})
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

func clientFromMeta(meta interface{}) (*client.Client, error) {
	tsClient, ok := meta.(*client.Client)
	if !ok {
		return nil, errors.New("meta does not contain a task Client")
	}
	return tsClient, nil
}

type findTaskRequest struct {
	ID        string
	urgency   string
	startDate string
	endDate   string
}

func findTaskInTaskManager(ctx context.Context, tsclient *client.Client, req findTaskRequest) ([]*client.Task, error) {
	searchReq := &pb.SearchTaskReq{}
	if req.urgency != "" {
		searchReq.Urgency = []pb.Urgency{parseUrgency(req.urgency)}
	}
	if req.startDate != "" {
		searchReq.DueRange = &pb.DateRange{
			Start: parseDate(req.startDate),
			End:   parseDate(req.endDate),
		}
	}
	resp, err := tsclient.SearchTasks(ctx, searchReq)
	if err != nil {
		log.Fatalf("Failed to search tasks: %v", err)
	}
	var tasks []*client.Task
	for task := range resp {
		task := task
		if task.Error() != nil {
			log.Fatalf("cloud not fetch tasks error %v", task.Error())
		}
		if req.ID == "" || task.Id == req.ID {
			tasks = append(tasks, &task)
		}
	}
	return tasks, nil
}
