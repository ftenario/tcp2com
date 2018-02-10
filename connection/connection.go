package connection

import (
  "github.com/go-ini/ini"
  "os"
  "strconv"
  "runtime"
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
    if runtime.GOOS == "darwin" {
      s.port = cfg.Section("darwin").Key("serial1").String()
      s.baud = cfg.Section("darwin").Key("baud1").Value()
    }
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
