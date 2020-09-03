package main

import (
	"time"

	"github.com/pulumi/pulumi-azure/sdk/v3/go/azure/containerservice"
	"github.com/pulumi/pulumi-azure/sdk/v3/go/azure/core"
	"github.com/pulumi/pulumi-docker/sdk/v2/go/docker"
	"github.com/pulumi/pulumi/sdk/v2/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create an Azure Resource Group
		resourceGroup, err := core.NewResourceGroup(ctx, "resourceGroup", &core.ResourceGroupArgs{
			Location: pulumi.String("WestUS"),
		})
		if err != nil {
			return err
		}

		// customImage is a folder containing a Dockerfile.
		customImage := "my-app"
		registry, err := containerservice.NewRegistry(ctx, "myregistry", &containerservice.RegistryArgs{
			AdminEnabled:      pulumi.BoolPtr(true),
			ResourceGroupName: resourceGroup.Name,
			Sku:               pulumi.StringPtr("Basic"),
		})
		if err != nil {
			return err
		}

		// TODO: REMOVE!!
		time.Sleep(10 * time.Second)
		//

		_, err = docker.NewImage(ctx, "myimage", &docker.ImageArgs{
			ImageName: pulumi.Sprintf("%s/%s:v1.0.0", registry.LoginServer, customImage),
			Registry: docker.ImageRegistryArgs{
				Server:   registry.LoginServer,
				Username: registry.AdminUsername,
				Password: registry.AdminPassword,
			},
			Build: docker.DockerBuildArgs{
				Context: pulumi.Sprintf("./%s", customImage),
			},
		})
		if err != nil {
			return err
		}

		return nil
	})
}
