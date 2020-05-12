package main

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/sirupsen/logrus"
)

//ProductType product type
type ProductType string

//const available value for enum
const (
	IT       ProductType = "it"
	Decorate ProductType = "decorate"
	Etc      ProductType = "etc"
)

//Value validate enum when set to database
func (t ProductType) Value() (driver.Value, error) {
	switch t {
	case IT, Decorate, Etc: //valid case
		return string(t), nil
	}
	return nil, errors.New("Inalid product type value") //else is invalid

}

//Scan validate enum on read from data base
func (t *ProductType) Scan(value interface{}) error {
	var pt ProductType
	if value == nil {
		*t = ""
		return nil
	}
	st, ok := value.(string)
	if !ok {
		return errors.New("Invalid data for product type")
	}
	pt = ProductType(st) //convert type from string to ProductType

	switch pt {
	case IT, Decorate, Etc: //valid case
		*t = pt
		return nil
	}
	return fmt.Errorf("Inalid product type value :%s", st) //else is invalid
}

//Product product struct
type Product struct {
	gorm.Model
	Type  ProductType
	Code  string
	Price uint
}

func main() {

	logrus.SetLevel(logrus.DebugLevel)
	db, err := gorm.Open("sqlite3", "test.db")
	// db.LogMode(true)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	err = db.AutoMigrate(&Product{}).Error
	if err != nil {
		logrus.Errorf("Error on run migrate: %+v", err)
	} else {
		logrus.Infof("Migreate success")
	}

	// Create
	err = db.Create(&Product{Type: Decorate, Code: "L1212", Price: 1000}).Error
	if err != nil {
		logrus.Errorf("Error on run create: %+v", err)
	} else {
		logrus.Infof("Create success")
	}

	// List
	var products []Product
	err = db.Find(&products, "code = ?", "L1212").Error
	if err != nil {
		logrus.Errorf("Error on run first: %+v", err)
	} else {
		logrus.Infof("Find success: %+v", products)
	}

	// Read
	var product Product
	err = db.First(&product, "code = ?", "L1212").Error
	if err != nil {
		logrus.Errorf("Error on run first: %+v", err)
	} else {
		logrus.Infof("First success: %+v", product)
	}

	// Update - update product's price to 2000
	product.Type = "ib"
	product.Price = 2000
	err = db.Save(&product).Error
	if err != nil {
		logrus.Errorf("Error on run update: %+v", err)
	} else {
		logrus.Infof("Update success")
	}

	// Delete - delete product
	// err = db.Delete(&product).Error
	// if err != nil {
	// 	logrus.Errorf("Error on run delete: %+v", err)
	// } else {
	// 	logrus.Infof("Delete success")
	// }
}
