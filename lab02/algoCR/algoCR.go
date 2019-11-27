package algoCR

import (
	"math"
	"strconv"
)

type algoCR struct{
	id       int         // id du site
	n        int         // nombre de processus
	h	     int         // valeur courante de l'estampille
	demCours bool        // faux/vrai quand le processus moi demande l'accès en SC
	sc		 bool        // faux/vrai quand le processus moi est en SC
	hDem	 int         // estampille de soumission de cette demande
	pDiff[]  int         // ensemble des numéros des processus pour lesquels on a différé le OK
	pAtt[]   int         // ensemble des processus desquels je dois obtenir une permission
	comm     chan string // communication (envoie)
}

func New(nbreSite int, c chan string, id int) algoCR{
	acr := algoCR{id,0,0,false,false,0,[]int{},[]int{}, c}
	return acr
}

//mettre recepteur si sa marche pas
func sendMsg(msgChanel chan string, msg string){
	msgChanel <- msg
}


func (acr algoCR) MsgHandle(msg string){

	hi , _ := strconv.Atoi(msg[1:3])
	acr.h = int(math.Max(float64(acr.h) , float64(hi))) + 1


	if acr.demCours == false {
		m := "O" + strconv.Itoa(acr.h) + msg[3:4]
		sendMsg(acr.comm, m)



	}
}


