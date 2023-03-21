package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	List() ([]Product, error)
	Get(id int) (*Product, error)
	Create(Product) error
	Update(Product) error
	Delete(id int) error
}

var _ Repository = &ProductRepository{}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(conn string) (*ProductRepository, error) {
	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Product{})
	if err != nil {
		return nil, err
	}
	return &ProductRepository{db: db}, nil
}

func (pr *ProductRepository) List() ([]Product, error) {
	var res []Product
	err := pr.db.Find(&Product{}).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (pr *ProductRepository) Get(id int) (*Product, error) {
	var res *Product
	err := pr.db.First(res, id).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (pr *ProductRepository) Update(prod Product) error {
	return pr.db.Save(prod).Error
}

func (pr *ProductRepository) Delete(id int) error {
	return pr.db.Delete(&Product{}, id).Error
}

func (pr *ProductRepository) Create(prod Product) error {
	return pr.db.Create(prod).Error
}
