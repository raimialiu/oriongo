package entities

type Migration struct {
	Name      string `gorm:"type:varchar(255);primary_key"`
	Timestamp int64  `gorm:"type:bigint"`
}
