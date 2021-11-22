package spoon

import (
	"net/http"
	"strconv"
	"sync"
)

func Run() {
	a := &app{
		clients:      make(map[*clientConnection]struct{}),
		clientsMutex: &sync.Mutex{},
	}

	http.HandleFunc("/ws/", clientHandle(a))
	http.HandleFunc("/api/push", a.push)
	http.HandleFunc("/api/begin", a.begin)
	http.HandleFunc("/api/end", a.end)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.ListenAndServe("127.0.0.1:80", nil)
}

type app struct {
	clients      map[*clientConnection]struct{}
	clientsMutex *sync.Mutex
}

func (a *app) push(w http.ResponseWriter, req *http.Request) {

	name := req.URL.Query().Get("name")
	if name == "" {
		return
	}

	group := req.URL.Query().Get("group")
	if group == "" {
		return
	}
	start, err := strconv.ParseInt(req.URL.Query().Get("start"), 10, 64)
	if err != nil {
		return
	}

	end, err := strconv.ParseInt(req.URL.Query().Get("end"), 10, 64)
	if err != nil {
		return
	}

	a.pushToClient(&PushMsg{
		Name:  name,
		Group: group,
		Start: int64(start),
		End:   end,
	})
}

func (a *app) begin(w http.ResponseWriter, req *http.Request) {

}

func (a *app) end(w http.ResponseWriter, req *http.Request) {

}

func (a *app) clientConnected(c *clientConnection) {
	a.clientsMutex.Lock()
	defer a.clientsMutex.Unlock()

	a.clients[c] = struct{}{}
}

func (a *app) clientDisconnected(c *clientConnection) {
	a.clientsMutex.Lock()
	defer a.clientsMutex.Unlock()

	delete(a.clients, c)
}

func (a *app) pushToClient(msg *PushMsg) {
	a.clientsMutex.Lock()
	defer a.clientsMutex.Unlock()

	msg.Command = "push"
	for c, _ := range a.clients {
		c.ws.WriteJSON(msg)
	}
}
