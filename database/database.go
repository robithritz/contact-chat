package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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
	RoomName        string `gorm:"type:varchar(100)" json:"roomName"`
	RoomType        string `gorm:"type:varchar(20)" json:"roomType"`
	RoomImageURL    string `json:"roomImageURL"`
	RoomDescription string `gorm:"type:varchar(150)"`
	CreatorID       uint
	Users           Users     `gorm:"foreignKey:CreatorID"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedBy       sql.NullInt32
	Users2          Users          `gorm:"foreignKey:UpdatedBy"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type MRoomParticipants struct {
	ID            uint           `gorm:"primaryKey;autoIncrement:true"`
	RoomId        uint           ``
	Room          MRooms         `gorm:"foreignKey:RoomId"`
	ParticipantId uint           ``
	Users         Users          `gorm:"foreignKey:ParticipantId"`
	CreatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type TChats struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoomID         uint      `gorm:"index"`
	Room           MRooms    `gorm:"foreignKey:RoomID"`
	SenderID       uint      `gorm:"index"`
	Users          Users     `gorm:"foreignKey:SenderID"`
	MessageType    string    `gorm:"type:varchar(50);index"`
	MessageContent string
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type TChatReaders struct {
	ID           uint      `gorm:"primaryKey;autoIncrement:true"`
	MessageId    uuid.UUID `gorm:"type:uuid;index"`
	Chat         TChats    `gorm:"foreignKey:MessageId"`
	TargetID     uint      `gorm:"index"`
	ReaderTarget Users     `gorm:"foreignKey:TargetID"`
	ReadAt       time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type TChatSents struct {
	ID         uint      `gorm:"primaryKey;autoIncrement:true"`
	MessageId  uuid.UUID `gorm:"type:uuid;index"`
	Chat       TChats    `gorm:"foreignKey:MessageId"`
	TargetID   uint      `gorm:"index"`
	SentTarget Users     `gorm:"foreignKey:TargetID"`
	SentAt     time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type LMQ struct {
	ID          uint           `gorm:"primaryKey;autoIncrement:true"`
	MQID        uuid.UUID      `gorm:"type:uuid"`
	TargetTopic string         `gorm:"type:text"`
	TargetID    uint           `gorm:"index"`
	Payload     string         `gorm:"type:text"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
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

	DB.AutoMigrate(&MRooms{}, &Users{}, &TChats{}, &TChatReaders{}, &TChatSents{}, &MRoomParticipants{}, &LMQ{})
	return true, nil
}
