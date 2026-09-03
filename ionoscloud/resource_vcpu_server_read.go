package ionoscloud

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	ionoscloud "github.com/ionos-cloud/sdk-go/v6"

	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/bundleclient"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapifirewall"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/cloudapinic"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/services/cloudapi/nsg"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils"
	diagutil "github.com/ionos-cloud/terraform-provider-ionoscloud/v6/utils/diags"
)

// The ionoscloud_vcpu_server resource has its own read/import/state-writer, deliberately NOT shared
// with ionoscloud_server. This keeps the two schemas decoupled: a change to the enterprise
// server's state-writer can no longer set a key that is missing from the VCPU schema (the 6.7.36
// enabled_features regression). Leaf helpers (SetVolumeProperties, cloudapinic/cloudapifirewall/nsg
// services, LabelsService, serverIsConfidential) remain shared - only the top-level per-type
// state-writer is duplicated.

func resourceVCPUServerRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	location := d.Get("location").(string)
	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return diag.FromErr(err)
	}

	dcID := d.Get("datacenter_id").(string)
	serverID := d.Id()

	server, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, dcID, serverID).Depth(4).Execute()
	logApiRequestTime(apiResponse)
	if err != nil {
		if httpNotFound(apiResponse) {
			tflog.Debug(ctx, "cannot find vcpu server by id", map[string]any{"server_id": serverID})
			d.SetId("")
			return nil
		}
		return diagutil.ToDiags(d, fmt.Errorf("error occurred while fetching a vcpu server: %w", err), nil)
	}
	if err := setResourceVCPUServerData(ctx, client, d, &server); err != nil {
		return diagutil.ToDiags(d, err, nil)
	}

	return nil
}

func resourceVCPUServerImport(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	importID := d.Id()
	location, parts := splitImportID(importID, "/")

	if len(parts) < 2 {
		return nil, diagutil.ToError(d, fmt.Errorf(
			"invalid import identifier: expected one of <location>:<datacenter-id>/<server-id> "+
				"or <datacenter-id>/<server-id>, got: %s", importID,
		), nil)
	}

	if err := validateImportIDParts(parts); err != nil {
		return nil, diagutil.ToError(d, fmt.Errorf("failed validating import identifier %q: %w", importID, err), nil)
	}

	datacenterID := parts[0]
	serverID := parts[1]

	client, err := meta.(bundleclient.SdkBundle).NewCloudAPIClient(ctx, location)
	if err != nil {
		return nil, err
	}

	server, apiResponse, err := client.ServersApi.DatacentersServersFindById(ctx, datacenterID, serverID).Depth(3).Execute()
	logApiRequestTime(apiResponse)

	if err != nil {
		if httpNotFound(apiResponse) {
			d.SetId("")
			return nil, diagutil.ToError(d, fmt.Errorf("unable to find vcpu server %q", serverID), nil)
		}
		return nil, diagutil.ToError(d, fmt.Errorf("error occurred while fetching a vcpu server ID %s %w", importID, err), nil)
	}
	var primaryNic ionoscloud.Nic
	d.SetId(*server.Id)
	primaryNicID := ""
	// first we try to get primary nic from parts, then if that fails, we get it from entities.
	if len(parts) > 2 {
		primaryNicID = parts[2]
		if err := d.Set("primary_nic", primaryNicID); err != nil {
			return nil, diagutil.ToError(d, fmt.Errorf("error setting primary_nic id %w", err), nil)
		}
	} else {
		if server.Entities != nil && server.Entities.Nics != nil && len(*server.Entities.Nics.Items) > 0 {
			primaryNic = (*server.Entities.Nics.Items)[0]
		}
	}
	if primaryNicID != "" {
		if server.Entities != nil && server.Entities.Nics != nil && server.Entities.Nics.Items != nil {
			for _, nic := range *server.Entities.Nics.Items {
				if *nic.Id == primaryNicID {
					primaryNic = nic
					if primaryNic.Properties != nil && *nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
						tflog.Debug(ctx, "setting primary_ip", map[string]any{"primary_ip": (*primaryNic.Properties.Ips)[0]})
						if err := d.Set("primary_ip", (*primaryNic.Properties.Ips)[0]); err != nil {
							return nil, diagutil.ToError(d, fmt.Errorf("error while setting primary ip: %w", err), nil)
						}
					}
					break
				}
			}
		}
	}

	if err := d.Set("datacenter_id", datacenterID); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}
	if err := d.Set("location", location); err != nil {
		return nil, err
	}

	if len(parts) > 3 {
		var rules []string
		rules = append(rules, parts[3])
		if err = cloudapifirewall.SetIdsInSchema(d, rules); err != nil {
			return nil, diagutil.ToError(d, err, nil)
		}

	}

	if err := setResourceVCPUServerData(ctx, client, d, &server); err != nil {
		return nil, diagutil.ToError(d, err, nil)
	}

	return []*schema.ResourceData{d}, nil
}

