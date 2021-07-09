package controllers

import (
	"context"
	"time"

	"github.com/Azure/go-autorest/autorest/to"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("AzureClusterMutator", func() {
	BeforeEach(func() {})
	AfterEach(func() {})

	Context("Default an AzureCluster", func() {
		It("should set the default security rules to allow GiantSwarm services like our monitoring", func() {
			ctx := context.TODO()

			azureCluster := &capz.AzureCluster{
				ObjectMeta: metav1.ObjectMeta{Name: "my-cluster", Namespace: "default"},
				Spec: capz.AzureClusterSpec{
					NetworkSpec: capz.NetworkSpec{
						Subnets: capz.Subnets{
							&capz.SubnetSpec{
								CIDRBlocks:    []string{"10.0.0.0/16"},
								Role:          "control-plane",
								Name:          "control-plane-subnet",
								SecurityGroup: capz.SecurityGroup{},
							},
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, azureCluster)).To(Succeed())
			defer func() {
				err := k8sClient.Delete(ctx, azureCluster)
				Expect(err).NotTo(HaveOccurred())
			}()

			Eventually(func() error {
				return k8sClient.Get(ctx, client.ObjectKey{Name: azureCluster.Name, Namespace: azureCluster.Namespace}, azureCluster)
			}, time.Second*5, time.Millisecond*500).Should(BeNil())

			cpSubnet := azureCluster.Spec.NetworkSpec.GetControlPlaneSubnet()
			Expect(cpSubnet).ShouldNot(BeNil())

			Expect(cpSubnet.SecurityGroup.IngressRules).Should(ContainElement(&capz.IngressRule{
				Name:             "allow_load_balancer",
				Description:      "Allow all TCP traffic from LB to master instance",
				Priority:         3902,
				Protocol:         capz.SecurityGroupProtocolTCP,
				Source:           to.StringPtr("AzureLoadBalancer"),
				SourcePorts:      to.StringPtr("*"),
				Destination:      to.StringPtr(cpSubnet.CIDRBlocks[0]),
				DestinationPorts: to.StringPtr("*"),
			}))
		})
	})
})
