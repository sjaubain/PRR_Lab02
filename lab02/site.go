package main

import (

	"encoding/json"
	"log"
	"os"
	"net"
	"strconv"
	"fmt"
)

type Conf struct {
	NB_SITES int
	SITES_ADDR []string
}
	
var (
	site_id int
	conf Conf
	connectedTo []bool
)

/**
 * each time a process is strated, it tries to connect to each other
 * and then listen to incoming connections. If a new site just connect,
 * it connect back to him (to allow bidirectionnal exchange)
 * This is managed whith the boolean tab connectedTo[]
 */
func main() {
	
	// Load configuration from json config file
	loadConfiguration()
	
	// list of sites on wich the process is connected
	connectedTo = make([]bool, conf.NB_SITES)
	
	// parse command line args
	if len(os.Args) == 1 {
		log.Println("you have to provide a site id")
		return
	} else {
		site_id, _ = strconv.Atoi(os.Args[1])
		if !(0 <= site_id && site_id <= conf.NB_SITES) {
			log.Println("invalid site id")
			return
		}	
	}	

	lookUp()
	for {
	}
	// wait for all sites to be running (connected)
	// TODO : applicative goroutine to listen to user commands
}

func loadConfiguration() {
	file, _ := os.Open("conf/conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&conf)
}

// try to connect to each site
func lookUp() {
	for id, site_addr := range conf.SITES_ADDR {
		// should not connect to myself
		if id != site_id {
			conn, err := net.Dial("tcp", site_addr)
			if err == nil {
				connectedTo[id] = true
				
				// send its id
				log.Println("i am connected to site " + strconv.Itoa(id))
				fmt.Fprintln(conn, strconv.Itoa(site_id))
			}
		}
	}
	go listen()
}

func listen() {
	listener, _ := net.Listen("tcp", conf.SITES_ADDR[site_id])
	
	for {
		conn, _ := listener.Accept()
		
		// receive id
		buf := make([]byte, 1) 
		_, _ = conn.Read(buf)
		id, _ := strconv.Atoi(string(buf[0]))

		if !connectedTo[id] {
			_, _ = net.Dial("tcp", conf.SITES_ADDR[id]) // note : ignoring errors
			log.Println("i am connected to site " + strconv.Itoa(id))
			// go handleConn
		}
	}
}

func handleConn(conn net.Conn) {
	
}

/*
func connect(site_addr string) {
	conn, err := net.Dial("tcp", site_addr)
	if err == nil {
		// send id so that all sites listening can
		// connect back
		fmt.Fprintln(conn, strconv.Itoa(site_id))
		go handleConn(conn)
	}
}

func listen() {
	listener, err := net.Listen("tcp", conf.SITES_ADDR[site_id])
	if err != nil {
		log.Fatal(err)
	}
	
	for i := 0; i < conf.NB_SITES; i++ {
		conn, err := listener.Accept()

		// new site connected
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	// receive id
	buf := make([]byte, 256) 
	_, _ = conn.Read(buf)
	id, _ := strconv.Atoi(string(buf[0]))
	
	connectedSites[id] = true
	log.Printf("site %d connected, waiting for following sites :\n%s", id, getAwaitingSites())
	// connect back
	connect(conf.SITES_ADDR[id])
}

func getAwaitingSites() string {
	var ret string = "[ "
	for i := 0; i < conf.NB_SITES; i++ {
		if !connectedSites[i] {ret += strconv.Itoa(i) + " "}
	}
	return ret + "]"
}
*/