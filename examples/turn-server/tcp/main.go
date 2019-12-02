package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/pion/turn"
)

// Allow one user 'foo' with the password 'bar'
func authCallback(username string, srcAddr net.Addr) (string, bool) {
	if username == "foo" {
		return "bar", true
	}

	return "", false
}

func main() {
	var err error
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	tcpListener, err := net.Listen("tcp4", "0.0.0.0:3478")
	if err != nil {
		log.Panicf("Failed to create TURN server listener: %s", err)
	}

	s, err := turn.NewServer(turn.ServerConfig{
		Realm:       "pion.ly",
		AuthHandler: authCallback,
		ListenerConfigs: []turn.ListenerConfig{
			{
				Listener: tcpListener,
				RelayAddressGenerator: &turn.RelayAddressGeneratorStatic{
					RelayAddress: net.ParseIP("127.0.0.1"),
					Address:      "0.0.0.0",
				},
			},
		},
	})
	if err != nil {
		log.Panic(err)
	}

	<-sigs
	if err = s.Close(); err != nil {
		log.Panic(err)
	}
}
