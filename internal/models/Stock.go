package models
import "github.com/google/uuid"

type Stock struct {
    ID 			uuid.UUID		`gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    Symbol		string 		    `gorm:"type:varchar(64);"`
    Name		string		    `gorm:"type:varchar(128);"`
    Sector		string		    `gorm:"type:varchar(64);"`
    Price		int		        `gorm:"not null;"`
    IsArchive 	bool 		    `gorm:"not null;default:false;"`
}
