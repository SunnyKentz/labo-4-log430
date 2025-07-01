import { test, expect } from '@playwright/test';

test('load testing', async ({ page }) => {
    await page.goto('http://localhost:8080/');
    await page.getByRole('textbox', { name: 'Votre nom' }).fill('Bob');
    await page.getByRole('textbox', { name: 'Votre mot de passe' }).fill('password');
    await page.getByRole('textbox', { name: 'Votre Magasin' }).fill('Magasin 1');
    await page.locator('#caisse').fill('Caisse 1');
    await page.getByRole('button', { name: 'Se connecter' }).click();

});

test('load testing 2', async ({ page }) => {
    await page.goto('http://localhost:8090/');
    await page.getByRole('textbox', { name: 'Votre nom' }).fill('Bob');
    await page.getByRole('textbox', { name: 'Votre mot de passe' }).fill('password');
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.getByRole('link', { name: 'Rapports' }).click();
    await page.getByRole('button', { name: 'Générer' }).click();

})

test('load testing 3', async ({ page }) => {
    await page.goto('http://localhost:8090/');
    await page.getByRole('textbox', { name: 'Votre nom' }).fill('Bob');
    await page.getByRole('textbox', { name: 'Votre mot de passe' }).fill('password');
    await page.getByRole('button', { name: 'Se connecter' }).click();
    await page.getByRole('link', { name: 'Produits' }).click();
    await page.getByRole('textbox', { name: 'Rechercher un produit...' }).fill('s');
    await page.getByRole('textbox', { name: 'Rechercher un produit...' }).press('Enter');
    await page.locator('.btn-primary').first().click();
    await page.locator('input[name="nom"]').fill('Cff');
    await page.getByRole('spinbutton').fill('12.00');
    await page.locator('#editPriceForm input[name="description"]').fill('new descp');
    page.once('dialog', dialog => {
        console.log(`Dialog message: ${dialog.message()}`);
        dialog.dismiss().catch(() => { });
    });
    await page.getByRole('button', { name: 'Mettre à jour' }).click();
})