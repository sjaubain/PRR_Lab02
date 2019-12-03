# Programmation Répartie - Laboratoire II

Auteurs: Jobin Simon, Teklehaimanot Robel

## Énoncé du problème

Partager une donnée parmi un ensemble de processus est un problème qui peut se résoudre par le biais d'un algorithme d'exclusion mutuelle. Dans ce laboratoire, nous allons utiliser l’algorithme de Carvalho et Roucairol, une optimisation de l’algorithme de Ricart et Agrawala, comme algorithme d’exclusion mutuelle.  Chaque processus détient une variable entière qui doit être cohérente. Les tâches applicatives peuvent faire 2 opérations sur cette variable : consulter sa valeur, et modifier sa valeur. La consultation revient à obtenir la valeur la plus récente. Par contre, une modification se passe en 3 étapes :

- **obtenir l'exclusion mutuelle sur la variable**
- **modifier la valeur en section critique, par exemple l’incrémenter**
- **informer tous les autres processus de la nouvelle valeur**
- **libérer la section critique**

## Fonctionnement

Le programme fonctionne en deux temps : tout d'abord, tous les sites doivent attendre d'être interconnectés. Pour se faire, au démarrage d'un site, celui-ci commence par faire un lookup sur tous les autres sites pour tenter de se connecter à ceux déjà lancés tout en restant à l'écoute des connexions entrantes. Lorsqu'un site accepte une connexion d'un autre site, il s'y connecte immédiatement en retour. Deux goroutines de lecture et d'écriture sur le socket permettent aux sites de communiquer correctement.

Une fois que la première étape est terminée, l'utilisateur a la possibilité d'entrer les commandes `W` (Write) et `R` (Read) pour soit changer la valeur partagée soit la consulter. Le site intègre une partie gérant l'exclusion mutuelle et implémente l'algorithme de Carvalho et Roucairol pour la gestion de la section critique. L'affichage de l'ancienne valeur et de la nouvelle, entrée par l'utilisateur, ainsi qu'une visualisation de l'entrée et sortie de la section critique (avec l'ajout d'un temps long de traitement artificiel pour des raisons de clarté) permettent d'observer le bon fonctionnement du traitement. Nous avons aussi mis dans l'affichage les messages reçus(REQ) et envoyés(OK) par les différents site.

Dans le cadre de ce laboratoire, l'état en section critique consiste uniquement en un changement de la valeur d'une variable cohérente et partagée par tous les sites, par le biais de la fonction *setInputValue* du fichier `algoCR.go` mais on pourrait très bien remplacer cela par un traitement plus complexe dans le cadre d'une autre application.

## Configuration

Un fichier Json contenant le nombre de sites ainsi que leurs adresses permet de configurer comme on le souhaite l'architecture des differents processus. Par défaut, toutes les addresses IP sont *localhost*, puisqu'on fait tourner tous les sites localement. Il est donc possible de changer ce fichier de configuration et d'ajouter d'autres sites (local ou remote) selon leur répartition.

Pour ce laboratoire, nous avons décidé d'attribuer les identifiants des sites avec des valeurs comprises entre *0* et *n-1*, où *n* est le nombre de sites. On pourrait très bien attribuer des identifiants quelconques, puisque les structures de données utilisées pour traiter les sites dans `algoCR.go` sont des map et non des tableaux statiques. Comme on suppose dans la donnée que tous les sites sont lancés une et une seule fois et que le réseau est fiable, notre représentation est adéquate.

Nous avons par conséquent choisi de construire tous les payloads envoyés de la manière suivante :

* **byte 1** : opcode `R`, `W` ou `V`
* **bytes 2 à 5** : estampille (sans perte de généralité, choix d'une estampille sur 4 digits)
* **byte 6** : id

Notons que le `V` est utilisé lorsque l'on désire informer les autres sites d'un changement de la valeur partagée.

## Utilisation

Tout d'abord, cloner le repository https://github.com/sjaubain/PRR_Lab02 quelque part depuis le GOPATH. Pour lancer les différents sites, il faut ouvrir un terminal où se situe le fichier `site.go` puis lancer la commande suivante pour construire l'exécutable :

```bash
go build site.go
```

Ouvrir ensuite autant de terminal qu'il y a de sites configurés (par défaut 3) et les lancer avec leur identifiant en argument. Le numéro du site doit être entre *0* et *n-1*, où *n* est le nombre de sites. Par exemple :

```bash
./site 0
```

Une fois que tous sites seront lancés et interconnectés, il suffit comme cité précédemment d'entrer les commandes `W` ou `R` pour changer la valeur de la variable partagée ou la consulter.
