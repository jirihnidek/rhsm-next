package main

import (
	_ "embed"
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/jirihnidek/rhsm2"
	"github.com/rs/zerolog/log"
)

func main() {
	// Try to connect to system bus
	conn, err := dbus.ConnectSystemBus()
	if err != nil {
		panic(err)
	} else {
		log.Debug().Msg("connected to system bus")
	}
	defer func(conn *dbus.Conn) {
		err := conn.Close()
		if err != nil {
			log.Error().Msgf("unable to close D-Bus connection: %v", err)
		}
	}(conn)

	rhsm2.SetUserAgentCmd("rhsm2.service")

	err = initConsumerInterface(conn)
	if err != nil {
		panic(err)
	}

	err = initRegisterInterface(conn)
	if err != nil {
		panic(err)
	}

	err = initUnregisterInterface(conn)
	if err != nil {
		panic(err)
	}

	reply, err := conn.RequestName(
		"com.redhat.RHSM2",
		dbus.NameFlagDoNotQueue)
	if err != nil {
		panic(fmt.Errorf("unable to request D-Bus name 'com.redhat.RHSM2': %s", err))
	}
	if reply != dbus.RequestNameReplyPrimaryOwner {
		panic(fmt.Errorf("name 'com.redhat.RHSM2' already taken"))
	}

	// "Block forever"
	<-make(chan struct{})
}
