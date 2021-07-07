package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Azure/go-autorest/autorest/to"
	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	cadvisorPort         = "4194"
	etcdPort             = "2379"
	kubeletPort          = "10250"
	nodeExporterPort     = "10300"
	kubeStateMetricsPort = "10301"
)

// +kubebuilder:rbac:groups=exp-infrastructure.cluster-x.k8s.io/v1alpha3,resources=azureclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:webhook:path=/mutate-exp-infrastructure-cluster-x-k8s-io-v1alpha3-azurecluster,mutating=true,failurePolicy=fail,groups="infrastructure.cluster.x-k8s.io",resources=azureclusters,verbs=create;update,versions=v1alpha3,name=azureclustermutator.giantswarm.io

type AzureClusterMutator struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (a *AzureClusterMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	azureCluster := &capz.AzureCluster{}

	err := a.decoder.Decode(req, azureCluster)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// We need to pass the management cluster CIDR somehow, like passing arguments when installing this app.
	managementClusterCIDR := ""

	// @todo: Document which priority range we use so that customers can override our rules.
	cpSubnet := azureCluster.Spec.NetworkSpec.GetControlPlaneSubnet()
	cpSubnet.SecurityGroup.IngressRules = append(cpSubnet.SecurityGroup.IngressRules, &capz.IngressRule{
		Name:             "allow_load_balancer",
		Description:      "Allow all TCP traffic from LB to master instance",
		Priority:         3902,
		Protocol:         capz.SecurityGroupProtocolTCP,
		Source:           to.StringPtr("AzureLoadBalancer"),
		SourcePorts:      to.StringPtr("*"),
		Destination:      to.StringPtr(cpSubnet.CIDRBlocks[0]),
		DestinationPorts: to.StringPtr("*"),
	})
	cpSubnet.SecurityGroup.IngressRules = append(cpSubnet.SecurityGroup.IngressRules, &capz.IngressRule{
		Name:             "allow_apiserver",
		Description:      "Allow K8s API Server",
		Priority:         2201,
		Protocol:         capz.SecurityGroupProtocolTCP,
		Source:           to.StringPtr("*"),
		SourcePorts:      to.StringPtr("*"),
		Destination:      to.StringPtr(cpSubnet.CIDRBlocks[0]),
		DestinationPorts: to.StringPtr(strconv.Itoa(int(azureCluster.Spec.ControlPlaneEndpoint.Port))),
	})
	cpSubnet.SecurityGroup.IngressRules = append(cpSubnet.SecurityGroup.IngressRules, &capz.IngressRule{
		Name:             "allow_cadvisor",
		Description:      "Allow host cluster Prometheus to reach Cadvisors",
		Priority:         3500,
		Protocol:         capz.SecurityGroupProtocolTCP,
		Source:           to.StringPtr(managementClusterCIDR),
		SourcePorts:      to.StringPtr("*"),
		Destination:      to.StringPtr("VirtualNetwork"),
		DestinationPorts: to.StringPtr(cadvisorPort),
	})
	cpSubnet.SecurityGroup.IngressRules = append(cpSubnet.SecurityGroup.IngressRules, &capz.IngressRule{
		Name:             "allow_kubelets",
		Description:      "Allow host cluster Prometheus to reach Kubelets",
		Priority:         3500,
		Protocol:         capz.SecurityGroupProtocolTCP,
		Source:           to.StringPtr(managementClusterCIDR),
		SourcePorts:      to.StringPtr("*"),
		Destination:      to.StringPtr("VirtualNetwork"),
		DestinationPorts: to.StringPtr(kubeletPort),
	})
	cpSubnet.SecurityGroup.IngressRules = append(cpSubnet.SecurityGroup.IngressRules, &capz.IngressRule{
		Name:             "allow_node_exporters",
		Description:      "Allow host cluster Prometheus to reach node-exporters",
		Priority:         3500,
		Protocol:         capz.SecurityGroupProtocolTCP,
		Source:           to.StringPtr(managementClusterCIDR),
		SourcePorts:      to.StringPtr("*"),
		Destination:      to.StringPtr("VirtualNetwork"),
		DestinationPorts: to.StringPtr(nodeExporterPort),
	})
	cpSubnet.SecurityGroup.IngressRules = append(cpSubnet.SecurityGroup.IngressRules, &capz.IngressRule{
		Name:             "allow_kube_state_metrics",
		Description:      "Allow host cluster Prometheus to reach kube-state-metrics",
		Priority:         3500,
		Protocol:         capz.SecurityGroupProtocolTCP,
		Source:           to.StringPtr(managementClusterCIDR),
		SourcePorts:      to.StringPtr("*"),
		Destination:      to.StringPtr("VirtualNetwork"),
		DestinationPorts: to.StringPtr(kubeStateMetricsPort),
	})
	cpSubnet.SecurityGroup.IngressRules = append(cpSubnet.SecurityGroup.IngressRules, &capz.IngressRule{
		Name:             "allow_ssh",
		Description:      "Allow SSH from management cluster",
		Priority:         2200,
		Protocol:         capz.SecurityGroupProtocolTCP,
		Source:           to.StringPtr(managementClusterCIDR),
		SourcePorts:      to.StringPtr("*"),
		Destination:      to.StringPtr("VirtualNetwork"),
		DestinationPorts: to.StringPtr("22"),
	})
	cpSubnet.SecurityGroup.IngressRules = append(cpSubnet.SecurityGroup.IngressRules, &capz.IngressRule{
		Name:             "etcd_load_balancer_rule_host",
		Description:      "Allow management cluster nodes to reach the etcd loadbalancer",
		Priority:         3901,
		Protocol:         capz.SecurityGroupProtocolTCP,
		Source:           to.StringPtr(managementClusterCIDR),
		SourcePorts:      to.StringPtr("*"),
		Destination:      to.StringPtr("VirtualNetwork"),
		DestinationPorts: to.StringPtr(etcdPort),
	})

	marshaledAzureCluster, err := json.Marshal(azureCluster)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledAzureCluster)
}

// InjectDecoder injects the decoder.
func (a *AzureClusterMutator) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
