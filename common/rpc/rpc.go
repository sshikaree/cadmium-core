package rpc

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"

	"github.com/sshikaree/cadmium-core/common/database"
	"github.com/sshikaree/wswrapper"

	"github.com/gorilla/websocket"
	// "golang.org/x/net/websocket"
)

const (
	SOCKET_NAME = "/tmp/cadmium.socket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		// TODO: Think about this hack !!! Is it dangerous?? Maybe try to avoid
		// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
		CheckOrigin: func(r *http.Request) bool { // to allow cross-origin requests
			return true
		},
	}
)

type Response struct {
	Result json.RawMessage `json:"result"`
	Error  error           `json:"error"`
	Id     string          `json:"id"`
}

func (r *Response) ToJson() ([]byte, error) {
	return json.Marshal(r)
}

type ResultRoster struct {
	From     string    `json:"from"`
	Contacts []Contact `json:"contacts"`
}

type ResultMessage struct{}
type ResultPresence struct{}

type DatabaseMethods struct{}

func (d *DatabaseMethods) CreateAccount(acc *database.Account, reply *int) error {
	return database.DB.CreateAccount(acc)
}
func (d *DatabaseMethods) UpdateAccount(acc *database.Account, reply *int) error {
	return database.DB.UpdateAccount(acc)
}

type CommonMethods struct{}

// dummy func
func (c *CommonMethods) SetStatus(status string, reply *string) error {
	log.Println("Set status:", status)
	*reply = "Status changed"
	return nil
}

func init() {
	// Register remote procedures
	if err := rpc.Register(new(DatabaseMethods)); err != nil {
		log.Fatal(err)
	}
	if err := rpc.Register(new(CommonMethods)); err != nil {
		log.Fatal(err)
	}
	if err := rpc.Register(new(XMPP)); err != nil {
		log.Fatal(err)
	}

}

func StartWebsocketRPCServer(port string) {
	// http.Handle("/ws/v1", websocket.Handler(serveWS))
	http.HandleFunc("/ws/v1", wsHandler)
	log.Println("Starting WS server on port:", port)
	http.ListenAndServe(":"+port, nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	wrapped_conn := wswrapper.NewWsWrapper(conn)
	wsHub.Register(wrapped_conn)
	jsonrpc.ServeConn(wrapped_conn)
	log.Println("Websocket connection closed")
	wsHub.Unregister(wrapped_conn)
	return
}

// func serveWS(ws *websocket.Conn) {
// 	jsonrpc.ServeConn(ws)
// }

func StartLocalSocketRPCServer() {
	os.Remove(SOCKET_NAME)

	ln, err := net.ListenUnix(
		"unix",
		&net.UnixAddr{SOCKET_NAME, "unix"},
	)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(SOCKET_NAME)
	defer ln.Close()

	// run RPC server
	rpc.Accept(ln)

}
