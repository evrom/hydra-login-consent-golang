package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
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
		c.HTML(http.StatusOK, "login.html", gin.H{
			"challenge": challenge,
		})
	})

	router.POST("/login", func(c *gin.Context) {
		challenge := c.PostForm("challenge")
		adminURL, _ := url.Parse("http://127.0.0.1:4445")
		hydraAdmin := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path})
		subject := "foo@bar.com"
		result, err := hydraAdmin.Admin.AcceptLoginRequest(
			admin.NewAcceptLoginRequestParams().
				WithLoginChallenge(challenge).
				WithBody(&models.AcceptLoginRequest{
					Subject: &subject,
				}))
		if err != nil {
			fmt.Printf("%v", err)
			c.Redirect(http.StatusFound, "/login")
		} else {
			c.Redirect(http.StatusFound, *result.Payload.RedirectTo)
		}
	})

	router.GET("/logout", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	router.Run(":3000")
}
