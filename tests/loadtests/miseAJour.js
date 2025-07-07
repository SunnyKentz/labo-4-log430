const fetch = (...args) => import('node-fetch').then(({ default: fetch }) => fetch(...args));
let magasin = 'Magasin 1';
let caisse = 'Caisse 1';
loadTestMiseAJour()
async function loadTestMiseAJour() {
    //  Mise à jour de produits à forte fréquence. 
    let token = await loginAndGetToken('http://localhost/mere/api/v1/merelogin', magasin, caisse)

    const productData = {
        productId: 5,
        nom: "new nom",
        prix: 23.45,
        description: "description",
    };

    fetch('http://localhost/mere/api/v1/produit', {
        method: 'PUT',
        headers: { 'Authorization': `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify(productData),
    });
}

async function loginAndGetToken(loginEndpoint, magasin, caisse) {
    const loginPayload = {
        username: 'Bob',
        password: 'password',
        magasin: magasin,
        caisse: caisse
    };
    try {
        const res = await fetch(loginEndpoint, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(loginPayload)
        });
        if (!res.ok) {
            throw new Error(`Login failed: ${res.status} ${res.statusText}`);
        }
        const data = await res.json();
        return data.token || data.access_token || null;
    } catch (err) {
        console.error('Error during login:', err.message, loginEndpoint);
        return null;
    }
}