package database

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

type Users struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true"`
	Username    string `gorm:"unique"`
	Password    string `gorm:"not null"`
	FirstName   string `gorm:"type:varchar(100)"`
	LastName    string `gorm:"type:varchar(100)"`
	ImageURL    string
	PhoneNumber string
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type MRooms struct {
	ID              uint   `gorm:"primaryKey;autoIncrement:true"`
	RoomName        string `gorm:"type:varchar(100)"`
	RoomType        string `gorm:"type:varchar(20)"`
	RoomImageURL    string
	RoomDescription string `gorm:"type:varchar(150)"`
	CreatorID       uint
	Users           Users     `gorm:"foreignKey:CreatorID"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedBy       uint
	Users2          Users          `gorm:"foreignKey:UpdatedBy"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type TChats struct {
	ID             uint   `gorm:"primaryKey;autoIncrement:true"`
	RoomID         uint   `gorm:"index"`
	Room           MRooms `gorm:"foreignKey:RoomID"`
	SenderID       uint   `gorm:"index"`
	Users          Users  `gorm:"foreignKey:SenderID"`
	MessageType    string `gorm:"type:varchar(50);index"`
	MessageContent string
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type TChatReaders struct {
	ID           uint      `gorm:"primaryKey;autoIncrement:true"`
	ChatID       uint      `gorm:"index"`
	Chat         TChats    `gorm:"foreignKey:ChatID"`
	TargetID     uint      `gorm:"index"`
	ReaderTarget Users     `gorm:"foreignKey:TargetID"`
	ReadAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type TChatSents struct {
	ID         uint      `gorm:"primaryKey;autoIncrement:true"`
	ChatID     uint      `gorm:"index"`
	Chat       TChats    `gorm:"foreignKey:ChatID"`
	TargetID   uint      `gorm:"index"`
	SentTarget Users     `gorm:"foreignKey:TargetID"`
	SentAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func Connect() (bool, error) {
	dsn := "host=localhost user=postgres password=123 dbname=karyawandb port=5432"
	var error error
	DB, error = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if error != nil {
		return false, error
	}

	DB.AutoMigrate(&MRooms{}, &Users{}, &TChats{}, &TChatReaders{}, &TChatSents{})

	// robith := Karyawan{
	// 	Name: "Robith Syaukil Islam",
	// }

	// result := db.Create(&robith)

	// if result.Error != nil {
	// 	return false, result.Error
	// }

	return true, nil

}
