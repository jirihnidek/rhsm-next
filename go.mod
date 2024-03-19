module rhsm-next

go 1.19

require (
	github.com/fsnotify/fsnotify v1.7.0
	github.com/godbus/dbus/v5 v5.1.0
	github.com/jirihnidek/rhsm2 v0.0.0-20240318122427-35b0a55b6778
	github.com/rs/zerolog v1.32.0
)

require (
	github.com/creasty/defaults v1.7.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/henvic/httpretty v0.1.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/jirihnidek/rhsm2 => /home/jhnidek/github/jirihnidek/rhsm2
