package node

import (
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/gitferry/bamboo/config"
	"github.com/gitferry/bamboo/log"
	"github.com/gitferry/bamboo/message"
)

// http request header names
const (
	HTTPClientID  = "Id"
	HTTPCommandID = "Cid"
)

var totalRecevReq int
var totalCommitReq int

var ppFree = sync.Pool{
	New: func() interface{} {
		return make(chan message.TransactionReply, 1)
	},
}

// serve serves the http REST API request from clients
func (n *node) http() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", n.handleRoot)

	// http string should be in form of ":8080"
	ip, err := url.Parse(config.Configuration.HTTPAddrs[n.id])
	if err != nil {
		log.Fatal("http url parse error: ", err)
	}
	port := ":" + ip.Port()
	n.server = &http.Server{
		Addr:    port,
		Handler: mux,
	}
	log.Info("http server starting on ", port)
	log.Fatal(n.server.ListenAndServe())
}

func (n *node) handleRoot(w http.ResponseWriter, r *http.Request) {
	var req message.Transaction
	defer r.Body.Close()
	var err error

	//req.C = ppFree.Get().(chan message.TransactionReply)
	req.C = make(chan message.TransactionReply, 1)
	req.NodeID = n.id // TODO does this work when forward twice
	req.ID = r.RequestURI
	req.HasBroadcast = false
	req.HasReplied = false
	n.TxChan <- req

	reply := <-req.C

	if reply.Err != nil {
		http.Error(w, reply.Err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = io.WriteString(w, string(reply.Value))
	if err != nil {
		log.Error(err)
	}
}
