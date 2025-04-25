package models

import (
	"database/sql"
	"errors"
)

// Product represents a product in the garage
type Product struct {
	ID          int     `json:"id" example:"1"`
	Name        string  `json:"name" example:"Hammer"`
	Description string  `json:"description,omitempty" example:"A sturdy hammer for construction"`
	Price       float64 `json:"price" example:"29.99"`
	ImagePath   string  `json:"image_path,omitempty" example:"/images/hammer.jpg"`
	HTMLContent string  `json:"html_content,omitempty" example:"<p>Product details in HTML</p>"`
}

type ProductModel struct {
	DB *sql.DB
}

func (m ProductModel) GetAll() ([]Product, error) {
	stmt := `SELECT id, name, description, price, image_path, html_content FROM products`
	
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImagePath, &product.HTMLContent)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (m ProductModel) Get(id int) (*Product, error) {
	stmt := `SELECT id, name, description, price, image_path, html_content FROM products WHERE id = $1`
	
	var product Product
	err := m.DB.QueryRow(stmt, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImagePath, &product.HTMLContent)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return &product, nil
}

func (m ProductModel) Create(product *Product) error {
	stmt := `
		INSERT INTO products (name, description, price, image_path, html_content)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	return m.DB.QueryRow(stmt, product.Name, product.Description, product.Price, product.ImagePath, product.HTMLContent).Scan(&product.ID)
}

func (m ProductModel) Update(product *Product) error {
	stmt := `
		UPDATE products 
		SET name = $1, description = $2, price = $3, image_path = $4, html_content = $5
		WHERE id = $6`

	result, err := m.DB.Exec(stmt, product.Name, product.Description, product.Price, product.ImagePath, product.HTMLContent, product.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
}

func (m ProductModel) Delete(id int) error {
	stmt := `DELETE FROM products WHERE id = $1`

	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("product not found")
	}

	return nil
} 