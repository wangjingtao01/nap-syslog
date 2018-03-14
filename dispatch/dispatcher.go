package dispatch

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ekanite/ekanite/input"
	"github.com/songtianyi/rrframework/config"
)

// dispatch events
type Dispatcher interface {
	Listen(chan []*input.Event)
	beforeDispatch()
	do(*input.Event) error
	afterDispatch()
}

// create dispatcher from json configuration file
func NewDispatcher(path string) ([]Dispatcher, error) {
	//
	jc, err := rrconfig.LoadJsonConfigFromFile(path)
	if err != nil {
		return nil, err
	}
	// dump config
	du, _ := jc.Dump()
	log.Println(du)
	var cfg DispatcherConfig
	if err := json.Unmarshal(jc.GetBytes(), &cfg); err != nil {
		return nil, err
	}
	dispatchers := make([]Dispatcher, 0)
	for _, v := range cfg.Dispatcher {
		switch strings.ToLower(v.Type) {
		case "nap":
			d, err := NewNapDispatcher(v)
			if err != nil {
				return nil, fmt.Errorf("new nap dispatcher failed, %s", err)
			}
			dispatchers = append(dispatchers, d)
			break
		case "elasticsearch":
			break
		default:
			return nil, fmt.Errorf("dispatcher type %s not valid or not support yet", v.Type)
		}
	}
	return dispatchers, nil
}