// setResourceVCPUServerData is the VCPU-specific state-writer. It mirrors setResourceServerData but
// is intentionally independent so the two resource schemas cannot drift into each other. Unlike the
// enterprise writer it never sets ssh_key_path (VCPU volumes do not expose it).
func setResourceVCPUServerData(ctx context.Context, client *ionoscloud.APIClient, d *schema.ResourceData, server *ionoscloud.Server) error {
	if server.Id != nil {
		d.SetId(*server.Id)
	}
	// takes care of an upgrade from a version that does not have firewallrule_ids(pre 6.4.2)
	// to one that has it(>=6.4.2)
	if err := cloudapifirewall.SetFwRuleIdsInSchemaInCaseOfProviderUpdate(d); err != nil {
		return err
	}

	// takes care of an upgrade from a version that does not have inline_volume_ids(pre 6.4.0)
	// to one that has it(>6.4.0). GetOk cannot be used here since it also returns false when
	// inline_volume_ids is present in the state as an empty list; checking the raw state directly
	// ensures this only fires when the attribute is completely absent.
	if rawState := d.GetRawState(); !rawState.IsNull() && rawState.GetAttr("inline_volume_ids").IsNull() {
		if bootVolumeItf, ok := d.GetOk("boot_volume"); ok {
			bootVolume := bootVolumeItf.(string)
			var inlineVolumeIDs []string
			inlineVolumeIDs = append(inlineVolumeIDs, bootVolume)
			if err := d.Set("inline_volume_ids", inlineVolumeIDs); err != nil {
				return utils.GenerateSetError("vcpu server", "inline_volume_ids", err)
			}
		}
	}

	datacenterID := d.Get("datacenter_id").(string)
	if server.Properties != nil {
		if server.Properties.Name != nil {
			if err := d.Set("name", *server.Properties.Name); err != nil {
				return fmt.Errorf("error setting name %w", err)
			}
		}
		if server.Properties.Hostname != nil {
			if err := d.Set("hostname", *server.Properties.Hostname); err != nil {
				return fmt.Errorf("error setting hostname %w", err)
			}
		}
		if server.Properties.Cores != nil {
			if err := d.Set("cores", *server.Properties.Cores); err != nil {
				return fmt.Errorf("error setting cores %w", err)
			}
		}

		if server.Properties.Ram != nil {
			if err := d.Set("ram", *server.Properties.Ram); err != nil {
				return fmt.Errorf("error setting ram %w", err)
			}
		}

		if server.Properties.AvailabilityZone != nil {
			if err := d.Set("availability_zone", *server.Properties.AvailabilityZone); err != nil {
				return fmt.Errorf("error setting availability_zone %w", err)
			}
		}

		if server.Properties.CpuFamily != nil {
			if err := d.Set("cpu_family", *server.Properties.CpuFamily); err != nil {
				return fmt.Errorf("error setting cpu_family %w", err)
			}
		}

		if server.Properties.Type != nil {
			if err := d.Set("type", *server.Properties.Type); err != nil {
				return fmt.Errorf("error setting type %w", err)
			}
		}

		if server.Properties.VmState != nil {
			if err := d.Set("vm_state", *server.Properties.VmState); err != nil {
				return fmt.Errorf("error setting vm_state %w", err)
			}
		}

		// Shared with the enterprise writer so these crash-prone keys can't drift between the two.
		if err := setServerConfidentialVisibility(d, server); err != nil {
			return err
		}

		if server.Properties.BootCdrom != nil {
			if err := d.Set("boot_cdrom", *server.Properties.BootCdrom.Id); err != nil {
				return fmt.Errorf("error setting boot_cdrom %w", err)
			}
		} else {
			d.Set("boot_cdrom", nil)
		}

		if server.Properties.BootVolume != nil {
			if err := d.Set("boot_volume", *server.Properties.BootVolume.Id); err != nil {
				return fmt.Errorf("error setting bootVolume %w", err)
			}
		} else {
			d.Set("boot_volume", nil)
		}
		if server.Entities != nil {
			if server.Entities.Volumes != nil && server.Entities.Volumes.Items != nil && len(*server.Entities.Volumes.Items) > 0 &&
				(*server.Entities.Volumes.Items)[0].Properties != nil && (*server.Entities.Volumes.Items)[0].Properties.Image != nil {
				if err := d.Set("boot_image", *(*server.Entities.Volumes.Items)[0].Properties.Image); err != nil {
					return fmt.Errorf("error setting boot_image %w", err)
				}
			}
			if server.Entities.Securitygroups != nil && server.Entities.Securitygroups.Items != nil {
				if err := nsg.SetNSGInResourceData(d, server.Entities.Securitygroups.Items); err != nil {
					return err
				}
			}
		}
		if server.Properties.NicMultiQueue != nil {
			if err := d.Set("nic_multi_queue", *server.Properties.NicMultiQueue); err != nil {
				return fmt.Errorf("error setting nic_multi_queue: %w", err)
			}
		}
	}

	if server.Entities == nil {
		return fmt.Errorf("vcpu server entities cannot be empty for %s", d.Id())
	}

	inlineVolumeIDs := d.Get("inline_volume_ids")
	if inlineVolumeIDs != nil {
		inlineVolumeIDs := inlineVolumeIDs.([]any)
		var volumes []any
		for i, volumeID := range inlineVolumeIDs {
			volume, apiResponse, err := client.ServersApi.DatacentersServersVolumesFindById(ctx, datacenterID, d.Id(), volumeID.(string)).Execute()
			logApiRequestTime(apiResponse)
			if err != nil {
				if apiResponse.HttpNotFound() {
					tflog.Info(ctx, "inline volume not found", map[string]any{"volume_id": volumeID.(string), "datacenter_id": datacenterID, "server_id": d.Id()})
					continue
				}
				return fmt.Errorf("error retrieving inline volume %w", err)
			}
			volumePath := fmt.Sprintf("volume.%d.", i)
			entry := SetVolumeProperties(volume)
			userData := d.Get(volumePath + "user_data")
			entry["user_data"] = userData
			backupUnit := d.Get(volumePath + "backup_unit_id")
			entry["backup_unit_id"] = backupUnit
			volumes = append(volumes, entry)
		}
		if err := d.Set("volume", volumes); err != nil {
			return fmt.Errorf("error setting inline volumes %w", err)
		}
	}

	// take nic and firewall from schema if set is used in resource read, else take it from entities
	var nicID string
	firewallRuleIDs := d.Get("firewallrule_ids").([]any)

	if nicIntf, primaryNicOk := d.GetOk("primary_nic"); primaryNicOk {
		nicID = nicIntf.(string)
		ns := cloudapinic.Service{Client: client, Meta: nil, D: d}
		nic, apiResponse, err := ns.Get(ctx, datacenterID, d.Id(), nicID, 2)
		if err != nil {
			// fixes #467
			if apiResponse.HttpNotFound() {
				tflog.Debug(ctx, "primary nic not found, clearing primary_nic/primary_ip/nic", map[string]any{"nic_id": nicID})
				if err := d.Set("primary_nic", ""); err != nil {
					return err
				}
				if err := d.Set("primary_ip", ""); err != nil {
					return err
				}
				if err := d.Set("nic", nil); err != nil {
					return err
				}
			} else {
				return err
			}
		}
		var nicEntry map[string]any
		var fwRulesEntries []map[string]any

		if nic != nil && nic.Properties != nil {
			// fixes #467
			if nic.Properties.Ips != nil && len(*nic.Properties.Ips) > 0 {
				if err := d.Set("primary_ip", (*nic.Properties.Ips)[0]); err != nil {
					return err
				}
			}
			nicEntry = cloudapinic.SetNetworkProperties(*nic)
			nicEntry["id"] = *nic.Id
			fs := cloudapifirewall.Service{Client: client, D: d}

			for _, id := range firewallRuleIDs {
				firewallEntry, err := fs.AddToMapIfRuleExists(ctx, datacenterID, d.Id(), nicID, id.(string))
				if err != nil {
					return err
				}
				if firewallEntry != nil && len(firewallEntry) != 0 {
					fwRulesEntries = append(fwRulesEntries, firewallEntry)
				}
			}
		}
		if nic != nil && nic.Entities != nil && nic.Entities.Securitygroups != nil && nic.Entities.Securitygroups.Items != nil {
			nsgIDs := make([]string, 0)
			for _, group := range *nic.Entities.Securitygroups.Items {
				if group.Id != nil {
					id := *group.Id
					nsgIDs = append(nsgIDs, id)
				}
			}
			utils.SetPropWithNilCheck(nicEntry, "security_groups_ids", nsgIDs)
		}
		nics := []map[string]any{}
		if fwRulesEntries != nil {
			nicEntry["firewall"] = fwRulesEntries
		}
		if len(nicEntry) > 0 {
			nics = []map[string]any{nicEntry}
		}
		if err := d.Set("nic", nics); err != nil {
			return fmt.Errorf("error settings nics %w", err)
		}
	}
	if len(firewallRuleIDs) == 0 {
		if err := d.Set("firewallrule_id", ""); err != nil {
			return err
		}
	}
	if err := d.Set("firewallrule_ids", firewallRuleIDs); err != nil {
		return err
	}

	// Labels logic
	ls := LabelsService{ctx: ctx, client: client}
	labels, err := ls.datacentersServersLabelsGet(datacenterID, d.Id(), false)
	if err != nil {
		return err
	}
	if err := d.Set("label", labels); err != nil {
		return err
	}

	return nil
}
