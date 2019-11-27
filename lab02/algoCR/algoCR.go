package algoCR

import (
	"math"
	"strconv"
)

type algoCR struct{
	nbreSite int         // nombre de sites de l'infrasctructure
	id       int         // id du site
	n        int         // nombre de processus
	h	     int         // valeur courante de l'estampille
	demCours bool        // faux/vrai quand le processus moi demande l'accès en SC
	sc		 bool        // faux/vrai quand le processus moi est en SC
	hDem	 int         // estampille de soumission de cette demande
	pDiff[]  int         // ensemble des numéros des processus pour lesquels on a différé le OK
	pAtt[]   int         // ensemble des processus desquels je dois obtenir une permission
}

func New(nbreSite int, id int) algoCR {
	acr := algoCR{nbreSite,id,0,0,false,false,0,[]int{},[]int{}}
	return acr
}

func (acr *algoCR) Init(deltaT int) {

}

func (acr algoCR) SendMsg(msgChannel chan <- string, msg string) {
	msgChannel <- msg
}

func (acr algoCR) MsgHandle(msgChannel chan <- string, msg string) {

	hi , _ := strconv.Atoi(msg[1:3])
	acr.h = int(math.Max(float64(acr.h) , float64(hi))) + 1


	if acr.demCours == false {
		m := "O" + strconv.Itoa(acr.h) + msg[3:4]
		acr.SendMsg(msgChannel, m)
	}
}


