package main

import (
  //"fmt"
  "net/http"
  "github.com/go-zoo/bone"
  "github.com/codegangsta/negroni"
  "github.com/gorilla/websocket"
  "log"
  //"os"
  //"tcp2com/serial"
  "github.com/tarm/serial"
  "fmt"
  "time"
  "bufio"
)
var rxChan = make(chan string)
var txChan = make(chan string)

func ConnectSerial(txChan chan string, rxChan chan string) {

  //create a serial object
  serConn := &serial.Config{Name: "/dev/cu.usbserial", Baud: 9600, ReadTimeout: time.Millisecond * 25}
  //serConn, err := serial.Open()
  s, err := serial.OpenPort(serConn)
  if err != nil {
          log.Fatal(err)
          fmt.Println("Error opeing serial port...")
  } else {
    fmt.Println("Serial device connected..")
  }
  defer s.Close()

  go func() {
    for {
      time.Sleep(time.Millisecond)
      m := <- txChan
      fmt.Printf("Sending: %s\n", m)
      s.Write([]byte(m))
    }
  }()

  //go routine to read data from serial port
  go func() {
    serial := bufio.NewReader(s)
    for {
        time.Sleep(50 * time.Millisecond)

        //read until newline
        recv,_ := serial.ReadBytes('\x0a')
        if len(string(recv)) > 0 {
          fmt.Printf("rx: %s\n", string(recv))
        }
        rxChan <- string(recv)
      }
  }()
}

var upgrader = websocket.Upgrader{
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
}

//Main entry point
func main() {

  ConnectSerial(txChan, rxChan)

  //create a multiplexer
  mux := bone.New()

  //create the http endpoints
  mux.Get("/", http.HandlerFunc(Home))
  mux.Get("/ws", http.HandlerFunc(WebSocket))
  //mux.Get("/ws", http.HandlerFunc(WebSocket))

  //create a middleware
  n:=negroni.Classic()
  n.UseHandler(mux)
  n.Run(":9000")

}

/*
  Home - handler for serving the index.html file
*/
func Home(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.Error(w, "Method not found", 404)
    return
  }

  if r.Method != "GET" {
    http.Error(w, "Method not allowed", 405)
    return
  }

  http.ServeFile(w, r, "index.html")
}

/*
  WebSocket - handler for the for websocket
*/
func WebSocket(w http.ResponseWriter, r *http.Request) {
  ws, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    log.Println("upgrade: ", err)
    return
  }

  defer ws.Close()

  for {
    _, message, err := ws.ReadMessage()
    if err != nil {
        log.Println("error: ", err)
        break
    }

    log.Printf("%s", message)
    if string(message) != "" {
      txChan <- "\n"
    }
    msg := <- rxChan
    //log.Printf("%s", msg)
    if len(msg) > 0 {
      ws.WriteMessage(websocket.TextMessage, []byte(msg))
    }

  }


}
