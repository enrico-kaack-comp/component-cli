// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gardener/component-cli/ociclient (interfaces: Client)

// Package mock_ociclient is a generated GoMock package.
package mock_ociclient

import (
	context "context"
	io "io"
	reflect "reflect"

	ociclient "github.com/gardener/component-cli/ociclient"
	oci "github.com/gardener/component-cli/ociclient/oci"
	gomock "github.com/golang/mock/gomock"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Fetch mocks base method.
func (m *MockClient) Fetch(arg0 context.Context, arg1 string, arg2 v1.Descriptor, arg3 io.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Fetch indicates an expected call of Fetch.
func (mr *MockClientMockRecorder) Fetch(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockClient)(nil).Fetch), arg0, arg1, arg2, arg3)
}

// GetManifest mocks base method.
func (m *MockClient) GetManifest(arg0 context.Context, arg1 string) (*v1.Manifest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManifest", arg0, arg1)
	ret0, _ := ret[0].(*v1.Manifest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManifest indicates an expected call of GetManifest.
func (mr *MockClientMockRecorder) GetManifest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManifest", reflect.TypeOf((*MockClient)(nil).GetManifest), arg0, arg1)
}

// GetOCIArtifact mocks base method.
func (m *MockClient) GetOCIArtifact(arg0 context.Context, arg1 string) (*oci.Artifact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOCIArtifact", arg0, arg1)
	ret0, _ := ret[0].(*oci.Artifact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOCIArtifact indicates an expected call of GetOCIArtifact.
func (mr *MockClientMockRecorder) GetOCIArtifact(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOCIArtifact", reflect.TypeOf((*MockClient)(nil).GetOCIArtifact), arg0, arg1)
}

// GetRawManifest mocks base method.
func (m *MockClient) GetRawManifest(arg0 context.Context, arg1 string) (v1.Descriptor, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRawManifest", arg0, arg1)
	ret0, _ := ret[0].(v1.Descriptor)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetRawManifest indicates an expected call of GetRawManifest.
func (mr *MockClientMockRecorder) GetRawManifest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRawManifest", reflect.TypeOf((*MockClient)(nil).GetRawManifest), arg0, arg1)
}

// PushBlob mocks base method.
func (m *MockClient) PushBlob(arg0 context.Context, arg1 string, arg2 v1.Descriptor, arg3 ...ociclient.PushOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PushBlob", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// PushBlob indicates an expected call of PushBlob.
func (mr *MockClientMockRecorder) PushBlob(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushBlob", reflect.TypeOf((*MockClient)(nil).PushBlob), varargs...)
}

// PushManifest mocks base method.
func (m *MockClient) PushManifest(arg0 context.Context, arg1 string, arg2 *v1.Manifest, arg3 ...ociclient.PushOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PushManifest", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// PushManifest indicates an expected call of PushManifest.
func (mr *MockClientMockRecorder) PushManifest(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushManifest", reflect.TypeOf((*MockClient)(nil).PushManifest), varargs...)
}

// PushOCIArtifact mocks base method.
func (m *MockClient) PushOCIArtifact(arg0 context.Context, arg1 string, arg2 *oci.Artifact, arg3 ...ociclient.PushOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PushOCIArtifact", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// PushOCIArtifact indicates an expected call of PushOCIArtifact.
func (mr *MockClientMockRecorder) PushOCIArtifact(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushOCIArtifact", reflect.TypeOf((*MockClient)(nil).PushOCIArtifact), varargs...)
}

// PushRawManifest mocks base method.
func (m *MockClient) PushRawManifest(arg0 context.Context, arg1 string, arg2 v1.Descriptor, arg3 []byte, arg4 ...ociclient.PushOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PushRawManifest", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// PushRawManifest indicates an expected call of PushRawManifest.
func (mr *MockClientMockRecorder) PushRawManifest(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushRawManifest", reflect.TypeOf((*MockClient)(nil).PushRawManifest), varargs...)
}

// Resolve mocks base method.
func (m *MockClient) Resolve(arg0 context.Context, arg1 string) (string, v1.Descriptor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resolve", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(v1.Descriptor)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Resolve indicates an expected call of Resolve.
func (mr *MockClientMockRecorder) Resolve(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resolve", reflect.TypeOf((*MockClient)(nil).Resolve), arg0, arg1)
}
