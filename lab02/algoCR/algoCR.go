package algoCR

import (
	"fmt"
	"math"
	"strconv"
	"strings"
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
	endAsk   chan bool
	value    int
}

func New() algoCR {
	acr := algoCR{0, 0, 0, false, false, 0, make(map[int]bool), make(map[int]bool), make(chan bool),make (chan bool), 0}
	sitesChannels = make(map[int]*chan<- string)
	return acr
}

// todo : trouver une meilleure facon d'initialiser
// avec l'id et lancer la goroutine Wait a la creation du ACR
func (acr *algoCR) Start(id int) {
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
	<- acr.endAsk
}

func (acr *algoCR) SetInputValue(){
	fmt.Println("Old Value : " + strconv.Itoa(acr.value) + " Please enter the new value :")
	for{
		_ , err := fmt.Scan(&acr.value)

		if err != nil {
			fmt.Println("Enter a number :")
		} else {
			fmt.Println("New Value : " + strconv.Itoa(acr.value))
			break
		}
	}
}
func (acr *algoCR) SetValue(v int){
	acr.value = v
}

func (acr *algoCR) GetValue() int{
	return acr.value
}

// goroutine attente passive sur SC
func (acr *algoCR) WaitSC() {

	for {
		<- acr.askingSC
		if len(acr.pAtt) == 0 { // pAtt vide
			
			//SC
			//Do something...
			acr.sc = true
			fmt.Println("\n\n===================== ENTER SC =====================\n\n")
			acr.SetInputValue()
			time.Sleep(5 * time.Second)
			fmt.Println("\n\n===================== LEAVE SC =====================\n\n")

			for i := range sitesChannels {
				acr.ChangeValueSites(i, acr.id, acr.value)
			}
			
			for i := range acr.pDiff {
				acr.Ok(i, acr.id)
			}
			
			// pAtt = pDiff, first clear pAtt then set pDiff finally clear pDiff
			for j := range acr.pAtt {
				delete(acr.pAtt, j)
			}
			
			for j := range acr.pDiff {
				acr.pAtt[j] = true
				delete(acr.pDiff, j)
			}

			acr.h = acr.h + 1
			acr.sc = false
			acr.demCours = false
			
			// si on le fait pas, il revient au main et l input sera incorrect dans le main
			acr.endAsk <- true
		}
	}
}

func (acr *algoCR) CheckSC() {
	acr.askingSC <- true
}

// convertit l'heure locale en un string de 4 digit
// pour la transmission 
func (acr *algoCR) intToString(h int) string {
	var ret string
	ret += strconv.Itoa(h / 1000)
	ret += strconv.Itoa(h % 1000 / 100)
	ret += strconv.Itoa(h % 100 / 10)
	ret += strconv.Itoa(h % 10)
	return ret
}

//Changement de valeur, il faut l'annoncer au autre site pour qu'ils puissent la changer
func (acr *algoCR) ChangeValueSites(idTo int, idFrom int, value int ) {
	msg := "V" + strconv.Itoa(value)
	acr.SendMsg(*sitesChannels[idTo],msg)
}

// OK
func (acr *algoCR) Ok(idTo int, idFrom int) {
	msg := "O" + acr.intToString(acr.h) + strconv.Itoa(idFrom)
	acr.SendMsg(*sitesChannels[idTo], msg)
}

// REQ
func (acr *algoCR) Req(idTo int, idFrom int) {
	msg := "R" + acr.intToString(acr.hDem) + strconv.Itoa(idFrom)
	acr.SendMsg(*sitesChannels[idTo], msg)
}

func (acr *algoCR) SendMsg(msgChannel chan<- string, msg string) {
	msgChannel <- msg
	//fmt.Println("sent : " + msg)
}

func (acr *algoCR) MsgHandle(msg string) {

	// extrait les parametres du message
	op := msg[0] // op : R ou O ou V

	// check le op
	if op == 'R' || op == 'O'{
	
		hi, _ := strconv.Atoi(string(msg[1:5]))
		i, _  := strconv.Atoi(string(msg[5]))

		// mise a jour de l'estampille
		acr.h = int(math.Max(float64(acr.h) , float64(hi))) + 1
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
	} else if op == 'V' {
		val,_ := strconv.Atoi(strings.Trim(msg, "V"))
		acr.SetValue(val)
	}
}


