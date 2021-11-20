package vos

type (
	AccountID string
)

func (accID AccountID) String() string {
	return string(accID)
}
