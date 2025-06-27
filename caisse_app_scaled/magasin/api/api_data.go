package api

import (
	"caisse-app-scaled/caisse_app_scaled/magasin/caissier"
	. "caisse-app-scaled/caisse_app_scaled/utils"

	_ "caisse-app-scaled/docs/swagger/magasin"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// @title Magasin API
// @version 1.0
// @description This is the API for the Magasin service.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @securityDefinitions.apikey Magasin
// @in header
// @name C-Mag
// @securityDefinitions.apikey Caisse
// @in header
// @name C-Caisse
func newDataApi() *fiber.App {
	api := fiber.New(fiber.Config{})
	api.Get("/swagger/*", swagger.HandlerDefault)
	api.Get("/metrics", prometheusHandler)
	api.Use(metricsMiddleware)
	api.Post("/login", loginHandler)
	api.Get("/produits", authMiddleWare, getProductsHandler)
	api.Get("/produits/:nom", authMiddleWare, findProductHandler)
	api.Post("/cart/:id", authMiddleWare, addToCartHandler)
	api.Get("/cart", authMiddleWare, getCartItemsHandler)
	api.Delete("/cart/:id", authMiddleWare, removeFromCartHandler)
	api.Post("/vendre", authMiddleWare, makeSaleHandler)
	api.Get("/transactions", authMiddleWare, getTransactionsHandler)
	api.Post("/rembourser/:id", authMiddleWare, refundTransactionHandler)
	api.Post("/produit/:id", authMiddleWare, requestRestockHandler)
	api.Put("/produit/:id", updateProductHandler)
	api.Put("/produit/:id/:qt", restockProductHandler)
	return api
}

// prometheusHandler handles the /metrics endpoint for Prometheus
func prometheusHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	// Get the prometheus handler and serve metrics directly
	handler := promhttp.Handler()

	// Create a simple HTTP request
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		return err
	}

	// Create a response writer that writes to Fiber response
	writer := &fiberResponseWriter{ctx: c}

	// Serve the prometheus metrics
	handler.ServeHTTP(writer, req)

	return nil
}

// fiberResponseWriter implements http.ResponseWriter for Fiber
type fiberResponseWriter struct {
	ctx *fiber.Ctx
}

func (w *fiberResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (w *fiberResponseWriter) Write(data []byte) (int, error) {
	w.ctx.Response().AppendBody(data)
	return len(data), nil
}

func (w *fiberResponseWriter) WriteHeader(statusCode int) {
	w.ctx.Status(statusCode)
}

// metricsMiddleware records HTTP request metrics
func metricsMiddleware(c *fiber.Ctx) error {
	start := time.Now()

	// Process request
	err := c.Next()

	// Record metrics
	duration := time.Since(start).Seconds()
	RecordHTTPRequestDuration(c.Method(), c.Path(), duration)

	status := strconv.Itoa(c.Response().StatusCode())
	RecordHTTPRequest(c.Method(), c.Path(), status)

	return err
}

// @Summary Login
// @Description Authenticate an employee with the Mere system
// @Tags Authentication
// @Accept application/json
// @Produce application/json
// @Param body body object{username=string,password=string,caisse=string,magasin=string} true "Login Credentials"
// @Success 200 {object} object{token=string} "JWT Token"
// @Failure 400 {object} models.ApiError "Bad Request"
// @Failure 403 {object} models.ApiError "Forbidden"
// @Router /api/v1/login [post]
func loginHandler(c *fiber.Ctx) error {
	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Caisse   string `json:"caisse"`
		Magasin  string `json:"magasin"`
	}

	if err := c.BodyParser(&requestBody); err != nil {
		return GetApiError(c, SYNTAX_ERR("body", "json"), http.StatusBadRequest)
	}
	employe := requestBody.Username
	pw := requestBody.Password
	caisse := requestBody.Caisse
	magasin := requestBody.Magasin
	jwt, err := caissier.Login(employe, pw)
	if err != nil {
		return GetApiError(c, "Failed to login, verifier la caisse, le nom, le mdp", http.StatusForbidden)
	}
	if !caissier.InitialiserPOS(employe, caisse, magasin) {
		return GetApiError(c, FAILURE_ERR(errors.New("echec d'ouverture de la caisse")), http.StatusInternalServerError)
	}
	return c.Status(200).JSON(fiber.Map{
		"token": jwt,
	})
}

