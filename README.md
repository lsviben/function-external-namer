# function-external-namer

A [Crossplane] Composition Function template, for Go.

## What does it do?

The function just goes through the desired resources, and sets their
[external-resource-annotation] to the name of the resource. This is useful
if users want to control the name of the resource created in the external
system.

If the resource already has an external name annotation, it will not be
overwritten. 

If the resource does not have a name, it will be skipped and Crossplane
will generate a name for it.

## Developing this Function

This template doesn't use the typical Crossplane build submodule and Makefile,
since we'd like Functions to have a less heavyweight developer experience.
It mostly relies on regular old Go tools:

```shell
# Run code generation - see input/generate.go
$ go generate ./...

# Run tests
$ go test -cover ./...
?       github.com/crossplane/function-template-go/input/v1beta1      [no test files]
ok      github.com/crossplane/function-template-go    0.006s  coverage: 25.8% of statements

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

First you'll need to create a `functions.yaml` file. This tells `xrender` what
Functions to run, and how. In this case we want to run the Function you're
developing in 'Development mode'. That pretty much means you'll run the Function
manually and tell `xrender` where to find it.

```yaml
---
apiVersion: pkg.crossplane.io/v1beta1
kind: Function
metadata:
  name: function-external-namer # Use your Function's name!
  annotations:
    # xrender will try to talk to your Function at localhost:9443
    xrender.crossplane.io/runtime: Development
    xrender.crossplane.io/runtime-development-target: localhost:9443
```

Next, run the Function locally:

```shell
# Run your Function in insecure mode
go run . --insecure --debug
```

Once your Function is running, in another window you can use `xrender`.

```shell
# Install xrender
$ go install github.com/crossplane-contrib/xrender@latest

# Run it! See the xrender repo for these examples.
$ xrender examples/xr.yaml examples/composition.yaml examples/functions.yaml
---
apiVersion: nopexample.org/v1
kind: XBucket
metadata:
  name: test-xrender
status:
  bucketRegion: us-east-2
---
apiVersion: s3.aws.upbound.io/v1beta1
kind: Bucket
metadata:
  annotations:
    crossplane.io/composition-resource-name: my-bucket
  generateName: test-xrender-
  labels:
    crossplane.io/composite: test-xrender
  ownerReferences:
  - apiVersion: nopexample.org/v1
    blockOwnerDeletion: true
    controller: true
    kind: XBucket
    name: test-xrender
    uid: ""
spec:
  forProvider:
    region: us-east-2
```

You can see an example Composition above. There's also some examples in the
`xrender` repo's examples directory.


[Crossplane]: https://crossplane.io
[external-resource-annotation]: https://docs.crossplane.io/v1.13/concepts/managed-resources/#naming-external-resources
[xrender]: https://github.com/crossplane-contrib/xrender