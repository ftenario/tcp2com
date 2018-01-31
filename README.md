# Tcp2Com

A web/serial application DEMO using Golang and WebSocket.

This is still a work-in-progress

## TO DO LIST
* use parameter.ini for serial config
* Add multiple ports in parameters.ini
* Add bootstrap css for the UI
* Add TCP port for automation
* Use css Tab for multiple port logs

## Building for different platforms:

For Linux (arm):
Build:
```
$ end GOOS=linux GOARCH=arm go build goterm
```

For MacOS:
Build:
```
$ go build goterm
```

For Windows:
Build:
```
$ end GOOS=windows go build goterm
```
