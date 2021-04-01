// Copyright 2021 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpaas

import (
	"context"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/sirupsen/logrus"
	"github.com/tsuru/tsuru/cmd"
	"istio.io/pkg/log"

	rpaas_client "github.com/tsuru/rpaas-operator/pkg/rpaas/client"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Description: "Target to tsuru API",
				Optional:    true,
			},
			"token": {
				Type:        schema.TypeString,
				Description: "Token to authenticate on tsuru API (optional)",
				Optional:    true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"rpaas_autoscale": resourceRpaasAutoscale(),
			"rpaas_block":     resourceRpaasBlock(),
			"rpaas_route":     resourceRpaasRoute(),
		},
	}
	p.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		return providerConfigure(ctx, d, p.TerraformVersion)
	}
	return p
}

type rpaasProvider struct {
	RpaasClient rpaas_client.Client
	Log         *logrus.Logger
}

func providerConfigure(ctx context.Context, d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	logger := logrus.New()
	file, err := os.OpenFile("/tmp/rpaas-provider.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	host := d.Get("host").(string)
	if host == "" {
		target, err := cmd.GetTarget()
		if err != nil {
			return nil, diag.FromErr(err)
		}
		if target == "" {
			return nil, diag.Errorf("Tsuru target is empty")
		}
	} else {
		os.Setenv("TSURU_TARGET", host)
	}

	token := d.Get("token").(string)
	if token == "" {
		t, err := cmd.ReadToken()
		if err != nil {
			return nil, diag.FromErr(err)
		}
		if t == "" {
			return nil, diag.Errorf("Tsuru token is empty")
		}
		token = t
	}

	cli, err := rpaas_client.NewClientThroughTsuruWithOptions(
		host,
		token,
		"unset",
		rpaas_client.ClientOptions{
			Timeout: 10 * time.Second,
		},
	)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	p := &rpaasProvider{
		Log:         logger,
		RpaasClient: cli,
	}

	return p, nil
}
