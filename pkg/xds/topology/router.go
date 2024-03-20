package topology

import (
	mesh_proto "github.com/apache/dubbo-kubernetes/api/mesh/v1alpha1"
	core_mesh "github.com/apache/dubbo-kubernetes/pkg/core/resources/apis/mesh"
	core_xds "github.com/apache/dubbo-kubernetes/pkg/core/xds"
)

// BuildDestinationMap creates a map of selectors to match other dataplanes reachable from a given one
// via given routes.
func BuildDestinationMap(dataplane *core_mesh.DataplaneResource) core_xds.DestinationMap {
	destinations := core_xds.DestinationMap{}
	for _, oface := range dataplane.Spec.Networking.GetOutbound() {
		serviceName := oface.GetService()
		// TODO: traffic route is not implemented yet
		destinations[serviceName] = destinations[serviceName].Add(mesh_proto.MatchService(serviceName))
	}
	return destinations
}
