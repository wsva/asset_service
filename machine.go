package main

type IPNote struct {
	IP   string `json:"ip"`
	Note string `json:"note"`
}

func (r *IPNote) Process() {
	r.Note = cleanNote(r.Note)
}
