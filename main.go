
package main

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
    "github.com/mindmorass/terraform-cloud-vault/example"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: example.Provider,
    })
}
