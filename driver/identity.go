/*
   Copyright (c) 2025 Saglabs SA. All Rights Reserved.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package driver

import (
	"context"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/protobuf/ptypes/wrappers"
)

type IdentityService struct {
	csi.UnimplementedIdentityServer
	Service

	MountPoints []*mfsHandler
}

var _ csi.IdentityServer = &IdentityService{}

func (is *IdentityService) GetPluginInfo(ctx context.Context, req *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	log.Infof("GetPluginInfo")

	return &csi.GetPluginInfoResponse{
		Name:          driverName,
		VendorVersion: driverVersion,
	}, nil
}

func (is *IdentityService) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	log.Infof("GetPluginCapabilities")

	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
			{
				Type: &csi.PluginCapability_VolumeExpansion_{
					VolumeExpansion: &csi.PluginCapability_VolumeExpansion{
						Type: csi.PluginCapability_VolumeExpansion_ONLINE,
					},
				},
			},
		},
	}, nil
}

func (is *IdentityService) Probe(ctx context.Context, req *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	log.Infof("Probe")

	if len(is.MountPoints) == 0 {
		log.Warn("Probe: No mount points configured")
		return &csi.ProbeResponse{
			Ready: &wrappers.BoolValue{Value: false},
		}, nil
	}

	// Validate actual mount paths exist
	for _, mp := range is.MountPoints {
		if mp == nil || mp.hostMountPath == "" {
			log.Warn("Probe: Skipping invalid mount point reference")
			return &csi.ProbeResponse{
				Ready: &wrappers.BoolValue{Value: false},
			}, nil
		}
		log.Debugf("Probe: Checking mount point: %s", mp.hostMountPath)
		if _, err := os.Stat(mp.hostMountPath); err != nil {
			log.Warnf("Probe: Mount path does not exist or is inaccessible: %s", mp.hostMountPath)
			return &csi.ProbeResponse{
				Ready: &wrappers.BoolValue{Value: false},
			}, nil
		}
	}

	// All mount points were valid
	return &csi.ProbeResponse{
		Ready: &wrappers.BoolValue{Value: true},
	}, nil
}
