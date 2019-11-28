package algoCR

import (
	"math"
	"strconv"
	"fmt"
	"time"
)

type siteChannel chan<- string

var sitesChannels  map[int]*chan<- string

type algoCR struct {
	id int
	n        int         // nombre de processus
	h	     int         // valeur courante de l'estampille
	demCours bool        // faux/vrai quand le processus moi demande l'accès en SC
	sc		 bool        // faux/vrai quand le processus moi est en SC
	hDem	 int         // estampille de soumission de cette demande
	pDiff  map[int]bool      // ensemble des numéros des processus pour lesquels on a différé le OK
	pAtt   map[int]bool         // ensemble des processus desquels je dois obtenir une permission
	askingSC chan bool
}

func New() algoCR {
	acr := algoCR{0, 0, 0, false, false, 0, make(map[int]bool), make(map[int]bool), make(chan bool)}
	sitesChannels = make(map[int]*chan<- string)
	return acr
}

// todo : trouver une meilleure facon d'initialiser
// avec l'id et lancer la goroutine Wait a la creation du ACR
func (acr *algoCR) SetId(id int) {
	acr.id = id
	go acr.WaitSC()
}

func (acr *algoCR) AddChannel(ch chan<- string, id int) {
	sitesChannels[id] = &ch
	acr.pAtt[id] = true
}

func (acr *algoCR) Ask() {
	acr.h = acr.h + 1
	acr.demCours = true
	acr.hDem = acr.h
	
	for i := range acr.pAtt {
		acr.Req(i, acr.id)
	}
	
	acr.CheckSC()
}

// goroutine attente passive sur SC
func (acr *algoCR) WaitSC() {
	
	for {
		<- acr.askingSC
		if len(acr.pAtt) == 0 { // pAtt vide
			
			//SC
			//Do something...
			acr.sc = true
			fmt.Println("ENTER SC*****************************")
			time.Sleep(5 * time.Second)
			fmt.Println("\n\n\n\n\nLEAVE SC*****************************")
			
			acr.h = acr.h + 1
			acr.sc = false
			acr.demCours = false
			
			for i := range acr.pDiff {
				acr.Ok(i, acr.id)
			}
			
			acr.pAtt = acr.pDiff
			for j := range acr.pDiff {
				delete(acr.pDiff, j)
			}
		}
	}
}

func (acr *algoCR) CheckSC() {
	acr.askingSC <- true
}

// OK
func (acr *algoCR) Ok(idTo int, idFrom int) {
	msg := "O" + strconv.Itoa(acr.h) + strconv.Itoa(idFrom)
	acr.SendMsg(*sitesChannels[idTo], msg)
}

// REQ
func (acr *algoCR) Req(idTo int, idFrom int) {
	msg := "R" + strconv.Itoa(acr.hDem) + strconv.Itoa(idFrom)
	acr.SendMsg(*sitesChannels[idTo], msg)
}

func (acr *algoCR) SendMsg(msgChannel chan<- string, msg string) {
	msgChannel <- msg
	fmt.Println("sent : " + msg)
}

func (acr *algoCR) MsgHandle(msg string) {

	op    := msg[0] // op : R ou O
	hi, _ := strconv.Atoi(string(msg[1]))
	i, _  := strconv.Atoi(string(msg[2]))
	
	// mise a jour de l'estampille
	acr.h = int(math.Max(float64(acr.h) , float64(hi))) + 1

	// check le op 
	if op == 'R' {
	
		if !acr.demCours {
			acr.pAtt[i] = true
			acr.Ok(i, acr.id)
		} else {
			
			if acr.sc || (acr.hDem < hi) || (acr.hDem == hi && acr.id < i) {
				acr.pDiff[i] = true
			} else {
				acr.Ok(i, acr.id)
				acr.pAtt[i] = true
				acr.Req(i, acr.id)
			}
		}
		
	} else if op == 'O' {
		delete(acr.pAtt, i)
		acr.CheckSC()
	}
}


