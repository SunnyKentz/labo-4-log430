package api

import (
	"caisse-app-scaled/caisse_app_scaled/centre_logistique/logistics"
	. "caisse-app-scaled/caisse_app_scaled/utils"

	_ "caisse-app-scaled/docs/swagger/logistique"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// @title Centre Logistique API
// @version 1.0
// @description This is the API for the Centre Logistique service.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8091
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /api/v1
func newDataApi() *fiber.App {
	api := fiber.New(fiber.Config{})
	api.Get("/swagger/*", swagger.HandlerDefault)
	api.Get("/metrics", prometheusHandler)
	api.Use(metricsMiddleware)
	api.Post("/login", loginHandler)
	api.Get("/commands", authMiddleWare, getAllCommandsHandler)
	api.Post("/commande/:magasin/:id", createCommandHandler)
	api.Put("/commande/:id", authMiddleWare, acceptCommandHandler)
	api.Delete("/commande/:id", authMiddleWare, refuseCommandHandler)
	api.Get("/produits/:nom", findProductHandler)
	api.Get("/produits/id/:id", findProductByIDHandler)
	api.Put("/produit/:id", updateProductHandler)

	return api
}

// @Summary Login
// @Description Authenticates an employee with the Mere system
// @Tags auth
// @Accept json
// @Produce json
// @Param body body object{username=string,password=string} true "Login Credentials"
// @Success 200 {object} object{token=string}
// @Failure 400 {object} models.ApiError
// @Failure 403 {object} models.ApiError
// @Router /api/v1/login [post]
func loginHandler(c *fiber.Ctx) error {
	var requestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return GetApiError(c, SYNTAX_ERR("body", "json"), http.StatusBadRequest)
	}
	employe := requestBody.Username
	pw := requestBody.Password
	jwt, err := logistics.Login(employe, pw)
	if err != nil {
		return GetApiError(c, "Failed to login", http.StatusForbidden)
	}
	return c.Status(200).JSON(fiber.Map{
		"token": jwt,
	})
}

// @Summary Get All Commands
// @Description Get all commands
// @Tags commands
// @Produce json
// @Security BearerAuth
// @Success 200 {array} logistics.Commande
// @Router /api/v1/commands [get]
func getAllCommandsHandler(c *fiber.Ctx) error {
	commands := logistics.GetAllCommands()

	// Update pending commands metric (all commands are considered pending)
	UpdateCommandsPending(len(commands))

	return c.JSON(commands)
}

// @Summary Create Command
// @Description Create a new command
// @Tags commands
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param magasin path string true "Magasin"
// @Param id path int true "Command ID"
// @Param body body object{host=string} true "Host"
// @Success 200 {object} models.ApiSuccess
// @Failure 400 {object} models.ApiError
// @Router /api/v1/commande/{magasin}/{id} [post]
func createCommandHandler(c *fiber.Ctx) error {
	mag, err := url.QueryUnescape(c.Params("magasin"))
	if err != nil {
		RecordError("invalid_parameter")
		return GetApiError(c, SYNTAX_ERR("magasin", c.Params("magasin")), http.StatusBadRequest)
	}

	var body struct {
		Host string `json:"host"`
	}
	if err := c.BodyParser(&body); err != nil {
		RecordError("invalid_parameter")
		return GetApiError(c, SYNTAX_ERR("body", "json"), http.StatusBadRequest)
	}
	id, err1 := strconv.Atoi(c.Params("id"))
	if err1 != nil {
		RecordError("invalid_parameter")
		return GetApiError(c, SYNTAX_ERR("id", c.Params("id")), http.StatusBadRequest)
	}
	logistics.AjouterUneCommande(id, mag, body.Host)

	// Record command metric
	RecordCommand()

	return GetApiSuccess(c, 200)
}

