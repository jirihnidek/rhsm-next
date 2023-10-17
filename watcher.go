package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/godbus/dbus/v5"
	"github.com/rs/zerolog/log"
)

// handleWatcher tries to handle watcher events
func handleWatcher(rhsm2Connection *RHSM2Consumer, watcher *fsnotify.Watcher, conn *dbus.Conn) error {
	for {
		log.Debug().Msg("handling watcher events...")
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			log.Debug().Msgf("event: %s, file: %s\n", event.Op.String(), event.Name)

			// Reload data only in the case, when consumer certificate was changed,
			// because we read data only form consumer certificate
			if event.Name == rhsm2Connection.consumerCertFileName {
				err := rhsm2Connection.readConsumerData()
				if err != nil {
					log.Debug().Msgf("error during reading data: %s", err)
				}
			}

			// Send D-Bus signal only in the case, when consumer cert or key was changed
			if event.Name == rhsm2Connection.consumerKeyFileName || event.Name == rhsm2Connection.consumerCertFileName {
				log.Debug().Msg("emitting signal com.redhat.RHSM2Consumer.Consumer.ConsumerChanged ...")
				err := conn.Emit(
					"/com/redhat/RHSM2Consumer/Consumer",
					"com.redhat.RHSM2Consumer.Consumer.ConsumerChanged",
					event.Name,
					event.Op.String())
				if err != nil {
					return err
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Error().Msgf("error: %s", err)
		}
	}
}

// initConsumerWatcher tries to initialize and trigger inotify watcher
// watching directory with consumer certificate and key
func initConsumerWatcher(conn *dbus.Conn, rhsm2Consumer *RHSM2Consumer) error {
	// Create watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.Error().Msgf("unable to close watcher: %v", err)
		}
	}(watcher)

	go func() {
		err := handleWatcher(rhsm2Consumer, watcher, conn)
		if err != nil {
			log.Error().Msgf("watcher terminated with error: %v", err)
		}
	}()

	// Add a path.
	consumerCertDirPath := "/etc/pki/consumer"
	log.Debug().Msgf("monitoring: %s", consumerCertDirPath)
	err = watcher.Add(consumerCertDirPath)
	if err != nil {
		return err
	}
	return nil
}
