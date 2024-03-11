package linode

import (
	"context"
	cryptoRand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/linode/linodego"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/linode/linode-cloud-controller-manager/cloud/annotations"
	"github.com/linode/linode-cloud-controller-manager/cloud/linode/firewall"
)

const testCert string = `-----BEGIN CERTIFICATE-----
MIIFITCCAwkCAWQwDQYJKoZIhvcNAQELBQAwUjELMAkGA1UEBhMCQVUxEzARBgNV
BAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0
ZDELMAkGA1UEAwwCY2EwHhcNMTkwNDA5MDkzNjQyWhcNMjMwNDA4MDkzNjQyWjBb
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMRQwEgYDVQQDDAtsaW5vZGUudGVzdDCCAiIw
DQYJKoZIhvcNAQEBBQADggIPADCCAgoCggIBANUC0KStr84PLnM1dTYuEtk4HOTc
ufb6pMHyttJv5oYxCAJaN5AI9QXPqJpUFI6GlS1oDpjRe9RQghXso/IihD9eoEP1
zkHcHJyb6TXThofatxX5jLUM9TgmTIrYH+1KyKraBO6iMz2UQkbJq04BZWI9wADq
ffn1Cw6RueDe4QdqXpv/M9d/PetsIQLjjNAFHo87gYIkw838DMyTNikIweg8tRSS
6hivBVLLF0WB7p4ZARic8t+VqEFz0xl9AANE3OYMcsZCYacHxMBnX/OpHgEMxVkZ
GZ/5ikb6HJNnK/OintBlTqmGJK77fwSYXeO/5Zn6HpakfsNf6ZWSXsWRaatRvwL7
RD45RqSUpx0GALhxXTlQWv4F0cEn5MJSZX9uTJbFTuTYqC5NrB/M33hcUWy5N/L8
fz8GOxLRmrAthZ//dW4GBASOHdwMJOPz0Hb7DwNP5tSi74o7k+vCNuAHW8c8KCno
EIOS5Z6VNc252KVWZ0Y7gz7/w1Jk+cepNmpTRWzQAWc1RRYgRvAfKwXCFZpE5y6T
iu9LYtH0eKp55MBdWJ44lBu2iXc/rzcWNo0jDeHkBevS0prBxIgH377WVq/GoPRW
g3uVC6nGczHEGq1j1u6q3JKU97JSVznXIJssZLCQ4NYxtuZtmqcfEUDictq1W2Lh
upOn8Y/XQtI8gdb1AgMBAAEwDQYJKoZIhvcNAQELBQADggIBAB1Se+wlSOsRlII3
zk5VYSwiuvWc3pBYHShbSjdOFo4StZ4MRFyKu+gBssNZ7ZyM5B1oDOjslwm31nWP
j5NnlCeSeTJ2LGIkn1AFsZ4LK/ffHnxRVSUZCTUdW9PLbwDf7oDUxdtfrLdsC39F
RBn22oXTto4SNAqNQJGSkPrVT5a23JSplsPWu8ZwruaslvCtC8MRwpUp+A8EKdau
8BeYgzJWY/QkJom159//crgvt4tDZA0ekByS/SOZ4YtIFckm5XMo7ToQCkoNNu6Y
JYfNBi9ryQMEiS0yUNghhJHxCMQp4cHISrftlPAsyv1yvf69FSoy2+RFa+KIyohK
7m6oCwCYl7I43em10kle3j8rNABEU2RCin2G92PKuweUYyabsOV8sgJpCn+r5tDJ
bIRgmSWyodP4tiu6xn1zfcK2aAQYl8PhoWIY9aSmFPKIPuxTkWu/dyNhZ2R0Ii/3
+2wU9j4bLc4ZrMROYAiQ5++EUaLIQRSVuuvJqGlfdUffJF7c6rjXHLyTKCmo079B
pCLzKBQTXQmeIWJue3/GcA8RLzcGtaTtQTJcAwNZp4V6exA869uDwFzbZA/z9jHJ
mmccdLY3hP1Ozwikm5Pecysk+bdx9rbzHbA6xLz8fp5oJYUbyyaqnWLdTZvubpur
2/6vm/KHkJHqFcF/LtIxgaZFnGYR
-----END CERTIFICATE-----`

const testKey string = `-----BEGIN RSA PRIVATE KEY-----
MIIJKAIBAAKCAgEA1QLQpK2vzg8uczV1Ni4S2Tgc5Ny59vqkwfK20m/mhjEIAlo3
kAj1Bc+omlQUjoaVLWgOmNF71FCCFeyj8iKEP16gQ/XOQdwcnJvpNdOGh9q3FfmM
tQz1OCZMitgf7UrIqtoE7qIzPZRCRsmrTgFlYj3AAOp9+fULDpG54N7hB2pem/8z
138962whAuOM0AUejzuBgiTDzfwMzJM2KQjB6Dy1FJLqGK8FUssXRYHunhkBGJzy
35WoQXPTGX0AA0Tc5gxyxkJhpwfEwGdf86keAQzFWRkZn/mKRvock2cr86Ke0GVO
qYYkrvt/BJhd47/lmfoelqR+w1/plZJexZFpq1G/AvtEPjlGpJSnHQYAuHFdOVBa
/gXRwSfkwlJlf25MlsVO5NioLk2sH8zfeFxRbLk38vx/PwY7EtGasC2Fn/91bgYE
BI4d3Awk4/PQdvsPA0/m1KLvijuT68I24AdbxzwoKegQg5LlnpU1zbnYpVZnRjuD
Pv/DUmT5x6k2alNFbNABZzVFFiBG8B8rBcIVmkTnLpOK70ti0fR4qnnkwF1YnjiU
G7aJdz+vNxY2jSMN4eQF69LSmsHEiAffvtZWr8ag9FaDe5ULqcZzMcQarWPW7qrc
kpT3slJXOdcgmyxksJDg1jG25m2apx8RQOJy2rVbYuG6k6fxj9dC0jyB1vUCAwEA
AQKCAgAJEXOcbyB63z6U/QOeaNu4j6D7RUJNd2IoN5L85nKj59Z1cy3GXftAYhTF
bSrq3mPfaPymGNTytvKyyD46gqmqoPalrgM33o0BRcnp1rV1dyQwNU1+L60I1OiR
SJ4jVfmw/FMVbaZMytD/fnpiecC9K+/Omiz+xSXRWvbU0eg2jpq0fWrRk8MpEJNf
Mhy+hllEs73Rsor7a+2HkATQPmUy49K5q393yYuqeKbm+J8V7+6SA6x7RD3De5DT
FvU3LmlRCdqhAhZyK+x+XGhDUUHLvaVxI5Zprw/p8Z/hzpSabKPiL03n/aP2JxLD
OVFV7sdxhKpks2AKJT0mdvK96nDbHFSn6cWvcwI9vprtfp3L+hk1OcYCpnjgphZf
Br6jTxIGOVVgzWGJQv89h17j1zYTY/VX0RZD+wSfewvjzm1lBdUWIZKvi5nhsoqd
4qjIeJnpBOVE0G4rY7hWlzPYk/JAPaXnD1Vj1u37CgodRGGWQjqtcoEPPQNI8HTU
wPPPJBrW9bSCywjupBPOZz+1gmwRKbyQgBGLQPJqn1BB3LsNpPervUa9udoTrelA
+c36EBlo9eAt5h2U11Q9yuLsyoUFWkndRWdHpJKPwt5tVOVQd8nnVZFGHvZhCt7M
XGy1jKL3CWpQavAtuSoX7YChQnQYM7TWTI/RtMdD62m8bbhgCQKCAQEA+YI8UvFm
6AZ4om8c3IwfGBeDpug4d2Dl1Uvmp5Zzaexp6UMKE8OgxFeyw5THjtjco6+IfDbm
lyxvUoDMxIWdBl8IuYpNZw5b8eW2SACTda7Sc8DeAuGg2VQcVYXUFzsUJiKhZLwc
CVfVVDoaMOC5T9M9cr/0dQ/AGk+dkdhx/IDRMSISNfZPwxEQvh43tciqpnme+eIg
CVqa+vfyUU4OC2kNpJj9m2bePkncRKUog+3exv+D4CPECXXF1a5qwFToXv6JiK3q
AlDPoVHz/MtZBw6PYiJau9gOV54bT+xdWSII4MO62bsvDM0GUppIMVpc3CgmDRcm
gnC/BIwcAvIBPwKCAQEA2o1/yEqniln6UfNbl8/AFFisZW9t+gXEHI0C1iYG588U
4NqpJqyFx62QlOgIgyfyE6Fk9M42LsW9CPoP+X9rdmqhnSVhbQgKbqI8ayeBCABu
oTbfh72MuFd0cco1P1Q/2XMGeQMAMMASSjyLe9xWHOGBnE5q1VfRz4yCA37+Zxo1
55eIbCfmYtu5S5GZLzTvFhpodDgC9qOBgWenXkYZor6AhopZU33Yr3a1Anp3VTfF
hMneGl6OVRyOhorphCG4yYS6hAL71ylLyqQRP0SPiSic/ipfdxT/Egs4Sov2f7cI
Lj8Sa5B7+vh4R4zsTAoeErpNZuMUo3y24rX+BzSmywKCAQB+BS6Mwgq01FfnyvEr
38XwuCexjIbAnPtYoQ5txMqkTFkuDMMxOlSf9p9+s02bs6K1NfpcqqoK3tGXPSCv
fcDSr/tLIzR3AcSkx94qPcg830DCYD6B/A3u1tG8zGxUE23Y2RLlOzF58pf4A6So
3UgbrljR9Wv2GC9x2pZ+THE+FJ4UD95czPx6TMtFCyQeN60hijomgfSmZNH0Qnls
YV0snDHc2bz12Z4Und+X+EcfY2xq3DFyav4fvRFgHMkkPX5kRHGYzCZuZvyHwUnX
e6mKq+r1qN5lE/oifOPUmVCIrW0IgTOFt0pLT96KqAwgiUBvngOiBvhXV7TTCiU3
w52nAoIBABie7jFLL7qnTkrjJoNgtRvVrX4z4mjTM3ef7xze5dJBgvGd0IZ50wxe
ojYUOblEy8GoYe4uOO5l+ljDiv8pepq5goFoj6QvzrUN886Cgce7/LqOqvnowayW
tZiIFh2PSS4fBjClxOS5DpZsYa5PcSgJw4cvUlu8a/d8tbzdFp3Y1w/DA2xjxlGG
vUYlHeOyi+iqiu/ky3irjNBeM/2r2gF6gpIljdCZEcsajWO9Fip0gPznnOzNkC1I
bUn85jercNzK5hQvHd3sWgx3FTZSa/UgrSb48Q5CQEXxG6NSRy+2F+bV1iZl/YGV
cj9lQc2DKkYj1MptdIrCZvv9UqPPK6cCggEBAO3uGtkCjbhiy2hZsfIybRBVk+Oz
/ViSe9xRTMO5UQYn7TXGUk5GwMIoBUSwujiLBPwPoAAlh26rZtnOfblLS74siBZu
sagVhoN02tqN5sM/AhUEVieGNb/WQjgeyd2bL8yIs9vyjH4IYZkljizp5+VLbEcR
o/aoxqmE0mN1lyCPOa9UP//LlsREkWVKI3+Wld/xERtzf66hjcH+ilsXDxxpMEXo
+jczfFY/ivf7HxfhyYAMMUT50XaQuN82ZcSdZt8fNwWL86sLtKQ3wugk9qsQG+6/
bSiPJQsGIKtQvyCaZY2szyOoeUGgOId+He7ITlezxKrjdj+1pLMESvAxKeo=
-----END RSA PRIVATE KEY-----`

