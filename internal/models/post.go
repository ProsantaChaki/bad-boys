package models

import "time"

type Post struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null"`
	Title         string    `json:"title" gorm:"not null"`
	Description   string    `json:"description" gorm:"not null"`
	Address       string    `json:"address" gorm:"not null"`
	ContactName   string    `json:"contact_name" gorm:"not null"`
	MobileNumber  string    `json:"mobile_number" gorm:"not null"`
	IncidentDate  Date      `json:"incident_date" gorm:"not null"`
	Status        string    `json:"status" gorm:"not null;default:'active'"`
	IsAnonymous   bool      `json:"is_anonymous" gorm:"default:false"`
	Visibility    string    `json:"visibility" gorm:"not null;default:'public'"`
	AllowComments bool      `json:"allow_comments" gorm:"default:true"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	User          User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type PostHistory struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	PostID        uint      `json:"post_id" gorm:"not null"`
	UserID        uint      `json:"user_id" gorm:"not null"`
	Title         string    `json:"title" gorm:"not null"`
	Description   string    `json:"description" gorm:"not null"`
	Address       string    `json:"address" gorm:"not null"`
	ContactName   string    `json:"contact_name" gorm:"not null"`
	MobileNumber  string    `json:"mobile_number" gorm:"not null"`
	IncidentDate  Date      `json:"incident_date" gorm:"not null"`
	Status        string    `json:"status" gorm:"not null"`
	IsAnonymous   bool      `json:"is_anonymous"`
	Visibility    string    `json:"visibility" gorm:"not null"`
	AllowComments bool      `json:"allow_comments"`
	CreatedAt     time.Time `json:"created_at"`
}

type Report struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	PostID     uint      `json:"post_id" gorm:"not null"`
	ReporterID uint      `json:"reporter_id" gorm:"not null"`
	Reason     string    `json:"reason" gorm:"not null"`
	Status     string    `json:"status" gorm:"not null;default:'pending'"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Post       Post      `json:"post,omitempty" gorm:"foreignKey:PostID"`
	Reporter   User      `json:"reporter,omitempty" gorm:"foreignKey:ReporterID"`
}

type CreatePostRequest struct {
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description" binding:"required"`
	Address       string `json:"address" binding:"required"`
	ContactName   string `json:"contact_name" binding:"required"`
	MobileNumber  string `json:"mobile_number" binding:"required"`
	IncidentDate  Date   `json:"incident_date" binding:"required"`
	IsAnonymous   bool   `json:"is_anonymous"`
	Visibility    string `json:"visibility" binding:"required,oneof=public private"`
	AllowComments bool   `json:"allow_comments"`
}

type UpdatePostRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Address       string `json:"address"`
	ContactName   string `json:"contact_name"`
	MobileNumber  string `json:"mobile_number"`
	IncidentDate  Date   `json:"incident_date"`
	IsAnonymous   bool   `json:"is_anonymous"`
	Visibility    string `json:"visibility" binding:"omitempty,oneof=public private"`
	AllowComments bool   `json:"allow_comments"`
}

type CreateReportRequest struct {
	Reason string `json:"reason" binding:"required"`
}
