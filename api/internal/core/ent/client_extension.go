package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
)

func (c *Client) Ping(ctx context.Context) error {
	driver, ok := c.driver.(*sql.Driver)
	if !ok {
		return fmt.Errorf("driver is not SQL driver")
	}
	return driver.DB().PingContext(ctx)
}