// @Summary Get Products
// @Description Get all products
// @Tags products
// @Produce json
// @Security BearerAuth
// @Security Magasin
// @Security Caisse
// @Success 200 {array} models.Produit
// @Failure 500 {object} models.ApiError
// @Router /api/v1/produits [get]
func getProductsHandler(c *fiber.Ctx) error {
	produits, err := caissier.AfficherProduits()
	if err != nil {
		RecordError("database_error")
		return GetApiError(c, FAILURE_ERR(err), http.StatusInternalServerError)
	}
	// Update products total metric
	UpdateProductsTotal(len(produits))

	return c.JSON(produits)
}

// @Summary Find Product
// @Description Find a product by name
// @Tags products
// @Produce json
// @Security BearerAuth
// @Security Magasin
// @Security Caisse
// @Param nom path string true "Product Name"
// @Success 200 {array} models.Produit
// @Failure 400 {object} models.ApiError
// @Router /api/v1/produits/{nom} [get]
func findProductHandler(c *fiber.Ctx) error {
	nomParam := c.Params("nom")
	nom, err := url.QueryUnescape(nomParam)
	if err != nil {
		return GetApiError(c, SYNTAX_ERR("nom", nomParam), http.StatusBadRequest)
	}
	produits, err := caissier.TrouverProduit(nom)
	if err != nil {
		return GetApiError(c, NOTFOUND_ERR("nom", nom), http.StatusNotFound)
	}
	return c.JSON(produits)
}

// @Summary Add to Cart
// @Description Add a product to the cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Security Magasin
// @Security Caisse
// @Param id path int true "Product ID"
// @Success 200 {object} models.ApiSuccess
// @Failure 400 {object} models.ApiError
// @Failure 404 {object} models.ApiError
// @Router /api/v1/cart/{id} [post]
func addToCartHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		RecordError("invalid_parameter")
		return GetApiError(c, SYNTAX_ERR("id", idParam), http.StatusBadRequest)
	}
	err = caissier.AjouterALaCart(id)
	if err != nil {
		RecordError("product_not_found")
		return GetApiError(c, NOTFOUND_ERR("id", id), http.StatusNotFound)
	}

	// Record cart operation metric
	RecordCartOperation("add")

	return GetApiSuccess(c, 200)
}

// @Summary Get Cart Items
// @Description Get all items in the cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Security Magasin
// @Security Caisse
// @Success 200 {array} models.Produit
// @Failure 500 {object} models.ApiError
// @Router /api/v1/cart [get]
func getCartItemsHandler(c *fiber.Ctx) error {
	items, err := caissier.GetCartItems()
	if err != nil {
		return GetApiError(c, FAILURE_ERR(err), http.StatusInternalServerError)
	}
	return c.JSON(items)
}

// @Summary Remove from Cart
// @Description Remove a product from the cart
// @Tags cart
// @Produce json
// @Security BearerAuth
// @Security Magasin
// @Security Caisse
// @Param id path int true "Product ID"
// @Success 200 {object} models.ApiSuccess
// @Failure 400 {object} models.ApiError
// @Router /api/v1/cart/{id} [delete]
func removeFromCartHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		RecordError("invalid_parameter")
		return GetApiError(c, SYNTAX_ERR("id", idParam), http.StatusBadRequest)
	}
	caissier.RetirerDeLaCart(id)

	// Record cart operation metric
	RecordCartOperation("remove")

	return GetApiSuccess(c, 200)
}

// @Summary Make a Sale
// @Description Finalize the sale of items in the cart
// @Tags sales
// @Produce json
// @Security BearerAuth
// @Security Magasin
// @Security Caisse
// @Success 200 {object} models.ApiSuccess
// @Failure 500 {object} models.ApiError
// @Router /api/v1/vendre [post]
func makeSaleHandler(c *fiber.Ctx) error {
	err := caissier.FaireUneVente()
	caissier.ViderLaCart()
	if err != nil {
		RecordError("sale_failed")
		return GetApiError(c, FAILURE_ERR(err), http.StatusInternalServerError)
	}

	// Record transaction and sale metrics
	RecordTransaction()
	RecordSale()

	return GetApiSuccess(c, 200)
}

