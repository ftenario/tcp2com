package main

import (
  //"fmt"
  "net/http"
  "github.com/go-zoo/bone"
  "github.com/codegangsta/negroni"
  "github.com/gorilla/websocket"
  "log"
  //"os"
  Conn "tcp2com/connection"
  "github.com/tarm/serial"
  "fmt"
  "time"
  "bufio"
)

var rxChan = make(chan string)
var txChan = make(chan string)

type readOp struct {
  resp chan string
}

type writeOp struct {
  resp chan string
}

// var rxChan = make(chan *readOp)
// var txChan = make(chan *writeOp)

func connectSerial()  {

  //create a serial parameter
  ser := Conn.Serial{}
  //initialize the parameter. This read the pramaeters.ini file and populate the struct
  ser.Init()

  //create a serial object
  serConn := &serial.Config{Name: ser.GetPort(), Baud: ser.GetBaudRate(), ReadTimeout: time.Millisecond * 25}
  //serConn := &serial.Config{Name: "/dev/cu.usbserial", Baud: 9600, ReadTimeout: time.Millisecond * 25}
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
      //fmt.Printf("Sending: %s", m)
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
        //recv,err := serial.ReadString('\n')
        if len(string(recv)) > 0 {
          //fmt.Printf("%s", recv)
          rxChan <- string(recv)
        }
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

  //connect to serial automaticaly. This needs to be change
  //to allow user to select serial port to connect to
  connectSerial()
  txChan <- "\n"

  mux := bone.New()
  //create the http endpoints
  mux.Get("/", http.HandlerFunc(Home))
  mux.Get("/ws", http.HandlerFunc(WebSocket))
  mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
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
  wsSendMsg - Send message to the WebSocket
  params - ws connection
*/
func wsSendMsg(ws *websocket.Conn) {
  for {
      select {

        case m := <- rxChan:
          ws.WriteMessage(websocket.TextMessage, []byte(m))

        case <- time.After(60000 * time.Millisecond):
          m := "console timeout...\n"
          ws.WriteMessage(websocket.TextMessage, []byte(m))
          txChan <- "exit\n"

      }
  }
}

/*
  wsGetMasg - Read message from the websocket
  params - ws connection
*/
func wsGetMsg(ws *websocket.Conn) {
  for {
    time.Sleep(50 * time.Millisecond)
    _, msg, err := ws.ReadMessage()

    if err != nil {
      fmt.Printf("wsGetMsg Error: %s\n", err)
      break;
    }
    //fmt.Printf("%d\n", len(string(msg)))
    txChan <- string(msg)
  }
  defer ws.Close()
}

/*
  WebSocket - handler for the websocket
*/
func WebSocket(w http.ResponseWriter, r *http.Request) {
  ws, err := upgrader.Upgrade(w, r, nil)

  if err != nil {
    log.Println("Upgrade Error: ", err)
    return
  }

  //create goroutines for incoming/outgoing websocket
  go wsSendMsg(ws)
  go wsGetMsg(ws)
}
