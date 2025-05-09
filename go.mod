module rhsm-next

go 1.19

require (
	github.com/fsnotify/fsnotify v1.9.0
	github.com/godbus/dbus/v5 v5.1.0
	github.com/jirihnidek/rhsm2 v0.0.0-20240424135902-3305cdd26a04
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/creasty/defaults v1.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/henvic/httpretty v0.1.3 // indirect
	github.com/jeandeaual/go-locale v0.0.0-20240223122105-ce5225dcaa49 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.19.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
)

//replace github.com/jirihnidek/rhsm2 => /home/jhnidek/github/jirihnidek/rhsm2
