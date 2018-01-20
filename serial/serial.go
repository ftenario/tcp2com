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
    ser_port = cfg.Section("params").Key("SerialPort").String()
    baud = cfg.Section("params").Key("Baud").String()
  }

}

func Open() {

  c := &serial.Config{Name: "/dev/cu.usbserial", Baud: 9600}
  s, err := serial.OpenPort(c)
  if err != nil {
    log.Fatal(err)
  }

  n, err := s.Write([]byte("\n\n\n\n\n"))
  if err != nil {
    log.Fatal(err)
  }

  buf := make([]byte, 128)
  n, err = s.Read(buf)
  if err != nil {
    log.Fatal(err)
  }

  log.Println("%q", buf[:n])


  defer s.Close()
}

func GetPorts() {

}

func SendCmd() {

}
