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
	// process events
	const restartDelay = 250 * time.Millisecond
	restartLast := time.Time{}
	restartPattern := regexp.MustCompile(`(.+\.go|input.+\.txt)`)
	cache := x.FileCache{}
	for {
		select {
		case event, ok := <-w.Events:
			if !ok {
				continue
			}
			if !restartPattern.MatchString(event.Name) {
				continue
			}
			if time.Since(restartLast) < restartDelay {
				continue
			}
			if !cache.Changed(event.Name) {
				continue
			}
			x.Logf("file changed %s", event.Name)
			events <- restart
			restartLast = time.Now()
		case err, ok := <-w.Errors:
			if ok {
				return fmt.Errorf("watcher error: %w", err)
			}
		}
	}
}
