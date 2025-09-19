package handlers

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

//go:embed swagger-ui/*
var swaggerFiles embed.FS

// DocsHandler serves API documentation
type DocsHandler struct{}

func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

// SwaggerUI serves the Swagger UI interface
func (h *DocsHandler) SwaggerUI() fiber.Handler {
	// Get the swagger-ui subdirectory from embedded files
	swaggerUI, err := fs.Sub(swaggerFiles, "swagger-ui")
	if err != nil {
		// Fallback to serving a simple documentation page
		return func(c *fiber.Ctx) error {
			return c.Type("html").SendString(`
<!DOCTYPE html>
<html>
<head>
    <title>Idea Collision Engine API</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .container { max-width: 800px; margin: 0 auto; }
        .header { background: #1f2937; color: white; padding: 20px; border-radius: 8px; }
        .content { background: #f9fafb; padding: 20px; border-radius: 8px; margin-top: 20px; }
        .endpoint { background: white; padding: 15px; margin: 10px 0; border-radius: 5px; border-left: 4px solid #3b82f6; }
        .method { font-weight: bold; color: #059669; }
        .method.post { color: #dc2626; }
        .method.put { color: #d97706; }
        .method.delete { color: #dc2626; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Idea Collision Engine API</h1>
            <p>Creative productivity API for generating unexpected idea combinations</p>
        </div>
        
        <div class="content">
            <h2>📚 API Documentation</h2>
            <p>The OpenAPI specification is available at: <a href="/docs/openapi.yaml">/docs/openapi.yaml</a></p>
            
            <h3>🔗 Key Endpoints</h3>
            
            <div class="endpoint">
                <span class="method">GET</span> <strong>/health</strong><br>
                <small>Service health check</small>
            </div>
            
            <div class="endpoint">
                <span class="method post">POST</span> <strong>/api/auth/register</strong><br>
                <small>Register a new user account</small>
            </div>
            
            <div class="endpoint">
                <span class="method post">POST</span> <strong>/api/auth/login</strong><br>
                <small>Authenticate and get access token</small>
            </div>
            
            <div class="endpoint">
                <span class="method post">POST</span> <strong>/api/collisions/generate</strong><br>
                <small>Generate idea collision (requires authentication)</small>
            </div>
            
            <div class="endpoint">
                <span class="method">GET</span> <strong>/api/collisions/history</strong><br>
                <small>Get collision generation history</small>
            </div>
            
            <div class="endpoint">
                <span class="method">GET</span> <strong>/api/domains/basic</strong><br>
                <small>Get available collision domains for basic users</small>
            </div>
            
            <div class="endpoint">
                <span class="method">GET</span> <strong>/api/subscriptions/plans</strong><br>
                <small>Get available subscription plans</small>
            </div>
            
            <h3>🔐 Authentication</h3>
            <p>Most endpoints require a Bearer token obtained via <code>/api/auth/login</code>.</p>
            <p>Include the token in the Authorization header: <code>Authorization: Bearer &lt;token&gt;</code></p>
            
            <h3>📊 Rate Limiting</h3>
            <p>Free users are limited to 10 collision generations per minute.</p>
            <p>Premium users have no rate limits.</p>
            
            <h3>📈 Usage Limits</h3>
            <p>Free users: 50 collisions per week</p>
            <p>Pro users: Unlimited collisions</p>
            <p>Team users: Unlimited collisions + premium domains</p>
        </div>
    </div>
</body>
</html>`)
		}
	}

	return filesystem.New(filesystem.Config{
		Root:       http.FS(swaggerUI),
		PathPrefix: "/docs",
		Browse:     true,
	})
}

// OpenAPISpec serves the OpenAPI YAML specification
func (h *DocsHandler) OpenAPISpec(c *fiber.Ctx) error {
	return c.SendFile("./docs/openapi.yaml")
}