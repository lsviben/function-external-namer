# function-external-namer

A [Crossplane] Composition Function template, for Go.

## What does it do?

The function just goes through the desired resources, and sets the 
[external-name-annotation] to the `metadata.name` of the resource.

If the resource already has an external name annotation, it will not be
overwritten. 

If the resource does not have a `metadata.name`, it will be skipped and 
Crossplane will generate a name for it.

## Future work

For now, this function only works with the `metadata.name` field. In the future,
it would be good to support other fields, or even fields from different 
resources.

## Developing this Function

This template doesn't use the typical Crossplane build submodule and Makefile,
since we'd like Functions to have a less heavyweight developer experience.
It mostly relies on regular old Go tools:

```shell

# Run tests
go test -cover ./...
ok  	github.com/crossplane/function-external-namer	0.019s	coverage: 58.3% of statements

# Lint the code
$ docker run --rm -v $(pwd):/app -v ~/.cache/golangci-lint/v1.54.2:/root/.cache -w /app golangci/golangci-lint:v1.54.2 golangci-lint run

# Build a Docker image - see Dockerfile
$ docker build .
```

This Function can be pushed to any Docker registry. To push to xpkg.upbound.io\
use `docker push` and `docker-credential-up` from
https://github.com/upbound/up/.

## Testing this Function

You can try your function out locally using [`xrender`][xrender]. With `xrender`
you can run a Function pipeline on your laptop.

```shell
# Install xrender
$ go install github.com/crossplane-contrib/xrender@latest

# Run it! 
$ xrender manifests/definition.yaml manifests/composition.yaml manigests/functions.yaml
---
apiVersion: apiextensions.crossplane.io/v1
kind: CompositeResourceDefinition
metadata:
  name: xnopresources.nop.example.org
---
apiVersion: nop.crossplane.io/v1alpha1
kind: NopResource
metadata:
  annotations:
    crossplane.io/composition-resource-name: no-name
  generateName: xnopresources.nop.example.org-
  labels:
    crossplane.io/composite: xnopresources.nop.example.org
  ownerReferences:
  - apiVersion: apiextensions.crossplane.io/v1
    blockOwnerDeletion: true
    controller: true
    kind: CompositeResourceDefinition
    name: xnopresources.nop.example.org
    uid: ""
---
apiVersion: nop.crossplane.io/v1alpha1
kind: NopResource
metadata:
  annotations:
    crossplane.io/composition-resource-name: annotated
    crossplane.io/external-name: annotated
  generateName: xnopresources.nop.example.org-
  labels:
    crossplane.io/composite: xnopresources.nop.example.org
  name: named
  ownerReferences:
  - apiVersion: apiextensions.crossplane.io/v1
    blockOwnerDeletion: true
    controller: true
    kind: CompositeResourceDefinition
    name: xnopresources.nop.example.org
    uid: ""
---
apiVersion: nop.crossplane.io/v1alpha1
kind: NopResource
metadata:
  annotations:
    crossplane.io/composition-resource-name: named
    crossplane.io/external-name: named
  generateName: xnopresources.nop.example.org-
  labels:
    crossplane.io/composite: xnopresources.nop.example.org
  name: named
  ownerReferences:
  - apiVersion: apiextensions.crossplane.io/v1
    blockOwnerDeletion: true
    controller: true
    kind: CompositeResourceDefinition
    name: xnopresources.nop.example.org
    uid: ""
```


[Crossplane]: https://crossplane.io
[external-name-annotation]: https://docs.crossplane.io/v1.13/concepts/managed-resources/#naming-external-resources
[xrender]: https://github.com/crossplane-contrib/xrender