package storage

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"github.com/kperanovic/epam-systems/api/v1/types"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *mySQLStorage

func setDefaultEnv() {
	viper.SetDefault("DB_USER", "user")
	viper.SetDefault("DB_PWD", "pass")
	viper.SetDefault("DB_HOST", "127.0.0.1:3306")
	viper.SetDefault("DB_NAME", "epam")

	viper.AutomaticEnv()
}

func Clear(db *gorm.DB) {
	db.Raw("TRUNCATE companies")
}

func TestMain(m *testing.M) {
	setDefaultEnv()
	pool, err := dockertest.NewPool("")
	pool.MaxWait = time.Minute * 2
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "mysql",
		Tag:        "5.7",
		Env: []string{
			"MYSQL_ROOT_PASSWORD=pass",
			fmt.Sprintf("MYSQL_DATABASE=%s", viper.GetString("DB_NAME")),
			fmt.Sprintf("MYSQL_USER=%s", viper.GetString("DB_USER")),
			fmt.Sprintf("MYSQL_PASSWORD=%s", viper.GetString("DB_PWD")),
		},
		ExposedPorts: []string{"3306"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"3306": {
				{HostIP: "0.0.0.0", HostPort: "3306"},
			},
		},
	}
	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("could not start resource %s", err.Error())
	}

	if err = pool.Retry(func() error {
		db = NewMySQLStorage()
		err := db.Connect()
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestMySQLStorage_SaveCompany(t *testing.T) {
	setDefaultEnv()

	company := &types.Company{
		ID:          uuid.New(),
		Name:        "test-title",
		Description: "description",
		Employees:   1,
		Registered:  false,
		CompanyType: 1,
	}

	err := db.SaveCompany(company)
	assert.Equal(t, err, nil)

	got, err := db.GetCompany(company.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, got, company)

	Clear(db.conn)
}

func TestMySQLStorage_UpdateCompany(t *testing.T) {
	setDefaultEnv()

	company := &types.Company{
		ID:          uuid.New(),
		Name:        "test-title",
		Description: "description",
		Employees:   1,
		Registered:  false,
		CompanyType: 1,
	}

	// Insert new row
	err := db.SaveCompany(company)
	assert.Equal(t, err, nil)

	company.Description = "changed description"

	err = db.UpdateCompany(company.ID, company)
	assert.Equal(t, err, nil)

	// Check if update is correct
	got, err := db.GetCompany(company.ID)
	assert.Equal(t, err, nil)
	assert.Equal(t, company, got)

	Clear(db.conn)
}

func TestMySQLStorage_DeleteCompany(t *testing.T) {
	setDefaultEnv()

	company := &types.Company{
		ID:          uuid.New(),
		Name:        "test-title",
		Description: "description",
		Employees:   1,
		Registered:  false,
		CompanyType: 1,
	}

	// Insert new row
	err := db.SaveCompany(company)
	assert.Equal(t, err, nil)

	err = db.DeleteCompany(company.ID)
	assert.Equal(t, err, nil)

	// Check if update is correct
	got, err := db.GetCompany(company.ID)
	assert.Equal(t, err, gorm.ErrRecordNotFound)
	assert.Equal(t, got, nil)

	Clear(db.conn)
}
