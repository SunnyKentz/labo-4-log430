# Rapport du laboratoire 4

Github : https://github.com/SunnyKentz/labo-4-log430.git

## Préambule

Dans le cadre du laboratoire 4, j'ai ajouté deux nouvelles dépendances à mon architecture : Grafana et Prometheus.
Ces deux dépendances augmentent l'observabilité de mon application.
Sur Grafana, je surveille les 4 Golden Signals.

## Test de charge initial et observabilité de base

#### Outils choisis : 
J'ai créé mon propre outil sous `tests/loadtests/test.js`. Cet outil me permet d'envoyer des requêtes fetch selon le nombre de requêtes par seconde mis en argument. Je fais ces trois requêtes :
- Consultation simultanée des stocks de plusieurs magasins.
- Génération de rapports consolidés.
- Mise à jour de produits à forte fréquence.

#### Résultats :
Mon application fonctionne bien à 1 RPS (requête par seconde) mais tombe en défaillance à 30 RPS. Le nombre de requêtes sur le service logistique est considérablement plus élevé que les autres à cause des rapports consolidés.

### ici on peut voir que le système est stable à 1rps
![ScreenShot Graphana](./imgs/1rps.png)

### le service logistique tombe en premier en retournans des erreurs de la base donnée à 30rps
![ScreenShot Graphana](./imgs/30rps.png)

### ensuite à 100rps je perd la connection ssh et la base de donner mere tombe
![ScreenShot Graphana](./imgs/100rps.png)

#### Analysez les points faibles de l’architecture:

1. il y a trops d'appel ver le service logistique, il faut donc cache les réponses du service
2. il faut aussi cache les résultats de la base de donné lors du listing des produits dans les magasins.
3. il fau aussi cache les résultat de la db quand on liste les transactions pour la base de donne mere

## Ajout d’un Load Balancer et résilience

J'ai décidé d'utiliser Traefik pour mon load balancing car, traefik est plus récent et est construit avec go and se configure sur docker-compose facilement. et l'orchestration via Docker Compose.

#### Analyse par nombre d:instance :
Je decide de scale avec un nombre N le service de logistique car ce service est celui qui reçois le plus de requête
Dans ces testes je perd la connection ssh à 30 rps

#### N = 1 instance

| RPS = 1 | RPS = 15 | RPS = 30 |
|---------|----------|----------|
| ![ScreenShot Graphana](./imgs/1rps-load.png) | ![ScreenShot Graphana](./imgs/15rps-load.png) | ![ScreenShot Graphana](./imgs/30rps-load.png) |

#### N = 2 instances

| RPS = 1 | RPS = 15 | RPS = 30 |
|---------|----------|----------|
| ![ScreenShot Graphana](./imgs/1rps-load-2.png) | ![ScreenShot Graphana](./imgs/15rps-load-2.png) | ![ScreenShot Graphana](./imgs/30rps-load-2.png) |

#### N = 3 instances

| RPS = 1 | RPS = 15 | RPS = 30 |
|---------|----------|----------|
| ![ScreenShot Graphana](./imgs/1rps-load-3.png) | ![ScreenShot Graphana](./imgs/15rps-load-3.png) | ![ScreenShot Graphana](./imgs/30rps-load-3.png) |

#### N = 4 instances

| RPS = 1 | RPS = 15 | RPS = 30 |
|---------|----------|----------|
| ![ScreenShot Graphana](./imgs/1rps-load-4.png) | ![ScreenShot Graphana](./imgs/15rps-load-4.png) | ![ScreenShot Graphana](./imgs/30rps-load-4.png) |

