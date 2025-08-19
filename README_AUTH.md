# Guide d'Intégration de l'Authentification JWT

Ce guide explique comment lier votre API d'authentification Node.js avec votre API d'annonces Go en utilisant le même token JWT.

## Configuration

### 1. API Node.js (Authentification)
Assurez-vous que votre API Node.js utilise le même secret JWT que dans le middleware Go :

```javascript
// Dans votre fichier .env
JWT_SECRET=mon_projet_agrobloc_jwt
```

### 2. API Go (Annonces)
Le middleware d'authentification est déjà configuré dans `middleware/auth.go` avec le même secret.

## Utilisation

### Obtenir un token depuis l'API Node.js
```bash
curl -X POST http://localhost:3000/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

Réponse :
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "user_id": 123,
    "profil_id": 456
  }
}
```

### Utiliser le token avec l'API Go
Incluez le token dans l'en-tête Authorization pour toutes les routes protégées :

```bash
# Créer une annonce de vente
curl -X POST http://localhost:8080/annonces_vente \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "statut": "active",
    "description": "Annonce de test",
    "type_culture_id": "uuid-du-type",
    "parcelle_id": "uuid-de-la-parcelle",
    "quantite": "100",
    "prix_kg": "5.50"
  }'
```

### Routes protégées
Les routes suivantes nécessitent un token valide :
- POST /annonces_vente
- PUT /annonces_vente/:id
- DELETE /annonces_vente/:id
- POST /annonces_achat
- PUT /annonces_achat/:id
- DELETE /annonces_achat/:id
- POST /annonces_pref
- PUT /annonces_pref/:id
- DELETE /annonces_pref/:id

### Routes publiques
Les routes suivantes sont accessibles sans authentification :
- GET /annonces_vente
- GET /annonces_vente/:id
- GET /annonces_achat
- GET /annonces_achat/:id
- GET /annonces_pref
- GET /annonces_pref/:id

## Structure du Token JWT
Le token JWT doit contenir :
```json
{
  "user_id": 123,
  "profil_id": 456,
  "iat": 1234567890,
  "exp": 1234567890
}
```

## Exemple d'intégration côté client
```javascript
// JavaScript/React exemple
const createAnnonce = async (annonceData) => {
  const token = localStorage.getItem('authToken');
  
  const response = await fetch('http://localhost:8080/annonces_vente', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`
    },
    body: JSON.stringify(annonceData)
  });
  
  return response.json();
};
```

## Dépannage
- **Token invalide** : Vérifiez que le secret JWT est identique dans les deux APIs
- **Format de token** : Assurez-vous d'utiliser le format `Bearer <token>`
- **Expiration** : Vérifiez la date d'expiration du token
- **User ID** : Le middleware récupère automatiquement l'ID utilisateur depuis le token
