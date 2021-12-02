package vos

type (
	Currency float64
)

//ConvertToCents converts currency unit to cents to solve floating point math problem
func (c *Currency) ConvertToCents() {
	*c = *c * 100.0
}

//ConvertFromCents converts cents to currency unit
func (c *Currency) ConvertFromCents() {
	*c = *c / 100.0
}
