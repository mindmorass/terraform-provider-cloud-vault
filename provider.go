
package example

import (
    "context"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    secretmanager "cloud.google.com/go/secretmanager/apiv1"
    "google.golang.org/api/option"
    "github.com/hashicorp/terraform-provider-google/google"
)

func Provider() *schema.Provider {
    return &schema.Provider{
        ResourcesMap: map[string]*schema.Resource{
            "example_secret": resourceSecret(),
        },
        ConfigureContextFunc: providerConfigure,
    }
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
    var diags diag.Diagnostics

    // Get the Google provider configuration from context
    config := google.Config{
        Credentials: d.Get("credentials").(string),
    }

    client, err := secretmanager.NewClient(ctx, option.WithCredentialsJSON([]byte(config.CredentialsJSON)))
    if err != nil {
        return nil, diag.FromErr(err)
    }

    return client, diags
}
