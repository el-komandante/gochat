package routes

// import (
//   "net/http"
//   "log"
//   "github.com/gorilla/websocket"
// )
//
// var upgrader = websocket.Upgrader{
//   ReadBufferSize:  1024,
//   WriteBufferSize: 1024,
//   CheckOrigin: func(r *http.Request) bool {
//         return true
//       },
// }
//
// // type Message struct {
// //   // Email string    `json:email`
// //   // Username string `json:username`
// //   Message string  `json:message`
// // }
//
// type Hub struct {
//   clients       map[*Client]bool
//   broadcast     chan []byte
//   addClient     chan *Client
//   removeClient  chan *Client
// }
//
// var hub = Hub{
//   broadcast:    make(chan []byte),
//   addClient:    make(chan *Client),
//   removeClient: make(chan *Client),
//   clients:      make(map[*Client]bool),
// }
//
// func (h *Hub) start() {
//   for {
//         // one of these fires when a channel
//         // receives data
//         select {
//         case conn := <-hub.addClient:
//             // add a new client
//             hub.clients[conn] = true
//         case conn := <-hub.removeClient:
//             // remove a client
//             if _, ok := hub.clients[conn]; ok {
//                 delete(hub.clients, conn)
//                 close(conn.send)
//             }
//         case message := <-hub.broadcast:
//             // broadcast a message to all clients
//             for conn := range hub.clients {
//                 select {
//                 case conn.send <- message:
//                 default:
//                     close(conn.send)
//                     delete(hub.clients, conn)
//                 }
//             }
//         }
//     }
// }
//
// type Client struct {
//     ws *websocket.Conn
//     // Hub passes broadcast messages to this channel
//     send chan []byte
// }
//
// // Hub broadcasts a new message and this fires
// func (c *Client) write() {
//     // make sure to close the connection incase the loop exits
//     defer func() {
//         c.ws.Close()
//     }()
//
//     for {
//         select {
//         case message, ok := <-c.send:
//             if !ok {
//                 c.ws.WriteMessage(websocket.CloseMessage, []byte{})
//                 return
//             }
//
//             c.ws.WriteMessage(websocket.TextMessage, message)
//         }
//     }
// }
//
// // New message received so pass it to the Hub
// func (c *Client) read() {
//     defer func() {
//         hub.removeClient <- c
//         c.ws.Close()
//     }()
//
//     for {
//         _, message, err := c.ws.ReadMessage()
//         if err != nil {
//             hub.removeClient <- c
//             c.ws.Close()
//             break
//         }
//
//         hub.broadcast <- message
//     }
// }
//
// func connectionHandler(w http.ResponseWriter, r *http.Request) {
//   ws, err := upgrader.Upgrade(w, r, nil)
//   if err != nil {
//     log.Fatal(err)
//     return
//   }
//   client := &Client{
//     ws: ws,
//     send: make(chan []byte),
//   }
//   hub.addClient <- client
//
//   go client.write()
//   go client.read()
// }
//
// func addWsRoutes(w http.ResponseWriter, req *http.Request) {
//   r.HandleFunc("/users", connectionHandler)
//   return r
// }
