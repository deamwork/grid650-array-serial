package spotify

import (
	"fmt"
	"net/http"
	"time"

	"github.com/XSAM/go-hybrid/log"
	spotifyclient "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"go.uber.org/zap"
)

const (
	redirectURL = "http://localhost:8008/oauth/spotify/callback"
)

var state = fmt.Sprintf("spo_%d", time.Now().Unix())

type authContext struct {
	cred          credentials
	authenticator *spotifyauth.Authenticator
	clientChan    chan *spotifyclient.Client
}

type credentials struct {
	clientID, clientSecret string
}

func (ac *authContext) Auth() string {
	ac.getAuth()
	url := ac.authenticator.AuthURL(state)
	log.BgLogger().Info("spotify.auth", zap.String("login_url", url))
	return url
}

func (ac *authContext) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	token, err := ac.authenticator.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "fail to fetch token", http.StatusForbidden)
		log.BgLogger().Error("spotify.auth", zap.Error(err))
		return
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.BgLogger().Error("spotify.auth", zap.String("err", "state mismatch"))
		return
	}

	ac.clientChan <- spotifyclient.New(ac.authenticator.Client(r.Context(), token))
	log.BgLogger().Info("spotify.auth", zap.Bool("login", true))
	return
}

func (ac *authContext) getAuth() {
	ac.authenticator = spotifyauth.New(
		spotifyauth.WithClientID(ac.cred.clientID),
		spotifyauth.WithClientSecret(ac.cred.clientSecret),
		spotifyauth.WithRedirectURL(redirectURL),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopeUserReadCurrentlyPlaying,
			spotifyauth.ScopeUserReadPlaybackState))

	return
}
