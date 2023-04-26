package main

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	googleOauthConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "http://localhost:3000/google/redirect",
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func main() {
	r := gin.Default()

	r.GET("/google", func(c *gin.Context) {
		url := googleOauthConfig.AuthCodeURL("state-token")
		c.Redirect(http.StatusTemporaryRedirect, url)
	})

	r.GET("/google/redirect", func(c *gin.Context) {
		code := c.Query("code")
		token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to exchange token",
				"error":   err.Error(),
			})
			return
		}

		client := googleOauthConfig.Client(oauth2.NoContext, token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to get user info",
				"error":   err.Error(),
			})
			return
		}
		defer resp.Body.Close()

		var userInfo map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to decode user info",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":     "User information from Google",
			"user":        userInfo,
			"accessToken": token.AccessToken,
		})
	})

	r.Run(":3000")
}
