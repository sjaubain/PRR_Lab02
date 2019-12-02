# PRR_Lab02

## Énoncé du problème

Partager une donnée parmi un ensemble de processus est un problème qui peut se résoudre par le biais d'un algorithme d'exclusion mutuelle. Dans ce laboratoire, nous allons utiliser l’algorithme de Carvalho et Roucairol, une optimisation de l’algorithme de Ricart et Agrawala, comme algorithme d’exclusion mutuelle.  Chaque processus détient une variable entière qui doit être cohérente. Les tâches applicatives peuvent faire 2 opérations sur cette variable : consulter sa valeur, et modifier sa valeur. La consultation revient à obtenir la valeur la plus récente. Par contre, une modification se passe en 3 étapes :

1. obtenir l'exclusion mutuelle sur la variable ;
2. modifier la valeur en section critique, par exemple l’incrémenter ;
3. informer tous les autres processus de la nouvelle valeur ;
4. libérer la section critique.

## Fonctionnement

Le programme fonctionne en deux temps : tout d'abord, tous les sites doivent attendre d'être interconnectés. Pour se faire, au démarrage d'un site, celui-ci commence par faire un lookup sur tous les autres sites pour tenter de se connecter à ceux déjà lancés tout en restant à l'écoute des connexions entrantes. Lorsqu'un site accepte une connexion d'un autre site, il s'y connecte immédiatement en retour. Deux goroutines de lecture et d'écriture sur le socket permettent aux sites de communiquer correctement.

Une fois que la première étape est terminée, l'utilisateur a la possibilité d'entrer les commandes `W` (Write) et `R` (Read) pour soit changer la valeur partagée soit la consulter. Le site intègre une partie gérant l'exclusion mutuelle et implémente l'algorithme de Carvalho et Roucairol pour la gestion de la section critique. L'affichage de l'ancienne valeur et de la nouvelle, entrée par l'utilisateur, ainsi qu'une visualisation de l'entrée et sortie de la section critique (avec l'ajout d'un temps long de traitement artificiel pour des raisons de clarté) permettent d'observer le bon fonctionnement du traitement.

Dans le cadre de ce laboratoire, l'état en section critique consiste uniquement en un changement de la valeur d'une variable cohérente et partagée par tous les sites, par le biais de la fonction *setInputValue* du fichier `algoCR.go` mais on pourrait très bien remplacer cela par un traitement plus complexe dans le cadre d'une autre application

## Configuration

Un fichier Json contenant le nombre de sites ainsi que leurs adresses permet de configurer comme on le souhaite l'architecture des differents processus. Par défaut, toutes les addresses IP sont *localhost*, puisqu'on fait tourner tous les sites localement.

## Utilisation

Tout d'abord, cloner le repository [https://github.com/sjaubain/PRR_Lab02] quelque part depuis le GOPATH. Pour lancer les différents sites, il faut ouvrir un terminal où se situe le fichier `site.go` puis de lancer la commande suivante pour construire l'exécutable :

```bash
go build site.go
```

Ouvrir ensuite autant de terminal qu'il y a de sites configurés (par défaut 3) et les lancer avec leur identifiant en argument. Par exemple :

```bash
./site 0
```

Une fois que tous sites seront lancés et interconnectés, il suffit comme cité précédemment d'entrer les commandes `W` ou `R` pour changer la valeur de la variable partagée ou de changer sa valeur. 
