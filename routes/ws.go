package routes

import (
  "net/http"
  "log"
  "encoding/json"
  "bytes"
  "errors"

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
    send chan []byte
    userID uint
}

type Hub struct {
  clients       map[*Client]bool
  broadcast     chan []byte
  addClient     chan *Client
  removeClient  chan *Client
}

var hub = Hub{
  broadcast:    make(chan []byte),
  addClient:    make(chan *Client),
  removeClient: make(chan *Client),
  clients:      make(map[*Client]bool),
}

func (h *Hub) start() {
  for {
        // one of these fires when a channel
        // receives data
        select {
        case conn := <-hub.addClient:
            // add a new client
            hub.clients[conn] = true
            log.Printf("%v", *conn)
        case conn := <-hub.removeClient:
            // remove a client
            if _, ok := hub.clients[conn]; ok {
                delete(hub.clients, conn)
                close(conn.send)
            }
        case message := <-hub.broadcast:

            // broadcast a message to all clients
            for conn := range hub.clients {
                select {
                case conn.send <- message:
                default:
                    close(conn.send)
                    delete(hub.clients, conn)
                }
            }
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
                log.Printf("%v", message)
                return
            }
            log.Printf("%v", message)
            c.ws.WriteMessage(websocket.TextMessage, message)
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
        log.Printf("%v", message)
        if err != nil {
            hub.removeClient <- c
            c.ws.Close()
            break
        }
        var (
            msg models.Message
            to models.User
            m map[string]interface{}
            ok bool
        )
        decoder := json.NewDecoder(bytes.NewReader(message))
        decoder.Decode(&m)
        if models.DB.Where("username = ?", m["to"]).First(&to).RecordNotFound() {
            err := errors.New("Recipient not found.")
            log.Fatal(err)
            break
        }
        msg.To = to.ID
        msg.From = c.userID
        msg.Text, ok = m["text"].(string)
        if !ok {
            log.Fatal(errors.New("Message empty"))
        }

        log.Printf("%v", msg)
        models.DB.Create(&msg)
        hub.broadcast <- message
    }
}

func connectionHandler(w http.ResponseWriter, r *http.Request) {
    var (
        u models.User
    )
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
    if err != nil {
        log.Fatal(err)
        return
    }

    client := &Client{
        ws: ws,
        send: make(chan []byte),
        userID: u.ID,
    }

    hub.addClient <- client
    go client.write()
    go client.read()
}

func addWsRoutes(r *mux.Router) *mux.Router {
    r.Handle("/ws", /*use(*/http.HandlerFunc(connectionHandler)/*, authenticate)*/)
    go hub.start()
    return r
}
