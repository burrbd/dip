package slack

// Notifier implements session.Notifier by posting notifications to the game
// channel via the Slack Channel adapter.
type Notifier struct {
	ch *Channel
}

// NewNotifier returns a Notifier that sends messages via ch.
func NewNotifier(ch *Channel) *Notifier {
	return &Notifier{ch: ch}
}

// Notify posts message to channelID.
func (n *Notifier) Notify(channelID, message string) error {
	return n.ch.Post(channelID, message)
}
