package client

import (
	"math/rand"
	"net"
	"sort"
	"strings"
	"sync"
)

type serverlist struct {
	e  endpoints
	mu sync.RWMutex
}

func newServerList() *serverlist {
	return &serverlist{}
}

// set the server list to a new list. The new list will be shuffled and sorted
// by priority.
func (s *serverlist) set(newe endpoints) {
	s.mu.Lock()
	s.e = newe
	s.mu.Unlock()
}

// all returns a copy of the full server list, shuffled and then sorted by
// priority
func (s *serverlist) all() endpoints {
	s.mu.RLock()
	out := make(endpoints, len(s.e))
	copy(out, s.e)
	s.mu.RUnlock()

	// Randomize the order
	for i, j := range rand.Perm(len(out)) {
		out[i], out[j] = out[j], out[i]
	}

	// Sort by priority
	sort.Sort(out)
	return out
}

// failed servers get deprioritized
func (s *serverlist) failed(e *endpoint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < len(s.e); i++ {
		if s.e[i].equal(e) {
			e.priority++
			return
		}
	}
}

// good servers get promoted to the highest priority
func (s *serverlist) good(e *endpoint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := 0; i < len(s.e); i++ {
		if s.e[i].equal(e) {
			e.priority = 0
			return
		}
	}
}

func (e endpoints) Len() int {
	return len(e)
}

func (e endpoints) Less(i int, j int) bool {
	// Sort only by priority as endpoints should be shuffled and ordered
	// only by priority
	return e[i].priority < e[j].priority
}

func (e endpoints) Swap(i int, j int) {
	e[i], e[j] = e[j], e[i]
}

type endpoints []*endpoint

func (e endpoints) String() string {
	names := make([]string, 0, len(e))
	for _, endpoint := range e {
		names = append(names, endpoint.name)
	}
	return strings.Join(names, ",")
}

type endpoint struct {
	name string
	addr net.Addr

	// 0 being the highest priority
	priority int
}

// equal returns true if the name and addr match between two endpoints.
// Priority is ignored because the same endpoint may be added by discovery and
// heartbeating with different priorities.
func (e *endpoint) equal(o *endpoint) bool {
	return e.name == o.name && e.addr == o.addr
}
