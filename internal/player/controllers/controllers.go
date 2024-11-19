package controllers

type PlayerController interface {
	SeekTo(float64) error
	ApplySpotifyWorkaround() error
}

var PlayerControllers map[string]PlayerController = map[string]PlayerController{
	"mpris": &MprisPlayerController{},
}
