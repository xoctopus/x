package cleanup

import (
	"errors"
	"sync"
)

type Collector interface {
	Collect(f func() error)
	JoinTo(dst *error) error
}

func NewCollector() Collector {
	return &collector{}
}

type collector struct {
	fns   []func() error
	once  sync.Once
	final error
	mu    sync.Mutex
}

func (c *collector) JoinTo(dst *error) error {
	c.once.Do(func() {
		errs := make([]error, 0, len(c.fns))

		for i := len(c.fns) - 1; i >= 0; i-- {
			if c.fns[i] != nil {
				if err := c.fns[i](); err != nil {
					errs = append(errs, err)
				}
			}
		}
		c.fns = nil

		var err error
		if dst != nil && *dst != nil {
			err = *dst
		}

		if err != nil {
			errs = append([]error{err}, errs...)
		}
		if len(errs) == 0 {
			c.final = err
		} else {
			c.final = errors.Join(errs...)
		}

		if dst != nil {
			*dst = c.final
		}
	})
	return c.final
}

func (c *collector) Collect(f func() error) {
	if f != nil {
		c.mu.Lock()
		c.fns = append(c.fns, f)
		c.mu.Unlock()
	}
}
