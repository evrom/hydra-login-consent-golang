package main

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

func radius_login(username string, password string) bool {
	packet := radius.New(radius.CodeAccessRequest, []byte(`secret`))
	rfc2865.UserName_SetString(packet, username)
	rfc2865.UserPassword_SetString(packet, password)
	response, err := radius.Exchange(context.Background(), packet, "localhost:1812")
	if err != nil {
		log.Printf("%v", err)
		return false
	}
	log.Println("Code:", response.Code)

	if response.Code == radius.CodeAccessAccept {
		return true
	} else {
		return false
	}
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	adminURL, _ := url.Parse("http://127.0.0.1:4445")
	hydraAdmin := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{Schemes: []string{adminURL.Scheme}, Host: adminURL.Host, BasePath: adminURL.Path})

	router.GET("/consent", func(c *gin.Context) {
		challenge := c.Query("consent_challenge")
		result, err := hydraAdmin.Admin.GetConsentRequest(
			admin.NewGetConsentRequestParams().
				WithConsentChallenge(challenge))
		if err != nil {
			log.Printf("%v", err)
		} else {
			log.Printf("%v", result)
		}

		c.HTML(http.StatusOK, "consent.html", gin.H{
			"challenge":       challenge,
			"RequestedScopes": result.Payload.RequestedScope,
			"ClientID":        result.Payload.Client.ClientID,
			"User":            result.Payload.Subject,
		})
	})

	router.POST("/consent", func(c *gin.Context) {
		challenge := c.PostForm("challenge")
		grant_scope := c.PostForm("grant_scope")
		result, err := hydraAdmin.Admin.AcceptConsentRequest(
			admin.NewAcceptConsentRequestParams().
				WithConsentChallenge(challenge).
				WithBody(&models.AcceptConsentRequest{
					GrantScope: *&models.StringSlicePipeDelimiter{
						grant_scope,
					},
				}))
		if err != nil {
			log.Printf("%v", err)
			c.Redirect(http.StatusFound, "/consent?consent_challenge="+challenge)
		} else {
			c.Redirect(http.StatusFound, *result.Payload.RedirectTo)
		}

	})

	router.GET("/login", func(c *gin.Context) {
		challenge := c.Query("login_challenge")
		c.HTML(http.StatusOK, "login.html", gin.H{
			"challenge": challenge,
		})
	})

	router.POST("/login", func(c *gin.Context) {
		challenge := c.PostForm("challenge")
		subject := c.PostForm("username")

		if radius_login(subject, c.PostForm("password")) == false {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"challenge": challenge,
				"error":     "wrong username or password",
			})
		} else {

			result, err := hydraAdmin.Admin.AcceptLoginRequest(
				admin.NewAcceptLoginRequestParams().
					WithLoginChallenge(challenge).
					WithBody(&models.AcceptLoginRequest{
						Subject: &subject,
					}))
			if err != nil {
				log.Printf("%v", err)
				c.Redirect(http.StatusFound, "/login")
			} else {
				c.Redirect(http.StatusFound, *result.Payload.RedirectTo)
			}
		}
	})

	router.GET("/logout", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	router.Run(":3000")
}
