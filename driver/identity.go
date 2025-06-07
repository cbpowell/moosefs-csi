/*
   Copyright (c) 2023 Saglabs SA. All Rights Reserved.

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
	Service Service
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

	ready := true

	switch svc := is.Service.(type) {
	case *NodeService:
		if len(svc.mountPoints) == 0 {
			log.Warn("Probe - no mount points found in NodeService")
			ready = false
		} else {
			for _, mp := range svc.mountPoints {
				log.Debugf("Checking Stat of mount point: %s", mp.hostMountPath)
				if _, err := os.Stat(mp.hostMountPath); err != nil {
					log.Warnf("Probe - mount point not ready: %s (err: %v)", mp.hostMountPath, err)
					ready = false
					break
				}
			}
		}

	case *ControllerService:
		if svc.ctlMount == nil {
			log.Warn("Probe - ControllerService ctlMount is nil")
			ready = false
		} else if _, err := os.Stat(svc.ctlMount.hostMountPath); err != nil {
			log.Warnf("Probe - ctlMount not ready: %s (err: %v)", svc.ctlMount.hostMountPath, err)
			ready = false
		}

	default:
		log.Warnf("Probe - unknown service type %T; assuming not ready", svc)
		ready = false
	}

	return &csi.ProbeResponse{
		Ready: &wrappers.BoolValue{Value: ready},
	}, nil
}
