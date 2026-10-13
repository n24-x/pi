package lifecycle

type appState int

const (
	stopped appState = iota
	starting
	incompleteStart
	started
	stopping
)

func (as appState) String() string {
	switch as {
	case stopped:
		return "stopped"
	case starting:
		return "starting"
	case incompleteStart:
		return "incompleteStart"
	case started:
		return "started"
	case stopping:
		return "stopping"
	default:
		return "invalidState"
	}
}
