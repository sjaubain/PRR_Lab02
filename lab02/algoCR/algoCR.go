package algoCR

import (
	"math"
	"strconv"
)

type siteChannel chan<- string

type algoCR struct {
	id int
	n        int         // nombre de processus
	h	     int         // valeur courante de l'estampille
	demCours bool        // faux/vrai quand le processus moi demande l'accès en SC
	sc		 bool        // faux/vrai quand le processus moi est en SC
	hDem	 int         // estampille de soumission de cette demande
	pDiff  map[int]bool      // ensemble des numéros des processus pour lesquels on a différé le OK
	pAtt   map[int]bool         // ensemble des processus desquels je dois obtenir une permission
	sitesChannels map[int]siteChannel
	askingSC chan bool
}

func New() algoCR {
	acr := algoCR{0, 0, 0, false, false, 0, make(map[int]bool), make(map[int]bool), make(map[int]siteChannel), make(chan bool)}
	go acr.WaitSC()
	return acr
}

func (acr *algoCR) SetId(id int) {
	acr.id = id
}

func (acr *algoCR) AddChannel(ch chan<- string, id int) {
	acr.sitesChannels[id] = ch
	acr.pAtt[id] = true
}

func (acr algoCR) Ask() {
	acr.h = acr.h + 1
	acr.demCours = true
	acr.hDem = acr.h
	
	for i := range acr.pAtt {
		acr.Req(i)
	}
}

// goroutine attente passive sur SC
func (acr algoCR) WaitSC() {
	
	for {
		<- acr.askingSC
		if len(acr.pAtt) == 0 { // pAtt vide
			
			//SC
			//Do something...
		}
	}
}

// OK
func (acr algoCR) Ok(i int) {
	msg := "O" + strconv.Itoa(acr.h) + strconv.Itoa(i)
	acr.SendMsg(acr.sitesChannels[i], msg)
}

// REQ
func (acr algoCR) Req(i int) {
	msg := "R" + strconv.Itoa(acr.hDem) + strconv.Itoa(i)
	acr.SendMsg(acr.sitesChannels[i], msg)
}

func (acr algoCR) SendMsg(msgChannel chan <- string, msg string) {
	msgChannel <- msg
}

func (acr algoCR) MsgHandle(msg string) {

	op := int(msg[0]) // op : R ou O
	hi := int(msg[1])
	i  := int(msg[2])
	
	// mise a jour de l'estampille
	acr.h = int(math.Max(float64(acr.h) , float64(hi))) + 1

	// check le op 
	if op == 'R' {
	
		if !acr.demCours {
			acr.Ok(i)
			acr.pAtt[i] = true
		} else {
			
			if acr.sc || (acr.hDem < hi) || (acr.hDem == hi && acr.id < i) {
				acr.pDiff[i] = true
			} else {
				acr.Ok(i)
				acr.pAtt[i] = true
				acr.Req(i)
			}
		}
		
	} else if op == 'O' {
	}
}


