// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/linode/linode-cloud-controller-manager/cloud/linode (interfaces: Client)

// Package linode is a generated GoMock package.
package linode

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	linodego "github.com/linode/linodego"
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

// CreateFirewall mocks base method.
func (m *MockClient) CreateFirewall(arg0 context.Context, arg1 linodego.FirewallCreateOptions) (*linodego.Firewall, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFirewall", arg0, arg1)
	ret0, _ := ret[0].(*linodego.Firewall)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFirewall indicates an expected call of CreateFirewall.
func (mr *MockClientMockRecorder) CreateFirewall(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFirewall", reflect.TypeOf((*MockClient)(nil).CreateFirewall), arg0, arg1)
}

// CreateFirewallDevice mocks base method.
func (m *MockClient) CreateFirewallDevice(arg0 context.Context, arg1 int, arg2 linodego.FirewallDeviceCreateOptions) (*linodego.FirewallDevice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFirewallDevice", arg0, arg1, arg2)
	ret0, _ := ret[0].(*linodego.FirewallDevice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFirewallDevice indicates an expected call of CreateFirewallDevice.
func (mr *MockClientMockRecorder) CreateFirewallDevice(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFirewallDevice", reflect.TypeOf((*MockClient)(nil).CreateFirewallDevice), arg0, arg1, arg2)
}

// CreateNodeBalancer mocks base method.
func (m *MockClient) CreateNodeBalancer(arg0 context.Context, arg1 linodego.NodeBalancerCreateOptions) (*linodego.NodeBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNodeBalancer", arg0, arg1)
	ret0, _ := ret[0].(*linodego.NodeBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNodeBalancer indicates an expected call of CreateNodeBalancer.
func (mr *MockClientMockRecorder) CreateNodeBalancer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNodeBalancer", reflect.TypeOf((*MockClient)(nil).CreateNodeBalancer), arg0, arg1)
}

// CreateNodeBalancerConfig mocks base method.
func (m *MockClient) CreateNodeBalancerConfig(arg0 context.Context, arg1 int, arg2 linodego.NodeBalancerConfigCreateOptions) (*linodego.NodeBalancerConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNodeBalancerConfig", arg0, arg1, arg2)
	ret0, _ := ret[0].(*linodego.NodeBalancerConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNodeBalancerConfig indicates an expected call of CreateNodeBalancerConfig.
func (mr *MockClientMockRecorder) CreateNodeBalancerConfig(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNodeBalancerConfig", reflect.TypeOf((*MockClient)(nil).CreateNodeBalancerConfig), arg0, arg1, arg2)
}

// DeleteFirewall mocks base method.
func (m *MockClient) DeleteFirewall(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFirewall", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFirewall indicates an expected call of DeleteFirewall.
func (mr *MockClientMockRecorder) DeleteFirewall(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFirewall", reflect.TypeOf((*MockClient)(nil).DeleteFirewall), arg0, arg1)
}

// DeleteFirewallDevice mocks base method.
func (m *MockClient) DeleteFirewallDevice(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFirewallDevice", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFirewallDevice indicates an expected call of DeleteFirewallDevice.
func (mr *MockClientMockRecorder) DeleteFirewallDevice(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFirewallDevice", reflect.TypeOf((*MockClient)(nil).DeleteFirewallDevice), arg0, arg1, arg2)
}

// DeleteNodeBalancer mocks base method.
func (m *MockClient) DeleteNodeBalancer(arg0 context.Context, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNodeBalancer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNodeBalancer indicates an expected call of DeleteNodeBalancer.
func (mr *MockClientMockRecorder) DeleteNodeBalancer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNodeBalancer", reflect.TypeOf((*MockClient)(nil).DeleteNodeBalancer), arg0, arg1)
}

// DeleteNodeBalancerConfig mocks base method.
func (m *MockClient) DeleteNodeBalancerConfig(arg0 context.Context, arg1, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteNodeBalancerConfig", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteNodeBalancerConfig indicates an expected call of DeleteNodeBalancerConfig.
func (mr *MockClientMockRecorder) DeleteNodeBalancerConfig(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteNodeBalancerConfig", reflect.TypeOf((*MockClient)(nil).DeleteNodeBalancerConfig), arg0, arg1, arg2)
}

// GetFirewall mocks base method.
func (m *MockClient) GetFirewall(arg0 context.Context, arg1 int) (*linodego.Firewall, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFirewall", arg0, arg1)
	ret0, _ := ret[0].(*linodego.Firewall)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFirewall indicates an expected call of GetFirewall.
func (mr *MockClientMockRecorder) GetFirewall(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFirewall", reflect.TypeOf((*MockClient)(nil).GetFirewall), arg0, arg1)
}

// GetInstance mocks base method.
func (m *MockClient) GetInstance(arg0 context.Context, arg1 int) (*linodego.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstance", arg0, arg1)
	ret0, _ := ret[0].(*linodego.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstance indicates an expected call of GetInstance.
func (mr *MockClientMockRecorder) GetInstance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstance", reflect.TypeOf((*MockClient)(nil).GetInstance), arg0, arg1)
}

// GetInstanceIPAddresses mocks base method.
func (m *MockClient) GetInstanceIPAddresses(arg0 context.Context, arg1 int) (*linodego.InstanceIPAddressResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInstanceIPAddresses", arg0, arg1)
	ret0, _ := ret[0].(*linodego.InstanceIPAddressResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInstanceIPAddresses indicates an expected call of GetInstanceIPAddresses.
func (mr *MockClientMockRecorder) GetInstanceIPAddresses(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInstanceIPAddresses", reflect.TypeOf((*MockClient)(nil).GetInstanceIPAddresses), arg0, arg1)
}

// GetNodeBalancer mocks base method.
func (m *MockClient) GetNodeBalancer(arg0 context.Context, arg1 int) (*linodego.NodeBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeBalancer", arg0, arg1)
	ret0, _ := ret[0].(*linodego.NodeBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeBalancer indicates an expected call of GetNodeBalancer.
func (mr *MockClientMockRecorder) GetNodeBalancer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeBalancer", reflect.TypeOf((*MockClient)(nil).GetNodeBalancer), arg0, arg1)
}

// ListFirewallDevices mocks base method.
func (m *MockClient) ListFirewallDevices(arg0 context.Context, arg1 int, arg2 *linodego.ListOptions) ([]linodego.FirewallDevice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFirewallDevices", arg0, arg1, arg2)
	ret0, _ := ret[0].([]linodego.FirewallDevice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFirewallDevices indicates an expected call of ListFirewallDevices.
func (mr *MockClientMockRecorder) ListFirewallDevices(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFirewallDevices", reflect.TypeOf((*MockClient)(nil).ListFirewallDevices), arg0, arg1, arg2)
}

// ListInstanceConfigs mocks base method.
func (m *MockClient) ListInstanceConfigs(arg0 context.Context, arg1 int, arg2 *linodego.ListOptions) ([]linodego.InstanceConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListInstanceConfigs", arg0, arg1, arg2)
	ret0, _ := ret[0].([]linodego.InstanceConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListInstanceConfigs indicates an expected call of ListInstanceConfigs.
func (mr *MockClientMockRecorder) ListInstanceConfigs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListInstanceConfigs", reflect.TypeOf((*MockClient)(nil).ListInstanceConfigs), arg0, arg1, arg2)
}

// ListInstances mocks base method.
func (m *MockClient) ListInstances(arg0 context.Context, arg1 *linodego.ListOptions) ([]linodego.Instance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListInstances", arg0, arg1)
	ret0, _ := ret[0].([]linodego.Instance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListInstances indicates an expected call of ListInstances.
func (mr *MockClientMockRecorder) ListInstances(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListInstances", reflect.TypeOf((*MockClient)(nil).ListInstances), arg0, arg1)
}

// ListNodeBalancerConfigs mocks base method.
func (m *MockClient) ListNodeBalancerConfigs(arg0 context.Context, arg1 int, arg2 *linodego.ListOptions) ([]linodego.NodeBalancerConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNodeBalancerConfigs", arg0, arg1, arg2)
	ret0, _ := ret[0].([]linodego.NodeBalancerConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNodeBalancerConfigs indicates an expected call of ListNodeBalancerConfigs.
func (mr *MockClientMockRecorder) ListNodeBalancerConfigs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNodeBalancerConfigs", reflect.TypeOf((*MockClient)(nil).ListNodeBalancerConfigs), arg0, arg1, arg2)
}

// ListNodeBalancerFirewalls mocks base method.
func (m *MockClient) ListNodeBalancerFirewalls(arg0 context.Context, arg1 int, arg2 *linodego.ListOptions) ([]linodego.Firewall, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNodeBalancerFirewalls", arg0, arg1, arg2)
	ret0, _ := ret[0].([]linodego.Firewall)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNodeBalancerFirewalls indicates an expected call of ListNodeBalancerFirewalls.
func (mr *MockClientMockRecorder) ListNodeBalancerFirewalls(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNodeBalancerFirewalls", reflect.TypeOf((*MockClient)(nil).ListNodeBalancerFirewalls), arg0, arg1, arg2)
}

// ListNodeBalancers mocks base method.
func (m *MockClient) ListNodeBalancers(arg0 context.Context, arg1 *linodego.ListOptions) ([]linodego.NodeBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNodeBalancers", arg0, arg1)
	ret0, _ := ret[0].([]linodego.NodeBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNodeBalancers indicates an expected call of ListNodeBalancers.
func (mr *MockClientMockRecorder) ListNodeBalancers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNodeBalancers", reflect.TypeOf((*MockClient)(nil).ListNodeBalancers), arg0, arg1)
}

// RebuildNodeBalancerConfig mocks base method.
func (m *MockClient) RebuildNodeBalancerConfig(arg0 context.Context, arg1, arg2 int, arg3 linodego.NodeBalancerConfigRebuildOptions) (*linodego.NodeBalancerConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RebuildNodeBalancerConfig", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*linodego.NodeBalancerConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RebuildNodeBalancerConfig indicates an expected call of RebuildNodeBalancerConfig.
func (mr *MockClientMockRecorder) RebuildNodeBalancerConfig(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RebuildNodeBalancerConfig", reflect.TypeOf((*MockClient)(nil).RebuildNodeBalancerConfig), arg0, arg1, arg2, arg3)
}

// UpdateFirewallRules mocks base method.
func (m *MockClient) UpdateFirewallRules(arg0 context.Context, arg1 int, arg2 linodego.FirewallRuleSet) (*linodego.FirewallRuleSet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFirewallRules", arg0, arg1, arg2)
	ret0, _ := ret[0].(*linodego.FirewallRuleSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFirewallRules indicates an expected call of UpdateFirewallRules.
func (mr *MockClientMockRecorder) UpdateFirewallRules(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFirewallRules", reflect.TypeOf((*MockClient)(nil).UpdateFirewallRules), arg0, arg1, arg2)
}

// UpdateInstanceConfigInterface mocks base method.
func (m *MockClient) UpdateInstanceConfigInterface(arg0 context.Context, arg1, arg2, arg3 int, arg4 linodego.InstanceConfigInterfaceUpdateOptions) (*linodego.InstanceConfigInterface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateInstanceConfigInterface", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*linodego.InstanceConfigInterface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateInstanceConfigInterface indicates an expected call of UpdateInstanceConfigInterface.
func (mr *MockClientMockRecorder) UpdateInstanceConfigInterface(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateInstanceConfigInterface", reflect.TypeOf((*MockClient)(nil).UpdateInstanceConfigInterface), arg0, arg1, arg2, arg3, arg4)
}

// UpdateNodeBalancer mocks base method.
func (m *MockClient) UpdateNodeBalancer(arg0 context.Context, arg1 int, arg2 linodego.NodeBalancerUpdateOptions) (*linodego.NodeBalancer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNodeBalancer", arg0, arg1, arg2)
	ret0, _ := ret[0].(*linodego.NodeBalancer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateNodeBalancer indicates an expected call of UpdateNodeBalancer.
func (mr *MockClientMockRecorder) UpdateNodeBalancer(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNodeBalancer", reflect.TypeOf((*MockClient)(nil).UpdateNodeBalancer), arg0, arg1, arg2)
}
