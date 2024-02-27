package main

import (
	"encoding/json"
	"envi/apps/api/initializers"
	model "envi/apps/api/models"
	envs "envi/packages/env"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func RegisterMiddlewares(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339}, remote_ip=${remote_ip}, method=${method}, " +
			"path=${path}, status=${status}, took=${response_time}, sent=t=${response_size} bytes\n",
	}))
}
func RegisterRoutes(e *echo.Echo) {
	e.GET("/health", healthHandler)
	e.GET("/oauth/authorize/google", oAuthAuthorizeGoogleHandler)
	e.GET("/oauth/callback/google", oAuthCallbackGoogleHandler)
}

const (
	googleOAuthURL = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL = "https://oauth2.googleapis.com/token"
)

func oAuthAuthorizeGoogleHandler(c echo.Context) error {
	u, _ := url.Parse(googleOAuthURL)
	values := u.Query()
	values.Set("client_id", envs.Env.GOOGLE_CLIENT_ID)
	values.Set("redirect_uri", "https://antoniel.zipper.ngrok.app/oauth/callback/google")
	values.Set("response_type", "code")
	values.Set("scope", "https://www.googleapis.com/auth/userinfo.email https://www.googleapis.com/auth/userinfo.profile")
	values.Set("access_type", "offline")
	u.RawQuery = values.Encode()
	return c.Redirect(http.StatusFound, u.String())
}

type googleTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	IDToken     string `json:"id_token"`
}

func oAuthCallbackGoogleHandler(c echo.Context) error {
	u, _ := url.Parse(googleTokenURL)
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
	jsonResp := new(googleTokenResponse)
	err = json.Unmarshal(bodyResp, jsonResp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	accountInfo, err := getTokenInfo(jsonResp.IDToken)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	user := model.User{
		Name:  accountInfo.Name,
		Email: accountInfo.Email,
		Image: accountInfo.Picture,
	}
	resultCreateUser := initializers.DB.Create(&user)
	if resultCreateUser.Error != nil {
		return c.JSON(http.StatusInternalServerError, resultCreateUser.Error)
	}

	account := model.Account{
		CompoundID:         strings.Join([]string{accountInfo.Email, "google"}, ""),
		UserID:             int(user.ID),
		PoviderType:        "google",
		RefreshToken:       jsonResp.AccessToken,
		AccessToken:        jsonResp.AccessToken,
		AccessTokenExpires: time.Now().Add(time.Duration(jsonResp.ExpiresIn) * time.Second),
		PoviderID:          "",
		PoviderAccountID:   "",
	}
	resultAccount := initializers.DB.Create(&account)
	if resultAccount.Error != nil {
		return c.JSON(http.StatusInternalServerError, resultAccount.Error)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"OK": true})
}

type tokenInfoResponse struct {
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Typ           string `json:"typ"`
}

func getTokenInfo(token string) (*tokenInfoResponse, error) {
	u, _ := url.Parse("https://oauth2.googleapis.com/tokeninfo")
	values := u.Query()
	values.Set("id_token", token)
	u.RawQuery = values.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	tokenInfo := new(tokenInfoResponse)
	err = json.Unmarshal(bodyResp, tokenInfo)
	if err != nil {
		return nil, err
	}
	return tokenInfo, nil
}

func healthHandler(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
