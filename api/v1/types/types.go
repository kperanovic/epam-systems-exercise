package types

import "github.com/google/uuid"

// Company represents the data structure for REST API.
// Also structure is used as a schema for storage table.
type Company struct {
	ID          uuid.UUID `json:"uuid" binding:"required" gorm:"primaryKey"`
	Name        string    `json:"name" binding:"required,max=15" gorm:"size:15"`
	Description string    `json:"description" binding:"max=3000" gorm:"size:3000"`
	Employees   int       `json:"employees" binding:"required" gorm:"type:int"`
	Registered  bool      `json:"registered" binding:"required"`
	CompanyType int       `json:"companyType" binding:"required" gorm:"size:1"`
}

type CompanyType struct {
	ID   int    `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"size:50"`
	// NOTE: I was unable to add the foreign key wih the gorm.Automigrate(). Needs more attention.
	// Company Company `gorm:"foreignKey:CompanyType;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
}
