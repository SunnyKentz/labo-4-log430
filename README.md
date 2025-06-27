## Labo 4

### Explication de l'application:

l'application sur les ports 8080 : magasin, 8090: mere, 8091 : centre logistique

On peut login avec Bob en tant que manager et Alice en tant que commis
le nom du magasin peut etre nimporte quoi
la caisse doit etre Caisse 1, Caisse 2 ou Caisse 3

le mdp est toujour : password

### Monitoring et Métriques

L'application inclut maintenant des métriques Prometheus pour le monitoring de tous les services :

#### Magasin API (Port 8080)
- **Endpoint des métriques** : `http://localhost:8080/api/v1/metrics`
- **Métriques** : Requêtes HTTP, produits, transactions, ventes, remboursements, opérations panier
- **Documentation** : `caisse_app_scaled/magasin/api/README_METRICS.md`

#### Maison Mere API (Port 8090)
- **Endpoint des métriques** : `http://localhost:8090/api/v1/metrics`
- **Métriques** : Requêtes HTTP, transactions, ventes, remboursements, magasins, notifications, alertes
- **Documentation** : `caisse_app_scaled/maison_mere/api/README_METRICS.md`

#### Centre Logistique API (Port 8091)
- **Endpoint des métriques** : `http://localhost:8091/api/v1/metrics`
- **Métriques** : Requêtes HTTP, produits, commandes, acceptations, refus, commandes en attente
- **Documentation** : `caisse_app_scaled/centre_logistique/api/README_METRICS.md`

### Comment run :
```
    make dev-setup
    make run
```

### Comment tester :
```
    make test
```

### Comment generer les docs :
```
    make docs
```

### Explication du CI
Après avoir fait un push, Github action check le linting du push,<br> execute les testes et si tous passe, le push vers dockerhub
