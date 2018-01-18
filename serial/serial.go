package serial

import (
  "github.com/jacobsa/go-serial/serial"
  "log"
)

options := serial.Options{
  PortName: "/dev/tty.USB",
  BaudRate: 19200,
  DataBits: 8,
  StopBits: 1,
  MinimumReadSize: 4,
}

func Open() {
  port, err := serial.Open(options)

  if err != nil {
    log.Fatal("serial.Open: %v\n", err)
  }

  defer port.Close()
}
