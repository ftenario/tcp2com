package main

import (
  //"fmt"
  "net/http"
  "github.com/go-zoo/bone"
  "github.com/codegangsta/negroni"
  "github.com/gorilla/websocket"
  "log"
  "os"
)

var upgrader = websocket.Upgrader{
  ReadBufferSize: 1024,
  WriteBufferSize: 1024,
}

func main() {

  mux := bone.New()
  mux.Get("/", http.HandlerFunc(Home))
  mux.Get("/ws", http.HandlerFunc(WebSocket))

  n:=negroni.Classic()
  n.UseHandler(mux)
  n.Run(":9000")
}

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

    log.Printf("rx: %s", message)

  }


}
