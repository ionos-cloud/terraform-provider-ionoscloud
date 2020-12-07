package profitbricks

import (
	"fmt"
	"net/url"
	"path"
)

func safeJoin(str ...string) string {
	for _, s := range str {
		if s == "" {
			panic(fmt.Sprintf("path contains empty element: %v", str))
		}
	}
	return path.Join(str...)
}

// datacentersPath: "datacenters"
func datacentersPath() string {
	return "datacenters"
}

// datacenterPath: "datacenters/<datacenterID>"
func datacenterPath(datacenterID string) string {
	return safeJoin(datacentersPath(), url.QueryEscape(datacenterID))
}

// imagesPath: "images"
func imagesPath() string {
	return "images"
}

// imagePath: "images/<imageID>"
func imagePath(imageID string) string {
	return safeJoin(imagesPath(), url.QueryEscape(imageID))
}

// ipblocksPath: "ipblocks"
func ipblocksPath() string {
	return "ipblocks"
}

//  ipblockPath: "ipblocks/<ipblockID>"
func ipblockPath(ipblockID string) string {
	return safeJoin(ipblocksPath(), url.QueryEscape(ipblockID))
}

// locationsPath: "locations"
func locationsPath() string {
	return "locations"
}

// locationRegionPath: "locations/<regionID>"
func locationRegionPath(regionID string) string {
	return safeJoin(locationsPath(), url.QueryEscape(regionID))
}

// locationPath: "locations/<regionID>/<locationID>"
func locationPath(regionID, locationID string) string {
	return safeJoin(locationRegionPath(regionID), url.QueryEscape(locationID))
}

// snapshotsPath: "snapshots"
func snapshotsPath() string {
	return "snapshots"
}

// snapshotPath: "snapshots/<snapshotID>"
func snapshotPath(snapshotID string) string {
	return safeJoin(snapshotsPath(), url.QueryEscape(snapshotID))
}

// lansPath: "datacenters/<datacenterID>/lans"
func lansPath(datacenterID string) string {
	return safeJoin(datacenterPath(datacenterID), "lans")
}

// lanPath: "datacenters/<datacenterID>/lans/<lanID>"
func lanPath(datacenterID, lanID string) string {
	return safeJoin(lansPath(datacenterID), url.QueryEscape(lanID))
}

//  loadbalancersPath: "datacenters/<datacenterID>/loadbalancers"
func loadbalancersPath(datacenterID string) string {
	return safeJoin(datacenterPath(datacenterID), "loadbalancers")
}

// loadbalancerPath: "datacenters/<datacenterID>/loadbalancers/<loadbalancerID>"
func loadbalancerPath(datacenterID, loadbalancerID string) string {
	return safeJoin(loadbalancersPath(datacenterID), url.QueryEscape(loadbalancerID))
}

// serversPath: "datacenters/<datacenterID>/servers"
func serversPath(datacenterID string) string {
	return safeJoin(datacenterPath(datacenterID), "servers")
}

// serverPath: "datacenters/<datacenterID>/servers/<serverID>"
func serverPath(datacenterID, serverID string) string {
	return safeJoin(serversPath(datacenterID), url.QueryEscape(serverID))
}

// serverStartPath: "datacenters/<datacenterID>/servers/<serverID>/start"
func serverStartPath(datacenterID, serverID string) string {
	return safeJoin(serverPath(datacenterID, serverID), "start")
}

// serverStopPath: "datacenters/<datacenterID>/servers/<serverID>/stop"
func serverStopPath(datacenterID, serverID string) string {
	return safeJoin(serverPath(datacenterID, serverID), "stop")
}

// serverRebootPath: "datacenters/<datacenterID>/servers/<serverID>/reboot"
func serverRebootPath(datacenterID, serverID string) string {
	return safeJoin(serverPath(datacenterID, serverID), "reboot")
}

// volumesPath: "datacenters/<datacenterID>/volumes"
func volumesPath(datacenterID string) string {
	return safeJoin(datacenterPath(datacenterID), "volumes")
}

// volume_path "datacenters/<datacenterID>/volumes/<volumeID>"
func volumePath(datacenterID, volumeID string) string {
	return safeJoin(volumesPath(datacenterID), url.QueryEscape(volumeID))
}

// createSnapshotPath: "datacenters/<datacenterID>/volumes/<volumeID>/create-snapshot"
func createSnapshotPath(datacenterID, volumeID string) string {
	return safeJoin(volumePath(datacenterID, volumeID), "create-snapshot")
}

