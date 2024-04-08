package linode

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/linode/linodego"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	cloudprovider "k8s.io/cloud-provider"
)

func nodeWithProviderID(providerID string) *v1.Node {
	return &v1.Node{Spec: v1.NodeSpec{
		ProviderID: providerID,
	}}
}

func nodeWithName(name string) *v1.Node {
	return &v1.Node{ObjectMeta: metav1.ObjectMeta{Name: name}}
}

func TestInstanceExists(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := NewMockClient(ctrl)

	t.Run("should return false if linode does not exist (by providerID)", func(t *testing.T) {
		instances := newInstances(client)
		node := nodeWithProviderID(providerIDPrefix + "123")
		client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{}, nil)

		exists, err := instances.InstanceExists(ctx, node)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("should return true if linode exists (by provider)", func(t *testing.T) {
		instances := newInstances(client)
		node := nodeWithProviderID(providerIDPrefix + "123")
		client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{
			{
				ID:     123,
				Label:  "mock",
				Region: "us-east",
				Type:   "g6-standard-2",
			},
		}, nil)

		exists, err := instances.InstanceExists(ctx, node)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("should return true if linode exists (by name)", func(t *testing.T) {
		instances := newInstances(client)
		name := "some-name"
		node := nodeWithName(name)

		client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{
			{ID: 123, Label: name},
		}, nil)

		exists, err := instances.InstanceExists(ctx, node)
		assert.NoError(t, err)
		assert.True(t, exists)
	})
}

func TestMetadataRetrieval(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := NewMockClient(ctrl)

	t.Run("uses name over IP for finding linode", func(t *testing.T) {
		instances := newInstances(client)
		publicIP := net.ParseIP("172.234.31.123")
		privateIP := net.ParseIP("192.168.159.135")
		expectedInstance := linodego.Instance{Label: "expected-instance", ID: 12345, IPv4: []*net.IP{&publicIP, &privateIP}}
		client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{{Label: "wrong-instance", ID: 3456, IPv4: []*net.IP{&publicIP, &privateIP}}, expectedInstance}, nil)
		name := "expected-instance"
		node := nodeWithName(name)

		meta, err := instances.InstanceMetadata(ctx, node)
		assert.Nil(t, err)
		assert.Equal(t, providerIDPrefix+strconv.Itoa(expectedInstance.ID), meta.ProviderID)
	})

	t.Run("fails when linode does not exist (by provider)", func(t *testing.T) {
		instances := newInstances(client)
		id := 456302
		providerID := providerIDPrefix + strconv.Itoa(id)
		node := nodeWithProviderID(providerID)
		client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{}, nil)
		meta, err := instances.InstanceMetadata(ctx, node)

		assert.ErrorIs(t, err, cloudprovider.InstanceNotFound)
		assert.Nil(t, meta)
	})

	t.Run("should return data when linode is found (by name)", func(t *testing.T) {
		instances := newInstances(client)
		id := 123
		name := "mock-instance"
		node := nodeWithName(name)
		publicIPv4 := net.ParseIP("45.76.101.25")
		privateIPv4 := net.ParseIP("192.168.133.65")
		linodeType := "g6-standard-1"
		region := "us-east"
		client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{
			{ID: id, Label: name, Type: linodeType, Region: region, IPv4: []*net.IP{&publicIPv4, &privateIPv4}},
		}, nil)

		meta, err := instances.InstanceMetadata(ctx, node)
		assert.NoError(t, err)
		assert.Equal(t, providerIDPrefix+strconv.Itoa(id), meta.ProviderID)
		assert.Equal(t, region, meta.Region)
		assert.Equal(t, linodeType, meta.InstanceType)
		assert.Equal(t, meta.NodeAddresses, []v1.NodeAddress{
			{
				Type:    v1.NodeHostName,
				Address: name,
			},
			{
				Type:    v1.NodeExternalIP,
				Address: publicIPv4.String(),
			},
			{
				Type:    v1.NodeInternalIP,
				Address: privateIPv4.String(),
			},
		})
	})

	ipTests := []struct {
		name            string
		inputIPs        []string
		outputAddresses []v1.NodeAddress
		expectedErr     error
	}{
		{"no IPs", nil, nil, instanceNoIPAddressesError{192910}},
		{
			"one public, one private",
			[]string{"32.74.121.25", "192.168.121.42"},
			[]v1.NodeAddress{{Type: v1.NodeExternalIP, Address: "32.74.121.25"}, {Type: v1.NodeInternalIP, Address: "192.168.121.42"}},
			nil,
		},
		{
			"one public, no private",
			[]string{"32.74.121.25"},
			[]v1.NodeAddress{{Type: v1.NodeExternalIP, Address: "32.74.121.25"}},
			nil,
		},
		{
			"one private, no public",
			[]string{"192.168.121.42"},
			[]v1.NodeAddress{{Type: v1.NodeInternalIP, Address: "192.168.121.42"}},
			nil,
		},
		{
			"two public addresses",
			[]string{"32.74.121.25", "32.74.121.22"},
			[]v1.NodeAddress{{Type: v1.NodeExternalIP, Address: "32.74.121.25"}, {Type: v1.NodeExternalIP, Address: "32.74.121.22"}},
			nil,
		},
		{
			"two private addresses",
			[]string{"192.168.121.42", "10.0.2.15"},
			[]v1.NodeAddress{{Type: v1.NodeInternalIP, Address: "192.168.121.42"}, {Type: v1.NodeInternalIP, Address: "10.0.2.15"}},
			nil,
		},
	}

	for _, test := range ipTests {
		t.Run(fmt.Sprintf("addresses are retrieved - %s", test.name), func(t *testing.T) {
			instances := newInstances(client)
			id := 192910
			name := "my-instance"
			providerID := providerIDPrefix + strconv.Itoa(id)
			node := nodeWithProviderID(providerID)

			ips := make([]*net.IP, 0, len(test.inputIPs))
			for _, ip := range test.inputIPs {
				parsed := net.ParseIP(ip)
				if parsed == nil {
					t.Fatalf("cannot parse %v as an ipv4", ip)
				}
				ips = append(ips, &parsed)
			}

			linodeType := "g6-standard-1"
			region := "us-east"
			client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{
				{ID: id, Label: name, Type: linodeType, Region: region, IPv4: ips},
			}, nil)

			meta, err := instances.InstanceMetadata(ctx, node)

			assert.Equal(t, err, test.expectedErr)
			if test.expectedErr == nil {
				assert.Equal(t, region, meta.Region)
				assert.Equal(t, linodeType, meta.InstanceType)
				addresses := append([]v1.NodeAddress{
					{Type: v1.NodeHostName, Address: name},
				}, test.outputAddresses...)
				assert.Equal(t, meta.NodeAddresses, addresses)
			}
		})

		getByIPTests := []struct {
			name          string
			nodeAddresses []v1.NodeAddress
			expectedErr   error
		}{
			{name: "gets linode by External IP", nodeAddresses: []v1.NodeAddress{{
				Type:    "ExternalIP",
				Address: "172.234.31.123",
			}, {
				Type:    "InternalIP",
				Address: "192.168.159.135",
			}}},
			{
				name: "returns error on node with only internal IP", nodeAddresses: []v1.NodeAddress{{
					Type:    "ExternalIP",
					Address: "123.2.1.23",
				}, {
					Type:    "InternalIP",
					Address: "192.168.159.135",
				}},
				expectedErr: cloudprovider.InstanceNotFound,
			},
			{
				name: "returns error on no matching nodes by IP", nodeAddresses: []v1.NodeAddress{{
					Type:    "ExternalIP",
					Address: "123.2.1.23",
				}, {
					Type:    "InternalIP",
					Address: "192.168.10.10",
				}},
				expectedErr: cloudprovider.InstanceNotFound,
			},
			{
				name: "returns error on no node IPs", nodeAddresses: []v1.NodeAddress{},
				expectedErr: fmt.Errorf("no IP address found on node test-node-1"),
			},
		}

		for _, test := range getByIPTests {
			t.Run(fmt.Sprintf("gets lindoe by IP - %s", test.name), func(t *testing.T) {
				instances := newInstances(client)

				publicIP := net.ParseIP("172.234.31.123")
				privateIP := net.ParseIP("192.168.159.135")
				wrongIP := net.ParseIP("1.2.3.4")
				expectedInstance := linodego.Instance{Label: "expected-instance", ID: 12345, IPv4: []*net.IP{&publicIP, &privateIP}}
				client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{{ID: 3456, IPv4: []*net.IP{&wrongIP}}, expectedInstance}, nil)
				node := v1.Node{ObjectMeta: metav1.ObjectMeta{Name: "test-node-1"}, Status: v1.NodeStatus{Addresses: test.nodeAddresses}}
				meta, err := instances.InstanceMetadata(ctx, &node)
				if test.expectedErr != nil {
					assert.Nil(t, meta)
					assert.Equal(t, err, test.expectedErr)
				} else {
					assert.Nil(t, err)
					assert.Equal(t, providerIDPrefix+strconv.Itoa(expectedInstance.ID), meta.ProviderID)
				}
			})
		}
	}
}