func TestCCMLoadBalancers(t *testing.T) {
	testCases := []struct {
		name string
		f    func(*testing.T, *linodego.Client, *fakeAPI)
	}{
		{
			name: "Get Load Balancer",
			f:    testGetLoadBalancer,
		},
		{
			name: "Create Load Balancer Without Firewall",
			f:    testCreateNodeBalancerWithOutFirewall,
		},
		{
			name: "Create Load Balancer With Valid Firewall ID",
			f:    testCreateNodeBalancerWithFirewall,
		},
		{
			name: "Create Load Balancer With Invalid Firewall ID",
			f:    testCreateNodeBalancerWithInvalidFirewall,
		},
		{
			name: "Create Load Balancer With Valid Firewall ACL - AllowList",
			f:    testCreateNodeBalancerWithAllowList,
		},
		{
			name: "Create Load Balancer With Valid Firewall ACL - DenyList",
			f:    testCreateNodeBalancerWithDenyList,
		},
		{
			name: "Create Load Balancer With Invalid Firewall ACL - Both Allow and Deny",
			f:    testCreateNodeBalanceWithBothAllowOrDenyList,
		},
		{
			name: "Create Load Balancer With Invalid Firewall ACL - NO Allow Or Deny",
			f:    testCreateNodeBalanceWithNoAllowOrDenyList,
		},
		{
			name: "Update Load Balancer - Add Annotation",
			f:    testUpdateLoadBalancerAddAnnotation,
		},
		{
			name: "Update Load Balancer - Add Port Annotation",
			f:    testUpdateLoadBalancerAddPortAnnotation,
		},
		{
			name: "Update Load Balancer - Add TLS Port",
			f:    testUpdateLoadBalancerAddTLSPort,
		},
		{
			name: "Update Load Balancer - Add Tags",
			f:    testUpdateLoadBalancerAddTags,
		},
		{
			name: "Update Load Balancer - Specify NodeBalancerID",
			f:    testUpdateLoadBalancerAddNodeBalancerID,
		},
		{
			name: "Update Load Balancer - Proxy Protocol",
			f:    testUpdateLoadBalancerAddProxyProtocol,
		},
		{
			name: "Update Load Balancer - Add new Firewall ID",
			f:    testUpdateLoadBalancerAddNewFirewall,
		},
		{
			name: "Update Load Balancer - Update Firewall ID",
			f:    testUpdateLoadBalancerUpdateFirewall,
		},
		{
			name: "Update Load Balancer - Delete Firewall ID",
			f:    testUpdateLoadBalancerDeleteFirewall,
		},
		{
			name: "Update Load Balancer - Update Firewall ACL",
			f:    testUpdateLoadBalancerUpdateFirewallACL,
		},
		{
			name: "Update Load Balancer - Remove Firewall ID & Add ACL",
			f:    testUpdateLoadBalancerUpdateFirewallRemoveIDaddACL,
		},
		{
			name: "Update Load Balancer - Remove Firewall ACL & Add ID",
			f:    testUpdateLoadBalancerUpdateFirewallRemoveACLaddID,
		},
		{
			name: "Update Load Balancer - Add a new Firewall ACL",
			f:    testUpdateLoadBalancerAddNewFirewallACL,
		},
		{
			name: "Build Load Balancer Request",
			f:    testBuildLoadBalancerRequest,
		},
		{
			name: "Ensure Load Balancer Deleted",
			f:    testEnsureLoadBalancerDeleted,
		},
		{
			name: "Ensure Load Balancer Deleted - Preserve Annotation",
			f:    testEnsureLoadBalancerPreserveAnnotation,
		},
		{
			name: "Ensure Existing Load Balancer",
			f:    testEnsureExistingLoadBalancer,
		},
		{
			name: "Ensure New Load Balancer",
			f:    testEnsureNewLoadBalancer,
		},
		{
			name: "Ensure New Load Balancer with NodeBalancerID",
			f:    testEnsureNewLoadBalancerWithNodeBalancerID,
		},
		{
			name: "getNodeBalancerForService - NodeBalancerID does not exist",
			f:    testGetNodeBalancerForServiceIDDoesNotExist,
		},
		{
			name: "makeLoadBalancerStatus",
			f:    testMakeLoadBalancerStatus,
		},
		{
			name: "makeLoadBalancerStatusEnvVar",
			f:    testMakeLoadBalancerStatusEnvVar,
		},
		{
			name: "Cleanup does not call the API unless Service annotated",
			f:    testCleanupDoesntCall,
		},
		{
			name: "Update Load Balancer - No Nodes",
			f:    testUpdateLoadBalancerNoNodes,
		},
	}

	for _, tc := range testCases {
		fake := newFake(t)
		ts := httptest.NewServer(fake)

		linodeClient := linodego.NewClient(http.DefaultClient)
		linodeClient.SetBaseURL(ts.URL)

		t.Run(tc.name, func(t *testing.T) {
			defer ts.Close()
			tc.f(t, &linodeClient, fake)
		})
	}
}

func stubService(fake *fake.Clientset, service *v1.Service) {
	_, _ = fake.CoreV1().Services("").Create(context.TODO(), service, metav1.CreateOptions{})
}

func testCreateNodeBalancer(t *testing.T, client *linodego.Client, _ *fakeAPI, annMap map[string]string) error {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: randString(),
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeThrottle:         "15",
				annotations.AnnLinodeLoadBalancerTags: "fake,test,yolo",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(8080),
					NodePort: int32(30001),
				},
			},
		},
	}
	for key, value := range annMap {
		svc.Annotations[key] = value
	}
	lb := &loadbalancers{client, "us-west", nil}
	nodes := []*v1.Node{
		{ObjectMeta: metav1.ObjectMeta{Name: "node-1"}},
	}
	nb, err := lb.buildLoadBalancerRequest(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		return err
	}

	if nb.Region != lb.zone {
		t.Error("unexpected nodebalancer region")
		t.Logf("expected: %s", lb.zone)
		t.Logf("actual: %s", nb.Region)
	}

	configs, err := client.ListNodeBalancerConfigs(context.TODO(), nb.ID, nil)
	if err != nil {
		return err
	}

	if len(configs) != len(svc.Spec.Ports) {
		t.Error("unexpected nodebalancer config count")
		t.Logf("expected: %v", len(svc.Spec.Ports))
		t.Logf("actual: %v", len(configs))
	}

	nb, err = client.GetNodeBalancer(context.TODO(), nb.ID)
	if err != nil {
		return err
	}

	if nb.ClientConnThrottle != 15 {
		t.Error("unexpected ClientConnThrottle")
		t.Logf("expected: %v", 15)
		t.Logf("actual: %v", nb.ClientConnThrottle)
	}

	expectedTags := []string{"linodelb", "fake", "test", "yolo"}
	if !reflect.DeepEqual(nb.Tags, expectedTags) {
		t.Error("unexpected Tags")
		t.Logf("expected: %v", expectedTags)
		t.Logf("actual: %v", nb.Tags)
	}

	_, ok := annMap[annotations.AnnLinodeCloudFirewallACL]
	if ok {
		// a firewall was configured for this
		firewalls, err := client.ListNodeBalancerFirewalls(context.TODO(), nb.ID, &linodego.ListOptions{})
		if err != nil {
			t.Errorf("Expected nil error, got %v", err)
		}

		if len(firewalls) == 0 {
			t.Errorf("Expected 1 firewall, got %d", len(firewalls))
		}
	}

	defer func() { _ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc) }()
	return nil
}

func testCreateNodeBalancerWithOutFirewall(t *testing.T, client *linodego.Client, f *fakeAPI) {
	err := testCreateNodeBalancer(t, client, f, nil)
	if err != nil {
		t.Fatalf("expected a nil error, got %v", err)
	}
}

