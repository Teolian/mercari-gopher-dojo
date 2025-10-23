package omikuji

import (
	"math/rand"
	"time"
)

// Fortune represents the level of luck in omikuji.
type Fortune string

const (
	DaiKichi  Fortune = "Dai-kichi"  // Great blessing
	Kichi     Fortune = "Kichi"      // Blessing
	ChuuKichi Fortune = "Chuu-kichi" // Middle blessing
	ShoKichi  Fortune = "Sho-kichi"  // Small blessing
	SueKichi  Fortune = "Sue-kichi"  // Future blessing
	Kyo       Fortune = "Kyo"        // Curse
	DaiKyo    Fortune = "Dai-kyo"    // Great curse
)

// AllFortunes is a list of all possible fortune values.
var AllFortunes = []Fortune{
	DaiKichi,
	Kichi,
	ChuuKichi,
	ShoKichi,
	SueKichi,
	Kyo,
	DaiKyo,
}

// Response represents the JSON response structure for omikuji API.
type Response struct {
	Fortune   Fortune `json:"fortune"`
	Health    string  `json:"health"`
	Residence string  `json:"residence"`
	Travel    string  `json:"travel"`
	Study     string  `json:"study"`
	Love      string  `json:"love"`
}

// Clock is a function type for getting the current time.
// This allows for dependency injection in tests.
type Clock func() time.Time

// DefaultClock returns the current local time.
var DefaultClock Clock = time.Now

// GenerateFortune generates a random fortune response.
// During New Year (Jan 1-3), it always returns Dai-kichi.
func GenerateFortune(clock Clock) Response {
	now := clock()
	month, day := now.Month(), now.Day()

	var fortune Fortune
	if month == time.January && day >= 1 && day <= 3 {
		// New Year special: always Dai-kichi
		fortune = DaiKichi
	} else {
		// Random fortune
		fortune = AllFortunes[rand.Intn(len(AllFortunes))]
	}

	return Response{
		Fortune:   fortune,
		Health:    "You will fully recover, but stay attentive after you do.",
		Residence: "You will have good fortune with a new house.",
		Travel:    "When traveling, you may find something to treasure.",
		Study:     "Things will be better. It may be worth aiming for a school in a different area.",
		Love:      "The person you are looking for is very close to you.",
	}
}

func init() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())
}
