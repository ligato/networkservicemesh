// Copyright (c) 2020 Doc.ai and/or its affiliates.
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

package wireguard

import (
	"fmt"
	"net"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/ipc"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection"
	"github.com/networkservicemesh/networkservicemesh/controlplane/api/connection/mechanisms/wireguard"
)

// DeviceManager manages wireguard devices
type DeviceManager struct {
	devices map[string]*device.Device
	mutex   sync.Mutex
	wrapper func(tun.Device) tun.Device
}

// NewWireguardDeviceManager creates new wireguard manager. Can create wireguard devices with L2 mode
func NewWireguardDeviceManager(allowL2Traffic bool) *DeviceManager {
	r := &DeviceManager{
		devices: make(map[string]*device.Device),
	}
	r.wrapper = func(t tun.Device) tun.Device {
		if allowL2Traffic {
			return &tunL2Adapter{original: t}
		}
		return t
	}
	return r
}

// CreateWireguardDevice creates a new wireguard device
func (w *DeviceManager) CreateWireguardDevice(ifaceName string, remoteConnection *connection.Connection, incoming bool) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	mechanism := wireguard.ToMechanism(remoteConnection.GetMechanism())

	/* Create interface - host namespace */
	var localPrivateKeyStr string
	var remotePublicKeyStr string
	var localPort int
	var remotePort int
	var dstIPStr string
	var err error
	if incoming {
		if localPrivateKeyStr, err = mechanism.DstPrivateKey(); err != nil {
			return err
		}
		if remotePublicKeyStr, err = mechanism.SrcPublicKey(); err != nil {
			return err
		}
		if dstIPStr, err = mechanism.SrcIP(); err != nil {
			return err
		}
		if localPort, err = mechanism.DstPort(); err != nil {
			return err
		}
		if remotePort, err = mechanism.SrcPort(); err != nil {
			return err
		}
	} else {
		if localPrivateKeyStr, err = mechanism.SrcPrivateKey(); err != nil {
			return err
		}
		if remotePublicKeyStr, err = mechanism.DstPublicKey(); err != nil {
			return err
		}
		if dstIPStr, err = mechanism.DstIP(); err != nil {
			return err
		}
		if localPort, err = mechanism.SrcPort(); err != nil {
			return err
		}
		if remotePort, err = mechanism.DstPort(); err != nil {
			return err
		}
	}

	dstIP := net.ParseIP(dstIPStr)
	localPrivateKey, err := wgtypes.ParseKey(localPrivateKeyStr)
	if err != nil {
		return errors.Errorf("failed to parse local private key: %v", err)
	}
	remotePublicKey, err := wgtypes.ParseKey(remotePublicKeyStr)
	if err != nil {
		return errors.Errorf("failed to parse remote public key: %v", err)
	}

	wgDevice, err := w.createWireguardDevice(ifaceName)
	if err != nil {
		return errors.Errorf("Wireguard error: %v", err)
	}

	w.devices[ifaceName] = wgDevice
	api, err := startWireguardAPI(ifaceName, wgDevice)

	if err != nil {
		return errors.Errorf("Wireguard error: %v", err)
	}
	defer func() {
		if uapiErr := api.Close(); uapiErr != nil {
			logrus.Errorf("Wireguard error: failed to close API client %v", uapiErr)
		}
	}()

	err = configureWireguardDevice(ifaceName, localPrivateKey, remotePublicKey, localPort, remotePort, dstIP)
	if err != nil {
		wgDevice.Close()
		return errors.Errorf("Wireguard error: %v", err)
	}
	return nil
}

// DeleteWireguardDevice creates a wireguard device
func (w *DeviceManager) DeleteWireguardDevice(ifaceName string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if wgDevice, ok := w.devices[ifaceName]; ok {
		wgDevice.Close()
		delete(w.devices, ifaceName)
	}
}

func (w *DeviceManager) createWireguardDevice(ifaceName string) (*device.Device, error) {
	tunIface, err := tun.CreateTUN(ifaceName, device.DefaultMTU)
	if err != nil {
		return nil, errors.Errorf("failed to create tun: %v", err)
	}

	logger := device.NewLogger(device.LogLevelDebug, fmt.Sprintf("Wireguard Error (%s): ", ifaceName))
	return device.NewDevice(w.wrapper(tunIface), logger), nil
}

func startWireguardAPI(ifaceName string, wgDevice *device.Device) (net.Listener, error) {
	fileUAPI, err := ipc.UAPIOpen(ifaceName)
	if err != nil {
		return nil, err
	}

	uapi, err := ipc.UAPIListen(ifaceName, fileUAPI)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			conn, err := uapi.Accept()
			if err != nil {
				return
			}
			go wgDevice.IpcHandle(conn)
		}
	}()

	return uapi, nil
}

func configureWireguardDevice(ifaceName string, localPrivateKey, remotePublicKey wgtypes.Key, localPort, remotePort int, dstIP net.IP) error {
	client, err := wgctrl.New()
	if err != nil {
		return errors.Errorf("failed to create configuration client: %v", err)
	}
	defer func() {
		if clientErr := client.Close(); clientErr != nil {
			logrus.Errorf("Wireguard error (%v): failed to close configuration client: %v", ifaceName, clientErr)
		}
	}()

	_, ipnet, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		return errors.Errorf("failed to configure device: %v", err)
	}
	err = client.ConfigureDevice(ifaceName, wgtypes.Config{
		ListenPort: &localPort,
		PrivateKey: &localPrivateKey,
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey: remotePublicKey,
				AllowedIPs: []net.IPNet{
					*ipnet,
				},
				Endpoint: &net.UDPAddr{
					IP:   dstIP,
					Port: remotePort,
				},
			},
		},
	})

	return errors.Wrapf(err, "failed to configure device: %v", err)
}
