# Tcp2Com

A web/serial application DEMO using Golang and WebSocket.
This is similar to a terminal serial application but is
using the web browser as the user interface

This is still a work-in-progress

## TO DO LIST
* use parameter.ini for serial config - done
* Add multiple ports in parameters.ini
* Add bootstrap css for the UI - done
* Add TCP port for automation
* Use css Tab for multiple port logs
* Add devices info in the UI
* Add connect to devices in the UI

## Building for different platforms:

For Linux (arm):
Build:
```
$ end GOOS=linux GOARCH=arm go build tcp2com
```

For MacOS:
Build:
```
$ go build tcp2com
```

For Windows:
Build:
```
$ end GOOS=windows go build tcp2com
```
