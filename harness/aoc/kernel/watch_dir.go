package kernel

import (
	"fmt"
	"regexp"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jpillora/puzzler/internal/pzlr/x"
)

func watch(events chan event) error {
	// code.go exists, check other files too
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer w.Close()
	// watch directory
	if err = w.Add("."); err != nil {
		return err
	}
	cache := x.FileCache{}
	// trigger restart event
	const restartDelay = 250 * time.Millisecond
	restartLast := time.Time{}
	restartPattern := regexp.MustCompile(`(.+\.go|input.+\.txt)`)
	handleRestart := func(event fsnotify.Event) {
		if time.Since(restartLast) < restartDelay {
			return
		}
		if !cache.Changed(event.Name) {
			return
		}
		x.Logf("file changed %s", event.Name)
		events <- restart
		restartLast = time.Now()
	}
	// trigger fetch README.md event
	// TODO: on delete
	for {
		select {
		case event, ok := <-w.Events:
			if ok && restartPattern.MatchString(event.Name) {
				handleRestart(event)
			}
		case err, ok := <-w.Errors:
			if ok {
				return fmt.Errorf("watcher error: %w", err)
			}
		}
	}
}
