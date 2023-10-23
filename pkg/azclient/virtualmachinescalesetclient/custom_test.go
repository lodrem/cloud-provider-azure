// /*
// Copyright The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

// Code generated by client-gen. DO NOT EDIT.
package virtualmachinescalesetclient

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	armcompute "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5"
	armnetwork "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	networkClientFactory  *armnetwork.ClientFactory
	virtualNetworksClient *armnetwork.VirtualNetworksClient
	vNet                  *armnetwork.VirtualNetwork
)

func init() {
	additionalTestCases = func() {
		When("get a non existing resource", func() {
			It("should return error", func(ctx context.Context) {
				_, err := realClient.Get(ctx, resourceGroupName, resourceName, nil)
				Expect(err).To(HaveOccurred())
			})
		})
	}

	beforeAllFunc = func(ctx context.Context) {
		networkClientFactory, err := armnetwork.NewClientFactory(recorder.SubscriptionID(), recorder.TokenCredential(), &arm.ClientOptions{
			ClientOptions: azcore.ClientOptions{
				Transport: recorder.HTTPClient(),
			},
		})
		Expect(err).NotTo(HaveOccurred())
		virtualNetworksClient = networkClientFactory.NewVirtualNetworksClient()
		vnetpoller, err := virtualNetworksClient.BeginCreateOrUpdate(ctx, resourceGroupName, "vnet1", armnetwork.VirtualNetwork{
			Location: to.Ptr(location),
			Properties: &armnetwork.VirtualNetworkPropertiesFormat{
				AddressSpace: &armnetwork.AddressSpace{
					AddressPrefixes: []*string{
						to.Ptr("10.1.0.0/16"),
					},
				},
				Subnets: []*armnetwork.Subnet{
					{
						Name: to.Ptr("subnet1"),
						Properties: &armnetwork.SubnetPropertiesFormat{
							AddressPrefix: to.Ptr("10.1.0.0/24"),
						},
					},
				},
			},
		}, nil)
		Expect(err).NotTo(HaveOccurred())

		vnetresp, err := vnetpoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 1 * time.Second,
		})
		Expect(err).NotTo(HaveOccurred())
		vNet = &vnetresp.VirtualNetwork
		newResource = &armcompute.VirtualMachineScaleSet{
			Location: to.Ptr(location),
			SKU: &armcompute.SKU{
				Name:     to.Ptr("Basic_A0"), //armcompute.VirtualMachineSizeTypesBasicA0
				Capacity: to.Ptr[int64](1),
			},
			Properties: &armcompute.VirtualMachineScaleSetProperties{
				Overprovision: to.Ptr(false),
				UpgradePolicy: &armcompute.UpgradePolicy{
					Mode: to.Ptr(armcompute.UpgradeModeManual),
					AutomaticOSUpgradePolicy: &armcompute.AutomaticOSUpgradePolicy{
						EnableAutomaticOSUpgrade: to.Ptr(false),
						DisableAutomaticRollback: to.Ptr(false),
					},
				},
				VirtualMachineProfile: &armcompute.VirtualMachineScaleSetVMProfile{
					OSProfile: &armcompute.VirtualMachineScaleSetOSProfile{
						ComputerNamePrefix: to.Ptr("vmss"),
						AdminUsername:      to.Ptr("sample-user"),
						AdminPassword:      to.Ptr("Password01!@#"),
					},
					StorageProfile: &armcompute.VirtualMachineScaleSetStorageProfile{
						ImageReference: &armcompute.ImageReference{
							Offer:     to.Ptr("WindowsServer"),
							Publisher: to.Ptr("MicrosoftWindowsServer"),
							SKU:       to.Ptr("2019-Datacenter"),
							Version:   to.Ptr("latest"),
						},
					},
					NetworkProfile: &armcompute.VirtualMachineScaleSetNetworkProfile{
						NetworkInterfaceConfigurations: []*armcompute.VirtualMachineScaleSetNetworkConfiguration{
							{
								Name: to.Ptr(resourceName),
								Properties: &armcompute.VirtualMachineScaleSetNetworkConfigurationProperties{
									Primary:            to.Ptr(true),
									EnableIPForwarding: to.Ptr(true),
									IPConfigurations: []*armcompute.VirtualMachineScaleSetIPConfiguration{
										{
											Name: to.Ptr(resourceName),
											Properties: &armcompute.VirtualMachineScaleSetIPConfigurationProperties{
												Subnet: &armcompute.APIEntityReference{
													ID: vNet.Properties.Subnets[0].ID,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}
	}
	afterAllFunc = func(ctx context.Context) {
		vnetPoller, err := virtualNetworksClient.BeginDelete(ctx, resourceGroupName, *vNet.Name, nil)
		Expect(err).NotTo(HaveOccurred())
		_, err = vnetPoller.PollUntilDone(ctx, &runtime.PollUntilDoneOptions{
			Frequency: 1 * time.Second,
		})
		Expect(err).NotTo(HaveOccurred())
	}
}