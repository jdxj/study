package full

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/peer"
	"github.com/btcsuite/btcd/wire"
	"net"
	"time"
)

// mockRemotePeer creates a basic inbound peer listening on the simnet port for
// use with Example_peerConnection.  It does not return until the listner is
// active.
func MockRemotePeer() {
	// Configure peer to act as a simnet node that offers no services.
	peerCfg := &peer.Config{
		UserAgentName:    "peer",  // User agent name to advertise.
		UserAgentVersion: "1.0.0", // User agent version to advertise.
		ChainParams:      &chaincfg.SimNetParams,
		TrickleInterval:  time.Second * 10,
	}

	// Accept connections on the simnet port.
	listener, err := net.Listen("tcp", "127.0.0.1:18555")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Accept: error %v\n", err)
			continue
		}

		go func() {
			// Create and start the inbound peer.
			p := peer.NewInboundPeer(peerCfg)
			p.AssociateConnection(conn)
		}()
	}
}

func Local() {
	// Create an outbound peer that is configured to act as a simnet node
	// that offers no services and has listeners for the version and verack
	// messages.  The verack listener is used here to signal the code below
	// when the handshake has been finished by signalling a channel.
	verack := make(chan struct{})
	peerCfg := &peer.Config{
		UserAgentName:    "peer",  // User agent name to advertise.
		UserAgentVersion: "1.0.0", // User agent version to advertise.
		ChainParams:      &chaincfg.SimNetParams,
		Services:         0,
		TrickleInterval:  time.Second * 10,
		Listeners: peer.MessageListeners{
			OnVersion: func(p *peer.Peer, msg *wire.MsgVersion) *wire.MsgReject {
				fmt.Println("outbound: received version")
				return nil
			},
			OnVerAck: func(p *peer.Peer, msg *wire.MsgVerAck) {
				verack <- struct{}{}
			},
		},
	}
	p, err := peer.NewOutboundPeer(peerCfg, "127.0.0.1:18555")
	if err != nil {
		fmt.Printf("NewOutboundPeer: error %v\n", err)
		return
	}

	// Establish the connection to the peer address and mark it connected.
	conn, err := net.Dial("tcp", p.Addr())
	if err != nil {
		fmt.Printf("net.Dial: error %v\n", err)
		return
	}
	fmt.Println("port:", conn.LocalAddr())
	p.AssociateConnection(conn)

	// Wait for the verack message or timeout in case of failure.
	select {
	case <-verack:
	case <-time.After(time.Second * 10):
		fmt.Printf("Example_peerConnection: verack timeout")
	}

	// Disconnect the peer.
	p.Disconnect()
	p.WaitForDisconnect()

}
