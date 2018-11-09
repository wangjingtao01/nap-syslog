package dispatch

import (
	"log"
	"regexp"
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
		strings.Contains(strings.ToLower(event.Parsed["message"].(string)), "attribute configured") /*fortinet*/ ||
		strings.Contains(strings.ToLower(event.Parsed["message"].(string)), "Configured from") /*NX-OS || IOS*/ {
		//config updated
		log.Println("[nap]", event)
		s.responser.Send(*event, s.triggers["config-updated"], "net.skycloud.nap.messaging.model.LogEvent")
		return nil
	}
	nexusUpDownPattern := `Interface ([\w\d\/]+) is (up|down)`
	m := regexp.MustCompile(nexusUpDownPattern).FindStringSubmatch(event.Parsed["message"].(string))
	if len(m) > 0 {
		event.Parsed["interface"] = m[1]
		event.Parsed["state"] = m[2]
		s.responser.Send(*event, s.triggers["interface-up-down"], "net.skycloud.nap.messaging.model.LogEvent")
		return nil
	}
	iosUpDownPattern := `Line protocol on Interface ([\w\d\/]+), changed state to (up|down)`
	n := regexp.MustCompile(iosUpDownPattern).FindStringSubmatch(event.Parsed["message"].(string))
	if len(n) > 0 {
		event.Parsed["interface"] = m[1]
		event.Parsed["state"] = m[2]
		s.responser.Send(*event, s.triggers["interface-up-down"], "net.skycloud.nap.messaging.model.LogEvent")
		return nil
	}
	// ignore
	return nil
}

func (s *nap) beforeDispatch() {

}

func (s *nap) afterDispatch() {

}
