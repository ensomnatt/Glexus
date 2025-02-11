package wsserver

import (
	"glexus/backend/internal/config"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Server struct {
  r *http.ServeMux
  srv *http.Server
  wsUpgrader *websocket.Upgrader
  clients *clients
  broadcast chan *wsMessage
  config *config.Config
}

type clients struct {
  mu *sync.Mutex
  wsClients map[*websocket.Conn]struct{}
  names map[string]string
}

func NewServer(addr string) *Server {
  r := http.NewServeMux()
  return &Server{
    r: r,
    srv: &http.Server{
      Addr: addr,
      Handler: r,
    },
    wsUpgrader: &websocket.Upgrader{},
    clients: &clients{
      mu: &sync.Mutex{},
      wsClients: map[*websocket.Conn]struct{}{},
      names: map[string]string{},
    },
    broadcast: make(chan *wsMessage),
    config: config.NewConfig(),
  }
}

func (s *Server) Start() error {
  s.r.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/src"))))
  s.r.Handle("/home/", http.StripPrefix("/home/", http.FileServer(http.Dir("/home"))))
  s.r.HandleFunc("/", s.watch)
  s.r.HandleFunc("/ws", s.wsHandler)
  s.r.HandleFunc("GET /api/videofiles", s.sendVideoFiles)
  go s.writeToTheClients()
  return s.srv.ListenAndServe()
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
  conn, err := s.wsUpgrader.Upgrade(w, r, nil)
  if err != nil {
    logrus.Errorf("error with creating connection: %v", err)
    w.WriteHeader(http.StatusBadRequest)
    return
  }
  logrus.Infof("new connection from: %s", conn.RemoteAddr().String())
  s.clients.mu.Lock()
  s.clients.wsClients[conn] = struct{}{}
  s.clients.names[conn.RemoteAddr().String()] = "arbuzik"
  s.clients.mu.Unlock()

  go s.updateUserList()
  go s.readFromClient(conn)
}

func (s *Server) readFromClient(conn *websocket.Conn) {
  for {
    msg := new(wsMessage)
    err := conn.ReadJSON(msg)
    if err != nil {
      wsErr, ok := err.(*websocket.CloseError)
      if !ok || wsErr.Code != websocket.CloseGoingAway {
        logrus.Errorf("failed to read message: %v", err)
      }
      break
    }

    msg.IP = conn.RemoteAddr().String()
    msg.Time = time.Now().Format("20:35")
    s.broadcast <- msg
  }
  s.clients.mu.Lock()
  delete(s.clients.wsClients, conn)
  delete(s.clients.names, conn.RemoteAddr().String())
  s.clients.mu.Unlock()

  s.updateUserList()
  logrus.Infof("user from %s has been disconneted", conn.RemoteAddr().String())
}

func (s *Server) writeToTheClients() {
  for msg := range s.broadcast {
    s.clients.mu.Lock()
    for client := range s.clients.wsClients {
      go func() {
        err := client.WriteJSON(msg)
        if err != nil {
          logrus.Errorf("error with writing json: %v", err)
        }
      }()
    }
    s.clients.mu.Unlock()
  }
  logrus.Info("sent to the clients all messages")
}

func (s *Server) updateUserList() {
  var names []string
  for _, v := range s.clients.names {
    names = append(names, v)
  }

  msg := &wsMessage{
    Action: "updateUsers",
    UserNames: names,
  }
  s.broadcast <- msg
}