func TestMalformedProviders(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := NewMockClient(ctrl)

	t.Run("fails on non-numeric providerID", func(t *testing.T) {
		instances := newInstances(client)
		providerID := providerIDPrefix + "abc"
		node := nodeWithProviderID(providerID)
		client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{}, nil)
		meta, err := instances.InstanceMetadata(ctx, node)

		assert.ErrorIs(t, err, invalidProviderIDError{providerID})
		assert.Nil(t, meta)
	})
}

func TestInstanceShutdown(t *testing.T) {
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := NewMockClient(ctrl)

	t.Run("fails when instance not found (by provider)", func(t *testing.T) {
		instances := newInstances(client)
		id := 12345
		node := nodeWithProviderID(providerIDPrefix + strconv.Itoa(id))
		client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{}, nil)
		shutdown, err := instances.InstanceShutdown(ctx, node)

		assert.Error(t, err)
		assert.False(t, shutdown)
	})

	t.Run("fails when instance not found (by name)", func(t *testing.T) {
		instances := newInstances(client)
		name := "some-name"
		node := nodeWithName(name)
		client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{}, nil)
		shutdown, err := instances.InstanceShutdown(ctx, node)

		assert.Error(t, err)
		assert.False(t, shutdown)
	})

	t.Run("returns true when instance is shut down", func(t *testing.T) {
		instances := newInstances(client)
		id := 12345
		node := nodeWithProviderID(providerIDPrefix + strconv.Itoa(id))
		client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{
			{ID: id, Label: "offline-linode", Status: linodego.InstanceOffline},
		}, nil)
		shutdown, err := instances.InstanceShutdown(ctx, node)

		assert.NoError(t, err)
		assert.True(t, shutdown)
	})

	t.Run("returns true when instance is shutting down", func(t *testing.T) {
		instances := newInstances(client)
		id := 12345
		node := nodeWithProviderID(providerIDPrefix + strconv.Itoa(id))
		client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{
			{ID: id, Label: "shutting-down-linode", Status: linodego.InstanceShuttingDown},
		}, nil)
		shutdown, err := instances.InstanceShutdown(ctx, node)

		assert.NoError(t, err)
		assert.True(t, shutdown)
	})

	t.Run("returns false when instance is running", func(t *testing.T) {
		instances := newInstances(client)
		id := 12345
		node := nodeWithProviderID(providerIDPrefix + strconv.Itoa(id))
		client.EXPECT().ListInstances(gomock.Any(), nil).Times(1).Return([]linodego.Instance{
			{ID: id, Label: "running-linode", Status: linodego.InstanceRunning},
		}, nil)
		shutdown, err := instances.InstanceShutdown(ctx, node)

		assert.NoError(t, err)
		assert.False(t, shutdown)
	})
}
