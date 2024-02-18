package main

import (
	envs "envi/src/internal/env"
	"io"
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, remote_ip=${remote_ip}, method=${method}, " +
			"path=${path}, status=${status}, took=${response_time}, sent=t=${response_size} bytes\n",
	}))
	e.GET("/health", healthHandler)
	e.GET("/oauth/authorize/google", oAuthAuthorizeGoogleHandler)
	e.GET("/oauth/callback/google", oAuthCallbackGoogleHandler)
	e.GET("/oauth/exchange/google", oAuthExchangeTokenGoogleHandler)
	e.Logger.Fatal(e.Start(":1323"))
}

func oAuthExchangeTokenGoogleHandler(c echo.Context) error {
	u, err := url.Parse("https://oauth2.googleapis.com/token")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	code := c.QueryParam("code")
	values := u.Query()

	values.Set("client_id", envs.Env.GOOGLE_CLIENT_ID)
	values.Set("client_secret", envs.Env.GOOGLE_SECRET)
	values.Set("code", code)
	values.Set("grant_type", "authorization_code")
	values.Set("redirect_uri", "https://antoniel.zipper.ngrok.app/oauth/callback/google")
	u.RawQuery = values.Encode()

	req, err := http.NewRequest(http.MethodPost, u.String(), nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusUnauthorized, "Invalid code")
	}

	return c.JSON(http.StatusOK, resp)
}

func oAuthAuthorizeGoogleHandler(c echo.Context) error {
	u, _ := url.Parse("https://accounts.google.com/o/oauth2/v2/auth")
	values := u.Query()
	values.Set("client_id", envs.Env.GOOGLE_CLIENT_ID)
	values.Set("redirect_uri", "https://antoniel.zipper.ngrok.app/oauth/callback/google")
	values.Set("response_type", "code")
	values.Set("scope", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile")
	values.Set("access_type", "offline")
	u.RawQuery = values.Encode()
	return c.Redirect(http.StatusFound, u.String())
}

func oAuthCallbackGoogleHandler(c echo.Context) error {
	u, _ := url.Parse("https://oauth2.googleapis.com/token")
	resp, err := http.PostForm(u.String(), url.Values{
		"client_id":     {envs.Env.GOOGLE_CLIENT_ID},
		"client_secret": {envs.Env.GOOGLE_SECRET},
		"code":          {c.QueryParam("code")},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"https://antoniel.zipper.ngrok.app/oauth/callback/google"},
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	defer resp.Body.Close()
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSONBlob(resp.StatusCode, bodyResp)
}

func healthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
