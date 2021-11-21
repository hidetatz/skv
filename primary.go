package skv

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

type Primary struct {
	port   int
	server *http.Server
	st     *Store
	// tl              *TransactionLog
	// replicas        []*Replica
	latestOperation uint64

	m sync.Mutex
}

func NewPrimary(port int) *Primary {
	p := &Primary{
		port: port,
		st:   NewStore(),
	}

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}
	mux := http.NewServeMux()
	mux.Handle("/get", p.getHandler())
	mux.Handle("/set", p.setHandler())
	server.Handler = mux
	p.server = server

	return p
}

func (p *Primary) Start() error {
	fmt.Fprintf(os.Stdout, `{"msg":"server started at %d"}`+"\n", p.port)

	// TODO: run background heartbeat
	return p.server.ListenAndServe()
}

func (p *Primary) getHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("key")
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"key must not be empty"}`))
			return
		}
		v, err := p.st.Get(key)
		if err == ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error":"key not found"}`))
			return
		}

		v = strings.ReplaceAll(v, `"`, `\"`)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"value":"%s"}`, v)))
	})
}

func (p *Primary) setHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.FormValue("key")
		val := r.FormValue("value")
		if key == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"key must not be empty"}`))
			return
		}
		p.st.Set(key, val)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"ok"}`))
	})
}
