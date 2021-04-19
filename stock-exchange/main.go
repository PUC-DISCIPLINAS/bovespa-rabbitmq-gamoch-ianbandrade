package main

import (
  "encoding/json"
  "github.com/gorilla/websocket"
  "html/template"
  "log"
  "net/http"
  "stock-exchange/data/app"
  "stock-exchange/middlewares/webSocketMiddleware"
  "stock-exchange/services"
)

func main() {
  serveMux := http.NewServeMux()

  serveMux.HandleFunc("/", pageHandler)
  serveMux.HandleFunc("/events", eventsHandler)

  log.Fatal(http.ListenAndServe(":80", serveMux))
}

func pageHandler(w http.ResponseWriter, _ *http.Request) {
  htmlTemplate := template.Must(template.ParseFiles("web/index.gohtml"))
  appData := app.GetAppData()

  if err := htmlTemplate.Execute(w, appData); err != nil {
    log.Println("falha ao renderizar o template HTML")
  }
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
  upgrader := websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
  }

  connection, err := upgrader.Upgrade(w, r, nil)
  if err != nil {
    http.Error(w, "falha ao realizar o upgrade para WebSocket", http.StatusInternalServerError)
    return
  }

  webSocket := webSocketMiddleware.GetWebSocket()
  webSocket.AddConnection(connection)

  for {
    _, messageBody, err := connection.ReadMessage()
    if err != nil {
      webSocket.RemoveConnection(connection)
      break
    }

    var pageData app.PageData
    if err = json.Unmarshal(messageBody, &pageData); err != nil {
      log.Println("falha ao desserializar os dados da página")
      break
    }

    appData := app.GetAppData()
    appData.Data = services.AppHandler(pageData)

    pageDataJSON, err := json.Marshal(appData.Data)
    if err != nil {
      log.Println("falha ao serializar os dados da página")
      break
    }

    webSocket.Broadcast(pageDataJSON)
  }
}
