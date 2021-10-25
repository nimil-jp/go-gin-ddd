package context

func (c *ctx) Validate(request interface{}) (ok bool) {
	return c.verr.Validate(request)
}

func (c *ctx) FieldError(fieldName string, message string) {
	c.verr.Add(fieldName, message)
}

func (c *ctx) IsInValid() bool {
	return c.verr.IsInValid()
}

func (c *ctx) ValidationError() error {
	return c.verr
}
