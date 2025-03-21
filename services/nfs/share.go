package nfs

import (
	"context"

	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionos-cloud/sdk-go-bundle/products/nfs/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
)

// GetNFSShareByID returns a share given an ID
func (c *Client) GetNFSShareByID(ctx context.Context, clusterID, shareID string, location string) (nfs.ShareRead, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.NFS, location)
	share, apiResponse, err := c.sdkClient.SharesApi.ClustersSharesFindById(ctx, clusterID, shareID).Execute()
	apiResponse.LogInfo()
	return share, apiResponse, err
}

// ListNFSShares returns a list of all shares
func (c *Client) ListNFSShares(ctx context.Context, d *schema.ResourceData) (nfs.ShareReadList, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.NFS, d.Get("location").(string))
	shares, apiResponse, err := c.sdkClient.SharesApi.
		ClustersSharesGet(ctx, d.Get("cluster_id").(string)).Execute()
	apiResponse.LogInfo()
	return shares, apiResponse, err
}

// DeleteNFSShare deletes a share given an ID
func (c *Client) DeleteNFSShare(ctx context.Context, clusterID, shareID string, location string) (*shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.NFS, location)
	apiResponse, err := c.sdkClient.SharesApi.ClustersSharesDelete(ctx, clusterID, shareID).Execute()
	apiResponse.LogInfo()
	return apiResponse, err
}

// UpdateNFSShare updates an existing share
func (c *Client) UpdateNFSShare(ctx context.Context, d *schema.ResourceData) (nfs.ShareRead, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.NFS, d.Get("location").(string))
	share, apiResponse, err := c.sdkClient.SharesApi.
		ClustersSharesPut(ctx, d.Get("cluster_id").(string), d.Id()).ShareEnsure(*setShareEnsureRequest(d)).Execute()
	apiResponse.LogInfo()
	return share, apiResponse, err
}

// CreateNFSShare creates a new share
func (c *Client) CreateNFSShare(ctx context.Context, d *schema.ResourceData) (nfs.ShareRead, *shared.APIResponse, error) {
	c.overrideClientEndpoint(fileconfiguration.NFS, d.Get("location").(string))
	share, apiResponse, err := c.sdkClient.SharesApi.
		ClustersSharesPost(ctx, d.Get("cluster_id").(string)).ShareCreate(*setShareCreateRequest(d)).Execute()
	apiResponse.LogInfo()
	return share, apiResponse, err
}

// setShareCreateRequest returns a ShareCreate object
func setShareCreateRequest(d *schema.ResourceData) *nfs.ShareCreate {
	return nfs.NewShareCreate(setShareConfig(d))
}

// setShareEnsureRequest returns a ShareEnsure object
func setShareEnsureRequest(d *schema.ResourceData) *nfs.ShareEnsure {
	shareID := d.Id()
	share := setShareConfig(d)

	return nfs.NewShareEnsure(shareID, share)
}

// setShareConfig returns a Share object
func setShareConfig(d *schema.ResourceData) nfs.Share {
	name := d.Get("name").(string)
	quota := int32(d.Get("quota").(int))
	gid := int32(d.Get("gid").(int))
	uid := int32(d.Get("uid").(int))

	clientGroupsRaw := d.Get("client_groups").([]interface{})
	clientGroups := make([]nfs.ShareClientGroups, 0, len(clientGroupsRaw))
	for _, cgRaw := range clientGroupsRaw {
		cg := cgRaw.(map[string]interface{})
		description := cg["description"].(string)
		ipNetworksRaw := cg["ip_networks"].([]interface{})
		var ipNetworks []string
		for _, ip := range ipNetworksRaw {
			ipNetworks = append(ipNetworks, ip.(string))
		}
		hostsRaw := cg["hosts"].([]interface{})
		var hosts []string
		for _, host := range hostsRaw {
			hosts = append(hosts, host.(string))
		}
		nfsRaw := cg["nfs"].([]interface{})
		var squash string
		if len(nfsRaw) > 0 {
			nfsData := nfsRaw[0].(map[string]interface{})
			squash = nfsData["squash"].(string)
		}

		clientGroups = append(clientGroups, nfs.ShareClientGroups{
			Description: &description,
			IpNetworks:  ipNetworks,
			Hosts:       hosts,
			Nfs: &nfs.ShareClientGroupsNfs{
				Squash: &squash,
			},
		})
	}

	return nfs.Share{
		Name:         name,
		Quota:        &quota,
		Gid:          &gid,
		Uid:          &uid,
		ClientGroups: clientGroups,
	}
}

// SetNFSShareData sets data from the SDK share to the resource data
func (c *Client) SetNFSShareData(d *schema.ResourceData, share nfs.ShareRead) error {
	d.SetId(share.Id)
	if err := d.Set("name", share.Properties.Name); err != nil {
		return err
	}
	if err := d.Set("nfs_path", share.Metadata.NfsPath); err != nil {
		return err
	}
	if err := d.Set("quota", int(*share.Properties.Quota)); err != nil {
		return err
	}
	if err := d.Set("gid", int(*share.Properties.Gid)); err != nil {
		return err
	}
	if err := d.Set("uid", int(*share.Properties.Uid)); err != nil {
		return err
	}
	if err := d.Set("client_groups", flattenClientGroups(share.Properties.ClientGroups)); err != nil {
		return err
	}
	return nil
}

func flattenClientGroups(clientGroups []nfs.ShareClientGroups) []map[string]interface{} {
	result := make([]map[string]interface{}, len(clientGroups))
	for i, cg := range clientGroups {
		flattened := map[string]interface{}{
			"description": *cg.Description,
			"ip_networks": cg.IpNetworks,
			"hosts":       cg.Hosts,
			"nfs": []map[string]interface{}{
				{
					"squash": *cg.Nfs.Squash,
				},
			},
		}
		result[i] = flattened
	}
	return result
}
