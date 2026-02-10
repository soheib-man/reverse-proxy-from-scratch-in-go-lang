package loadbalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

var (
	baseURL = "http://localhost:808"
)

type loadBalancer struct {
	RevProxy httputil.ReverseProxy
}

type Endpoint struct {
	List []*url.URL
}

func (e *Endpoint) shuffle() {
	tmp := e.List[0]
	e.List = e.List[1:]
	e.List = append(e.List, tmp)
}

func MakeLoadBalancer(amount int) {
	// initiate objects

	var lb loadBalancer
	var ep Endpoint
	// the SERVER and ROUTER
	router := http.NewServeMux()
	server := http.Server{
		Addr:    ":8090",
		Handler: router,
	}
	//creating the endpoints

	for i := 0; i < amount; i++ {
		ep.List = append(ep.List, createEndpoint(baseURL, i))
	}

	router.HandleFunc("/loadbalancer", makeRequest(&lb, &ep))
	// listen and serve
	log.Fatal(server.ListenAndServe())
}

func makeRequest(lb *loadBalancer, ep *Endpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		lb.RevProxy = *httputil.NewSingleHostReverseProxy(ep.List[0])
		ep.shuffle()
		lb.RevProxy.ServeHTTP(w, r)
	}
}

func createEndpoint(endpoint string, idx int) *url.URL {
	link := endpoint + strconv.Itoa(idx)
	url, err := url.Parse(link)
	if err != nil {
		log.Fatal(err)
	}
	return url
}
