/**
algoCR.go
Author: Simon Jobin, Robel Teklehaimanot
Date  : 04.12.2019
Goal  : it's the part Mutex, use the Carvalho et Roucairol algorithm
**/


package algoCR

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type siteChannel chan<- string

// references to all other process channels
var sitesChannels map[int]*chan<- string

type algoCR struct {
	id       int          // current process id
	n        int          // number of process
	h        int          // current stamp value
	demCours bool         // true / false if I have asked for SC
	sc       bool         // true / false if I am in SC
	hDem     int          // SC asking stamp
	pDiff    map[int]bool // set of process for which we differ the OK
	pAtt     map[int]bool // set of process for which I need to wait the OK
	askingSC chan bool    // channel to communicate between SC goroutine and main context
	endAsk   chan bool    // channel to inform end of SC
	value    int          // global shared value
}

// Construct the objet mutex
func New() algoCR {
	acr := algoCR{0, 0, 0, false, false, 0,
		make(map[int]bool),
		make(map[int]bool),
		make(chan bool),
		make(chan bool), 0}

	sitesChannels = make(map[int]*chan<- string)
	return acr
}


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
	<-acr.endAsk
}

/**
 * function executed in critical section to
 * change global shared value
 */
func (acr *algoCR) SetInputValue() {
	fmt.Println("Old Value : " + strconv.Itoa(acr.value) + "\nPlease enter the new value :")
	for {
		_, err := fmt.Scan(&acr.value)

		if err != nil {
			fmt.Println("\nEnter a number :")
		} else {
			fmt.Println("\nNew Value : " + strconv.Itoa(acr.value))
			break
		}
	}
}

func (acr *algoCR) SetValue(v int) {
	acr.value = v
}

func (acr *algoCR) GetValue() int {
	return acr.value
}

/**
 * goroutine for passive waiting on SC
 */
func (acr *algoCR) WaitSC() {

	for {
		<-acr.askingSC // all site wait here until one of them ask to enter in SC
		if len(acr.pAtt) == 0 {

			acr.sc = true
			fmt.Println("\n\n===================== ENTER SC =====================\n\n")
			acr.SetInputValue()
			time.Sleep(5 * time.Second)
			fmt.Println("\n\n===================== LEAVE SC =====================\n\n")

			for i := range sitesChannels {
				acr.ChangeValueSites(i, acr.id, acr.value)
			}

			// pAtt = pDiff, first clear pAtt then set pDiff finally clear pDiff
			for j := range acr.pAtt {
				delete(acr.pAtt, j)
			}

			for j := range acr.pDiff {
				acr.pAtt[j] = true
				acr.Ok(j, acr.id)
				delete(acr.pDiff, j)
			}

			acr.h = acr.h + 1
			acr.sc = false
			acr.demCours = false

			// needed to wait to go back to main
			acr.endAsk <- true
		}
	}
}

func (acr *algoCR) CheckSC() {
	acr.askingSC <- true
}

/**
 * convert local time to 4 digit string
 * in order to transmit it in payloads
 */
func (acr *algoCR) intToString(h int) string {
	var ret string
	ret += strconv.Itoa(h / 1000)
	ret += strconv.Itoa(h % 1000 / 100)
	ret += strconv.Itoa(h % 100 / 10)
	ret += strconv.Itoa(h % 10)
	return ret
}

func (acr *algoCR) ChangeValueSites(idTo int, idFrom int, value int) {
	msg := "V" + strconv.Itoa(value)
	acr.SendMsg(*sitesChannels[idTo], msg)
}

func (acr *algoCR) Ok(idTo int, idFrom int) {
	msg := "O" + acr.intToString(acr.h) + strconv.Itoa(idFrom)
	acr.SendMsg(*sitesChannels[idTo], msg)
}

func (acr *algoCR) Req(idTo int, idFrom int) {
	msg := "R" + acr.intToString(acr.hDem) + strconv.Itoa(idFrom)
	acr.SendMsg(*sitesChannels[idTo], msg)
}

func (acr *algoCR) SendMsg(msgChannel chan<- string, msg string) {
	msgChannel <- msg
	fmt.Println("sent : " + msg)
}

func (acr *algoCR) MsgHandle(msg string) {

	fmt.Println("received : " + msg)
	// extract payload parameters
	op := msg[0] // op : R | O | V

	if op == 'R' || op == 'O' {

		hi, _ := strconv.Atoi(string(msg[1:5]))
		i, _ := strconv.Atoi(string(msg[5]))

		// update stamp
		acr.h = int(math.Max(float64(acr.h), float64(hi))) + 1
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
		val, _ := strconv.Atoi(strings.Trim(msg, "V"))
		acr.SetValue(val)
	}

}
