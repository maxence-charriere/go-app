package app

var (
	// UIChan is a channel which take a func as payload.
	// Every func are executed in a dedicated goroutine.
	// Component callbacks should be called through this channel.
	UIChan = make(chan func(), 255)
)

func init() {
	go startUIScheduler()
}

func startUIScheduler() {
	for f := range UIChan {
		f()
	}
}
