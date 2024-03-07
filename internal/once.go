package internal

// Once is non-atomic version of [[sync.Once]]
type Once struct {
	done bool
}

func (o *Once) Do(do func()) {
	if !o.done {
		o.done = true
		do()
	}
}
