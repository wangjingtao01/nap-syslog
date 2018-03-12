package input

import (
	"regexp"
	"strconv"
)

// RFC3164 represents a parser for RFC3164-compliant log messages
// BUT made some modifications to make it compatible with juniper syslog messages
type RFC3164 struct {
	matcher []*regexp.Regexp
}

func (p *Parser) newRFC3164Parser() {
	p.rfc = &RFC3164{}
	p.rfc.compileMatcher()
}

func (s *RFC3164) compileMatcher() {
	s.matcher = make([]*regexp.Regexp, 4)
	pri := `<([0-9]{1,3})>`
	ts := `([A-Za-z]+\s\d+(\s\d+)?\s\d+:\d+:\d+)`  // with year
	chost := `([:alnum:]]{1,}(-[[:alnum:]]+){0,})` // cisco hostname
	jhost := `([[:alnum:]._-]+)`                   // juniper hostname
	// uuid := `([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12})`
	uuid := `([^ =:]+)`
	capp := `(%\w+-\d+-\d+)` // %ASA-5-611103
	japp := `(\w+\[\d+\])`   // sshd[8144]
	msg := `(.+$)`
	// cisco host only
	s.matcher[0] = regexp.MustCompile(pri + ts + `\s` + chost + `\s:\s` + capp + `:\s` + msg)
	// cisco uuid only
	s.matcher[1] = regexp.MustCompile(pri + ts + `\s` + uuid + `\s:\s` + capp + `:\s` + msg)
	// juniper host+uuid
	s.matcher[2] = regexp.MustCompile(pri + ts + `\s` + jhost + `\s` + uuid + `:\s` + japp + `:\s` + msg)
	// juniper host only
	s.matcher[3] = regexp.MustCompile(pri + ts + `\s` + jhost + `:\s` + japp + `:\s` + msg)
}

func (s *RFC3164) parse(raw []byte, result *map[string]interface{}) {
	for i, v := range s.matcher {
		m := v.FindStringSubmatch(string(raw))
		if len(m) == 0 {
			continue
		}
		pri, _ := strconv.Atoi(m[1])

		*result = map[string]interface{}{
			"priority":  pri,
			"timestamp": m[2],
		}
		if i == 0 || i == 1 || i == 3 {
			(*result)["identifier"] = m[4]
			(*result)["app"] = m[5]
			(*result)["message"] = m[6]
		} else if i == 2 {
			(*result)["identifier"] = m[5]
			(*result)["app"] = m[6]
			(*result)["message"] = m[7]
		}
		stats.Add("rfc3164Parsed", 1)
		return
	}
	stats.Add("rfc3164Unparsed", 1)
}
