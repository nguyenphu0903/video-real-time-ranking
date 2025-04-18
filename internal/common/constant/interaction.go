package constant

type InteractionType string

const (
	View    InteractionType = "view"
	Like    InteractionType = "like"
	Comment InteractionType = "comment"
	Share   InteractionType = "share"
)

func (it *InteractionType) GetScore() float64 {
	switch *it {
	case View:
		return 1.0
	case Like:
		return 2.0
	case Comment:
		return 3.0
	case Share:
		return 4.0
	default:
		return 0.0
	}
}

const InteractionEventsChannel string = "interaction_events"
