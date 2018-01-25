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
var rxChan = make(chan string, 1024)
var txChan = make(chan string, 1024)
var wsRxChan = make(chan string, 1024)
var wsTxChan = make(chan string, 1024)

func connectSerial()  {

  //create a serial object
  serConn := &serial.Config{Name: "/dev/cu.usbserial", Baud: 9600, ReadTimeout: time.Millisecond * 25}
  //serConn, err := serial.Open()
  var err error
  s, err := serial.OpenPort(serConn)
  if err != nil {
          log.Fatal(err)
          fmt.Println("Error opeing serial port...")
  } else {
    fmt.Println("Serial device connected..")
  }

  go func() {
    for {
      time.Sleep(50 * time.Millisecond)
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
        recv,err := serial.ReadBytes('\x0a')
        if err != nil {
          //fmt.Printf("Rx Error: %s\n", err)
        }
        if len(string(recv)) > 0 {
          //fmt.Printf("%s", string(recv))
        }
        rxChan <- string(recv)
      }
  }()
//  defer s.Close()

}

var upgrader = websocket.Upgrader{
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
}

//Main entry point
func main() {

  connectSerial()
  txChan <- "\n"

  go func() {
    for{
      select {
        case tx := <-wsRxChan:
            txChan <- tx

          case r := <- rxChan:
            wsTxChan <- r
            //fmt.Printf("%s", r)

          case <- time.After(4000 * time.Millisecond):
            fmt.Println("timeout\n")
            //fmt.Println("\n->")
      }
    }
  }()
  RunServer()

  //mux := bone.New()

  // //create the http endpoints
  // mux.Get("/", http.HandlerFunc(Home))
  // mux.Get("/ws", http.HandlerFunc(WebSocket))
  //
  // //create a middleware
  // n:=negroni.Classic()
  // n.UseHandler(mux)
  // n.Run(":9000")
}

func RunServer() {
  //create a multiplexer
  mux := bone.New()

  //create the http endpoints
  mux.Get("/", http.HandlerFunc(Home))
  mux.Get("/ws", http.HandlerFunc(WebSocket))

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

  for {
    time.Sleep(25 * time.Millisecond)

    _, message, err := ws.ReadMessage()
    if err != nil {
        log.Println("error: ", err)
        break
    }
    if len(string(message)) > 0 {
      fmt.Printf("Ws rx: %s\n", message)
      wsRxChan <- string(message)
    }
    t := <- wsTxChan
    fmt.Printf("Rx: %s", t)
    ws.WriteMessage(websocket.TextMessage, []byte(t))
  }
  defer ws.Close()

}
