package clusterinfo

import (
	"context"

	configv1 "github.com/openshift/api/config/v1"
	hypershift "github.com/openshift/hypershift/api/v1beta1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterInfo is global information about where the ClusterLogForwarder is running.
type ClusterInfo struct {
	Version string // Version of the cluster.
	ID      string // Unique identifier of the cluster.
}

// GetClusterInfo gets cluster info for the cluster we are running in.
//
// If the namepsace contains a Hosted Control Plane (we are on the hypershift management cluster)
// this information describes the *guest* cluster, not the host cluster.
// We assume in this case that CLF is running on behalf of the guest cluster to collect API audit logs.
//
// FIXME review assumptions.
func GetClusterInfo(c client.Reader, cfg *rest.Config, namespace string) (*ClusterInfo, error) {
	// Use HCP info if exactly one HCP is present in the current namespace.
	hcps := &hypershift.HostedControlPlaneList{}
	if err := c.List(context.Background(), hcps, client.InNamespace(namespace)); err == nil && len(hcps.Items) == 1 {
		return &ClusterInfo{
			Version: hcps.Items[0].Status.VersionStatus.Desired.Version,
			ID:      hcps.Items[0].Spec.ClusterID,
		}, nil
	}
	// Use standalone cluster info.
	cv := &configv1.ClusterVersion{}
	if err := c.Get(context.Background(), client.ObjectKey{Name: "version"}, cv); err != nil {
		return nil, err
	}
	return &ClusterInfo{
		Version: cv.Spec.DesiredUpdate.Version,
		ID:      string(cv.Spec.ClusterID),
	}, nil
}
