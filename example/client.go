
package example

import (
    "context"
    secretmanager "cloud.google.com/go/secretmanager/apiv1"
    secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func CreateSecret(client *secretmanager.Client, projectID, secretID, secretValue string) error {
    ctx := context.Background()

    // Create the secret
    req := &secretmanagerpb.CreateSecretRequest{
        Parent:   "projects/" + projectID,
        SecretId: secretID,
        Secret: &secretmanagerpb.Secret{
            Replication: &secretmanagerpb.Replication{
                Replication: &secretmanagerpb.Replication_Automatic_{
                    Automatic: &secretmanagerpb.Replication_Automatic{},
                },
            },
        },
    }

    secret, err := client.CreateSecret(ctx, req)
    if err != nil {
        return err
    }

    // Add a secret version
    reqVersion := &secretmanagerpb.AddSecretVersionRequest{
        Parent: secret.Name,
        Payload: &secretmanagerpb.SecretPayload{
            Data: []byte(secretValue),
        },
    }

    _, err = client.AddSecretVersion(ctx, reqVersion)
    if err != nil {
        return err
    }

    return nil
}
