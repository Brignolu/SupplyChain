//imports
package main

import (
	"fmt"
	"math/rand"
	"time"
	. "github.com/pspaces/gospace"
)

//constantes définissant le nombres d'usines, de fournisseurs, d'entrepots et de revendeurs
const NPLANT = 5
const NSUPPL = 5
const NWARE = 5
const NRETAIL = 5
const NTRANS = 5

func main() {

	//On crée plusieurs espaces de tuples
	sync := NewSpace("tcp://localhost:31414/sync")
	propositions := NewSpace("tcp://localhost:31415/propositions")
	emballages := NewSpace("tcp://localhost:31416/emballages")
	transports := NewSpace("tcp://localhost:31147/transports")
	entrepots := NewSpace("tcp://localhost:31148/entrepots")

	client(&propositions,&emballages,&transports,&sync)
	go plants(&propositions,&sync)
	go suppliers(&emballages,&sync)
	go transporters(&transports,&sync)
	go warehouses(&entrepots,&propositions,&sync)
	go retailers(&entrepots,&sync)
	sync.Get("done")

}

//Agent Client qui va interargir avec l'utilisateur
func client(propositions *Space,emballages *Space,transports *Space, sync *Space){
	var productname,req string
	var time,qty,cost int
	fmt.Printf("Que desirez vous produire ?\n")
	fmt.Scanln(&productname)
	fmt.Printf("Caractéristiques Techniques ?\n")
	fmt.Scanln(&req)
	fmt.Printf("en combien de temps ?\n")
	fmt.Scanln(&time)
	fmt.Printf("Combien d'exemplaires désirez vous produire ?\n")
	fmt.Scanln(&qty)
	fmt.Printf("Budget souhaité (en euros) ?\n")
	fmt.Scanln(&cost)
	propositions.Put(productname,req,cost,time,qty)
	go first_appel(propositions,sync)
	go second_appel(emballages,sync)
	go trois_appel(transports,sync)
}

//Place de multiples propositions dans l'espace des propositions
//(produits,requirements,cost,time,qty)
func plants(propositions *Space, sync *Space) {
	var item,req string
	var time,qty,cost, rcost,rtime,rqty int
	t, _ := propositions.Query(&item,&req,&cost,&time,&qty)
	for i := 1; i < NPLANT+1; i++ {
		item = (t.GetFieldAt(0)).(string)
		req = (t.GetFieldAt(1)).(string)
		cost = (t.GetFieldAt(2)).(int)
		time = (t.GetFieldAt(3)).(int)
		qty = (t.GetFieldAt(4)).(int)
		rtime = time-(rand.Intn(time)-1)
		rqty = qty-(rand.Intn(qty)-1)
		rcost = cost-(rand.Intn(cost)-1)
		//design
		//workshop
		propositions.Put(i ,item, req, rcost, rtime, rqty)
		fmt.Printf("%d) Je propose %s, j'ai besoin de %s, je peux en fabriquer %d exemplaires en %d jours pour la modique somme de %d €		||	",i, item , req , rqty , rtime , rcost)
		fmt.Printf("out(%d, %s, %s, %d, %d, %d )\n\n",i, item , req , rcost , rtime , rqty)
}
sync.Put("appel1")
}

//Agent fournisseurs
func suppliers(emballages *Space, sync *Space){
	matiere := []string{
		"Carton",
		"Plastique",
		"Verre",
		"Papier",
	}
	_,_ = sync.Query("selection")
	for i:=1;i<NSUPPL+1;i++{
		n := rand.Int() % len(matiere)
		prix := rand.Int() % 100
		emballage := matiere[n]
		emballages.Put(i,emballage,prix)
		fmt.Printf("%d)Je propose un emballage en %s pour la modique somme de %d €		||	", i, emballage, prix)
		fmt.Printf("out(%d, %s, %d )\n\n", i, emballage, prix)
}
sync.Put("appel2")
}

//Agent transporteurs
func transporters(transports *Space, sync *Space){
	transport := []string{
		"Bateau",
		"Camion",
		"Avion",
	}
	_,_ = sync.Query("selection2")
	for i:=1;i<NTRANS+1;i++{
		n := rand.Int() % len(transport)
		prix := rand.Int() % 100
		transporteur := transport[n]
		transports.Put(i,transporteur,prix)
		fmt.Printf("%d)Je propose un transport avec %s pour la modique somme de %d €\n\n", i, transporteur, prix)
}
sync.Put("appel3")
}

