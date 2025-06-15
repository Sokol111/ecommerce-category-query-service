package category

import "time"

type Category struct {
	ID         string
	Version    int
	Name       string
	Enabled    bool
	CreatedAt  time.Time
	ModifiedAt time.Time
}
