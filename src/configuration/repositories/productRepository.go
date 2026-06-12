package repositories

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/qmessentials/qmessentials/configuration/models"
)

type ProductRepository interface {
	GetByPartNumber(ctx context.Context, partNumber string) (*models.Product, error)
}

type ProductRepositoryPG struct {
	db      *sql.DB
	typeMap *pgtype.Map
}

func NewProductRepositoryPG(db *sql.DB) *ProductRepositoryPG {
	return &ProductRepositoryPG{db, pgtype.NewMap()}
}

func (r *ProductRepositoryPG) GetByPartNumber(ctx context.Context, partNumber string) (*models.Product, error) {
	return r.getProductByPartNumber(ctx, partNumber)
}

func (r *ProductRepositoryPG) getProductByPartNumber(ctx context.Context, partNumber string) (*models.Product, error) {
	row := r.db.QueryRowContext(ctx, "select id, part_number, product_name, is_active, created_at, updated_at from products where part_number = $1", partNumber)
	var product models.Product
	err := row.Scan(&product.ID, &product.PartNumber, &product.ProductName, &product.IsActive, &product.CreatedAt, &product.UpdatedAt)
	if err != nil {
		return nil, err
	}
	product.ProductTestConfigurations, err = r.getProductTestConfigurations(ctx, product.ID)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryPG) getProductTestConfigurations(ctx context.Context, productId int) ([]models.ProductTestConfiguration, error) {
	rows, err := r.db.QueryContext(ctx, "select id, product_id, test_id, product_test_sequence, coalesce(specific_modifiers, '{}'::text[]) as specific_modifiers, unit, decimal_places, min_value, max_value, is_active, created_at, updated_at from product_test_configurations where product_id = $1", productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var productTestConfigurations []models.ProductTestConfiguration
	for rows.Next() {
		var productTestConfiguration models.ProductTestConfiguration
		err = rows.Scan(&productTestConfiguration.ID, &productTestConfiguration.ProductID, &productTestConfiguration.TestID, &productTestConfiguration.ProductTestSequence, r.typeMap.SQLScanner(&productTestConfiguration.SpecificModifiers), &productTestConfiguration.Unit, &productTestConfiguration.DecimalPlaces, &productTestConfiguration.MinValue, &productTestConfiguration.MaxValue, &productTestConfiguration.IsActive, &productTestConfiguration.CreatedAt, &productTestConfiguration.UpdatedAt)
		if err != nil {
			return nil, err
		}
		productTestConfigurations = append(productTestConfigurations, productTestConfiguration)
	}
	testIds := make([]int, 0)
	for _, productTestConfiguration := range productTestConfigurations {
		testIds = append(testIds, productTestConfiguration.TestID)
	}
	tests, err := r.getTests(ctx, testIds)
	if err != nil {
		return nil, err
	}
	testsById := make(map[int]models.Test)
	for _, test := range tests {
		testsById[test.ID] = test
	}
	for i := range productTestConfigurations {
		test := testsById[productTestConfigurations[i].TestID]
		productTestConfigurations[i].Test = &test
	}
	return productTestConfigurations, nil
}

func (r *ProductRepositoryPG) getTests(ctx context.Context, testIds []int) ([]models.Test, error) {
	rows, err := r.db.QueryContext(ctx, "select id, test_name, test_unit_category, coalesce(documentation_references, '{}'::text[]) as documentation_references, are_modifiers_allowed, is_active, created_at, updated_at from tests where id = any($1)", testIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tests []models.Test
	for rows.Next() {
		var test models.Test
		err = rows.Scan(&test.ID, &test.TestName, &test.TestUnitCategory, r.typeMap.SQLScanner(&test.DocumentationReferences), &test.AreModifiersAllowed, &test.IsActive, &test.CreatedAt, &test.UpdatedAt)
		if err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}
	return tests, nil
}
