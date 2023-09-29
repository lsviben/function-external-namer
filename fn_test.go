package main

import (
	"context"
	"testing"

	"github.com/crossplane/function-sdk-go/resource"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/crossplane/crossplane-runtime/pkg/logging"

	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/response"
)

func TestRunFunction(t *testing.T) {

	type args struct {
		ctx context.Context
		req *fnv1beta1.RunFunctionRequest
	}
	type want struct {
		rsp *fnv1beta1.RunFunctionResponse
		err error
	}

	cases := map[string]struct {
		reason string
		args   args
		want   want
	}{
		"SetExternalName": {
			reason: "Should set the external name annotation of the resource to its meta.name.",
			args: args{
				req: &fnv1beta1.RunFunctionRequest{
					Meta:     &fnv1beta1.RequestMeta{Tag: "hello"},
					Observed: &fnv1beta1.State{},
					Desired: &fnv1beta1.State{
						Resources: map[string]*fnv1beta1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestNamer",
									"metadata": {
										"name": "my-test-namer"
									}
								}`),
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1beta1.RunFunctionResponse{
					Meta: &fnv1beta1.ResponseMeta{Tag: "hello", Ttl: durationpb.New(response.DefaultTTL)},
					Desired: &fnv1beta1.State{
						Resources: map[string]*fnv1beta1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestNamer",
									"metadata": {
										"name": "my-test-namer",
                                        "annotations": {
 										   "crossplane.io/external-name": "my-test-namer"
										}	
									}
								}`),
							},
						},
					},
					Results: []*fnv1beta1.Result{
						{
							Severity: fnv1beta1.Severity_SEVERITY_NORMAL,
							Message:  "External names added successfully",
						},
					},
				},
			},
		},
		"DontSetExternalNameAlreadySet": {
			reason: "Should not set the external name annotation of the resource to its meta.name if it is already set.",
			args: args{
				req: &fnv1beta1.RunFunctionRequest{
					Meta:     &fnv1beta1.RequestMeta{Tag: "hello"},
					Observed: &fnv1beta1.State{},
					Desired: &fnv1beta1.State{
						Resources: map[string]*fnv1beta1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestNamer",
									"metadata": {
										"name": "my-test-namer",
                                        "annotations": {
 										   "crossplane.io/external-name": "existing"
										}	
									}
								}`),
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1beta1.RunFunctionResponse{
					Meta: &fnv1beta1.ResponseMeta{Tag: "hello", Ttl: durationpb.New(response.DefaultTTL)},
					Desired: &fnv1beta1.State{
						Resources: map[string]*fnv1beta1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestNamer",
									"metadata": {
										"name": "my-test-namer",
                                        "annotations": {
 										   "crossplane.io/external-name": "existing"
										}	
									}
								}`),
							},
						},
					},
					Results: []*fnv1beta1.Result{
						{
							Severity: fnv1beta1.Severity_SEVERITY_NORMAL,
							Message:  "External names added successfully",
						},
					},
				},
			},
		},
		"DontSetExternalNameEmpty": {
			reason: "Should not set the external name annotation of the resource to its meta.name the meta.name is empty",
			args: args{
				req: &fnv1beta1.RunFunctionRequest{
					Meta:     &fnv1beta1.RequestMeta{Tag: "hello"},
					Observed: &fnv1beta1.State{},
					Desired: &fnv1beta1.State{
						Resources: map[string]*fnv1beta1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestNamer",
									"metadata": {}
								}`),
							},
						},
					},
				},
			},
			want: want{
				rsp: &fnv1beta1.RunFunctionResponse{
					Meta: &fnv1beta1.ResponseMeta{Tag: "hello", Ttl: durationpb.New(response.DefaultTTL)},
					Desired: &fnv1beta1.State{
						Resources: map[string]*fnv1beta1.Resource{
							"ready-composed-resource": {
								Resource: resource.MustStructJSON(`{
									"apiVersion": "test.crossplane.io/v1",
									"kind": "TestNamer",
									"metadata": {}
								}`),
							},
						},
					},
					Results: []*fnv1beta1.Result{
						{
							Severity: fnv1beta1.Severity_SEVERITY_NORMAL,
							Message:  "External names added successfully",
						},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f := &Function{log: logging.NewNopLogger()}
			rsp, err := f.RunFunction(tc.args.ctx, tc.args.req)

			if diff := cmp.Diff(tc.want.rsp, rsp, protocmp.Transform()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want rsp, +got rsp:\n%s", tc.reason, diff)
			}

			if diff := cmp.Diff(tc.want.err, err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("%s\nf.RunFunction(...): -want err, +got err:\n%s", tc.reason, diff)
			}
		})
	}
}
