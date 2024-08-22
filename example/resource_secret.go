
package example

import (
    "context"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
    secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func resourceSecret() *schema.Resource {
    return &schema.Resource{
        Create: resourceSecretCreate,
        Read:   resourceSecretRead,
        Delete: resourceSecretDelete,

        Schema: map[string]*schema.Schema{
            "username": {
                Type:     schema.TypeString,
                Optional: true,
                Computed: true,
            },
            "project_id": {
                Type:     schema.TypeString,
                Required: true,
            },
            "secret_id": {
                Type:     schema.TypeString,
                Required: true,
            },
        },
    }
}

func resourceSecretCreate(d *schema.ResourceData, m interface{}) error {
    client := m.(*secretmanager.Client)
    projectID := d.Get("project_id").(string)
    secretID := d.Get("secret_id").(string)

    // Generate a random password
    password, err := bcrypt.GenerateFromPassword([]byte(uuid.New().String()), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    // Generate or use provided username
    username := d.Get("username").(string)
    if username == "" {
        username = "user_" + uuid.New().String()
        d.Set("username", username)
    }

    // Store the secret in GCP
    err = CreateSecret(client, projectID, secretID, string(password))
    if err != nil {
        return err
    }

    // Set resource ID (but not storing the password in state)
    d.SetId(secretID)

    return resourceSecretRead(d, m)
}

func resourceSecretRead(d *schema.ResourceData, m interface{}) error {
    // No sensitive data is read back into the state
    return nil
}

func resourceSecretDelete(d *schema.ResourceData, m interface{}) error {
    client := m.(*secretmanager.Client)
    projectID := d.Get("project_id").(string)
    secretID := d.Id()

    ctx := context.Background()
    req := &secretmanagerpb.DeleteSecretRequest{
        Name: "projects/" + projectID + "/secrets/" + secretID,
    }

    if err := client.DeleteSecret(ctx, req); err != nil {
        return err
    }

    // Remove the resource from the state
    d.SetId("")
    return nil
}
