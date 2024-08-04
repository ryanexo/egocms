package data

import (
	"runtime"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type DBConfig struct {
	Type    string `json:"type,omitempty"`
	Host    string `json:"host,omitempty"`
	Name    string `json:"name,omitempty"`
	User    string `json:"user,omitempty"`
	Pass    string `json:"pass,omitempty"`
	Charset string `json:"charset,omitempty"`
}

func DefaultConfig() *gorm.Config {
	return &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		FullSaveAssociations:   false,
		Logger:                 logger.Default.LogMode(logger.Silent),
		PrepareStmt:            true,
		AllowGlobalUpdate:      false,
		QueryFields:            true,
		SkipDefaultTransaction: true,

		DisableForeignKeyConstraintWhenMigrating: true,
	}
}

func buildDB(dialector gorm.Dialector, config *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, config)
	if err != nil {
		return nil, err
	}
	sqlDB, sqlDBErr := db.DB()
	if sqlDBErr != nil {
		return nil, err
	}

	sqlDB.SetConnMaxLifetime(time.Second * 30)
	sqlDB.SetConnMaxIdleTime(time.Minute)
	sqlDB.SetMaxOpenConns(runtime.NumCPU() * 2)

	return db, nil
}

func NewDB(config *DBConfig) (*gorm.DB, error) {
	driver, err := GetDriver(config)
	if err != nil {
		return nil, err
	}
	return buildDB(driver, DefaultConfig())
}

func NewDBMock() (*gorm.DB, sqlmock.Sqlmock) {
	sqlDB, sqlm, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	sqlm.ExpectQuery(`SELECT VERSION()`).WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("8.0.0"))
	db, err := buildDB(mysql.New(mysql.Config{Conn: sqlDB}), DefaultConfig())
	if err != nil {
		panic(err)
	}
	return db, sqlm
}
