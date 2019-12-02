package turn

import (
	"fmt"
	"net"
	"strconv"
)

// RelayAddressGeneratorStatic can be used to return static IP address each time a relay is created.
// This can be used when you have a single static IP address that you want to use
type RelayAddressGeneratorStatic struct {
	// RelayAddress is the IP returned to the user when the relay is created
	RelayAddress net.IP

	// Network, Address are the arguments passed to ListenPacket
	Address string
}

// Validate is caled on server startup and confirms the RelayAddressGenerator is properly configured
func (r *RelayAddressGeneratorStatic) Validate() error {
	switch {
	case r.RelayAddress == nil:
		return errRelayAddressInvalid
	case r.Address == "":
		return errListeningAddressInvalid
	default:
		return nil
	}
}

// AllocatePacketConn generates a new PacketConn to receive traffic on and the IP/Port to populate the allocation response with
func (r *RelayAddressGeneratorStatic) AllocatePacketConn(network string, dst net.IP, requestedPort int) (net.PacketConn, net.Addr, error) {
	conn, err := net.ListenPacket(network, r.Address+":"+strconv.Itoa(requestedPort))
	if err != nil {
		return nil, nil, err
	}

	// Replace actual listening IP with the user requested one of RelayAddressGeneratorStatic
	relayAddr := conn.LocalAddr().(*net.UDPAddr)
	relayAddr.IP = r.RelayAddress

	return conn, relayAddr, nil
}

// ReadDestinationAddress if we should poll the Dst address from the request
func (r *RelayAddressGeneratorStatic) ReadDestinationAddress() bool {
	return false
}

// AllocateConn generates a new Conn to receive traffic on and the IP/Port to populate the allocation response with
func (r *RelayAddressGeneratorStatic) AllocateConn(network string, dst net.IP, requestedPort int) (net.Conn, net.Addr, error) {
	return nil, nil, fmt.Errorf("TODO")
}
