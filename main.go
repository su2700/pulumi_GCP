package main

import (
    "github.com/pulumi/pulumi-gcp/sdk/v5/go/gcp/compute"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        vm, err := compute.NewInstance(ctx, "vm-instance", &compute.InstanceArgs{
            MachineType: pulumi.String("n2-standard-4"),
            Zone:        pulumi.String("us-central1-a"),
            BootDisk: &compute.InstanceBootDiskArgs{
                InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
                    Image: pulumi.String("ubuntu-os-cloud/ubuntu-2004-lts"),
                    Size:  pulumi.Int(256),
                },
            },
            NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
                &compute.InstanceNetworkInterfaceArgs{
                    Network: pulumi.String("default"),
                    AccessConfigs: compute.InstanceNetworkInterfaceAccessConfigArray{
                        &compute.InstanceNetworkInterfaceAccessConfigArgs{},
                    },
                },
            },
        })
        if err != nil {
            return err
        }

        // Handling Outputs Properly
        // ApplyT on the entire NetworkInterfaces array to safely handle potential nil values
        vm.NetworkInterfaces.ApplyT(func(interfaces []compute.InstanceNetworkInterface) error {
            if len(interfaces) > 0 && len(interfaces[0].AccessConfigs) > 0 && interfaces[0].AccessConfigs[0].NatIp != nil {
                // Safe dereference of NatIp, which is a pointer to a string
                ctx.Export("instanceIP", pulumi.String(*interfaces[0].AccessConfigs[0].NatIp))
            }
            return nil
        })

        return nil
    })
}
