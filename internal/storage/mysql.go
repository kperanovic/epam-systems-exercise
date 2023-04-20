package storage

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/kperanovic/epam-systems/api/v1/types"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mySQLStorage struct {
	conn *gorm.DB
}

func NewMySQLStorage() *mySQLStorage {
	return &mySQLStorage{}
}

func (m *mySQLStorage) Connect() error {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PWD"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_NAME"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	dbSql, err := db.DB()
	if err != nil {
		return err
	}

	dbSql.SetMaxIdleConns(1)
	dbSql.SetMaxOpenConns(10)

	// Migrate the database
	if err := db.Migrator().AutoMigrate(&types.CompanyType{}, &types.Company{}); err != nil {
		return err
	}

	m.conn = db

	return nil
}

func (m *mySQLStorage) SaveCompany(company *types.Company) error {
	res := m.conn.Create(company)

	return res.Error
}

func (m *mySQLStorage) GetCompany(id uuid.UUID) (*types.Company, error) {
	var company types.Company
	if err := m.conn.First(&company, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &company, nil
}

func (m *mySQLStorage) UpdateCompany(id uuid.UUID, company *types.Company) error {
	if err := m.conn.UpdateColumns(company).Error; err != nil {
		return err
	}

	return nil
}

func (m *mySQLStorage) DeleteCompany(id uuid.UUID) error {
	if err := m.conn.Where("id=?", id).Delete(&types.Company{}).Error; err != nil {
		return err
	}

	return nil
}
