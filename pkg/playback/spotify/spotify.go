package spotify

import (
	"context"
	"net/http"

	"github.com/XSAM/go-hybrid/errorw"
	"github.com/XSAM/go-hybrid/log"
	"github.com/deamwork/grid650-array-serial/pkg/httpserver"
	"github.com/deamwork/grid650-array-serial/pkg/playback/utils"
	spotifyclient "github.com/zmb3/spotify/v2"
	"go.uber.org/zap"
)

type Client struct {
	ctx    authContext
	client *spotifyclient.Client
}

func NewClient(clientID, clientSecret string, router httpserver.Router) *Client {
	client := &Client{}
	router.RegisterRouteHandlerFunc(http.MethodGet, "/oauth/spotify/callback", client.ctx.RedirectHandler)
	go client.OAuthCallback()
	client.OAuth(clientID, clientSecret)
	return client
}

func (c *Client) VendorName() string {
	return "spotify"
}

func (c *Client) OAuth(clientID, clientSecret string) string {
	c.setCreds(clientID, clientSecret)
	c.ctx.clientChan = make(chan *spotifyclient.Client, 1)
	return c.ctx.Auth()
}

func (c *Client) OAuthCallback() (state bool) {
	// wait for auth to complete
	c.client = <-c.ctx.clientChan

	// use the client to make calls that require authorization
	user, err := c.client.CurrentUser(context.Background())
	if err != nil {
		log.BgLogger().Error("spotify.auth", zap.Error(err))
		return false
	}

	log.BgLogger().Info("spotify.auth", zap.String("user_id", user.ID))

	return true
}

func (c *Client) CurrentPlaying(ctx context.Context) (trackInfo utils.TrackInfo, err error) {
	if c.client == nil {
		// client is not ready
		return utils.TrackInfo{}, errorw.NewMessage("client not ready")
	}
	trackInfo = utils.TrackInfo{}
	playing, err := c.client.PlayerCurrentlyPlaying(ctx)
	if err != nil {
		return
	}

	// not playing?
	if !playing.Playing {
		return trackInfo, errorw.NewMessage("not playing")
	}

	trackInfo = utils.TrackInfo{
		Album: playing.Item.Album.Name,
		Name:  playing.Item.Name,
		Artists: func() []string {
			var artists []string
			for _, artist := range playing.Item.Artists {
				artists = append(artists, artist.Name)
			}
			return artists
		}(),
	}

	return trackInfo, nil
}

func (c *Client) setCreds(clientID, clientSecret string) {
	c.ctx.cred.clientID = clientID
	c.ctx.cred.clientSecret = clientSecret
}
