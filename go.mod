module github.com/giantswarm/capz-azure-admission-controller

go 1.13

require (
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/onsi/ginkgo v1.14.2
	github.com/onsi/gomega v1.10.3
	k8s.io/api v0.17.14
	k8s.io/apimachinery v0.17.14
	k8s.io/client-go v0.17.14
	sigs.k8s.io/cluster-api-provider-azure v0.4.13
	sigs.k8s.io/controller-runtime v0.5.14
)
