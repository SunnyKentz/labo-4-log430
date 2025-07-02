# Rapport du laboratoire 3

Github : https://github.com/SunnyKentz/labo-3-log430.git

## Description des endpoints


## Authentification

Tous les endpoints de l'API nécessitent une authentification via JWT (JSON Web Token). Pour utiliser les endpoints protégés :

1. **Obtenir un token** : Utilisez d'abord l'endpoint de login correspondant à votre service
2. **Configurer l'autorisation** : Dans Swagger UI, cliquez sur le bouton "Authorize" en haut à droite
3. **Ajouter le token** : Dans le champ "BearerAuth", entrez votre token au format : `Bearer <votre_token_jwt>`
4. **Tester les endpoints** : Tous les endpoints protégés utiliseront automatiquement ce token


### API Magasin (Port 8080)

#### Authentification
- **POST** `/api/v1/login` - Authentification des employés du magasin
  - Body: `{"username": "string", "password": "string", "magasin": "string", "caisse": "string"}`
  - Retourne un token JWT pour l'authentification

#### Produits
- **GET** `/api/v1/produits` - Récupération de tous les produits
- **GET** `/api/v1/produits/:nom` - Recherche de produits par nom
- **GET** `/api/v1/produits/id/:id` - Recherche de produit par ID
- **PUT** `/api/v1/produit/:id` - Mise à jour d'un produit

#### Panier et Transactions
- **POST** `/api/v1/cart/:productId` - Ajout d'un produit au panier
- **GET** `/api/v1/cart` - Récupération du contenu du panier
- **POST** `/api/v1/checkout` - Finalisation d'une transaction
- **GET** `/api/v1/transactions` - Historique des transactions

#### Monitoring
- **GET** `/api/v1/metrics` - Métriques Prometheus du service magasin
  - Métriques disponibles : Requêtes HTTP, produits, transactions, ventes, remboursements, opérations panier

### API Maison Mère (Port 8090)

#### Authentification
- **POST** `/api/v1/merelogin` - Authentification des managers
  - Body: `{"username": "string", "password": "string"}`
  - Retourne un token JWT

#### Gestion des Magasins
- **GET** `/api/v1/magasins` - Liste de tous les magasins
- **GET** `/api/v1/magasins/:nom` - Informations d'un magasin spécifique

#### Transactions et Rapports
- **GET** `/api/v1/transactions` - Toutes les transactions
- **GET** `/api/v1/transactions/:id` - Transaction spécifique
- **POST** `/api/v1/vente` - Enregistrement d'une nouvelle vente

#### Notifications et Alertes
- **POST** `/api/v1/notify` - Réception de notifications
- **POST** `/api/v1/subscribe` - Souscription d'un nouveau magasin
- **GET** `/api/v1/alerts` - Récupération de toutes les alertes

#### Analytics
- **GET** `/api/v1/analytics/:mag` - Analytics pour un magasin spécifique
- **GET** `/api/v1/raport` - Rapport consolidé pour tous les magasins

#### Monitoring
- **GET** `/api/v1/metrics` - Métriques Prometheus du service maison mère
  - Métriques disponibles : Requêtes HTTP, transactions, ventes, remboursements, magasins, notifications, alertes

### API Centre Logistique (Port 8091)

#### Authentification
- **POST** `/api/v1/login` - Authentification des employés logistiques
  - Body: `{"username": "string", "password": "string"}`

#### Commandes
- **GET** `/api/v1/commands` - Toutes les commandes
- **POST** `/api/v1/commande/:magasin/:id` - Création d'une commande
- **PUT** `/api/v1/commande/:id` - Acceptation d'une commande
- **DELETE** `/api/v1/commande/:id` - Refus d'une commande

#### Produits
- **GET** `/api/v1/produits/:nom` - Recherche de produits par nom
- **GET** `/api/v1/produits/id/:id` - Recherche de produit par ID
- **PUT** `/api/v1/produit/:id` - Mise à jour d'un produit

#### Monitoring
- **GET** `/api/v1/metrics` - Métriques Prometheus du service centre logistique
  - Métriques disponibles : Requêtes HTTP, produits, commandes, acceptations, refus, commandes en attente

### Services de Monitoring

#### Prometheus (Port 9090)
- **Interface web** : `http://localhost:9090`
- **Configuration** : `./data/prometheus.yml`
- **Intervalle de scraping** : 5 secondes
- **Services surveillés** : magasin (8080), mere (8090), logistique (8091)

#### Grafana (Port 3000)
- **Interface web** : `http://localhost:3000`
- **Identifiants par défaut** : admin/admin
- **Source de données** : Prometheus
- **Fonctionnalités** : Dashboards, alertes, visualisations

### Screenshots Swagger/Postman

Pour utiliser Swagger et accéder à la documentation interactive des APIs, il faut naviguer vers `/api/v1/swagger` pour chaque service :

- **Magasin** : `http://localhost:8080/api/v1/swagger`
- **Maison Mère** : `http://localhost:8090/api/v1/swagger`  
- **Centre Logistique** : `http://localhost:8091/api/v1/swagger`

Ces URLs permettent d'accéder à l'interface Swagger UI qui documente toutes les endpoints disponibles pour chaque service, avec la possibilité de tester les APIs directement depuis l'interface web.

![ScreenShot Swagger](./snap.png)

### Vue mise à jour du modèle 4+1 (si n´ecessaire)

Pas de mise à jour pour le modèle 4+1



