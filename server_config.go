package turn

import (
	"net"
	"time"

	"github.com/pion/logging"
)

// RelayAddressGenerator is used to generate a RelayAddress when creating an allocation.
// You can use one of the provided ones or provide your own.
type RelayAddressGenerator interface {
	// Validate confirms that the RelayAddressGenerator is properly initialized
	Validate() error

	// Should we poll for DestinationAddress, not supported by Go on all platforms
	ReadDestinationAddress() bool

	// Allocate a PacketConn (UDP) RelayAddress
	AllocatePacketConn(network string, dst net.IP, requestedPort int) (net.PacketConn, net.Addr, error)

	// Allocate a Conn (TCP) RelayAddress
	AllocateConn(network string, dst net.IP, requestedPort int) (net.Conn, net.Addr, error)
}

// PacketConnConfig is a single net.PacketConn to listen/write on. This will be used for UDP listeners
type PacketConnConfig struct {
	PacketConn net.PacketConn

	// When an allocation is generated the RelayAddressGenerator
	// creates the net.PacketConn and returns the IP/Port it is available at
	RelayAddressGenerator RelayAddressGenerator
}

func (c *PacketConnConfig) validate() error {
	if c.PacketConn == nil {
		return errConnUnset
	}
	if c.RelayAddressGenerator == nil {
		return errRelayAddressGeneratorUnset
	}

	return nil
}

// ListenerConfig is a single net.Listener to accept connections on. This will be used for TCP, TLS and DTLS listeners
type ListenerConfig struct {
	Listener net.Listener

	// When an allocation is generated the RelayAddressGenerator
	// creates the net.PacketConn and returns the IP/Port it is available at
	RelayAddressGenerator RelayAddressGenerator
}

func (c *ListenerConfig) validate() error {
	if c.Listener == nil {
		return errListenerUnset
	}

	if c.RelayAddressGenerator == nil {
		return errRelayAddressGeneratorUnset
	}

	return nil
}

// AuthHandler is a callback used to handle incoming auth requests, allowing users to customize Pion TURN with custom behavior
type AuthHandler func(username string, srcAddr net.Addr) (password string, ok bool)

// ServerConfig configures the Pion TURN Server
type ServerConfig struct {
	// PacketConnConfigs and ListenerConfigs are a list of all the turn listeners
	// Each listener can have custom behavior around the creation of Relays
	PacketConnConfigs []PacketConnConfig
	ListenerConfigs   []ListenerConfig

	// LoggerFactory must be set for logging from this server.
	LoggerFactory logging.LoggerFactory

	// Realm sets the realm for this server
	Realm string

	// AuthHandler is a callback used to handle incoming auth requests, allowing users to customize Pion TURN with custom behavior
	AuthHandler AuthHandler

	// ChannelBindTimeout sets the lifetime of channel binding. Defaults to 10 minutes.
	ChannelBindTimeout time.Duration
}

func (s *ServerConfig) validate() error {
	if len(s.PacketConnConfigs) == 0 && len(s.ListenerConfigs) == 0 {
		return errNoAvailableConns
	}

	for _, s := range s.PacketConnConfigs {
		if err := s.validate(); err != nil {
			return err
		}
	}

	for _, s := range s.ListenerConfigs {
		if err := s.validate(); err != nil {
			return err
		}
	}

	return nil
}
