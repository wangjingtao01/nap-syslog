package input

type Delimiter interface {
	Push(b byte) (string, bool)
	Vestige() (string, bool)
}
