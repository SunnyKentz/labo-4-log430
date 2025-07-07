package api

import (
	"caisse-app-scaled/caisse_app_scaled/logger"
	"caisse-app-scaled/caisse_app_scaled/magasin/caissier"
	. "caisse-app-scaled/caisse_app_scaled/utils"
	"log"
	"net/http"
	"strings"
	"time"

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
		return c.Render("login", fiber.Map{
			"Title": "Login - Caisse App",
		})
	})

	app.Get("/home", func(c *fiber.Ctx) error {
		return c.Render("product", nil)
	})
	app.Get("/panier", func(c *fiber.Ctx) error {
		return c.Render("checkout", nil)
	})
	app.Get("/transactions", func(c *fiber.Ctx) error {
		return c.Render("transactions", nil)
	})

	port := ":8080"
	logger.Init("Magasin")
	caissier.InitialiserPOS("init", "Caisse 1", "NO MAGASIN")
	log.Fatal(app.Listen(port))
}

func loadBalancerMiddlerWare(c *fiber.Ctx) error {
	path := c.Path()
	if strings.HasPrefix(path, "/magasin") {
		p := strings.TrimPrefix(path, "/magasin")
		p = strings.TrimPrefix(p, "/")
		c.Path("/" + p)
	}
	return c.Next()
}

func authMiddleWare(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	magasin := c.Get("C-Mag")
	caisse := c.Get("C-Caisse")
	err := caissier.CheckLogedIn(authHeader, magasin, caisse)
	if err != nil {
		if strings.HasPrefix(c.Path(), "/api") {
			return GetApiError(c, "this action requires authentification", http.StatusUnauthorized)
		}
		return c.Redirect("/")
	}
	return c.Next()
}

var cache = make(map[string]any)
var cacheTime = make(map[string]time.Time)

func cacheMiddleware(c *fiber.Ctx) error {
	noCache := c.Get("no-cache", "false")
	if noCache != "false" {
		return c.Next()
	}
	path := c.Method() + " " + c.Path()
	now := time.Now()

	// Check if we have a cached value and it's not older than 30 seconds
	if t, ok := cacheTime["t - "+path]; ok {
		if now.Sub(t) < 30*time.Second {
			if val, ok := cache[path]; ok {
				return c.JSON(val)
			}
		}
	}
	err := c.Next()
	if err != nil {
		return err
	}
	return nil
}
