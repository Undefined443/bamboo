package main

import (
	"flag"
	"sync"

	"github.com/gitferry/zeitgeber"
	"github.com/gitferry/zeitgeber/config"
	"github.com/gitferry/zeitgeber/identity"
	"github.com/gitferry/zeitgeber/log"
)

var algorithm = flag.String("algorithm", "pacemaker", "synchronizer algorithm")
var id = flag.String("id", "", "NodeID of the node")
var simulation = flag.Bool("sim", false, "simulation mode")
var isByz = flag.Bool("isByz", false, "this is a Byzantine node")

func replica(id identity.NodeID, isByz bool) {
	log.Infof("node %v starting...", id)
	if isByz {
		log.Infof("node %v is Byzantine", id)
	}

	r := zeitgeber.NewReplica(id, isByz)

	switch *algorithm {
	case "pacemaker":
		r.Pacemaker = zeitgeber.NewBcb(r.Node, r.Election)
	default:
		r.Pacemaker = zeitgeber.NewBcb(r.Node, r.Election)
	}
	if r.IsLeader(id, 1) {
		log.Debugf("[%v] should kick off", id)
		go r.MakeProposal(1)
	}
	r.Run()
}

func main() {
	zeitgeber.Init()

	if *simulation {
		var wg sync.WaitGroup
		wg.Add(1)
		config.Simulation()
		for id := range config.GetConfig().Addrs {
			isByz := false
			if id.Node() <= config.GetConfig().ByzNo {
				isByz = true
			}
			n := id
			go replica(n, isByz)
		}
		wg.Wait()
	} else {
		replica(identity.NodeID(*id), *isByz)
	}
}