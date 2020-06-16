// Copyright (c) 2020 Cisco Systems, Inc.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package ipaddress provides networkservice chain elements that support setting ip addresses on kernel interfaces
package ipaddress

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/networkservicemesh/api/pkg/api/networkservice"
	"github.com/networkservicemesh/sdk/pkg/networkservice/core/next"

	"github.com/networkservicemesh/sdk-vppagent/pkg/tools/kernelctx"
)

type setIPKernelServer struct{}

// NewServer provides a NetworkServiceServer that sets the IP on a kernel interface
// It sets the IP Address on the *kernel* side of an interface plugged into the
// Endpoint.  Generally only used by privileged Endpoints like those implementing
// the Cross Connect Network Service for K8s (formerly known as NSM Forwarder).
//                                         Endpoint
//                              +---------------------------+
//                              |                           |
//                              |                           |
//                              |                           |
//                              |                           |
//                              |                           |
//                              |                           |
//                              |                           |
//          +-------------------+                           |
//  ipaddress.NewServer()       |                           |
//                              |                           |
//                              |                           |
//                              |                           |
//                              |                           |
//                              |                           |
//                              |                           |
//                              +---------------------------+
//
func NewServer() networkservice.NetworkServiceServer {
	return &setIPKernelServer{}
}

func (s *setIPKernelServer) Request(ctx context.Context, request *networkservice.NetworkServiceRequest) (*networkservice.Connection, error) {
	iface := kernelctx.ServerInterface(ctx)
	if iface != nil {
		srcIP := request.GetConnection().GetContext().GetIpContext().GetSrcIpAddr()
		if srcIP != "" {
			iface.IpAddresses = append(iface.GetIpAddresses(), srcIP)
		}
	}
	return next.Server(ctx).Request(ctx, request)
}

func (s *setIPKernelServer) Close(ctx context.Context, conn *networkservice.Connection) (*empty.Empty, error) {
	iface := kernelctx.ServerInterface(ctx)
	if iface != nil {
		srcIP := conn.GetContext().GetIpContext().GetSrcIpAddr()
		if srcIP != "" {
			iface.IpAddresses = append(iface.GetIpAddresses(), srcIP)
		}
	}
	return next.Server(ctx).Close(ctx, conn)
}
