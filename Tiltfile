# -*- mode: Python -*-

load('ext://restart_process', 'docker_build_with_restart')
load('ext://cert_manager', 'deploy_cert_manager')

deploy_cert_manager()
local('opsctl ensure crds -k ~/.kube/config -p azure --crds azureclusters.infrastructure.cluster.x-k8s.io', quiet=True)

local_resource(
  'manager',
  'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/manager ./',
  deps=["config", "controllers", "Dockerfile", "go.mod", "go.sum", "main.go"])

docker_build_with_restart(
  'giantswarm/capz-azure-admission-controller',
  '.',
  entrypoint=['/app/build/manager'],
  live_update=[
    sync('./build', '/app/build/'),
  ],
)

k8s_yaml(kustomize('./config/default'))

# Add task to create an AzureCluster CR to test the admission controller.
local_resource(
  'azureCluster with empty rules',
  'kubectl delete --ignore-not-found=true -f azurecluster.yaml && kubectl apply -f azurecluster.yaml && kubectl get -f azurecluster.yaml -o yaml',
  auto_init=False,
  trigger_mode=TRIGGER_MODE_MANUAL)
