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
	// filter
	// commandPattern := regexp.MustCompile(`(executed the '(.+)' command?`)
	if strings.Contains(event.Parsed["message"].(string), "'write memory' command") || // cisco
		strings.Contains(event.Parsed["message"].(string), "commit complete") || // juniper
		strings.Contains(strings.ToLower(event.Parsed["message"].(string)), "attribute configured") /*fortinet*/ {
		//config updated
		log.Println("[nap]", event)
		s.responser.Send(*event, s.triggers["config-updated"], "net.skycloud.nap.messaging.model.LogEvent")
	}
	return nil
}

func (s *nap) beforeDispatch() {

}

func (s *nap) afterDispatch() {

}
