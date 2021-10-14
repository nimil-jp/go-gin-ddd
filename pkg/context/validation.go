package context

func (c *ctx) FieldError(fieldName string, message string) {
	c.verr.Add(fieldName, message)
}

func (c *ctx) IsInValid() bool {
	return c.verr.IsInValid()
}

func (c *ctx) ValidationError() error {
	return c.verr
}
