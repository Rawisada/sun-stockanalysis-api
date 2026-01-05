package stock
type Stock struct {
    ID 			uint64		`gorm:"primaryKeyวautoIncrement;"`
    Symbol		string 		`gorm:"type:varchar(64);"`
    Name		string		`gorm:"type:varchar(128);"`
    Sector		string		`gorm:"type:varchar(64);"`
    Price		uint		`gorm:"not null;"`
    IsArchive 	bool 		`gorm:"not null;default:false;"`
}
