package playback

import (
	"context"

	"github.com/deamwork/grid650-array-serial/pkg/httpserver"
	"github.com/deamwork/grid650-array-serial/pkg/playback/utils"
)

type Provider interface {
	VendorName() string
	OAuth(clientID, clientSecret string) (authURL string)
	OAuthCallback() (state bool)
	CurrentPlaying(ctx context.Context) (result utils.TrackInfo, err error)
	//SetTrackInfo(album, name string, artists []string)
	//GetTrackInfo() (album, name string, artists []string)
}

type Player interface {
	Start(clientID, clientSecret string, register httpserver.Router, trackChan chan utils.TrackInfo)
	getTrack(ch chan utils.TrackInfo)
}