func testCreateNodeBalanceWithNoAllowOrDenyList(t *testing.T, client *linodego.Client, f *fakeAPI) {
	annotations := map[string]string{
		annotations.AnnLinodeCloudFirewallACL: `{}`,
	}

	err := testCreateNodeBalancer(t, client, f, annotations)
	if err == nil || !stderrors.Is(err, firewall.ErrInvalidFWConfig) {
		t.Fatalf("expected a %v error, got %v", firewall.ErrInvalidFWConfig, err)
	}
}

func testCreateNodeBalanceWithBothAllowOrDenyList(t *testing.T, client *linodego.Client, f *fakeAPI) {
	annotations := map[string]string{
		annotations.AnnLinodeCloudFirewallACL: `{
			"allowList": {
				"ipv4": ["2.2.2.2"]
			},
			"denyList": {
				"ipv4": ["2.2.2.2"]
			}
		}`,
	}

	err := testCreateNodeBalancer(t, client, f, annotations)
	if err == nil || !stderrors.Is(err, firewall.ErrInvalidFWConfig) {
		t.Fatalf("expected a %v error, got %v", firewall.ErrInvalidFWConfig, err)
	}
}

func testCreateNodeBalancerWithAllowList(t *testing.T, client *linodego.Client, f *fakeAPI) {
	annotations := map[string]string{
		annotations.AnnLinodeCloudFirewallACL: `{
			"allowList": {
				"ipv4": ["2.2.2.2"]
			}
		}`,
	}

	err := testCreateNodeBalancer(t, client, f, annotations)
	if err != nil {
		t.Fatalf("expected a non-nil error, got %v", err)
	}
}

func testCreateNodeBalancerWithDenyList(t *testing.T, client *linodego.Client, f *fakeAPI) {
	annotations := map[string]string{
		annotations.AnnLinodeCloudFirewallACL: `{
			"denyList": {
				"ipv4": ["2.2.2.2"]
			}
		}`,
	}

	err := testCreateNodeBalancer(t, client, f, annotations)
	if err != nil {
		t.Fatalf("expected a non-nil error, got %v", err)
	}
}

func testCreateNodeBalancerWithFirewall(t *testing.T, client *linodego.Client, f *fakeAPI) {
	annotations := map[string]string{
		annotations.AnnLinodeCloudFirewallID: "123",
	}
	err := testCreateNodeBalancer(t, client, f, annotations)
	if err != nil {
		t.Fatalf("expected a nil error, got %v", err)
	}
}

func testCreateNodeBalancerWithInvalidFirewall(t *testing.T, client *linodego.Client, f *fakeAPI) {
	annotations := map[string]string{
		annotations.AnnLinodeCloudFirewallID: "qwerty",
	}
	expectedError := "strconv.Atoi: parsing \"qwerty\": invalid syntax"
	err := testCreateNodeBalancer(t, client, f, annotations)
	if err.Error() != expectedError {
		t.Fatalf("expected a %s error, got %v", expectedError, err)
	}
}

func testUpdateLoadBalancerAddAnnotation(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: randString(),
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeThrottle: "15",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset

	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
	}()

	lbStatus, err := lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("EnsureLoadBalancer returned an error: %s", err)
	}
	svc.Status.LoadBalancer = *lbStatus

	stubService(fakeClientset, svc)
	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeThrottle: "10",
	})

	err = lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("UpdateLoadBalancer returned an error while updated annotations: %s", err)
	}

	nb, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	if nb.ClientConnThrottle != 10 {
		t.Errorf("unexpected ClientConnThrottle: expected %d, got %d", 10, nb.ClientConnThrottle)
		t.Logf("expected: %v", 10)
		t.Logf("actual: %v", nb.ClientConnThrottle)
	}
}

func testUpdateLoadBalancerAddPortAnnotation(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	targetTestPort := 80
	portConfigAnnotation := fmt.Sprintf("%s%d", annotations.AnnLinodePortConfigPrefix, targetTestPort)
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        randString(),
			UID:         "foobar123",
			Annotations: map[string]string{},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset

	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
	}()

	lbStatus, err := lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("EnsureLoadBalancer returned an error: %s", err)
	}
	svc.Status.LoadBalancer = *lbStatus
	stubService(fakeClientset, svc)

	svc.ObjectMeta.SetAnnotations(map[string]string{
		portConfigAnnotation: `{"protocol": "http"}`,
	})

	err = lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Fatalf("UpdateLoadBalancer returned an error while updated annotations: %s", err)
	}

	nb, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer by status: %v", err)
	}

	cfgs, errConfigs := client.ListNodeBalancerConfigs(context.TODO(), nb.ID, nil)
	if errConfigs != nil {
		t.Fatalf("error getting NodeBalancer configs: %v", errConfigs)
	}

	expectedPortConfigs := map[int]string{
		80: "http",
	}
	observedPortConfigs := make(map[int]string)

	for _, cfg := range cfgs {
		observedPortConfigs[cfg.Port] = string(cfg.Protocol)
	}

	if !reflect.DeepEqual(expectedPortConfigs, observedPortConfigs) {
		t.Errorf("NodeBalancer port mismatch: expected %v, got %v", expectedPortConfigs, observedPortConfigs)
	}
}

func testUpdateLoadBalancerAddTags(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        randString(),
			UID:         "foobar123",
			Annotations: map[string]string{},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset
	clusterName := "linodelb"

	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), clusterName, svc)
	}()

	lbStatus, err := lb.EnsureLoadBalancer(context.TODO(), clusterName, svc, nodes)
	if err != nil {
		t.Errorf("EnsureLoadBalancer returned an error: %s", err)
	}
	svc.Status.LoadBalancer = *lbStatus
	stubService(fakeClientset, svc)

	testTags := "test,new,tags"
	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeLoadBalancerTags: testTags,
	})

	err = lb.UpdateLoadBalancer(context.TODO(), clusterName, svc, nodes)
	if err != nil {
		t.Fatalf("UpdateLoadBalancer returned an error while updated annotations: %s", err)
	}

	nb, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer by status: %v", err)
	}

	expectedTags := append([]string{clusterName}, strings.Split(testTags, ",")...)
	observedTags := nb.Tags

	if !reflect.DeepEqual(expectedTags, observedTags) {
		t.Errorf("NodeBalancer tags mismatch: expected %v, got %v", expectedTags, observedTags)
	}
}

func testUpdateLoadBalancerAddTLSPort(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: randString(),
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeThrottle: "15",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	extraPort := v1.ServicePort{
		Name:     randString(),
		Protocol: "TCP",
		Port:     int32(443),
		NodePort: int32(30001),
	}

	lb := &loadbalancers{client, "us-west", nil}

	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
	}()

	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset
	addTLSSecret(t, lb.kubeClient)

	lbStatus, err := lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("EnsureLoadBalancer returned an error: %s", err)
	}
	svc.Status.LoadBalancer = *lbStatus

	stubService(fakeClientset, svc)
	svc.Spec.Ports = append(svc.Spec.Ports, extraPort)
	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodePortConfigPrefix + "443": `{ "protocol": "https", "tls-secret-name": "tls-secret"}`,
	})
	err = lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Fatalf("UpdateLoadBalancer returned an error while updated annotations: %s", err)
	}

	nb, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cfgs, errConfigs := client.ListNodeBalancerConfigs(context.TODO(), nb.ID, nil)
	if errConfigs != nil {
		t.Fatalf("error getting NodeBalancer configs: %v", errConfigs)
	}

	expectedPorts := map[int]struct{}{
		80:  {},
		443: {},
	}

	observedPorts := make(map[int]struct{})

	for _, cfg := range cfgs {
		nodes, errNodes := client.ListNodeBalancerNodes(context.TODO(), nb.ID, cfg.ID, nil)
		if errNodes != nil {
			t.Errorf("error getting NodeBalancer nodes: %v", errNodes)
		}

		if len(nodes) == 0 {
			t.Errorf("no nodes found for port %d", cfg.Port)
		}

		observedPorts[cfg.Port] = struct{}{}
	}

	if !reflect.DeepEqual(expectedPorts, observedPorts) {
		t.Errorf("NodeBalancer ports mismatch: expected %v, got %v", expectedPorts, observedPorts)
	}
}

func testUpdateLoadBalancerAddProxyProtocol(t *testing.T, client *linodego.Client, fakeAPI *fakeAPI) {
	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset

	for _, tc := range []struct {
		name                string
		proxyProtocolConfig linodego.ConfigProxyProtocol
		invalidErr          bool
	}{
		{
			name:                "with invalid Proxy Protocol",
			proxyProtocolConfig: "bogus",
			invalidErr:          true,
		},
		{
			name:                "with none",
			proxyProtocolConfig: linodego.ProxyProtocolNone,
		},
		{
			name:                "with v1",
			proxyProtocolConfig: linodego.ProxyProtocolV1,
		},
		{
			name:                "with v2",
			proxyProtocolConfig: linodego.ProxyProtocolV2,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			svc := &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:        randString(),
					UID:         "foobar123",
					Annotations: map[string]string{},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Name:     randString(),
							Protocol: "tcp",
							Port:     int32(80),
							NodePort: int32(8080),
						},
					},
				},
			}

			defer func() {
				_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
			}()
			nodeBalancer, err := client.CreateNodeBalancer(context.TODO(), linodego.NodeBalancerCreateOptions{
				Region: lb.zone,
			})
			if err != nil {
				t.Fatalf("failed to create NodeBalancer: %s", err)
			}

			svc.Status.LoadBalancer = *makeLoadBalancerStatus(svc, nodeBalancer)
			svc.ObjectMeta.SetAnnotations(map[string]string{
				annotations.AnnLinodeDefaultProxyProtocol: string(tc.proxyProtocolConfig),
			})

			stubService(fakeClientset, svc)
			if err = lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes); err != nil {
				expectedErrMessage := fmt.Sprintf("invalid NodeBalancer proxy protocol value '%s'", tc.proxyProtocolConfig)
				if tc.invalidErr && err.Error() == expectedErrMessage {
					return
				}
				t.Fatalf("UpdateLoadBalancer returned an unexpected error while updated annotations: %s", err)
				return
			}
			if tc.invalidErr {
				t.Fatal("expected UpdateLoadBalancer to return an error")
			}

			nodeBalancerConfigs, err := client.ListNodeBalancerConfigs(context.TODO(), nodeBalancer.ID, nil)
			if err != nil {
				t.Fatalf("failed to get NodeBalancer: %s", err)
			}

			for _, config := range nodeBalancerConfigs {
				proxyProtocol := config.ProxyProtocol
				if proxyProtocol != tc.proxyProtocolConfig {
					t.Errorf("expected ProxyProtocol to be %s; got %s", tc.proxyProtocolConfig, proxyProtocol)
				}
			}
		})
	}
}

