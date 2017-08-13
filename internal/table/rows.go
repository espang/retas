package table

type Rows struct {
	current, next []int
	nextFn        func() ([]int, error)
	err           error
}

func (r *Rows) Next() []int {
	r.current = r.next
	r.next, r.err = r.nextFn()
	return r.current
}

func (r *Rows) Err() error {
	return r.err
}