// @Summary Get Transactions
// @Description Get all transactions
// @Tags sales
// @Produce json
// @Security BearerAuth
// @Security Magasin
// @Security Caisse
// @Success 200 {array} models.Transaction
// @Failure 500 {object} models.ApiError
// @Router /api/v1/transactions [get]
func getTransactionsHandler(c *fiber.Ctx) error {
	transactions := caissier.AfficherTransactions()
	if transactions == nil {
		return GetApiError(c, FAILURE_ERR(errors.New("transaction list was null")), http.StatusInternalServerError)
	}
	return c.JSON(transactions)
}

// @Summary Refund Transaction
// @Description Refund a transaction by its ID
// @Tags sales
// @Produce json
// @Security BearerAuth
// @Security Magasin
// @Security Caisse
// @Param id path int true "Transaction ID"
// @Success 200 {object} models.ApiSuccess
// @Failure 400 {object} models.ApiError
// @Failure 404 {object} models.ApiError
// @Router /api/v1/rembourser/{id} [post]
func refundTransactionHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		RecordError("invalid_parameter")
		return GetApiError(c, SYNTAX_ERR("id", idParam), http.StatusBadRequest)
	}

	err = caissier.FaireUnRetour(id)
	if err != nil {
		RecordError("refund_failed")
		return GetApiError(c, NOTFOUND_ERR("id", id), http.StatusNotFound)
	}

	// Record refund metric
	RecordRefund()

	return GetApiSuccess(c, 200)
}

// @Summary Request Restock
// @Description Request to restock a product
// @Tags products
// @Produce json
// @Security BearerAuth
// @Security Magasin
// @Security Caisse
// @Param id path int true "Product ID"
// @Success 200 {object} models.ApiSuccess
// @Failure 400 {object} models.ApiError
// @Router /api/v1/produit/{id} [post]
func requestRestockHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return GetApiError(c, SYNTAX_ERR("id", idParam), http.StatusBadRequest)
	}
	caissier.DemmandeReapprovisionner(id)
	return GetApiSuccess(c, 200)
}

// @Summary Update Product
// @Description Update a product's details
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body object{nom=string,prix=float64,description=string} true "Product Data"
// @Success 200 {object} models.ApiSuccess
// @Failure 400 {object} models.ApiError
// @Failure 500 {object} models.ApiError
// @Router /api/v1/produit/{id} [put]
func updateProductHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return GetApiError(c, SYNTAX_ERR("id", idParam), http.StatusBadRequest)
	}
	var body struct {
		Nom         string  `json:"nom"`
		Prix        float64 `json:"prix"`
		Description string  `json:"description"`
	}

	if err := c.BodyParser(&body); err != nil {
		return GetApiError(c, SYNTAX_ERR("body", "json"), http.StatusBadRequest)
	}

	err = caissier.MiseAJourProduit(id, body.Nom, body.Prix, body.Description)
	if err != nil {
		return GetApiError(c, FAILURE_ERR(err), http.StatusInternalServerError)
	}
	return GetApiSuccess(c, 200)
}

// @Summary Restock Product
// @Description Restock a product with a given quantity
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Param qt path int true "Quantity"
// @Success 200 {object} models.ApiSuccess
// @Failure 400 {object} models.ApiError
// @Failure 500 {object} models.ApiError
// @Router /api/v1/produit/{id}/{qt} [put]
func restockProductHandler(c *fiber.Ctx) error {
	qtParam := c.Params("qt")
	qt, err := strconv.Atoi(qtParam)
	if err != nil {
		return GetApiError(c, SYNTAX_ERR("qt", qtParam), http.StatusBadRequest)
	}
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return GetApiError(c, SYNTAX_ERR("id", idParam), http.StatusBadRequest)
	}
	err = caissier.Reapprovisionner(id, qt)
	if err != nil {
		return GetApiError(c, FAILURE_ERR(err), http.StatusInternalServerError)
	}
	return GetApiSuccess(c, 200)
}