func testUpdateLoadBalancerAddNewFirewall(t *testing.T, client *linodego.Client, fakeAPI *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: randString(),
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeThrottle: "15",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset

	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
	}()

	lbStatus, err := lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("EnsureLoadBalancer returned an error: %s", err)
	}
	svc.Status.LoadBalancer = *lbStatus
	stubService(fakeClientset, svc)
	fwClient := firewall.NewFirewalls(client)
	fw, err := fwClient.CreateFirewall(context.TODO(), linodego.FirewallCreateOptions{
		Label: "test",
		Rules: linodego.FirewallRuleSet{Inbound: []linodego.FirewallRule{{
			Action:      "ACCEPT",
			Label:       "inbound-rule123",
			Description: "inbound rule123",
			Ports:       "4321",
			Protocol:    linodego.TCP,
			Addresses: linodego.NetworkAddresses{
				IPv4: &[]string{"0.0.0.0/0"},
			},
		}}, Outbound: []linodego.FirewallRule{}, InboundPolicy: "ACCEPT", OutboundPolicy: "ACCEPT"},
	})
	if err != nil {
		t.Errorf("CreatingFirewall returned an error: %s", err)
	}
	defer func() {
		_ = fwClient.DeleteFirewall(context.TODO(), fw)
	}()

	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeCloudFirewallID: strconv.Itoa(fw.ID),
	})

	err = lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("UpdateLoadBalancer returned an error while updated annotations: %s", err)
	}

	nb, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	firewalls, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nb.ID, &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("failed to List Firewalls %s", err)
	}

	if len(firewalls) == 0 {
		t.Fatalf("No attached firewalls found")
	}

	if firewalls[0].ID != fw.ID {
		t.Fatalf("Attached firewallID not matching with created firewall")
	}
}

// This will also test the firewall with >255 IPs
func testUpdateLoadBalancerAddNewFirewallACL(t *testing.T, client *linodego.Client, fakeAPI *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: randString(),
			UID:  "foobar123",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset

	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
	}()
	lbStatus, err := lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("EnsureLoadBalancer returned an error: %s", err)
	}
	svc.Status.LoadBalancer = *lbStatus
	stubService(fakeClientset, svc)

	nb, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	firewalls, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nb.ID, &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list nodeBalancer firewalls %s", err)
	}

	if len(firewalls) != 0 {
		t.Fatalf("Firewalls attached when none specified")
	}

	var ipv4s []string
	var ipv6s []string
	i := 0
	for i < 400 {
		ipv4s = append(ipv4s, fmt.Sprintf("%d.%d.%d.%d", 192, rand.Int31n(255), rand.Int31n(255), rand.Int31n(255)))
		i += 1
	}
	i = 0
	for i < 300 {
		ip := make([]byte, 16)
		if _, err := cryptoRand.Read(ip); err != nil {
			t.Fatalf("unable to read random bytes")
		}
		ipv6s = append(ipv6s, fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s",
			hex.EncodeToString(ip[0:2]),
			hex.EncodeToString(ip[2:4]),
			hex.EncodeToString(ip[4:6]),
			hex.EncodeToString(ip[6:8]),
			hex.EncodeToString(ip[8:10]),
			hex.EncodeToString(ip[10:12]),
			hex.EncodeToString(ip[12:14]),
			hex.EncodeToString(ip[14:16])))
		i += 1
	}
	acl := map[string]map[string][]string{
		"allowList": {
			"ipv4": ipv4s,
			"ipv6": ipv6s,
		},
	}
	aclString, err := json.Marshal(acl)
	if err != nil {
		t.Fatalf("unable to marshal json acl")
	}

	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeCloudFirewallACL: string(aclString),
	})

	err = lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("UpdateLoadBalancer returned an error: %s", err)
	}

	nbUpdated, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	firewallsNew, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nbUpdated.ID, &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("failed to List Firewalls %s", err)
	}

	if len(firewallsNew) == 0 {
		t.Fatalf("No firewalls found")
	}

	if firewallsNew[0].Rules.InboundPolicy != "DROP" {
		t.Errorf("expected DROP inbound policy, got %s", firewallsNew[0].Rules.InboundPolicy)
	}

	if len(firewallsNew[0].Rules.Inbound) != 4 {
		t.Errorf("expected 4 rules, got %d", len(firewallsNew[0].Rules.Inbound))
	}
}

func testUpdateLoadBalancerUpdateFirewallRemoveACLaddID(t *testing.T, client *linodego.Client, fakeAPI *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: randString(),
			UID:  "foobar123",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset

	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeCloudFirewallACL: `{
			"allowList": {
				"ipv4": ["2.2.2.2"]
			}
		}`,
	})

	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
	}()
	lbStatus, err := lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("EnsureLoadBalancer returned an error: %s", err)
	}
	svc.Status.LoadBalancer = *lbStatus
	stubService(fakeClientset, svc)

	nb, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	firewalls, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nb.ID, &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list nodeBalancer firewalls %s", err)
	}

	if len(firewalls) == 0 {
		t.Fatalf("No firewalls attached")
	}

	if firewalls[0].Rules.InboundPolicy != "DROP" {
		t.Errorf("expected DROP inbound policy, got %s", firewalls[0].Rules.InboundPolicy)
	}

	fwIPs := firewalls[0].Rules.Inbound[0].Addresses.IPv4
	if fwIPs == nil {
		t.Errorf("expected IP, got %v", fwIPs)
	}

	fwClient := firewall.NewFirewalls(client)
	fw, err := fwClient.CreateFirewall(context.TODO(), linodego.FirewallCreateOptions{
		Label: "test",
		Rules: linodego.FirewallRuleSet{Inbound: []linodego.FirewallRule{{
			Action:      "ACCEPT",
			Label:       "inbound-rule123",
			Description: "inbound rule123",
			Ports:       "4321",
			Protocol:    linodego.TCP,
			Addresses: linodego.NetworkAddresses{
				IPv4: &[]string{"0.0.0.0/0"},
			},
		}}, Outbound: []linodego.FirewallRule{}, InboundPolicy: "ACCEPT", OutboundPolicy: "ACCEPT"},
	})
	if err != nil {
		t.Errorf("Error creating firewall %s", err)
	}
	defer func() {
		_ = fwClient.DeleteFirewall(context.TODO(), fw)
	}()

	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeCloudFirewallID: strconv.Itoa(fw.ID),
	})

	err = lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("UpdateLoadBalancer returned an error: %s", err)
	}

	nbUpdated, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	firewallsNew, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nbUpdated.ID, &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("failed to List Firewalls %s", err)
	}

	if len(firewallsNew) == 0 {
		t.Fatalf("No attached firewalls found")
	}

	if firewallsNew[0].Rules.InboundPolicy != "ACCEPT" {
		t.Errorf("expected ACCEPT inbound policy, got %s", firewallsNew[0].Rules.InboundPolicy)
	}

	fwIPs = firewallsNew[0].Rules.Inbound[0].Addresses.IPv4
	if fwIPs == nil {
		t.Errorf("expected 2.2.2.2, got %v", fwIPs)
	}

	if firewallsNew[0].ID != fw.ID {
		t.Errorf("Firewall ID does not match what we created, something wrong.")
	}
}

func testUpdateLoadBalancerUpdateFirewallRemoveIDaddACL(t *testing.T, client *linodego.Client, fakeAPI *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: randString(),
			UID:  "foobar123",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset

	fwClient := firewall.NewFirewalls(client)
	fw, err := fwClient.CreateFirewall(context.TODO(), linodego.FirewallCreateOptions{
		Label: "test",
		Rules: linodego.FirewallRuleSet{Inbound: []linodego.FirewallRule{{
			Action:      "ACCEPT",
			Label:       "inbound-rule123",
			Description: "inbound rule123",
			Ports:       "4321",
			Protocol:    linodego.TCP,
			Addresses: linodego.NetworkAddresses{
				IPv4: &[]string{"0.0.0.0/0"},
			},
		}}, Outbound: []linodego.FirewallRule{}, InboundPolicy: "ACCEPT", OutboundPolicy: "ACCEPT"},
	})
	if err != nil {
		t.Errorf("Error creating firewall %s", err)
	}
	defer func() {
		_ = fwClient.DeleteFirewall(context.TODO(), fw)
	}()

	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeCloudFirewallID: strconv.Itoa(fw.ID),
	})

	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
	}()
	lbStatus, err := lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("EnsureLoadBalancer returned an error: %s", err)
	}
	svc.Status.LoadBalancer = *lbStatus
	stubService(fakeClientset, svc)

	nb, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	firewalls, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nb.ID, &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list nodeBalancer firewalls %s", err)
	}

	if len(firewalls) == 0 {
		t.Fatalf("No firewalls attached")
	}

	if firewalls[0].Rules.InboundPolicy != "ACCEPT" {
		t.Errorf("expected ACCEPT inbound policy, got %s", firewalls[0].Rules.InboundPolicy)
	}

	fwIPs := firewalls[0].Rules.Inbound[0].Addresses.IPv4
	if fwIPs == nil {
		t.Errorf("expected IP, got %v", fwIPs)
	}
	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeCloudFirewallACL: `{
			"allowList": {
				"ipv4": ["2.2.2.2"]
			}
		}`,
	})

	err = lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("UpdateLoadBalancer returned an error: %s", err)
	}

	nbUpdated, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	firewallsNew, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nbUpdated.ID, &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("failed to List Firewalls %s", err)
	}

	if len(firewallsNew) == 0 {
		t.Fatalf("No attached firewalls found")
	}

	if firewallsNew[0].Rules.InboundPolicy != "DROP" {
		t.Errorf("expected DROP inbound policy, got %s", firewallsNew[0].Rules.InboundPolicy)
	}

	fwIPs = firewallsNew[0].Rules.Inbound[0].Addresses.IPv4
	if fwIPs == nil {
		t.Errorf("expected 2.2.2.2, got %v", fwIPs)
	}

	if firewallsNew[0].ID != fw.ID {
		t.Errorf("Firewall ID does not match, something wrong.")
	}
}