// @Summary Accept Command
// @Description Accept a command
// @Tags commands
// @Produce json
// @Security BearerAuth
// @Param id path int true "Command ID"
// @Success 200 {object} models.ApiSuccess
// @Failure 400 {object} models.ApiError
// @Router /api/v1/commande/{id} [put]
func acceptCommandHandler(c *fiber.Ctx) error {
	id, err1 := strconv.Atoi(c.Params("id"))
	if err1 != nil {
		RecordError("invalid_parameter")
		return GetApiError(c, SYNTAX_ERR("id", c.Params("id")), http.StatusBadRequest)
	}
	if ok := logistics.AccepterUneCommande(id); !ok {
		RecordError("command_accept_failed")
		return GetApiError(c, FAILURE_ERR(errors.New("failed to accept command")), http.StatusInternalServerError)
	}

	// Record accepted command metric
	RecordCommandAccepted()

	return GetApiSuccess(c, 200)
}

// @Summary Refuse Command
// @Description Refuse a command
// @Tags commands
// @Produce json
// @Security BearerAuth
// @Param id path int true "Command ID"
// @Success 200 {object} models.ApiSuccess
// @Failure 400 {object} models.ApiError
// @Router /api/v1/commande/{id} [delete]
func refuseCommandHandler(c *fiber.Ctx) error {
	id, err1 := strconv.Atoi(c.Params("id"))
	if err1 != nil {
		RecordError("invalid_parameter")
		return GetApiError(c, SYNTAX_ERR("id", c.Params("id")), http.StatusBadRequest)
	}
	if ok := logistics.RefuserUneCommande(id); !ok {
		RecordError("command_refuse_failed")
		return GetApiError(c, FAILURE_ERR(errors.New("failed to refuse command")), http.StatusInternalServerError)
	}

	// Record refused command metric
	RecordCommandRefused()

	return GetApiSuccess(c, 200)
}

// @Summary Find Product
// @Description Find a product by name
// @Tags produits
// @Produce json
// @Param nom path string true "Product Name"
// @Success 200 {array} models.Produit
// @Failure 400 {object} models.ApiError
// @Failure 404 {object} models.ApiError
// @Router /api/v1/produits/{nom} [get]
func findProductHandler(c *fiber.Ctx) error {
	nom, err := url.QueryUnescape(c.Params("nom"))
	if err != nil {
		return GetApiError(c, SYNTAX_ERR("nom", c.Params("nom")), http.StatusBadRequest)
	}
	produits, err := logistics.TrouverProduit(nom)
	if err != nil {
		return GetApiError(c, NOTFOUND_ERR("nom", nom), http.StatusNotFound)
	}
	return c.JSON(produits)
}

// @Summary Find Product
// @Description Find a product by id
// @Tags produits
// @Produce json
// @Param id path string true "Product Name"
// @Success 200 {array} models.Produit
// @Failure 400 {object} models.ApiError
// @Failure 404 {object} models.ApiError
// @Router /api/v1/produits/id/{id} [get]
func findProductByIDHandler(c *fiber.Ctx) error {
	id, err1 := strconv.Atoi(c.Params("id"))
	if err1 != nil {
		return GetApiError(c, SYNTAX_ERR("id", c.Params("id")), http.StatusBadRequest)
	}
	prod, err := logistics.TrouverProduitParID(id)
	if err != nil {
		return GetApiError(c, NOTFOUND_ERR("id", id), http.StatusNotFound)
	}
	return c.JSON(prod)
}

// @Summary Update Product
// @Description Updates a product
// @Tags produits
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body object{nom=string,prix=float64,description=string} true "Product Update Data"
// @Success 200 {object} models.ApiSuccess
// @Failure 400 {object} models.ApiError
// @Failure 500 {object} models.ApiError
// @Router /api/v1/produit/{id} [put]
func updateProductHandler(c *fiber.Ctx) error {
	id, err1 := strconv.Atoi(c.Params("id"))
	if err1 != nil {
		return GetApiError(c, SYNTAX_ERR("id", c.Params("id")), http.StatusBadRequest)
	}
	var body struct {
		Nom         string  `json:"nom"`
		Prix        float64 `json:"prix"`
		Description string  `json:"description"`
	}

	if err := c.BodyParser(&body); err != nil {
		return GetApiError(c, SYNTAX_ERR("body", "json"), http.StatusBadRequest)
	}

	err := logistics.MiseAJourProduit(id, body.Nom, body.Prix, body.Description)
	if err != nil {
		return GetApiError(c, FAILURE_ERR(err), http.StatusInternalServerError)
	}
	return GetApiSuccess(c, 200)
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
