package main

import "loadbalancer/loadbalancer"

func main() {
	//servers.RunServers(5)
	loadbalancer.MakeLoadBalancer(5)

}
