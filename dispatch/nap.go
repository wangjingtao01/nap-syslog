package dispatch

import (
	"log"
	"strings"

	"github.com/ekanite/ekanite/input"
)

// NAP dispatcher
type nap struct {
	responser *Responser
	triggers  map[string]Trigger
}

func NewNapDispatcher(config DispatcherInstance) (Dispatcher, error) {
	responser, err := NewResponser(config)
	if err != nil {
		return nil, err
	}
	return &nap{
		responser: responser,
		triggers:  config.Triggers,
	}, nil
}

func (s *nap) Listen(c chan []*input.Event) {
	for {
		select {
		case e := <-c:
			for _, v := range e {
				s.do(v)
			}
			break
		}
	}
}

func (s *nap) do(event *input.Event) error {
	log.Println("[nap]", event)
	// filter
	// commandPattern := regexp.MustCompile(`(executed the '(.+)' command?`)
	if strings.Contains(event.Parsed["message"].(string), "'write memory' command") ||
		strings.Contains(event.Parsed["message"].(string), "commit complete") {
		//config updated
		s.responser.Send(*event, s.triggers["config-updated"])
	}
	return nil
}

func (s *nap) beforeDispatch() {

}

func (s *nap) afterDispatch() {

}
