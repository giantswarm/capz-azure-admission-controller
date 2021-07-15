# -*- mode: Python -*-

load('ext://restart_process', 'docker_build_with_restart')
load('ext://cert_manager', 'deploy_cert_manager')

def kubectl_apply(file):
    """Apply a file to the tilt k8s context"""
    local("kubectl --context {} apply -f {}".format(k8s_context(), file), quiet=True, echo_off=True)

deploy_cert_manager()

# Install CAPI/CAPZ CRDs
kubectl_apply('https://raw.githubusercontent.com/giantswarm/apiextensions/v3.27.1/helm/crds-common/templates/upstream.yaml')
kubectl_apply('https://raw.githubusercontent.com/giantswarm/apiextensions/v3.27.1/helm/crds-azure/templates/upstream.yaml')

# Install GiantSwarm app platform
local("apptestctl bootstrap --github-token {} --kubeconfig=\"$(kind get kubeconfig --name {})\"".format(os.getenv("OPSCTL_GITHUB_TOKEN"), k8s_context().lstrip("kind-")), quiet=True, echo_off=True)

# Install kyverno from upstream
kubectl_apply('https://raw.githubusercontent.com/kyverno/kyverno/main/definitions/release/install.yaml')

# Install GiantSwarm kyverno policies
kubectl_apply('kyverno-policies.yaml')

# Add task to create a cluster to test the admission controller.
local_resource(
  'azureCluster with empty rules',
  "kubectl --context {} delete --ignore-not-found=true -f azurecluster.yaml && kubectl --context {} apply -f azurecluster.yaml && kubectl --context {} get -f azurecluster.yaml -o yaml".format(k8s_context(), k8s_context(), k8s_context()),
  auto_init=False,
  trigger_mode=TRIGGER_MODE_MANUAL)

# Add task to run integration tests
test(
  'Run integration tests',
  'make test',
  auto_init=False,
  trigger_mode=TRIGGER_MODE_MANUAL)

local_resource(
  'compile',
  "make build",
  deps=["config", "controllers", "Dockerfile", "go.mod", "go.sum", "main.go"])

docker_build_with_restart(
  'giantswarm/capz-azure-admission-controller',
  '.',
  only=['./build/manager'],
  entrypoint=['/app/manager'],
  live_update=[
    sync('./build/manager', '/app/manager'),
  ],
)

k8s_yaml(kustomize('./config/dev'))
