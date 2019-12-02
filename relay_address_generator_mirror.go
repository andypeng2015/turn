package turn

import (
	"fmt"
	"net"
	"strconv"
)

// RelayAddressGeneratorMirror returns the IP Address that the request used as the destination
// This can be useful for a 1:1 NAT like AWS, you can listen on a private IP Address but your Public IP Address will be used
type RelayAddressGeneratorMirror struct {
	// Network, Address are the arguments passed to ListenPacket
	Address string
}

// Validate is caled on server startup and confirms the RelayAddressGenerator is properly configured
func (r *RelayAddressGeneratorMirror) Validate() error {
	switch {
	case r.Address == "":
		return errListeningAddressInvalid
	default:
		return nil
	}
}

// AllocatePacketConn generates a new PacketConn to receive traffic on and the IP/Port to populate the allocation response with
func (r *RelayAddressGeneratorMirror) AllocatePacketConn(network string, dst net.IP, requestedPort int) (net.PacketConn, net.Addr, error) {
	conn, err := net.ListenPacket(network, r.Address+":"+strconv.Itoa(requestedPort))
	if err != nil {
		return nil, nil, err
	}

	relayAddr := conn.LocalAddr().(*net.UDPAddr)
	relayAddr.IP = dst

	return conn, relayAddr, nil
}

// ReadDestinationAddress if we should poll the Dst address from the request
func (r *RelayAddressGeneratorMirror) ReadDestinationAddress() bool {
	return true
}

// AllocateConn generates a new Conn to receive traffic on and the IP/Port to populate the allocation response with
func (r *RelayAddressGeneratorMirror) AllocateConn(network string, dst net.IP, requestedPort int) (net.Conn, net.Addr, error) {
	return nil, nil, fmt.Errorf("TODO")
}
