package connection

import (
  "github.com/go-ini/ini"
  "os"
  "strconv"
)

// options := serial.options{
//   portname: "/dev/tty.usb",
//   baudrate: 19200,
//   databits: 8,
//   stopbits: 1,
//   minimumreadsize: 4,
// }

type Serial struct{
  port string
  baud string
}

func (s *Serial) Init() {

  if _,err := os.Stat("parameters.ini"); err == nil {
    cfg,_ := ini.Load("parameters.ini")
    s.port = cfg.Section("port1").Key("serial").String()
    s.baud = cfg.Section("port1").Key("baud").Value()
  }
}

func (s Serial) GetPort() string{
  return s.port
}

func (s Serial) GetBaudRate() int {
  r, _ := strconv.Atoi(s.baud)
  return r
}

func GetPorts() {

}

func SendCmd() {

}
