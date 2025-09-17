package constants

type EntityStatus int

const (
	ACTIVE   EntityStatus = 1
	INACTIVE EntityStatus = 0
	DELETED  EntityStatus = -1
	ARCHIVED EntityStatus = -2
)

func (s EntityStatus) Int() int {
	return int(s)
}
