package turn

import (
	"fmt"
	"net"
	"strconv"
)

// RelayAddressGeneratorNone returns the listener with no modifications
type RelayAddressGeneratorNone struct {
	Address string
}

// Validate is caled on server startup and confirms the RelayAddressGenerator is properly configured
func (r *RelayAddressGeneratorNone) Validate() error {
	switch {
	case r.Address == "":
		return errListeningAddressInvalid
	default:
		return nil
	}
}

// AllocatePacketConn generates a new PacketConn to receive traffic on and the IP/Port to populate the allocation response with
func (r *RelayAddressGeneratorNone) AllocatePacketConn(network string, dst net.IP, requestedPort int) (net.PacketConn, net.Addr, error) {
	conn, err := net.ListenPacket(network, r.Address+":"+strconv.Itoa(requestedPort))
	if err != nil {
		return nil, nil, err
	}

	return conn, conn.LocalAddr(), nil
}

// ReadDestinationAddress if we should poll the Dst address from the request
func (r *RelayAddressGeneratorNone) ReadDestinationAddress() bool {
	return true
}

// AllocateConn generates a new Conn to receive traffic on and the IP/Port to populate the allocation response with
func (r *RelayAddressGeneratorNone) AllocateConn(network string, dst net.IP, requestedPort int) (net.Conn, net.Addr, error) {
	return nil, nil, fmt.Errorf("TODO")
}