func testUpdateLoadBalancerUpdateFirewallACL(t *testing.T, client *linodego.Client, fakeAPI *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: randString(),
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeCloudFirewallACL: `{
					"allowList": {
						"ipv4": ["2.2.2.2"]
					}
				}`,
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset

	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
	}()
	lbStatus, err := lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("EnsureLoadBalancer returned an error: %s", err)
	}
	svc.Status.LoadBalancer = *lbStatus
	stubService(fakeClientset, svc)

	nb, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	firewalls, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nb.ID, &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list nodeBalancer firewalls %s", err)
	}

	if len(firewalls) == 0 {
		t.Fatalf("No firewalls attached")
	}

	if firewalls[0].Rules.InboundPolicy != "DROP" {
		t.Errorf("expected DROP inbound policy, got %s", firewalls[0].Rules.InboundPolicy)
	}

	fwIPs := firewalls[0].Rules.Inbound[0].Addresses.IPv4
	if fwIPs == nil {
		t.Errorf("expected 2.2.2.2, got %v", fwIPs)
	}

	fmt.Printf("got %v", fwIPs)

	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeCloudFirewallACL: `{
			"allowList": {
				"ipv4": ["2.2.2.2"],
				"ipv6": ["dead:beef::/128"]
			}
		}`,
	})

	err = lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("UpdateLoadBalancer returned an error: %s", err)
	}

	nbUpdated, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	firewallsNew, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nbUpdated.ID, &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("failed to List Firewalls %s", err)
	}

	if len(firewallsNew) == 0 {
		t.Fatalf("No attached firewalls found")
	}

	fwIPs = firewallsNew[0].Rules.Inbound[0].Addresses.IPv4
	if fwIPs == nil {
		t.Errorf("expected non nil IPv4, got %v", fwIPs)
	}

	if len(*fwIPs) != 1 {
		t.Errorf("expected one IPv4, got %v", fwIPs)
	}

	if firewallsNew[0].Rules.Inbound[0].Addresses.IPv6 == nil {
		t.Errorf("expected non nil IPv6, got %v", firewallsNew[0].Rules.Inbound[0].Addresses.IPv6)
	}

	if len(*firewallsNew[0].Rules.Inbound[0].Addresses.IPv6) != 1 {
		t.Errorf("expected one IPv6, got %v", firewallsNew[0].Rules.Inbound[0].Addresses.IPv6)
	}
}

func testUpdateLoadBalancerUpdateFirewall(t *testing.T, client *linodego.Client, fakeAPI *fakeAPI) {
	firewallCreateOpts := linodego.FirewallCreateOptions{
		Label: "test",
		Rules: linodego.FirewallRuleSet{Inbound: []linodego.FirewallRule{{
			Action:      "ACCEPT",
			Label:       "inbound-rule123",
			Description: "inbound rule123",
			Ports:       "4321",
			Protocol:    linodego.TCP,
			Addresses: linodego.NetworkAddresses{
				IPv4: &[]string{"0.0.0.0/0"},
			},
		}}, Outbound: []linodego.FirewallRule{}, InboundPolicy: "ACCEPT", OutboundPolicy: "ACCEPT"},
	}

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: randString(),
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeThrottle: "15",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset

	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
	}()

	fwClient := firewall.NewFirewalls(client)
	fw, err := fwClient.CreateFirewall(context.TODO(), firewallCreateOpts)
	if err != nil {
		t.Errorf("Error creating firewall %s", err)
	}
	defer func() {
		_ = fwClient.DeleteFirewall(context.TODO(), fw)
	}()

	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeCloudFirewallID: strconv.Itoa(fw.ID),
	})
	lbStatus, err := lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("EnsureLoadBalancer returned an error: %s", err)
	}
	svc.Status.LoadBalancer = *lbStatus
	stubService(fakeClientset, svc)

	nb, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	firewalls, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nb.ID, &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("Failed to list nodeBalancer firewalls %s", err)
	}

	if len(firewalls) == 0 {
		t.Fatalf("No firewalls attached")
	}

	if fw.ID != firewalls[0].ID {
		t.Fatalf("Attached firewallID not matching with created firewall")
	}

	firewallCreateOpts.Label = "test2"
	firewallNew, err := fwClient.CreateFirewall(context.TODO(), firewallCreateOpts)
	if err != nil {
		t.Fatalf("Error in creating firewall %s", err)
	}
	defer func() {
		_ = fwClient.DeleteFirewall(context.TODO(), firewallNew)
	}()

	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeCloudFirewallID: strconv.Itoa(firewallNew.ID),
	})

	err = lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("UpdateLoadBalancer returned an error: %s", err)
	}

	nbUpdated, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	firewallsNew, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nbUpdated.ID, &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("failed to List Firewalls %s", err)
	}

	if len(firewallsNew) == 0 {
		t.Fatalf("No attached firewalls found")
	}

	if firewallsNew[0].ID != firewallNew.ID {
		t.Fatalf("Attached firewallID not matching with created firewall")
	}
}

func testUpdateLoadBalancerDeleteFirewall(t *testing.T, client *linodego.Client, fakeAPI *fakeAPI) {
	firewallCreateOpts := linodego.FirewallCreateOptions{
		Label: "test",
		Rules: linodego.FirewallRuleSet{Inbound: []linodego.FirewallRule{{
			Action:      "ACCEPT",
			Label:       "inbound-rule123",
			Description: "inbound rule123",
			Ports:       "4321",
			Protocol:    linodego.TCP,
			Addresses: linodego.NetworkAddresses{
				IPv4: &[]string{"0.0.0.0/0"},
			},
		}}, Outbound: []linodego.FirewallRule{}, InboundPolicy: "ACCEPT", OutboundPolicy: "ACCEPT"},
	}

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: randString(),
			UID:  "foobar123",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset

	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
	}()

	fwClient := firewall.NewFirewalls(client)
	fw, err := fwClient.CreateFirewall(context.TODO(), firewallCreateOpts)
	if err != nil {
		t.Errorf("Error in creating firewall %s", err)
	}
	defer func() {
		_ = fwClient.DeleteFirewall(context.TODO(), fw)
	}()

	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeCloudFirewallID: strconv.Itoa(fw.ID),
	})

	lbStatus, err := lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("EnsureLoadBalancer returned an error: %s", err)
	}
	svc.Status.LoadBalancer = *lbStatus
	stubService(fakeClientset, svc)

	nb, err := lb.getNodeBalancerByStatus(context.TODO(), svc)
	if err != nil {
		t.Fatalf("failed to get NodeBalancer via status: %s", err)
	}

	firewalls, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nb.ID, &linodego.ListOptions{})
	if err != nil {
		t.Errorf("Error in listing firewalls %s", err)
	}

	if len(firewalls) == 0 {
		t.Fatalf("No firewalls attached")
	}

	if fw.ID != firewalls[0].ID {
		t.Fatalf("Attached firewallID not matching with created firewall")
	}

	svc.ObjectMeta.SetAnnotations(map[string]string{})

	err = lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("UpdateLoadBalancer returned an error: %s", err)
	}

	firewallsNew, err := lb.client.ListNodeBalancerFirewalls(context.TODO(), nb.ID, &linodego.ListOptions{})
	if err != nil {
		t.Fatalf("failed to List Firewalls %s", err)
	}

	if len(firewallsNew) != 0 {
		t.Fatalf("firewall's %d still attached", firewallsNew[0].ID)
	}
}

func testUpdateLoadBalancerAddNodeBalancerID(t *testing.T, client *linodego.Client, fakeAPI *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        randString(),
			UID:         "foobar123",
			Annotations: map[string]string{},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "http",
					Port:     int32(80),
					NodePort: int32(8080),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
	}()

	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset

	nodeBalancer, err := client.CreateNodeBalancer(context.TODO(), linodego.NodeBalancerCreateOptions{
		Region: lb.zone,
	})
	if err != nil {
		t.Fatalf("failed to create NodeBalancer: %s", err)
	}

	svc.Status.LoadBalancer = *makeLoadBalancerStatus(svc, nodeBalancer)

	newNodeBalancer, err := client.CreateNodeBalancer(context.TODO(), linodego.NodeBalancerCreateOptions{
		Region: lb.zone,
	})
	if err != nil {
		t.Fatalf("failed to create new NodeBalancer: %s", err)
	}

	stubService(fakeClientset, svc)
	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeNodeBalancerID: strconv.Itoa(newNodeBalancer.ID),
	})
	err = lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Errorf("UpdateLoadBalancer returned an error while updated annotations: %s", err)
	}

	lbStatus, _, err := lb.GetLoadBalancer(context.TODO(), svc.ClusterName, svc)
	if err != nil {
		t.Errorf("GetLoadBalancer returned an error: %s", err)
	}

	expectedLBStatus := makeLoadBalancerStatus(svc, newNodeBalancer)
	if !reflect.DeepEqual(expectedLBStatus, lbStatus) {
		t.Errorf("LoadBalancer status mismatch: expected %v, got %v", expectedLBStatus, lbStatus)
	}

	if !fakeAPI.didRequestOccur(http.MethodDelete, fmt.Sprintf("/nodebalancers/%d", nodeBalancer.ID), "") {
		t.Errorf("expected old NodeBalancer to have been deleted")
	}
}