func warehouses(warehouses *Space,propositions *Space,sync *Space){
	var item,req string
	var time,qty,cost int
	sync.Query("selection3")
	t,_:= propositions.Query("selected", &item, &req, &cost, &time, &qty)
	item = (t.GetFieldAt(1)).(string)
	time = (t.GetFieldAt(4)).(int)
	qty = (t.GetFieldAt(5)).(int)
	for i:=1;i<NWARE+1;i++{
	fmt.Printf("%d exemplaires de %s dans l'entrepot %d\n",qty/NWARE , item, i)
	warehouses.Put(qty/NWARE , item, i)


}
sync.Put("retailers")
}

func retailers(warehouses *Space,sync *Space){
	var item string
	var qty int
	sync.Query("retailers")
	fmt.Printf("\nTransport vers les détaillants en cours ")
	for i:=0;i<3;i++{
	 time.Sleep(1 * time.Second)
		fmt.Printf(".")
	}
	fmt.Printf("\n\n")
	for i:=1;i<NWARE+1;i++{
	t,_:= warehouses.Query(&qty,&item,i)
	item = (t.GetFieldAt(1)).(string)
	qty = (t.GetFieldAt(0)).(int)
	fmt.Printf("%d exemplaires de %s en vente dans le magasin %d\n",qty,item,i)
}
}

func first_appel(propositions *Space, sync *Space) {
	var item, req string
	var qty,cost,tim,id,selectedprop int
	_, _ = sync.Get("appel1")
	fmt.Printf("Tapez le numero de la proposition qui vous convient ou Enter pour renegocier\n")
	fmt.Scanln(&selectedprop)

	if(selectedprop != 0){
		t, _ := propositions.Get(selectedprop,&item, &req, &cost, &tim, &qty)
		id = (t.GetFieldAt(0)).(int)
		item = (t.GetFieldAt(1)).(string)
		req = (t.GetFieldAt(2)).(string)
		cost = (t.GetFieldAt(3)).(int)
		tim = (t.GetFieldAt(4)).(int)
		qty = (t.GetFieldAt(5)).(int)
		fmt.Printf("Vous avez selectionné ==> %d) Je propose %s, j'ai besoin de %s, je peux en fabriquer %d exemplaires en %d jours pour la modique somme de %d €\n\n",id, item , req , qty , tim , cost)
	propositions.Put("selected", item, req, cost, tim, qty)
	fmt.Printf("Fabrication en cours ")
	for i:=0;i<3;i++{
		time.Sleep(1 * time.Second)
		fmt.Printf(".")
	}
	fmt.Printf("\n")
	sync.Put("selection")
	} else {
		fmt.Printf("Renegociation en cours ...\n")
		go plants(propositions, sync)
		first_appel(propositions, sync)
	}
}


func second_appel(emballages *Space, sync *Space) {
	var emballage string
	var selectedprop, prix int
	_, _ = sync.Query("appel2")
	fmt.Printf("Selectionnez l'offre du fournisseur qui vous convient ou Enter pour renegocier\n")
	fmt.Scanln(&selectedprop)

	if(selectedprop != 0){
		t, _ := emballages.Get(selectedprop,&emballage,&prix)
		id := (t.GetFieldAt(0)).(int)
		emballage = (t.GetFieldAt(1)).(string)
		prix = (t.GetFieldAt(2)).(int)

		fmt.Printf("Vous avez selectionné ==> %d) Emballage en %s, pour la modique somme de %d €\n\n", id, emballage, prix)
	emballages.Put(1, emballage, prix)
	fmt.Printf("Emballage en cours ")
	for i:=0;i<3;i++{
		time.Sleep(1 * time.Second)
		fmt.Printf(".")
	}
	fmt.Printf("\n")
		sync.Put("selection2")
	} else {
		fmt.Printf("Renegociation en cours ...\n")
		go suppliers(emballages, sync)
		second_appel(emballages, sync)
	}
}

func trois_appel(transports *Space, sync *Space) {
	var transport string
	var selectedprop, prix int
	_, _ = sync.Query("appel3")
	fmt.Printf("Selectionnez l'offre du Transporteur qui vous convient ou Enter pour renegocier\n")
	fmt.Scanln(&selectedprop)

	if(selectedprop != 0){
		t, _ := transports.Get(selectedprop, &transport, &prix)
		id := (t.GetFieldAt(0)).(int)
		transport = (t.GetFieldAt(1)).(string)
		prix = (t.GetFieldAt(2)).(int)

		fmt.Printf("Vous avez selectionné ==> %d) Transport en %s, pour la modique somme de %d €\n\n", id, transport, prix)
	transports.Put(1, transport, prix)
	fmt.Printf("Transport en cours ")
	for i:=0;i<3;i++{
	 time.Sleep(1 * time.Second)
		fmt.Printf(".")
	}
	fmt.Printf("\n")
	sync.Put("selection3")
	} else {
		fmt.Printf("Renegociation en cours ...\n")
		go transporters(transports, sync)
		trois_appel(transports, sync)
	}
}