// restoreSnapshotPath: "datacenters/<datacenterID>/volumes/<volumeID>/restore-snapshot"
func restoreSnapshotPath(datacenterID, volumeID string) string {
	return safeJoin(volumePath(datacenterID, volumeID), "restore-snapshot")
}

//  balancedNicsPath: "datacenters/<datacenterID>/loadbalancers/<loadbalancerID>/balancednics"
func balancedNicsPath(datacenterID, loadbalancerID string) string {
	return safeJoin(loadbalancerPath(datacenterID, loadbalancerID), "balancednics")
}

//  balancedNicPath: "datacenters/<datacenterID>/loadbalancers/<loadbalancerID>/balancednics/<balancedNicID>"
func balancedNicPath(datacenterID, loadbalancerID, balancedNicID string) string {
	return safeJoin(balancedNicsPath(datacenterID, loadbalancerID), url.QueryEscape(balancedNicID))
}

// cdromsPath: "datacenters/<datacenterID>/servers/<serverID>/cdroms"
func cdromsPath(datacenterID, serverID string) string {
	return safeJoin(serverPath(datacenterID, serverID), "cdroms")
}

// cdromPath: "datacenters/<datacenterID>/servers/<serverID>/cdroms/<cdID>"
func cdromPath(datacenterID, serverID, cdID string) string {
	return safeJoin(cdromsPath(datacenterID, serverID), url.QueryEscape(cdID))
}

// attachedVolumesPath: "datacenters/<datacenterID>/servers/<serverID>/volumes"
func attachedVolumesPath(datacenterID, serverID string) string {
	return safeJoin(serverPath(datacenterID, serverID), "volumes")
}

// attachedVolumePath: "datacenters/<datacenterID>/servers/<serverID>/volumes/<volumeID>"
func attachedVolumePath(datacenterID, serverID, volumeID string) string {
	return safeJoin(attachedVolumesPath(datacenterID, serverID), url.QueryEscape(volumeID))
}

// nicsPath: "datacenters/<datacenterID>/servers/<serverID>/nics"
func nicsPath(datacenterID, serverID string) string {
	return safeJoin(serverPath(datacenterID, serverID), "nics")
}

// nicPath: "datacenters/<datacenterID>/servers/<serverID>/nics/<nicID>"
func nicPath(datacenterID, serverID, nicID string) string {
	return safeJoin(nicsPath(datacenterID, serverID), url.QueryEscape(nicID))
}

// firewallRulesPath: "datacenters/<datacenterID>/servers/<serverID>/nics/<nicID>/firewallrules"
func firewallRulesPath(datacenterID, serverID, nicID string) string {
	return safeJoin(nicPath(datacenterID, serverID, nicID), "firewallrules")
}

// firewallRulePath:
//  "datacenters/<datacenterID>/servers/<serverID>/nics/<nicID>/firewallrules/<firewallRuleID>"
func firewallRulePath(datacenterID, serverID, nicID, firewallRuleID string) string {
	return safeJoin(firewallRulesPath(datacenterID, serverID, nicID), url.QueryEscape(firewallRuleID))
}

// RequestsPath: "requests"
func RequestsPath() string {
	return "requests"
}

// RequestPath: "requests/<requestID>"
func RequestPath(requestID string) string {
	return safeJoin(RequestsPath(), url.QueryEscape(requestID))
}

// RequestStatusPath: "requests/<requestID>/status"
func RequestStatusPath(requestID string) string {
	return safeJoin(RequestPath(requestID), "status")
}

// contractsPath: "contracts"
func contractsPath() string {
	return "contracts"
}

// umPath: "um"
func um() string {
	return "um"
}

// groupsPath: "um/groups"
func groupsPath() string {
	return safeJoin(um(), "groups")
}

// groupPath: "um/groups/<groupID>"
func groupPath(groupID string) string {
	return safeJoin(groupsPath(), url.QueryEscape(groupID))
}

// sharesPath: "um/groups/<groupID>/shares"
func sharesPath(groupID string) string {
	return safeJoin(groupPath(groupID), "shares")
}

// sharePath: "um/groups/<groupID>/shares/<resourceID>"
func sharePath(groupID string, resourceID string) string {
	return safeJoin(sharesPath(groupID), url.QueryEscape(resourceID))
}

// groupUsersPath: "um/groups/<groupID>/users"
func groupUsersPath(groupID string) string {
	return safeJoin(groupPath(groupID), "users")
}

// groupUserPath: "um/groups/<groupID>/users/<userID>"
func groupUserPath(groupID string, userID string) string {
	return safeJoin(groupUsersPath(groupID), url.QueryEscape(userID))
}

