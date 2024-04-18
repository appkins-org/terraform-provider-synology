package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/appkins-org/terraform-provider-synology/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceHostItem() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostItemRead,
		Schema: map[string]*schema.Schema{
			"hosts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"free_cpu_core": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"free_ram_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_cpu_core": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_ram_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceHostItemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)
	service := HostGuestService{synologyClient: client}

	hostResponse, err := service.Read()
	if err != nil {
		return diag.FromErr(err)
	}

	var hosts []map[string]interface{}
	for _, h := range hostResponse.Hosts {
		m := map[string]interface{}{
			"free_cpu_core":  h.FreeCpuCore,
			"free_ram_size":  h.FreeRamSize,
			"host_id":        h.HostID,
			"host_name":      h.HostName,
			"status":         h.Status,
			"total_cpu_core": h.TotalCpuCore,
			"total_ram_size": h.TotalRamSize,
		}
		hosts = append(hosts, m)
	}

	if err := d.Set("hosts", hosts); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
