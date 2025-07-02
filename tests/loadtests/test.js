const fetch = (...args) => import('node-fetch').then(({ default: fetch }) => fetch(...args));
let magasin = 'Magasin 1';
let caisse = 'Caisse 1';

main()

function main() {
    const argv = process.argv;
    const rps = parseInt(argv[2], 10) || 10; // Default to 10 rps if not provided
    console.log(`Starting load test with ${rps} requests per second...`);
    loadTestConsultation(rps)
    loadTestMiseAJour(rps)
    loadTestRapport(rps)
}

async function loadTestConsultation(requestsPerSecond) {
    // Consultation simultanée des stocks de plusieurs magasins.
    let initial = requestsPerSecond
    let token = await loginAndGetToken('http://magasin.localhost/api/v1/login', magasin, caisse)
    setInterval(() => {
        while (requestsPerSecond > 0) {
            fetch('http://magasin.localhost/api/v1/produits', {
                method: 'GET',
                headers: { 'Authorization': `Bearer ${token}`, 'C-Mag': magasin, 'C-Caisse': caisse }
            })
            requestsPerSecond--;
        }
        requestsPerSecond = initial;
    }, 1000);
}

async function loadTestRapport(requestsPerSecond) {
    //  Génération de rapports consolidés.
    let token = await loginAndGetToken('http://mere.localhost/api/v1/merelogin', magasin, caisse)
    let initial = requestsPerSecond
    setInterval(() => {
        while (requestsPerSecond > 0) {
            fetch('http://mere.localhost/api/v1/raport', {
                method: 'GET',
                headers: { 'Authorization': `Bearer ${token}` }
            })
            requestsPerSecond--;
        }
        requestsPerSecond = initial;
    }, 1000);
}

async function loadTestMiseAJour(requestsPerSecond) {
    //  Mise à jour de produits à forte fréquence. 
    let token = await loginAndGetToken('http://mere.localhost/api/v1/merelogin', magasin, caisse)
    let initial = requestsPerSecond
    const productData = {
        productId: 5,
        nom: "new nom",
        prix: 23.45,
        description: "description",
    };
    setInterval(() => {
        while (requestsPerSecond > 0) {
            fetch('http://mere.localhost/api/v1/raport', {
                method: 'PUT',
                headers: { 'Authorization': `Bearer ${token}`,"Content-Type": "application/json" },
                body: JSON.stringify(productData),
            })
            requestsPerSecond--;
        }
        requestsPerSecond = initial;
    }, 1000);
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
        console.error('Error during login:', err.message,loginEndpoint);
        return null;
    }
}