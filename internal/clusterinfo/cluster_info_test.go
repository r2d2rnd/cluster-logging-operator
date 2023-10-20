package clusterinfo

import (
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	configv1 "github.com/openshift/api/config/v1"
	hypershift "github.com/openshift/hypershift/api/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const namespace = "testing"

var (
	// All clusters have a ClusterVersion
	clusterVersion = &configv1.ClusterVersion{
		Spec: configv1.ClusterVersionSpec{
			ClusterID:     "clusterVersion-id",
			DesiredUpdate: &configv1.Update{Version: "clusterVersion-version"},
		}}

	// Each HCP management namespace has a HostedControlPlane
	hostedControlPlane = &hypershift.HostedControlPlane{
		ObjectMeta: metav1.ObjectMeta{Namespace: "testing", Name: "foobar"},
		Spec: hypershift.HostedControlPlaneSpec{
			ClusterID: "hypershift-id",
		},
		Status: hypershift.HostedControlPlaneStatus{
			VersionStatus: &hypershift.ClusterVersionStatus{
				Desired: configv1.Release{
					Version: "hypershift-version",
				},
			},
		},
	}
)

var _ = DescribeTable("GetClusterInfo", func(id, version string, objs ...runtime.Object) {
	c := fake.NewFakeClient(objs...)
	info, err := GetClusterInfo(c, &rest.Config{}, namespace)
	Expect(err).To(Succeed())
	Expect(info).To(Equal(ClusterInfo{Version: version, ID: id}))
},
	Entry("Standalone cluster", "clusterVersion-id", "clusterVersion-version", clusterVersion),
	Entry("Hypershift cluster", "hypershift-id", "hypershift-version", clusterVersion, hostedControlPlane),
)
