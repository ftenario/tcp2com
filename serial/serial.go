package serial

import (
  //"github.com/jacobsa/go-serial/serial"
  "github.com/tarm/serial"
  "github.com/go-ini/ini"
  "log"
  "os"
  //"strconv"
)

// options := serial.options{
//   portname: "/dev/tty.usb",
//   baudrate: 19200,
//   databits: 8,
//   stopbits: 1,
//   minimumreadsize: 4,
// }
var ser_port = ""
var baud = ""

func init() {
  if _,err := os.Stat("parameters.ini"); err == nil {
    cfg,_ := ini.Load("parameters.ini")
    ser_port = cfg.Section("port1").Key("serial").String()
    baud = cfg.Section("port1").Key("Baud").String()
  }

}

func Open() *serial {

  c := &serial.Config{Name: "/dev/cu.usbserial", Baud: 9600}
  s, err := serial.OpenPort(c)
  if err != nil {
    log.Fatal(err)
  }
  defer s.Close()
  return s
}

func GetPorts() {

}

func SendCmd() {

}
