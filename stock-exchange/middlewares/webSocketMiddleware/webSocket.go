package webSocketMiddleware

import (
  "github.com/gorilla/websocket"
  "log"
  "sync"
)

type WebSocket struct {
  connections map[*websocket.Conn]bool
  broadcast   chan []byte
}

var once sync.Once

var server *WebSocket

func GetWebSocket() *WebSocket {
  once.Do(func() {
    server = &WebSocket{
      connections: make(map[*websocket.Conn]bool),
      broadcast:   make(chan []byte),
    }

    go func() {
      for {
        select {
        case message := <-server.broadcast:
          for connection := range server.connections {
            if err := connection.WriteMessage(1, message); err != nil {
              log.Println("falha ao escrever mensagem para o cliente")
            }
          }
        }
      }
    }()
  })

  return server
}

func (webSocket *WebSocket) AddConnection(connection *websocket.Conn) {
  webSocket.connections[connection] = true
}

func (webSocket *WebSocket) RemoveConnection(connection *websocket.Conn) {
  delete(webSocket.connections, connection)
}

func (webSocket *WebSocket) Broadcast(message []byte) {
  webSocket.broadcast <- message
}
