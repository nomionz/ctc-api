package repositories

import (
	"github.com/nomionz/ctc-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	List() ([]*models.Product, error)
	Get(id int) (*models.Product, error)
	Create(*models.Product) error
	Update(*models.Product) error
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
	err = db.AutoMigrate(&models.Product{})
	if err != nil {
		return nil, err
	}
	return &ProductRepository{db: db}, nil
}

func (pr *ProductRepository) List() ([]*models.Product, error) {
	var res []*models.Product
	err := pr.db.Model(&models.Product{}).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (pr *ProductRepository) Get(id int) (*models.Product, error) {
	var res *models.Product
	err := pr.db.Model(&models.Product{}).First(&res, id).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (pr *ProductRepository) Update(prod *models.Product) error {
	err := pr.db.Save(prod).Error
	if err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) Delete(id int) error {
	err := pr.db.Model(&models.Product{}).Delete(&models.Product{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (pr *ProductRepository) Create(prod *models.Product) error {
	err := pr.db.Model(&models.Product{}).Create(prod).Error
	if err != nil {
		return err
	}
	return nil
}
