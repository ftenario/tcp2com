package serial

import (
  "github.com/jacobsa/go-serial/serial"
  "github.com/go-ini/ini"
  "log"
  "os"
)

options := serial.Options{
  PortName: "/dev/tty.USB",
  BaudRate: 19200,
  DataBits: 8,
  StopBits: 1,
  MinimumReadSize: 4,
}

func Init() {
  var ser_port = ""
  var baud = ""
  var data = ""
  var stop = ""
  if _,err := os.Stat("parameters.ini"); err == nil {
    cfg,_ := ini.Load("parameters.ini")
    ser_port = cfg.Section("parameters").Key("SerialPort").string()
    baud = cfg.Section("parameters").Key("Baud").string()
    data = cfg.Section("parameters").Key("DataBits").string()
    stop = cfg.Section("parameters").Key("StopBits").string()
  }

}

func Open() {

  port, err := serial.Open(options)

  if err != nil {
    log.Fatal("serial.Open: %v\n", err)
  }

  defer port.Close()
}

func GetPorts() {

}

func SendCmd() {

}
