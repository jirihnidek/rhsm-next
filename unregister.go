package main

import (
	_ "embed"
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/jirihnidek/rhsm2"
	"github.com/rs/zerolog/log"
)

//go:embed com.redhat.RHSM2.Unregister.xml
var InterfaceUnregister string

// RHSM2Unregister is structure for methods implementing interface for com.redhat.RHSM2.Register
type RHSM2Unregister struct{}

// Unregister tries to unregister system
func (rhsm2Unregister *RHSM2Unregister) Unregister(
	locale string,
	dbusSender dbus.Sender,
) *dbus.Error {
	log.Debug().Msgf("Unregister() called with arguments: locale: %s, sender: %s",
		locale, dbusSender)

	rhsmClient, err := rhsm2.GetRHSMClient(nil)
	if err != nil {
		return dbus.MakeFailedError(err)
	}

	clientInfo := rhsm2.ClientInfo{}
	clientInfo.Locale = locale
	clientInfo.DBusSender = string(dbusSender) // TODO: get app name
	err = rhsmClient.Unregister(&clientInfo)
	if err != nil {
		return dbus.MakeFailedError(err)
	}

	return nil
}

// initRegisterInterface tries to initialize D-Bus interface for
// com.redhat.RHSM2.Register interface
func initUnregisterInterface(conn *dbus.Conn) error {
	rhsm2Unregister := RHSM2Unregister{}

	err := conn.Export(
		&rhsm2Unregister,
		"/com/redhat/RHSM2/Unregister",
		"com.redhat.RHSM2.Unregister")
	if err != nil {
		return fmt.Errorf("unable to export D-Bus interface 'com.redhat.RHSM.Unregister': %v", err)
	}

	err = conn.Export(
		introspect.Introspectable(InterfaceUnregister),
		"/com/redhat/RHSM2/Unregister",
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return fmt.Errorf("unable to export D-Bus interface 'org.freedesktop.DBus.Introspectable':%v", err)
	}

	log.Debug().Msg("interface 'com.redhat.RHSM2.Unregister' created")
	return nil
}
