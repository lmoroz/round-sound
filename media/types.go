package media

// StateMode represents playback state
type StateMode int

const (
	StatePlaying StateMode = 0
	StatePaused  StateMode = 1
	StateStopped StateMode = 2
)

// RepeatMode represents repeat mode
type RepeatMode int

const (
	RepeatNone RepeatMode = 1
	RepeatAll  RepeatMode = 2
	RepeatOne  RepeatMode = 4
)

// RatingSystem represents rating system type
type RatingSystem int

const (
	RatingNone        RatingSystem = 0
	RatingLike        RatingSystem = 1
	RatingLikeDislike RatingSystem = 2
	RatingScale       RatingSystem = 3
)

// Player represents media player state
type Player struct {
	ID              int          `json:"id"`
	Name            string       `json:"name"`
	Title           string       `json:"title"`
	Artist          string       `json:"artist"`
	Album           string       `json:"album"`
	Cover           string       `json:"cover"`
	CoverData       []byte       `json:"-"` // Binary cover data (not sent to frontend)
	State           StateMode    `json:"state"`
	Position        int          `json:"position"` // seconds
	Duration        int          `json:"duration"` // seconds
	Volume          int          `json:"volume"`   // 0-100
	Rating          int          `json:"rating"`   // 0=none, 1=dislike, 5=like
	Repeat          RepeatMode   `json:"repeat"`
	Shuffle         bool         `json:"shuffle"`
	RatingSystem    RatingSystem `json:"ratingSystem"`
	AvailableRepeat int          `json:"availableRepeat"`
	CanSetState     bool         `json:"canSetState"`
	CanSkipPrevious bool         `json:"canSkipPrevious"`
	CanSkipNext     bool         `json:"canSkipNext"`
	CanSetPosition  bool         `json:"canSetPosition"`
	CanSetVolume    bool         `json:"canSetVolume"`
	CanSetRating    bool         `json:"canSetRating"`
	CanSetRepeat    bool         `json:"canSetRepeat"`
	CanSetShuffle   bool         `json:"canSetShuffle"`
	CreatedAt       int64        `json:"createdAt"`
	UpdatedAt       int64        `json:"updatedAt"`
	ActiveAt        int64        `json:"activeAt"`
}

// Clone creates a copy of the player
func (p *Player) Clone() *Player {
	if p == nil {
		return nil
	}
	clone := *p
	return &clone
}
