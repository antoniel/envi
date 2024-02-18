package main

import (
	"io"
	"net/http"
	"net/url"
	envs "pkg/env"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	_, err := NewPostgresStorage()
	if err != nil {
		e.Logger.Fatal(err)
	}
	RegisterMiddlewares(e)
	RegisterRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
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