func Test_getConnectionThrottle(t *testing.T) {
	testcases := []struct {
		name     string
		service  *v1.Service
		expected int
	}{
		{
			"throttle not specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:        randString(),
					UID:         "abc123",
					Annotations: map[string]string{},
				},
			},
			0,
		},
		{
			"throttle value is a string",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeThrottle: "foo",
					},
				},
			},
			0,
		},
		{
			"throttle value is less than 0",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeThrottle: "-123",
					},
				},
			},
			0,
		},
		{
			"throttle value is valid",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeThrottle: "1",
					},
				},
			},
			1,
		},
		{
			"throttle value is too high",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeThrottle: "21",
					},
				},
			},
			20,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			connThrottle := getConnectionThrottle(test.service)

			if test.expected != connThrottle {
				t.Fatalf("expected throttle value (%d) does not match actual value (%d)", test.expected, connThrottle)
			}
		})
	}
}

func Test_getPortConfig(t *testing.T) {
	testcases := []struct {
		name               string
		service            *v1.Service
		expectedPortConfig portConfig
		err                error
	}{
		{
			"default no proxy protocol specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
				},
			},
			portConfig{Port: 443, Protocol: "tcp", ProxyProtocol: linodego.ProxyProtocolNone},
			nil,
		},
		{
			"default proxy protocol specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeDefaultProxyProtocol: string(linodego.ProxyProtocolV2),
					},
				},
			},
			portConfig{Port: 443, Protocol: "tcp", ProxyProtocol: linodego.ProxyProtocolV2},
			nil,
		},
		{
			"port specific proxy protocol specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeDefaultProxyProtocol:     string(linodego.ProxyProtocolV2),
						annotations.AnnLinodePortConfigPrefix + "443": fmt.Sprintf(`{"proxy-protocol": "%s"}`, linodego.ProxyProtocolV1),
					},
				},
			},
			portConfig{Port: 443, Protocol: "tcp", ProxyProtocol: linodego.ProxyProtocolV1},
			nil,
		},
		{
			"default invalid proxy protocol",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeDefaultProxyProtocol: "invalid",
					},
				},
			},
			portConfig{},
			fmt.Errorf("invalid NodeBalancer proxy protocol value '%s'", "invalid"),
		},
		{
			"default no protocol specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
				},
			},
			portConfig{Port: 443, Protocol: "tcp", ProxyProtocol: linodego.ProxyProtocolNone},

			nil,
		},
		{
			"default tcp protocol specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeDefaultProtocol: "tcp",
					},
				},
			},
			portConfig{Port: 443, Protocol: "tcp", ProxyProtocol: linodego.ProxyProtocolNone},
			nil,
		},
		{
			"default capitalized protocol specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeDefaultProtocol: "HTTP",
					},
				},
			},
			portConfig{Port: 443, Protocol: "http", ProxyProtocol: linodego.ProxyProtocolNone},
			nil,
		},
		{
			"default invalid protocol",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeDefaultProtocol: "invalid",
					},
				},
			},
			portConfig{},
			fmt.Errorf("invalid protocol: %q specified", "invalid"),
		},
		{
			"port config falls back to default",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeDefaultProtocol:          "http",
						annotations.AnnLinodePortConfigPrefix + "443": `{}`,
					},
				},
			},
			portConfig{Port: 443, Protocol: "http", ProxyProtocol: linodego.ProxyProtocolNone},
			nil,
		},
		{
			"port config capitalized protocol",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodePortConfigPrefix + "443": `{ "protocol": "HTTp" }`,
					},
				},
			},
			portConfig{Port: 443, Protocol: "http", ProxyProtocol: linodego.ProxyProtocolNone},
			nil,
		},
		{
			"port config invalid protocol",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: randString(),
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodePortConfigPrefix + "443": `{ "protocol": "invalid" }`,
					},
				},
			},
			portConfig{},
			fmt.Errorf("invalid protocol: %q specified", "invalid"),
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			testPort := 443
			portConfig, err := getPortConfig(test.service, testPort)

			if !reflect.DeepEqual(portConfig, test.expectedPortConfig) {
				t.Error("unexpected port config")
				t.Logf("expected: %q", test.expectedPortConfig)
				t.Logf("actual: %q", portConfig)
			}

			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %q", test.err)
				t.Logf("actual: %q", err)
			}
		})
	}
}

func Test_getHealthCheckType(t *testing.T) {
	testcases := []struct {
		name       string
		service    *v1.Service
		healthType linodego.ConfigCheck
		err        error
	}{
		{
			"no type specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:        randString(),
					UID:         "abc123",
					Annotations: map[string]string{},
				},
			},
			linodego.CheckConnection,
			nil,
		},
		{
			"http specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeHealthCheckType: "http",
					},
				},
			},
			linodego.CheckHTTP,
			nil,
		},
		{
			"invalid specified",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test",
					UID:  "abc123",
					Annotations: map[string]string{
						annotations.AnnLinodeHealthCheckType: "invalid",
					},
				},
			},
			"",
			fmt.Errorf("invalid health check type: %q specified in annotation: %q", "invalid", annotations.AnnLinodeHealthCheckType),
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			hType, err := getHealthCheckType(test.service)
			if !reflect.DeepEqual(hType, test.healthType) {
				t.Error("unexpected health check type")
				t.Logf("expected: %v", test.healthType)
				t.Logf("actual: %v", hType)
			}

			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}

func Test_getNodePrivateIP(t *testing.T) {
	testcases := []struct {
		name    string
		node    *v1.Node
		address string
	}{
		{
			"node internal ip specified",
			&v1.Node{
				Status: v1.NodeStatus{
					Addresses: []v1.NodeAddress{
						{
							Type:    v1.NodeInternalIP,
							Address: "127.0.0.1",
						},
					},
				},
			},
			"127.0.0.1",
		},
		{
			"node internal ip not specified",
			&v1.Node{
				Status: v1.NodeStatus{
					Addresses: []v1.NodeAddress{
						{
							Type:    v1.NodeExternalIP,
							Address: "127.0.0.1",
						},
					},
				},
			},
			"",
		},
		{
			"node internal ip annotation present",
			&v1.Node{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						annotations.AnnLinodeNodePrivateIP: "192.168.42.42",
					},
				},
				Status: v1.NodeStatus{
					Addresses: []v1.NodeAddress{
						{
							Type:    v1.NodeInternalIP,
							Address: "10.0.1.1",
						},
					},
				},
			},
			"192.168.42.42",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			ip := getNodePrivateIP(test.node)
			if ip != test.address {
				t.Error("unexpected certificate")
				t.Logf("expected: %q", test.address)
				t.Logf("actual: %q", ip)
			}
		})
	}
}

func testBuildLoadBalancerRequest(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeDefaultProtocol: "tcp",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}
	nodes := []*v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-1",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-2",
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-3",
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	nb, err := lb.buildLoadBalancerRequest(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(err, err) {
		t.Error("unexpected error")
		t.Logf("expected: %v", nil)
		t.Logf("actual: %v", err)
	}

	configs, err := client.ListNodeBalancerConfigs(context.TODO(), nb.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(configs) != len(svc.Spec.Ports) {
		t.Error("unexpected nodebalancer config count")
		t.Logf("expected: %v", len(svc.Spec.Ports))
		t.Logf("actual: %v", len(configs))
	}

	nbNodes, err := client.ListNodeBalancerNodes(context.TODO(), nb.ID, configs[0].ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(nbNodes) != len(nodes) {
		t.Error("unexpected nodebalancer nodes count")
		t.Logf("expected: %v", len(nodes))
		t.Logf("actual: %v", len(nbNodes))
	}
}

func testEnsureLoadBalancerPreserveAnnotation(t *testing.T, client *linodego.Client, fake *fakeAPI) {
	testServiceSpec := v1.ServiceSpec{
		Ports: []v1.ServicePort{
			{
				Name:     "test",
				Protocol: "TCP",
				Port:     int32(80),
				NodePort: int32(30000),
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	for _, test := range []struct {
		name        string
		deleted     bool
		annotations map[string]string
	}{
		{
			name:        "load balancer preserved",
			annotations: map[string]string{annotations.AnnLinodeLoadBalancerPreserve: "true"},
			deleted:     false,
		},
		{
			name:        "load balancer not preserved (deleted)",
			annotations: map[string]string{annotations.AnnLinodeLoadBalancerPreserve: "false"},
			deleted:     true,
		},
		{
			name:        "invalid value treated as false (deleted)",
			annotations: map[string]string{annotations.AnnLinodeLoadBalancerPreserve: "bogus"},
			deleted:     true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			svc := &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "test",
					UID:         types.UID("foobar" + randString()),
					Annotations: test.annotations,
				},
				Spec: testServiceSpec,
			}

			nb, err := lb.createNodeBalancer(context.TODO(), "linodelb", svc, []*linodego.NodeBalancerConfigCreateOptions{})
			if err != nil {
				t.Fatal(err)
			}

			svc.Status.LoadBalancer = *makeLoadBalancerStatus(svc, nb)
			err = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)

			didDelete := fake.didRequestOccur(http.MethodDelete, fmt.Sprintf("/nodebalancers/%d", nb.ID), "")
			if didDelete && !test.deleted {
				t.Fatal("load balancer was unexpectedly deleted")
			} else if !didDelete && test.deleted {
				t.Fatal("load balancer was unexpectedly preserved")
			}

			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
	}
}

func testEnsureLoadBalancerDeleted(t *testing.T, client *linodego.Client, fake *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeDefaultProtocol: "tcp",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}
	testcases := []struct {
		name        string
		clusterName string
		service     *v1.Service
		err         error
	}{
		{
			"load balancer delete",
			"linodelb",
			svc,
			nil,
		},
		{
			"load balancer not exists",
			"linodelb",
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "notexists",
					UID:  "notexists123",
					Annotations: map[string]string{
						annotations.AnnLinodeDefaultProtocol: "tcp",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Name:     "test",
							Protocol: "TCP",
							Port:     int32(80),
							NodePort: int32(30000),
						},
					},
				},
			},
			nil,
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	configs := []*linodego.NodeBalancerConfigCreateOptions{}
	_, err := lb.createNodeBalancer(context.TODO(), "linodelb", svc, configs)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc) }()

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			err := lb.EnsureLoadBalancerDeleted(context.TODO(), test.clusterName, test.service)
			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}

