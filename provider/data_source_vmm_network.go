package provider

import (
	"context"
	"strconv"
	"time"

	"github.com/appkins-org/terraform-provider-synology/client"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func datasourceNetworkItem() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkItemRead,
		Schema: map[string]*schema.Schema{
			"networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNetworkItemRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client := m.(client.SynologyClient)
	service := NetworkGuestService{synologyClient: client}

	networkResponse, err := service.Read()
	if err != nil {
		return diag.FromErr(err)
	}

	var networks []map[string]interface{}
	for _, n := range networkResponse.Networks {
		network := map[string]interface{}{
			"network_id":   n.NetworkID,
			"network_name": n.NetworkName,
		}
		networks = append(networks, network)
	}

	if err := d.Set("networks", networks); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
