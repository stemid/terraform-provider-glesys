package glesys

import (
	"context"

	"github.com/glesys/glesys-go/v6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGlesysLoadBalancerTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlesysLoadBalancerTargetCreate,
		ReadContext:   resourceGlesysLoadBalancerTargetRead,
		UpdateContext: resourceGlesysLoadBalancerTargetUpdate,
		DeleteContext: resourceGlesysLoadBalancerTargetDelete,

		Description: "Create a LoadBalancer Target for a `glesys_loadbalancer_backend`.",

		Schema: map[string]*schema.Schema{
			"backend": {
				Description: "Backend to associate with.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"enabled": {
				Description: "Enable or disable Target. `true`, `false`",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},

			"loadbalancerid": {
				Description: "LoadBalancer ID.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"name": {
				Description: "Target name.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},

			"port": {
				Description: "Target port to connect to.",
				Type:        schema.TypeInt,
				Required:    true,
			},

			"status": {
				Description: "Target status. `UP`, `DOWN`",
				Type:        schema.TypeString,
				Computed:    true,
			},

			"targetip": {
				Description: "Target IP.",
				Type:        schema.TypeString,
				Required:    true,
			},

			"weight": {
				Description: "Target weight. `1-256`. Higher weight gets more requests.",
				Type:        schema.TypeInt,
				Required:    true,
			},
		},
	}
}

func resourceGlesysLoadBalancerTargetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Add target to glesys_loadbalancer_backend resource
	client := m.(*glesys.Client)

	params := glesys.AddTargetParams{
		Backend:  d.Get("backend").(string),
		Name:     d.Get("name").(string),
		Port:     d.Get("port").(int),
		TargetIP: d.Get("targetip").(string),
		Weight:   d.Get("weight").(int),
	}

	loadbalancerID := d.Get("loadbalancerid").(string)

	_, err := client.LoadBalancers.AddTarget(ctx, loadbalancerID, params)
	if err != nil {
		return diag.Errorf("Error creating LoadBalancer Target: %s", err)
	}

	if !d.Get("enabled").(bool) {
		// Disable the target after creation
		targetParams := glesys.ToggleTargetParams{Name: d.Get("name").(string), Backend: d.Get("backend").(string)}
		_, err := client.LoadBalancers.DisableTarget(ctx, loadbalancerID, targetParams)

		if err != nil {
			return diag.Errorf("could not disable Target during creation: %s", err)
		}
	}

	d.SetId(d.Get("name").(string))

	return resourceGlesysLoadBalancerTargetRead(ctx, d, m)
}

func resourceGlesysLoadBalancerTargetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)
	lb, err := client.LoadBalancers.Details(ctx, loadbalancerid)
	if err != nil {
		diag.Errorf("loadbalancer not found: %s", err)
		d.SetId("")
		return nil
	}

	// iterate over all backends && targets of the loadbalancer
	for _, n := range lb.BackendsList {
		if n.Name == d.Get("backend").(string) {
			for _, t := range n.Targets {
				if t.Name == d.Get("name").(string) {
					d.Set("enabled", t.Enabled)
					d.Set("port", t.Port)
					d.Set("status", t.Status)
					d.Set("targetip", t.TargetIP)
					d.Set("weight", t.Weight)
				}
			}
		} else {
			diag.Errorf("loadbalancer Target not found: %s", d.Get("name").(string))
			d.SetId("")
			return nil
		}
	}

	return nil
}

func resourceGlesysLoadBalancerTargetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)

	params := glesys.EditTargetParams{
		Backend: d.Get("backend").(string),
		Name:    d.Get("name").(string),
	}

	if d.HasChange("name") {
		params.Name = d.Get("name").(string)
	}

	// If changed, toggle the enabled state of the target.
	if d.HasChange("enabled") {
		currentState, _ := d.GetChange("enabled")

		toggleParams := glesys.ToggleTargetParams{
			Backend: d.Get("backend").(string),
			Name:    d.Get("name").(string),
		}

		if currentState == true {
			_, err := client.LoadBalancers.DisableTarget(ctx, loadbalancerid, toggleParams)
			if err != nil {
				return diag.Errorf("error toggling LoadBalancer Target from: %s - %s", currentState, err)
			}
		} else {
			_, err := client.LoadBalancers.EnableTarget(ctx, loadbalancerid, toggleParams)
			if err != nil {
				return diag.Errorf("error toggling LoadBalancer Target from: %s - %s", currentState, err)
			}
		}
	}

	if d.HasChange("port") {
		params.Port = d.Get("port").(int)
	}

	if d.HasChange("targetip") {
		params.TargetIP = d.Get("targetip").(string)
	}

	if d.HasChange("weight") {
		params.Weight = d.Get("weight").(int)
	}

	_, err := client.LoadBalancers.EditTarget(ctx, loadbalancerid, params)

	if err != nil {
		return diag.Errorf("Error updating LoadBalancer Target: %s", err)
	}

	return resourceGlesysLoadBalancerTargetRead(ctx, d, m)
}

func resourceGlesysLoadBalancerTargetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*glesys.Client)

	loadbalancerid := d.Get("loadbalancerid").(string)

	params := glesys.RemoveTargetParams{
		Backend: d.Get("backend").(string),
		Name:    d.Get("name").(string),
	}

	err := client.LoadBalancers.RemoveTarget(ctx, loadbalancerid, params)
	if err != nil {
		return diag.Errorf("Error deleting LoadBalancer Target: %s", err)
	}

	d.SetId("")
	return nil
}
