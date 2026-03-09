package local

// Notifier implements session.Notifier using the in-memory Channel.
// Notifications are posted to the game channel as plain text prefixed with "[notify] ".
type Notifier struct {
	ch *Channel
}

// NewNotifier returns a Notifier that posts to ch.
func NewNotifier(ch *Channel) *Notifier {
	return &Notifier{ch: ch}
}

// Notify posts "[notify] <message>" to channelID.
func (n *Notifier) Notify(channelID, message string) error {
	return n.ch.Post(channelID, "[notify] "+message)
}
