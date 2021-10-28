package playback

import (
	"context"
	"time"

	"github.com/XSAM/go-hybrid/log"
	"github.com/deamwork/grid650-array-serial/pkg/httpserver"
	"github.com/deamwork/grid650-array-serial/pkg/playback/spotify"
	"github.com/deamwork/grid650-array-serial/pkg/playback/utils"
	"go.uber.org/zap"
)

type Spotify struct {
	c *spotify.Client
}

func NewSpotifyClient() *Spotify {
	return &Spotify{}
}

func (s *Spotify) Start(clientID, clientSecret string, register httpserver.Router, trackChan chan utils.TrackInfo) {
	s.c = spotify.NewClient(clientID, clientSecret, register)

	s.getTrack(trackChan)
}

func (s *Spotify) getTrack(ch chan utils.TrackInfo) {
	ctx := context.Background()
	for {
		var errCounter int
		track, err := s.c.CurrentPlaying(ctx)
		if err != nil {
			if err.Error() == "not playing" {
				log.Logger(ctx).Info("spotify.getTrack", zap.Error(err))
				ch <- utils.TrackInfo{}
			} else if err.Error() == "client not ready" {
				log.Logger(ctx).Error("spotify.getTrack", zap.Error(err), zap.String("status", "waiting for client"))
			} else {
				log.Logger(ctx).Error("spotify.getTrack", zap.Error(err))
				errCounter += 1
			}
			time.Sleep(time.Second * 5)
			continue
		} else if errCounter > 5 {
			log.Logger(ctx).Error("spotify.getTrack", zap.String("err", "error happened more then 5 times, exit"))
			return
		}

		ch <- track

		time.Sleep(time.Second)
	}
}
