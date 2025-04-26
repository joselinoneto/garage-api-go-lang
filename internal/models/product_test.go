package models

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestProductModel_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := ProductModel{DB: db}

	// Test case 1: Successful retrieval
	t.Run("successful retrieval", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "image_path", "html_content"}).
			AddRow(1, "Hammer", "A sturdy hammer", 29.99, "/images/hammer.jpg", "<p>Hammer details</p>").
			AddRow(2, "Screwdriver", "A useful tool", 19.99, "/images/screwdriver.jpg", "<p>Screwdriver details</p>")

		mock.ExpectQuery("SELECT id, name, description, price, image_path, html_content FROM products").
			WillReturnRows(rows)

		products, err := model.GetAll()
		assert.NoError(t, err)
		assert.Len(t, products, 2)
		assert.Equal(t, "Hammer", products[0].Name)
		assert.Equal(t, "Screwdriver", products[1].Name)
	})

	// Test case 2: Database error
	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, description, price, image_path, html_content FROM products").
			WillReturnError(sql.ErrConnDone)

		products, err := model.GetAll()
		assert.Error(t, err)
		assert.Nil(t, products)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductModel_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := ProductModel{DB: db}

	// Test case 1: Successful retrieval
	t.Run("successful retrieval", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "price", "image_path", "html_content"}).
			AddRow(1, "Hammer", "A sturdy hammer", 29.99, "/images/hammer.jpg", "<p>Hammer details</p>")

		mock.ExpectQuery("SELECT id, name, description, price, image_path, html_content FROM products WHERE id = \\$1").
			WithArgs(1).
			WillReturnRows(rows)

		product, err := model.Get(1)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, "Hammer", product.Name)
		assert.Equal(t, 29.99, product.Price)
	})

	// Test case 2: Product not found
	t.Run("product not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, description, price, image_path, html_content FROM products WHERE id = \\$1").
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		product, err := model.Get(999)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Equal(t, "product not found", err.Error())
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductModel_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := ProductModel{DB: db}

	// Test case 1: Successful creation
	t.Run("successful creation", func(t *testing.T) {
		product := &Product{
			Name:        "Hammer",
			Description: "A sturdy hammer",
			Price:       29.99,
			ImagePath:   "/images/hammer.jpg",
			HTMLContent: "<p>Hammer details</p>",
		}

		rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
		mock.ExpectQuery("INSERT INTO products").
			WithArgs(product.Name, product.Description, product.Price, product.ImagePath, product.HTMLContent).
			WillReturnRows(rows)

		err := model.Create(product)
		assert.NoError(t, err)
		assert.Equal(t, 1, product.ID)
	})

	// Test case 2: Database error
	t.Run("database error", func(t *testing.T) {
		product := &Product{
			Name:        "Hammer",
			Description: "A sturdy hammer",
			Price:       29.99,
		}

		mock.ExpectQuery("INSERT INTO products").
			WithArgs(product.Name, product.Description, product.Price, product.ImagePath, product.HTMLContent).
			WillReturnError(sql.ErrConnDone)

		err := model.Create(product)
		assert.Error(t, err)
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductModel_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := ProductModel{DB: db}

	// Test case 1: Successful update
	t.Run("successful update", func(t *testing.T) {
		product := &Product{
			ID:          1,
			Name:        "Updated Hammer",
			Description: "An updated hammer",
			Price:       39.99,
			ImagePath:   "/images/updated-hammer.jpg",
			HTMLContent: "<p>Updated hammer details</p>",
		}

		mock.ExpectExec("UPDATE products").
			WithArgs(product.Name, product.Description, product.Price, product.ImagePath, product.HTMLContent, product.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := model.Update(product)
		assert.NoError(t, err)
	})

	// Test case 2: Product not found
	t.Run("product not found", func(t *testing.T) {
		product := &Product{
			ID:          999,
			Name:        "Non-existent Product",
		}

		mock.ExpectExec("UPDATE products").
			WithArgs(product.Name, product.Description, product.Price, product.ImagePath, product.HTMLContent, product.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := model.Update(product)
		assert.Error(t, err)
		assert.Equal(t, "product not found", err.Error())
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestProductModel_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	model := ProductModel{DB: db}

	// Test case 1: Successful deletion
	t.Run("successful deletion", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := model.Delete(1)
		assert.NoError(t, err)
	})

	// Test case 2: Product not found
	t.Run("product not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM products WHERE id = \\$1").
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := model.Delete(999)
		assert.Error(t, err)
		assert.Equal(t, "product not found", err.Error())
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
} 