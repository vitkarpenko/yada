package sqlite

type Swear struct {
	ID   int64  `gorm:"primarykey"`
	Word string `gorm:"index:unique"`
}
