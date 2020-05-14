package main

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	sdb, mock, err := sqlmock.New()
	db, err := gorm.Open("mysql", sdb)
	if err != nil {
		t.Errorf("Cannot setup db")
	}
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "type", "code", "price", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, string(IT), "code01", 40, now, now, nil).
		AddRow(2, string(Decorate), "code02", 45, now, now, nil).
		AddRow(3, "ib", "code03", 50, now, now, nil) //invalid record

	mock.ExpectQuery("SELECT (.+) FROM `products`").WillReturnRows(rows)

	products := []Product{}
	err = db.Find(&products).Error
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Invalid product type value :ib")
}

func TestFirst(t *testing.T) {
	sdb, mock, err := sqlmock.New()
	db, err := gorm.Open("mysql", sdb)
	if err != nil {
		t.Errorf("Cannot setup db")
	}
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "type", "code", "price", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "ib", "code01", 50, now, now, nil) //invalid record

	mock.ExpectQuery("SELECT (.+) FROM `products`").WillReturnRows(rows)

	product := Product{}
	err = db.First(&product).Error
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Invalid product type value :ib")
}

func TestCreate(t *testing.T) {
	sdb, mock, err := sqlmock.New()
	db, err := gorm.Open("mysql", sdb)
	if err != nil {
		t.Errorf("Cannot setup db")
	}
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "type", "code", "price", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "it", "code01", 50, now, now, nil)

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT (.+) FROM `products`").WillReturnRows(rows)
	mock.ExpectExec("UPDATE `products").WillReturnResult(sqlmock.NewResult(0, 1))

	tx := db.Begin()
	product := Product{}
	err = tx.First(&product, "id=?", 1).Error
	assert.Equal(t, nil, err)

	product.Type = "ib"
	err = tx.Save(&product).Error

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "Invalid product type value")
}
