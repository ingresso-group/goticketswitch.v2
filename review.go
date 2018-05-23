package ticketswitch

import "time"

type Review struct {
	// review test.
	Body string
	// date and time of the review.
	DateTime time.Time
	// rating on a scale of 1-5, with 1 being the lowest rating and 5 being the highest rating.
	StarRating int
	// the IETF language tag for the review.
	Language string
	// a review title if available.
	Title string
	// the review was made by a user not a critic.
	IsUser bool
	// the authors name.
	Author string
	// the original url.
	URL string
}
