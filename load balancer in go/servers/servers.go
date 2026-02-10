package servers

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type ServerList struct {
	Ports []int
}

func (s *ServerList) Populate(amount int) {
	if amount >= 10 { // FIXED: changed 'for' to 'if' to avoid infinite loop
		log.Fatal("bzaf les port lhadj")
	}

	for x := 0; x < amount; x++ {
		s.Ports = append(s.Ports, x)
	}
}

func (s *ServerList) Pop() int {
	if len(s.Ports) == 0 { // FIXED: added bounds checking to prevent panic
		log.Fatal("No servers available")
	}
	port := s.Ports[0]
	s.Ports = s.Ports[1:]
	return port
}

func RunServers(amount int) {
	// server list object
	var MyserverList ServerList
	MyserverList.Populate(amount)
	//wait group to assure concurency

	wg := sync.WaitGroup{} // add the amount of servers to run to the wait group
	wg.Add(amount)         // adds the servers and runs them in the goroutines
	defer wg.Wait()        // run servers when the add function is done

	for x := 0; x < amount; x++ {
		go makeservers(&MyserverList, &wg) // FIXED: pass WaitGroup by pointer
	}

}

func makeservers(sl *ServerList, wg *sync.WaitGroup) { // FIXED: accept WaitGroup by pointer
	r := http.NewServeMux()
	defer wg.Done()  // make sure  the server is done
	port := sl.Pop() // get the port to run the server on
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Server %d", port) // FIXED: added missing port argument
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":808%d", port),
		Handler: r,
	}
	server.ListenAndServe()
}