// usersPath: "um/users"
func usersPath() string {
	return safeJoin(um(), "users")
}

// userPath: "um/users/<userID>"
func userPath(userID string) string {
	return safeJoin(usersPath(), url.QueryEscape(userID))
}

// resourcesPath: "um/resources"
func resourcesPath() string {
	return safeJoin(um(), "resources")
}

// resourcesTypePath: "um/resources/<resourceType>"
func resourcesTypePath(resourceType string) string {
	return safeJoin(resourcesPath(), url.QueryEscape(resourceType))
}

// resourcePath: "um/resources/<resourceType>/<resourceID>"
func resourcePath(resourceType string, resourceID string) string {
	return safeJoin(resourcesTypePath(resourceType), url.QueryEscape(resourceID))
}

// tokensPath: "/tokens"
func tokensPath() string {
	// comes with leading slash, as it is used for calls to a different api.
	return "/tokens"
}

// tokenPath: "tokens/<tokenID>"
func tokenPath(tokenID string) string {
	return safeJoin(tokensPath(), url.QueryEscape(tokenID))
}

// kubernetesClustersPath: "k8s"
func kubernetesClustersPath() string {
	return "k8s"
}

// kubernetesClusterPath: "k8s/<clusterID>"
func kubernetesClusterPath(clusterID string) string {
	return safeJoin(kubernetesClustersPath(), clusterID)
}

// kubeConfigPath: "k8s/<clusterID>/kubeconfig"
func kubeConfigPath(clusterID string) string {
	return safeJoin(kubernetesClusterPath(clusterID), "kubeconfig")
}

// kubernetesNodePoolsPath: "k8s/<clusterID>/nodepools"
func kubernetesNodePoolsPath(clusterID string) string {
	return safeJoin(kubernetesClusterPath(clusterID), "nodepools")
}

// kubernetesNodePoolPath: "k8s/<clusterID>/nodepools/<nodepoolID>"
func kubernetesNodePoolPath(clusterID, nodepoolID string) string {
	return safeJoin(kubernetesNodePoolsPath(clusterID), nodepoolID)
}

// kubernetesNodesPath: "k8s/<clusterID>/nodepools/<nodepoolID>/nodes"
func kubernetesNodesPath(clusterID, nodepoolID string) string {
	return safeJoin(kubernetesNodePoolPath(clusterID, nodepoolID), "nodes")
}

// kubernetesNodePath: "k8s/<clusterID>/nodepools/<nodepoolID>/nodes/<nodeID>"
func kubernetesNodePath(clusterID, nodepoolID, nodeID string) string {
	return safeJoin(kubernetesNodesPath(clusterID, nodepoolID), nodeID)
}

// kubernetesNodeReplacePath: "k8s/<clusterID>/nodepools/<nodepoolID>/nodes/<nodeID>/replace"
func kubernetesNodeReplacePath(clusterID, nodepoolID, nodeID string) string {
	return safeJoin(kubernetesNodePath(clusterID, nodepoolID, nodeID), "replace")
}

// backupUnitsPath: "backupunits"
func backupUnitsPath() string {
	return "backupunits"
}

// backupUnitsPath: "backupunits/<backupUnitID>"
func backupUnitPath(backupUnitID string) string {
	return safeJoin(backupUnitsPath(), backupUnitID)
}

// backupUnitSSOURLPath: "backupunits/backupUnitID/ssourl"
func backupUnitSSOURLPath(backupUnitID string) string {
	return safeJoin(backupUnitsPath(), backupUnitID, "ssourl")
}

// s3KeysPath: "um/users/<userID>/s3keys"
func s3KeysPath(userID string) string {
	return safeJoin(userPath(userID), "s3keys")
}

// s3KeysListPath: "um/users/<userID>/s3keys?depth=1"
func s3KeysListPath(userID string) string {
	return safeJoin(userPath(userID), "s3keys?depth=1")
}

// s3KeyPath: "um/users/<userID>/s3keys/<s3KeyID>"
func s3KeyPath(userID string, s3KeyID string) string {
	return safeJoin(s3KeysPath(userID), s3KeyID)
}

// PrivateCrossConnectsPath: "pccs"
func PrivateCrossConnectsPath() string {
	return "pccs"
}

// PrivateCrossConnectPath: "pccs/<pccID>"
func PrivateCrossConnectPath(pccID string) string {
	return safeJoin(PrivateCrossConnectsPath(), pccID)
}
