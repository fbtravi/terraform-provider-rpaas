// Copyright 2021 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/tsuru/terraform-provider-rpaas/internal/provider"
)

//go:generate tfplugindocs

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
