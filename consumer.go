package main

import (
	_ "embed"
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/jirihnidek/rhsm2"
	"github.com/rs/zerolog/log"
)

//go:embed com.redhat.RHSM2.Consumer.xml
var InterfaceConsumer string

type RHSM2Consumer struct {
	consumerCertFileName string
	consumerKeyFileName  string
	uuid                 string
	orgID                string
}

func (rhsm2Consumer *RHSM2Consumer) GetUuid(locale string, dbusSender dbus.Sender) (string, *dbus.Error) {
	log.Debug().Msgf("GetUuid() called with locale argument: \"%s\" from sender: %s\n", locale, dbusSender)
	err := rhsm2Consumer.readConsumerData()
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return rhsm2Consumer.uuid, nil
}

func (rhsm2Consumer *RHSM2Consumer) GetOrg(locale string, dbusSender dbus.Sender) (string, *dbus.Error) {
	log.Debug().Msgf("GetOrg() called with locale argument: \"%s\" from sender: %s\n", locale, dbusSender)
	err := rhsm2Consumer.readConsumerData()
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}
	return rhsm2Consumer.orgID, nil
}

// readConsumerData tries to read data from consumer certificate
// and set fields in rhsm2 structure
func (rhsm2Consumer *RHSM2Consumer) readConsumerData() error {
	rhsmClient, err := rhsm2.GetRHSMClient(nil)
	if err != nil {
		return fmt.Errorf("unable to create rhsm client: %s", err)
	}

	log.Debug().Msgf("reloading consumer certificate...")
	consumerUUID, err := rhsmClient.GetConsumerUUID()
	if err != nil {
		rhsm2Consumer.uuid = ""
	} else {
		rhsm2Consumer.uuid = *consumerUUID
	}

	ownerId, err := rhsmClient.GetOwner()
	if err != nil {
		rhsm2Consumer.orgID = ""
	} else {
		rhsm2Consumer.orgID = *ownerId
	}

	return nil
}

// initConsumerInterface tries to initialize D-Bus interface for
// com.redhat.RHSM2Consumer.Consumer interface
func initConsumerInterface(conn *dbus.Conn) error {
	rhsm2Consumer := RHSM2Consumer{
		consumerCertFileName: "",
		consumerKeyFileName:  "",
		uuid:                 "",
		orgID:                "",
	}

	err := rhsm2Consumer.readConsumerData()
	if err != nil {
		return err
	}

	err = conn.Export(
		&rhsm2Consumer,
		"/com/redhat/RHSM2/Consumer",
		"com.redhat.RHSM2.Consumer")
	if err != nil {
		return fmt.Errorf("unable to export D-Bus interface 'com.redhat.RHSM2.Consumer': %v", err)
	}

	err = conn.Export(
		introspect.Introspectable(InterfaceConsumer),
		"/com/redhat/RHSM2/Consumer",
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return fmt.Errorf("unable to export D-Bus interface 'org.freedesktop.DBus.Introspectable':%v", err)
	}

	err = initConsumerWatcher(conn, &rhsm2Consumer)
	if err != nil {
		return err
	}

	log.Debug().Msgf("interface 'com.redhat.RHSM2.Consumer' created")

	return nil
}
