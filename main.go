package main

import (
	"database/sql/driver"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	return nil, errors.New("Invalid product type value") //else is invalid
}

//Scan validate enum on read from data base
func (t *ProductType) Scan(value interface{}) error {
	var pt ProductType
	if value == nil {
		*t = ""
		return nil
	}
	st, ok := value.([]uint8) // if we declare db type as ENUM gorm will scan value as []uint8
	if !ok {
		return errors.New("Invalid data for product type")
	}
	pt = ProductType(string(st)) //convert type from string to ProductType

	switch pt {
	case IT, Decorate, Etc: //valid case
		*t = pt
		return nil
	}
	return fmt.Errorf("Invalid product type value :%s", st) //else is invalid
}

//Product product struct
type Product struct {
	gorm.Model
	Type  ProductType `gorm:"type:ENUM('it', 'decorate', 'etc')"`
	Code  string
	Price uint
}

func main() {

}
