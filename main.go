package main

import (
	"net/http"
	"github.com/gin-gonic/gin"

		"fmt"
	"net/url"
  "github.com/ory/hydra-client-go/client"
  "github.com/ory/hydra-client-go/client/admin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("/consent", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	router.GET("/login", func(c *gin.Context) {
		challenge := c.Query("login_challenge")
		adminURL, _ := url.Parse("http://127.0.0.1:4445")
		hydraAdmin := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path})
		result, err := hydraAdmin.Admin.GetLoginRequest(admin.NewGetLoginRequestParams().WithLoginChallenge(challenge))
		if err != nil {
			fmt.Printf("%v", err)
		} else {
			fmt.Printf("%v", result)
		}
		c.HTML(http.StatusOK, "login.html", gin.H{
			"challenge": "Main website",
		})
	})

	router.POST("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"challenge": "hi",
		})
	})

	router.GET("/logout", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	router.Run(":3000")
}
