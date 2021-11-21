package vos

type (
	Currency int64
)

//ConvertToCents converts currency unit to cents to solve floating point math problem
func (c Currency) ConvertToCents() int64 {
	return int64(c * 100)
}

//ConvertFromCents converts cents to currency unit
func (c *Currency) ConvertFromCents() {
	*c = *c / 100
}