func testEnsureExistingLoadBalancer(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testensure",
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeDefaultProtocol:           "tcp",
				annotations.AnnLinodePortConfigPrefix + "8443": `{ "protocol": "https", "tls-secret-name": "tls-secret"}`,
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(8443),
					NodePort: int32(30000),
				},
				{
					Name:     "test2",
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30001),
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	lb.kubeClient = fake.NewSimpleClientset()
	addTLSSecret(t, lb.kubeClient)

	configs := []*linodego.NodeBalancerConfigCreateOptions{}
	nb, err := lb.createNodeBalancer(context.TODO(), "linodelb", svc, configs)
	if err != nil {
		t.Fatal(err)
	}

	svc.Status.LoadBalancer = *makeLoadBalancerStatus(svc, nb)
	defer func() { _ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc) }()
	getLBStatus, exists, err := lb.GetLoadBalancer(context.TODO(), "linodelb", svc)
	if err != nil {
		t.Fatalf("failed to create nodebalancer: %s", err)
	}
	if !exists {
		t.Fatal("Node balancer not found")
	}

	testcases := []struct {
		name        string
		service     *v1.Service
		nodes       []*v1.Node
		clusterName string
		nbIP        string
		err         error
	}{
		{
			"update load balancer",
			svc,
			[]*v1.Node{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-1",
					},
					Status: v1.NodeStatus{
						Addresses: []v1.NodeAddress{
							{
								Type:    v1.NodeInternalIP,
								Address: "127.0.0.1",
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-2",
					},
					Status: v1.NodeStatus{
						Addresses: []v1.NodeAddress{
							{
								Type:    v1.NodeInternalIP,
								Address: "127.0.0.2",
							},
						},
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "node-3",
					},
					Status: v1.NodeStatus{
						Addresses: []v1.NodeAddress{
							{
								Type:    v1.NodeInternalIP,
								Address: "127.0.0.3",
							},
						},
					},
				},
			},
			"linodelb",
			getLBStatus.Ingress[0].IP,
			nil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			lbStatus, err := lb.EnsureLoadBalancer(context.TODO(), test.clusterName, test.service, test.nodes)
			if err != nil {
				t.Fatal(err)
			}
			if lbStatus.Ingress[0].IP != test.nbIP {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.nbIP)
				t.Logf("actual: %v", lbStatus.Ingress)
			}
			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}

func testMakeLoadBalancerStatus(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	ipv4 := "192.168.0.1"
	hostname := "nb-192-168-0-1.newark.nodebalancer.linode.com"
	nb := &linodego.NodeBalancer{
		IPv4:     &ipv4,
		Hostname: &hostname,
	}

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test",
			Annotations: make(map[string]string, 1),
		},
	}

	expectedStatus := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{{
			Hostname: hostname,
			IP:       ipv4,
		}},
	}
	status := makeLoadBalancerStatus(svc, nb)
	if !reflect.DeepEqual(status, expectedStatus) {
		t.Errorf("expected status for basic service to be %#v; got %#v", expectedStatus, status)
	}

	svc.Annotations[annotations.AnnLinodeHostnameOnlyIngress] = "true"
	expectedStatus.Ingress[0] = v1.LoadBalancerIngress{Hostname: hostname}
	status = makeLoadBalancerStatus(svc, nb)
	if !reflect.DeepEqual(status, expectedStatus) {
		t.Errorf("expected status for %q annotated service to be %#v; got %#v", annotations.AnnLinodeHostnameOnlyIngress, expectedStatus, status)
	}
}

func testMakeLoadBalancerStatusEnvVar(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	ipv4 := "192.168.0.1"
	hostname := "nb-192-168-0-1.newark.nodebalancer.linode.com"
	nb := &linodego.NodeBalancer{
		IPv4:     &ipv4,
		Hostname: &hostname,
	}

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test",
			Annotations: make(map[string]string, 1),
		},
	}

	expectedStatus := &v1.LoadBalancerStatus{
		Ingress: []v1.LoadBalancerIngress{{
			Hostname: hostname,
			IP:       ipv4,
		}},
	}
	status := makeLoadBalancerStatus(svc, nb)
	if !reflect.DeepEqual(status, expectedStatus) {
		t.Errorf("expected status for basic service to be %#v; got %#v", expectedStatus, status)
	}

	t.Setenv("LINODE_HOSTNAME_ONLY_INGRESS", "true")
	expectedStatus.Ingress[0] = v1.LoadBalancerIngress{Hostname: hostname}
	status = makeLoadBalancerStatus(svc, nb)
	if !reflect.DeepEqual(status, expectedStatus) {
		t.Errorf("expected status for %q annotated service to be %#v; got %#v", annotations.AnnLinodeHostnameOnlyIngress, expectedStatus, status)
	}

	t.Setenv("LINODE_HOSTNAME_ONLY_INGRESS", "false")
	expectedStatus.Ingress[0] = v1.LoadBalancerIngress{Hostname: hostname}
	status = makeLoadBalancerStatus(svc, nb)
	if reflect.DeepEqual(status, expectedStatus) {
		t.Errorf("expected status for %q annotated service to be %#v; got %#v", annotations.AnnLinodeHostnameOnlyIngress, expectedStatus, status)
	}

	t.Setenv("LINODE_HOSTNAME_ONLY_INGRESS", "banana")
	expectedStatus.Ingress[0] = v1.LoadBalancerIngress{Hostname: hostname}
	status = makeLoadBalancerStatus(svc, nb)
	if reflect.DeepEqual(status, expectedStatus) {
		t.Errorf("expected status for %q annotated service to be %#v; got %#v", annotations.AnnLinodeHostnameOnlyIngress, expectedStatus, status)
	}
	os.Unsetenv("LINODE_HOSTNAME_ONLY_INGRESS")
}

func testCleanupDoesntCall(t *testing.T, client *linodego.Client, fakeAPI *fakeAPI) {
	region := "us-west"
	nb1, err := client.CreateNodeBalancer(context.TODO(), linodego.NodeBalancerCreateOptions{Region: region})
	if err != nil {
		t.Fatal(err)
	}
	nb2, err := client.CreateNodeBalancer(context.TODO(), linodego.NodeBalancerCreateOptions{Region: region})
	if err != nil {
		t.Fatal(err)
	}

	svc := &v1.Service{ObjectMeta: metav1.ObjectMeta{Name: "test"}}
	svcAnn := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test",
			Annotations: map[string]string{annotations.AnnLinodeNodeBalancerID: strconv.Itoa(nb2.ID)},
		},
	}
	svc.Status.LoadBalancer = *makeLoadBalancerStatus(svc, nb1)
	svcAnn.Status.LoadBalancer = *makeLoadBalancerStatus(svcAnn, nb1)
	lb := &loadbalancers{client, region, nil}

	fakeAPI.ResetRequests()
	t.Run("non-annotated service shouldn't call the API during cleanup", func(t *testing.T) {
		if err := lb.cleanupOldNodeBalancer(context.TODO(), svc); err != nil {
			t.Fatal(err)
		}
		if len(fakeAPI.requests) != 0 {
			t.Fatalf("unexpected API calls: %v", fakeAPI.requests)
		}
	})

	fakeAPI.ResetRequests()
	t.Run("annotated service calls the API to load said NB", func(t *testing.T) {
		if err := lb.cleanupOldNodeBalancer(context.TODO(), svcAnn); err != nil {
			t.Fatal(err)
		}
		expectedRequests := map[fakeRequest]struct{}{
			{Path: "/nodebalancers", Body: "", Method: "GET"}:                            {},
			{Path: fmt.Sprintf("/nodebalancers/%v", nb2.ID), Body: "", Method: "GET"}:    {},
			{Path: fmt.Sprintf("/nodebalancers/%v", nb1.ID), Body: "", Method: "DELETE"}: {},
		}
		if !reflect.DeepEqual(fakeAPI.requests, expectedRequests) {
			t.Fatalf("expected requests %#v, got %#v instead", expectedRequests, fakeAPI.requests)
		}
	})
}

