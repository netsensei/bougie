package history

import "log"

type History struct {
	Entries  [20]string
	Position int
	Length   int
}

func (h *History) Forward() {
	h.Position++
	if h.Position > h.Length-1 {
		h.Position = h.Length - 1
	}
}

func (h *History) Backward() {
	h.Position--
	if h.Position < 0 {
		h.Position = 0
	}
}

func (h *History) Add(entry string) {
	if h.Position == h.Length-1 && h.Length < len(h.Entries) {
		h.Entries[h.Length] = entry
		h.Length++
		h.Position++
	} else if h.Position == h.Length-1 && h.Length == 20 {
		copy(h.Entries[:], h.Entries[1:])
		h.Entries[h.Length-1] = entry
	} else {
		h.Position++
		h.Length = h.Position + 1
		h.Entries[h.Position] = entry
	}
}

func (h *History) Current() string {
	if h.Length == 0 {
		return ""
	}

	log.Printf("entries: %v, position: %d", h.Entries, h.Position)
	return h.Entries[h.Position]
}
