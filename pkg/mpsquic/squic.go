// Copyright 2019 ETH Zurich
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// squic drop in
package mpsquic

import (
	"github.com/lucas-clemente/quic-go"

	"github.com/scionproto/scion/go/lib/addr"
	"github.com/scionproto/scion/go/lib/common"
	"github.com/scionproto/scion/go/lib/snet"
)

// DialSCION creates a single path QUIC connection over SCION.
// It returns a QUIC Session.
func DialSCION(network *snet.SCIONNetwork, laddr, raddr *snet.Addr,
	quicConfig *quic.Config) (quic.Session, error) {

	mpQuic, err := DialMPWithBindSVC(network, laddr, raddr, nil, nil, addr.SvcNone, quicConfig)
	if err != nil {
		return nil, err
	}
	return mpQuic, nil
}

// DialSCION creates a single path QUIC connection over SCION on the bind address baddr.
// It returns a QUIC Session.
func DialSCIONWithBindSVC(network *snet.SCIONNetwork, laddr, raddr, baddr *snet.Addr,
	svc addr.HostSVC, quicConfig *quic.Config) (quic.Session, error) {

	mpQuic, err := DialMPWithBindSVC(network, laddr, raddr, nil, baddr, svc, quicConfig)
	if err != nil {
		return nil, err
	}
	return mpQuic, nil
}

// ListenSCION listens for incoming QUIC connections over SCION addresses to the local SCION address laddr.
// It returns a QUIC Listener.
func ListenSCION(network *snet.SCIONNetwork, laddr *snet.Addr,
	quicConfig *quic.Config) (quic.Listener, error) {

	return ListenSCIONWithBindSVC(network, laddr, nil, addr.SvcNone, quicConfig)
}

// ListenSCION listens for incoming QUIC connections over SCION binding to the SCION address baddr.
// It returns a QUIC Listener.
func ListenSCIONWithBindSVC(network *snet.SCIONNetwork, laddr, baddr *snet.Addr,
	svc addr.HostSVC, quicConfig *quic.Config) (quic.Listener, error) {

	if len(srvTlsCfg.Certificates) == 0 {
		return nil, common.NewBasicError("mpsquic: No server TLS certificate configured", nil)
	}
	sconn, err := sListen(network, laddr, baddr, svc)
	if err != nil {
		return nil, err
	}
	return quic.Listen(sconn, srvTlsCfg, quicConfig)
}

func sListen(network *snet.SCIONNetwork, laddr, baddr *snet.Addr,
	svc addr.HostSVC) (snet.Conn, error) {

	if network == nil {
		network = snet.DefNetwork
	}
	return network.ListenSCIONWithBindSVC("udp4", laddr, baddr, svc, 0)
}
