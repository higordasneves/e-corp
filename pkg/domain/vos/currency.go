package vos

type (
	Currency int64
)

func (c Currency) ConvertToCents() int64 {
	return int64(c * 100)
}
