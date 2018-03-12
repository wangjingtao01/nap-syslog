package input

type RFC interface {
	compileMatcher()
	parse([]byte, *map[string]interface{})
}
