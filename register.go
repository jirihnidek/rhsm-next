package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/jirihnidek/rhsm2"
	"github.com/rs/zerolog/log"
)

//go:embed com.redhat.RHSM2.Register.xml
var InterfaceRegister string

// RHSM2Register is structure for methods implementing interface for com.redhat.RHSM2.Register
type RHSM2Register struct{}

// RegisterWithUsername tries to register system with org, username and password
// When registration is successful, then string with JSON representing consumer is
// returned
func (rhsm2Register *RHSM2Register) RegisterWithUsername(
	org string,
	username string,
	password string,
	options map[string]string,
	locale string,
	dbusSender dbus.Sender,
) (string, *dbus.Error) {
	log.Debug().Msgf("RegisterWithUsername() called with arguments: org: %s, username: %s, password: %s, options: %s, locale: %s, sender: %s",
		org, username, password, options, locale, dbusSender)

	rhsmClient, err := rhsm2.CreateRHSMClient(nil)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}

	consumerData, err := rhsmClient.RegisterUsernamePasswordOrg(&username, &password, &org)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}

	data, err := json.Marshal(consumerData)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}

	return string(data), nil
}

// RegisterWithActivationKeys tries to register system with org and activation keys
func (rhsm2Register *RHSM2Register) RegisterWithActivationKeys(
	org string,
	activationKeys []string,
	options map[string]string,
	locale string,
	dbusSender dbus.Sender,
) (string, *dbus.Error) {
	log.Debug().Msgf("RegisterWithActivationKeys() called with arguments: org: %s, activation_keys: %s, options: %s, locale: %s, sender: %s",
		org, activationKeys, options, locale, dbusSender)

	rhsmClient, err := rhsm2.CreateRHSMClient(nil)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}

	consumerData, err := rhsmClient.RegisterOrgActivationKeys(&org, activationKeys)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}

	data, err := json.Marshal(consumerData)
	if err != nil {
		return "", dbus.MakeFailedError(err)
	}

	return string(data), nil
}

// GetOrgs tries to get list of organizations which is user member of
func (rhsm2Register *RHSM2Register) GetOrgs(
	username string,
	password string,
	locale string,
	dbusSender dbus.Sender,
) ([]string, *dbus.Error) {
	log.Debug().Msgf("GetOrgs() calles with arguments: username: %s, password: %s, locale: %s, sender: %s",
		username, password, locale, dbusSender)
	var orgs []string
	var organizations []rhsm2.OrganizationData

	rhsmClient, err := rhsm2.CreateRHSMClient(nil)
	if err != nil {
		return orgs, dbus.MakeFailedError(err)
	}

	organizations, err = rhsmClient.GetOrgs(username, password)
	if err != nil {
		return orgs, dbus.MakeFailedError(err)
	}

	for _, organization := range organizations {
		org, err := json.Marshal(organization)
		if err != nil {
			log.Error().Msgf("unable to marshal organization JSON objet to text: %s", err)
			continue
		}
		orgs = append(orgs, string(org))
	}

	return orgs, nil
}

// initRegisterInterface tries to initialize D-Bus interface for
// com.redhat.RHSM2.Register interface
func initRegisterInterface(conn *dbus.Conn) error {
	rhsm2Register := RHSM2Register{}

	err := conn.Export(
		&rhsm2Register,
		"/com/redhat/RHSM2/Register",
		"com.redhat.RHSM2.Register")
	if err != nil {
		return fmt.Errorf("unable to export D-Bus interface 'com.redhat.RHSM.Register': %v", err)
	}

	err = conn.Export(
		introspect.Introspectable(InterfaceRegister),
		"/com/redhat/RHSM2/Register",
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		return fmt.Errorf("unable to export D-Bus interface 'org.freedesktop.DBus.Introspectable':%v", err)
	}

	log.Debug().Msg("interface 'com.redhat.RHSM2.Register' created")
	return nil
}
