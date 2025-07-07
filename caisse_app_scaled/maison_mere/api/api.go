package api

import (
	"caisse-app-scaled/caisse_app_scaled/logger"
	"caisse-app-scaled/caisse_app_scaled/maison_mere/db"
	"caisse-app-scaled/caisse_app_scaled/maison_mere/mere"
	. "caisse-app-scaled/caisse_app_scaled/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
)

func NewApp() {

	engine := html.New("./view", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return false
		},
	}))
	app.Use(loadBalancerMiddlerWare)
	app.Static("/static", "./view")
	app.Static("/js", "./commonjs")
	app.Mount("/api/v1", newDataApi())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("login", nil)
	})

	app.Get("/home", func(c *fiber.Ctx) error {
		return c.Render("analytics", nil)
	})

	app.Get("/rapport", func(c *fiber.Ctx) error {
		return c.Render("reports", nil)
	})

	app.Get("/produits", func(c *fiber.Ctx) error {
		return c.Render("products", nil)
	})

	db.Init()
	logger.Init("Mere")
	db.SetupLog()
	log.Fatal(app.Listen(":8090"))
}

func loadBalancerMiddlerWare(c *fiber.Ctx) error {
	path := c.Path()
	if strings.HasPrefix(path, "/mere") {
		p := strings.TrimPrefix(path, "/mere")
		p = strings.TrimPrefix(p, "/")
		c.Path("/" + p)
	}
	return c.Next()
}

func authMiddleWare(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	err := mere.CheckLogedIn(authHeader)
	if err != nil {
		if c.Path() == "api/v1/login" {
			return c.Next()
		}
		if strings.HasPrefix(c.Path(), "/api") {
			return GetApiError(c, "this action requires authentification", http.StatusUnauthorized)
		}
		return c.Redirect("/")
	}
	return c.Next()
}