func testUpdateLoadBalancerNoNodes(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        randString(),
			UID:         "foobar123",
			Annotations: map[string]string{},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     randString(),
					Protocol: "http",
					Port:     int32(80),
					NodePort: int32(8080),
				},
			},
		},
	}

	lb := &loadbalancers{client, "us-west", nil}
	defer func() {
		_ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc)
	}()

	fakeClientset := fake.NewSimpleClientset()
	lb.kubeClient = fakeClientset

	nodeBalancer, err := client.CreateNodeBalancer(context.TODO(), linodego.NodeBalancerCreateOptions{
		Region: lb.zone,
	})
	if err != nil {
		t.Fatalf("failed to create NodeBalancer: %s", err)
	}
	svc.Status.LoadBalancer = *makeLoadBalancerStatus(svc, nodeBalancer)
	stubService(fakeClientset, svc)
	svc.ObjectMeta.SetAnnotations(map[string]string{
		annotations.AnnLinodeNodeBalancerID: strconv.Itoa(nodeBalancer.ID),
	})

	// setup done, test ensure/update
	nodes := []*v1.Node{}

	if _, err = lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes); !stderrors.Is(err, errNoNodesAvailable) {
		t.Errorf("EnsureLoadBalancer should return %v, got %v", errNoNodesAvailable, err)
	}

	if err := lb.UpdateLoadBalancer(context.TODO(), "linodelb", svc, nodes); !stderrors.Is(err, errNoNodesAvailable) {
		t.Errorf("UpdateLoadBalancer should return %v, got %v", errNoNodesAvailable, err)
	}
}

func testGetNodeBalancerForServiceIDDoesNotExist(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	lb := &loadbalancers{client, "us-west", nil}
	bogusNodeBalancerID := "123456"

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeNodeBalancerID: bogusNodeBalancerID,
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(8443),
					NodePort: int32(30000),
				},
			},
		},
	}

	_, err := lb.getNodeBalancerForService(context.TODO(), svc)
	if err == nil {
		t.Fatal("expected getNodeBalancerForService to return an error")
	}

	nbid, _ := strconv.Atoi(bogusNodeBalancerID)
	expectedErr := lbNotFoundError{
		serviceNn:      getServiceNn(svc),
		nodeBalancerID: nbid,
	}
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error to be '%s' but got '%s'", expectedErr, err)
	}
}

func testEnsureNewLoadBalancerWithNodeBalancerID(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	lb := &loadbalancers{client, "us-west", nil}
	nodeBalancer, err := client.CreateNodeBalancer(context.TODO(), linodego.NodeBalancerCreateOptions{
		Region: lb.zone,
	})
	if err != nil {
		t.Fatalf("failed to create NodeBalancer: %s", err)
	}

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testensure",
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeNodeBalancerID: strconv.Itoa(nodeBalancer.ID),
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(8443),
					NodePort: int32(30000),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-1",
			},
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}

	defer func() { _ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc) }()

	if _, err = lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes); err != nil {
		t.Fatal(err)
	}
}

func testEnsureNewLoadBalancer(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "testensure",
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeDefaultProtocol:           "tcp",
				annotations.AnnLinodePortConfigPrefix + "8443": `{ "protocol": "https", "tls-secret-name": "tls-secret"}`,
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(8443),
					NodePort: int32(30000),
				},
				{
					Name:     "test2",
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30001),
				},
			},
		},
	}

	nodes := []*v1.Node{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name: "node-1",
			},
			Status: v1.NodeStatus{
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeInternalIP,
						Address: "127.0.0.1",
					},
				},
			},
		},
	}
	lb := &loadbalancers{client, "us-west", nil}
	lb.kubeClient = fake.NewSimpleClientset()
	addTLSSecret(t, lb.kubeClient)

	defer func() { _ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc) }()

	_, err := lb.EnsureLoadBalancer(context.TODO(), "linodelb", svc, nodes)
	if err != nil {
		t.Fatal(err)
	}
}

func testGetLoadBalancer(t *testing.T, client *linodego.Client, _ *fakeAPI) {
	lb := &loadbalancers{client, "us-west", nil}
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			UID:  "foobar123",
			Annotations: map[string]string{
				annotations.AnnLinodeDefaultProtocol: "tcp",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	configs := []*linodego.NodeBalancerConfigCreateOptions{}
	nb, err := lb.createNodeBalancer(context.TODO(), "linodelb", svc, configs)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = lb.EnsureLoadBalancerDeleted(context.TODO(), "linodelb", svc) }()

	lbStatus := makeLoadBalancerStatus(svc, nb)
	svc.Status.LoadBalancer = *lbStatus

	testcases := []struct {
		name        string
		service     *v1.Service
		clusterName string
		found       bool
		err         error
	}{
		{
			"Load balancer exists",
			svc,
			"linodelb",
			true,
			nil,
		},
		{
			"Load balancer not exists",

			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "notexists",
					UID:  "notexists123",
					Annotations: map[string]string{
						annotations.AnnLinodeDefaultProtocol: "tcp",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						{
							Name:     "test",
							Protocol: "TCP",
							Port:     int32(80),
							NodePort: int32(30000),
						},
					},
				},
			},
			"linodelb",
			false,
			nil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			_, found, err := lb.GetLoadBalancer(context.TODO(), test.clusterName, test.service)
			if found != test.found {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.found)
				t.Logf("actual: %v", found)
			}
			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}

func Test_getPortConfigAnnotation(t *testing.T) {
	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:     "test",
					Protocol: "TCP",
					Port:     int32(80),
					NodePort: int32(30000),
				},
			},
		},
	}

	testcases := []struct {
		name     string
		ann      map[string]string
		expected portConfigAnnotation
		err      string
	}{
		{
			name: "Test single port annotation",
			ann:  map[string]string{annotations.AnnLinodePortConfigPrefix + "443": `{ "tls-secret-name": "prod-app-tls", "protocol": "https" }`},
			expected: portConfigAnnotation{
				TLSSecretName: "prod-app-tls",
				Protocol:      "https",
			},
			err: "",
		},
		{
			name: "Test multiple port annotation",
			ann: map[string]string{
				annotations.AnnLinodePortConfigPrefix + "443": `{ "tls-secret-name": "prod-app-tls", "protocol": "https" }`,
				annotations.AnnLinodePortConfigPrefix + "80":  `{ "protocol": "http" }`,
			},
			expected: portConfigAnnotation{
				TLSSecretName: "prod-app-tls",
				Protocol:      "https",
			},
			err: "",
		},
		{
			name: "Test no port annotation",
			ann:  map[string]string{},
			expected: portConfigAnnotation{
				Protocol: "",
			},
			err: "",
		},
		{
			name: "Test invalid json",
			ann: map[string]string{
				annotations.AnnLinodePortConfigPrefix + "443": `{ "tls-secret-name": "prod-app-tls" `,
			},
			expected: portConfigAnnotation{},
			err:      "unexpected end of JSON input",
		},
	}
	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			svc.Annotations = test.ann
			ann, err := getPortConfigAnnotation(svc, 443)
			if !reflect.DeepEqual(ann, test.expected) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.expected)
				t.Logf("actual: %v", ann)
			}
			if test.err != "" && test.err != err.Error() {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}

func Test_getTLSCertInfo(t *testing.T) {
	kubeClient := fake.NewSimpleClientset()
	addTLSSecret(t, kubeClient)

	testcases := []struct {
		name       string
		portConfig portConfig
		cert       string
		key        string
		err        error
	}{
		{
			name: "Test valid Cert info",
			portConfig: portConfig{
				TLSSecretName: "tls-secret",
				Port:          8080,
			},
			cert: testCert,
			key:  testKey,
			err:  nil,
		},
		{
			name: "Test unspecified Cert info",
			portConfig: portConfig{
				Port: 8080,
			},
			cert: "",
			key:  "",
			err:  fmt.Errorf("TLS secret name for port 8080 is not specified"),
		},
		{
			name: "Test blank Cert info",
			portConfig: portConfig{
				TLSSecretName: "",
				Port:          8080,
			},
			cert: "",
			key:  "",
			err:  fmt.Errorf("TLS secret name for port 8080 is not specified"),
		},
		{
			name: "Test no secret found",
			portConfig: portConfig{
				TLSSecretName: "secret",
				Port:          8080,
			},
			cert: "",
			key:  "",
			err: errors.NewNotFound(schema.GroupResource{
				Group:    "",
				Resource: "secrets",
			}, "secret"), /*{}(`secrets "secret" not found`)*/
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			cert, key, err := getTLSCertInfo(context.TODO(), kubeClient, "", test.portConfig)
			if cert != test.cert {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.cert)
				t.Logf("actual: %v", cert)
			}
			if key != test.key {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.key)
				t.Logf("actual: %v", key)
			}
			if !reflect.DeepEqual(err, test.err) {
				t.Error("unexpected error")
				t.Logf("expected: %v", test.err)
				t.Logf("actual: %v", err)
			}
		})
	}
}

func addTLSSecret(t *testing.T, kubeClient kubernetes.Interface) {
	_, err := kubeClient.CoreV1().Secrets("").Create(context.TODO(), &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tls-secret",
		},
		Data: map[string][]byte{
			v1.TLSCertKey:       []byte(testCert),
			v1.TLSPrivateKeyKey: []byte(testKey),
		},
		StringData: nil,
		Type:       "kubernetes.io/tls",
	}, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("failed to add TLS secret: %s\n", err)
	}
}

func Test_LoadbalNodeNameCoercion(t *testing.T) {
	type testCase struct {
		nodeName       string
		padding        string
		expectedOutput string
	}
	testCases := []testCase{
		{
			nodeName:       "n",
			padding:        "z",
			expectedOutput: "zzn",
		},
		{
			nodeName:       "n",
			padding:        "node-",
			expectedOutput: "node-n",
		},
		{
			nodeName:       "n",
			padding:        "",
			expectedOutput: "xxn",
		},
		{
			nodeName:       "infra-logging-controlplane-3-atl1-us-prod",
			padding:        "node-",
			expectedOutput: "infra-logging-controlplane-3-atl",
		},
		{
			nodeName:       "node1",
			padding:        "node-",
			expectedOutput: "node1",
		},
	}

	for _, tc := range testCases {
		if out := coerceString(tc.nodeName, 3, 32, tc.padding); out != tc.expectedOutput {
			t.Fatalf("Expected loadbal backend name to be %s (got: %s)", tc.expectedOutput, out)
		}
	}
}
