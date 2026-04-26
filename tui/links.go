package tui

type Links struct {
	items  []map[int]string
	active int // index into items
}

func NewLinks(items []map[int]string) Links {
	return Links{items: items}
}

func (l Links) Len() int {
	return len(l.items)
}

func (l Links) ActiveLineNumber() int {
	if len(l.items) == 0 {
		return -1
	}
	for k := range l.items[l.active] {
		return k
	}
	return -1
}

func (l Links) ActiveURL() string {
	if len(l.items) == 0 {
		return ""
	}
	for k := range l.items[l.active] {
		return l.items[l.active][k]
	}
	return ""
}

func (l *Links) Forward() bool {
	if l.active < len(l.items)-1 {
		l.active++
		return true
	}
	return false
}

func (l *Links) Backward() bool {
	if l.active > 0 {
		l.active--
		return true
	}
	return false
}

func (l *Links) Reset() {
	l.active = 0
}
