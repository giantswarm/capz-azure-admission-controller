# capz-azure-admission-controller
Admission controller that contains GiantSwarm business logic for CAPI/CAPZ types.

## Development
### Tilt
You can start `Tilt` to run the controller and test it locally with this command:
```
make tilt-up
```

You can remove the `kind` cluster created for local testing with this command:
```
make kind-reset
```

### Tests
The automated tests rely on the CAPZ repository to be present on the filesystem so that we can register the CAPZ CRDs.
It needs to be at the same level than this repository under the name `cluster-api-provider-azure`.
