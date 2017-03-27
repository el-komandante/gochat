package routes

import (
  "net/http"
  "log"
  "encoding/json"
  "bytes"
  "errors"
  "unicode/utf8"

  "github.com/el-komandante/gochat/models"
  "github.com/gorilla/websocket"
  "github.com/gorilla/mux"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
            return true
        },
}

type Client struct {
    ws *websocket.Conn
    // Hub passes broadcast messages to this channel
    send chan models.Message
    userID uint
}

type Hub struct {
    clients       map[uint]*Client
    broadcast     chan models.Message
    addClient     chan *Client
    removeClient  chan *Client
}

var hub = Hub{
    broadcast:    make(chan models.Message),
    addClient:    make(chan *Client),
    removeClient: make(chan *Client),
    clients:      make(map[uint]*Client),
}

func (h *Hub) start() {
  for {
        // one of these fires when a channel
        // receives data
        select {
        case client := <-hub.addClient:
            // add a new client
            hub.clients[client.userID] = client
        case client := <-hub.removeClient:
            // remove a client
            if _, ok := hub.clients[client.userID]; ok {
                delete(hub.clients, client.userID)
                close(client.send)
            }
        case message := <-hub.broadcast:
            log.Printf("message at 60: %+v", message)

            // broadcast a message to all clients
            log.Printf("%v", hub.clients)
            client, ok := hub.clients[message.To]
            if !ok {
                log.Fatal("Client not found")
            }
            log.Printf("message at 64: %v", message)
            select {
                case client.send <- message:
                    // log.Printf("to id: %v %t, from id: %v %t", client.userID, client.userID, message.From, message.From)
                    log.Printf("%v %v", message.From, message.To)
                    if message.From != message.To {
                        receiver, ok := hub.clients[message.From]
                        if !ok {
                            log.Fatal("Client not found")
                        }
                        receiver.send <- message
                    }
                default:
                    close(client.send)
                    delete(hub.clients, client.userID)
            }
            // for conn := range hub.clients {
            //     select {
            //     case conn.send <- message:
            //     default:
            //         close(conn.send)
            //         delete(hub.clients, conn)
            //     }
            // }
        }
    }
}


// Hub broadcasts a new message and this fires
func (c *Client) write() {
    // make sure to close the connection incase the loop exits
    defer func() {
        c.ws.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            if !ok {
                c.ws.WriteMessage(websocket.CloseMessage, []byte{})
                // log.Printf("%v", message)
                return
            }
            var jsonBytes bytes.Buffer
            err := json.NewEncoder(&jsonBytes).Encode(&message)
            if err != nil {
                panic(err)
            }
            // jsonBytes := make([]byte, len(jsonString))
            c.ws.WriteMessage(websocket.TextMessage, jsonBytes.Bytes())
        }
    }
}

// New message received so pass it to the Hub
func (c *Client) read() {
    defer func() {
        hub.removeClient <- c
        c.ws.Close()
    }()

    for {
        _, message, err := c.ws.ReadMessage()
        if err != nil {
            hub.removeClient <- c
            c.ws.Close()
            break
        }
        var (
            msg models.Message
            to models.User
        )
        err = json.NewDecoder(bytes.NewReader(message)).Decode(&msg)
        if err != nil {
            panic(err)
        }
        log.Printf("message at 137 = %+v", msg)
        if utf8.RuneCountInString(msg.Text) < 1 {
            log.Fatal(errors.New("Message empty"))
        }
        if models.DB.Where("ID = ?", msg.To).First(&to).RecordNotFound() {
            err := errors.New("Recipient not found.")
            log.Fatal(err)
            break
        }
        msg.To = to.ID
        msg.From = c.userID
        log.Printf("message at 148: %+v", msg)
        models.DB.Create(&msg)
        log.Printf("message at 150: %+v", msg)
        hub.broadcast <- msg
    }
}

func connectionHandler(w http.ResponseWriter, r *http.Request) {
    var u models.User

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Origin", "localhost")
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    w.Header().Set("Access-Control-Allow-Headers",
      "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

    cookie, err := r.Cookie("session")
    if err != nil {
        log.Fatal(err)
        return
    }

    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatal(err)
        return
    }

    // s, err = s.FromKey(id)
    // if err != nil {
    //     log.Fatal(err)
    //     return
    // }
    u, err = u.FromSessionID(cookie.Value)
    log.Printf("user at 181: %+v", u)
    if err != nil {
        log.Fatal(err)
        return
    }

    client := &Client{
        ws: ws,
        send: make(chan models.Message),
        userID: u.ID,
    }

    hub.addClient <- client
    go client.write()
    go client.read()
}

func addWsRoutes(r *mux.Router) *mux.Router {
    r.Handle("/ws", use(http.HandlerFunc(connectionHandler), authenticate))
    go hub.start()
    return r
}
