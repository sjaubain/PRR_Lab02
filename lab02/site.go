/**
site.go
Author: Simon Jobin, Robel Teklehaimanot
Date  : 04.12.2019
Goal  : it's the part client and network, connect with others site. Listen, send message ...
**/

package main

import (
	"PRR_Lab02/lab02/algoCR"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

type Conf struct {
	NB_SITES   int
	SITES_ADDR []string
}

type siteChannel chan<- string

var (
	siteId      int
	conf        Conf
	connectedTo []bool
	connecting  = make(chan siteChannel)
	acr         = algoCR.New()
	newSite     = make(chan bool, 1) // to protect concurrent acces to connectedTo tab
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
		siteId, _ = strconv.Atoi(os.Args[1])
		if !(0 <= siteId && siteId <= conf.NB_SITES) {
			log.Println("invalid site id")
			return
		}
	}

	// set id and start main loop of acr
	acr.Start(siteId)

	// initialize chan newSite as an open mutex
	newSite <- true

	go lookUp()
	go listen()

	// wait for all sites to be running (connected)
	for i := 0; i < conf.NB_SITES-1; i++ {
		<-connecting
	}

	fmt.Println("\nall sites connected, now able to accept user commands")
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nEnter text: [R (Read) | W (Write)]\n")

		// read the input of the user
		cmd, _ := reader.ReadString('\n')

		// if W, the site do an ask
		if cmd == "W\n" {
			acr.Ask()
		} else if cmd == "R\n" {
			fmt.Println("value is " + strconv.Itoa(acr.GetValue()))
		} else {
			fmt.Println("unknown command " + cmd)
		}
	}
}

func loadConfiguration() {
	file, _ := os.Open("conf/conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	_ = decoder.Decode(&conf)
}

// try to connect to each site
func lookUp() {
	for id, _ := range conf.SITES_ADDR {

		// should not connect to myself
		<-newSite
		if id != siteId {
			connectToSite(id)
		}
		newSite <- true
	}
}

func listen() {

	// listen all site
	listener, _ := net.Listen("tcp", conf.SITES_ADDR[siteId])

	for {
		conn, _ := listener.Accept()

		// receive id
		buf := make([]byte, 256)
		_, _ = conn.Read(buf)
		id, _ := strconv.Atoi(string(buf[0]))

		<-newSite
		if !connectedTo[id] {
			connectToSite(id)
		}
		newSite <- true

		go reader(conn, id)
	}
}

func connectToSite(id int) {

	conn, err := net.Dial("tcp", conf.SITES_ADDR[id])

	if err == nil {
		connectedTo[id] = true

		log.Println("i am connected to site " + strconv.Itoa(id))

		// send its id
		fmt.Fprintln(conn, strconv.Itoa(siteId))

		writer(conn, id)
	}
}

func writer(conn net.Conn, id int) {

	ch := make(chan string)
	acr.AddChannel(ch, id)
	connecting <- ch

	go func() {
		for msg := range ch {
			fmt.Fprintln(conn, msg)
		}
	}()
}

func reader(conn net.Conn, id int) {

	input := bufio.NewScanner(conn)
	for input.Scan() {
		if input.Text() != "" {
			acr.MsgHandle(input.Text())
		}
	}
}
