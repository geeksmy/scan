package gorm

import (
	"regexp"
	"scan/pkg/libc/test_tool"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

type Product struct {
	ID              uuid.UUID `gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name            string    `gorm:"type:varchar(128);not null"`
	Description     string    `gorm:"type:varchar(500);not null"`
	CloudProtocol   string    `gorm:"type:varchar(36);not null"`
	GatewayProtocol string    `gorm:"type:varchar(36)"`
	ProjectID       string    `gorm:"type:varchar(36);not null"`
	Type            int
	CreateAt        *time.Time
	UpdateAt        *time.Time
}

func (p *Product) TableName() string {
	return "product"
}

func TestPagination_ParamResultError(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	tx := test_tool.MockedGORMDBForTest(t, db)
	page, err := Pagination(tx, 0, 0, nil)
	assert.Error(t, err, ErrPageParamError.Error())
	assert.Empty(t, page)
}

func TestPagination_Error(t *testing.T) {
	db, dbMock, err := sqlmock.New()
	assert.NoError(t, err)

	var (
		limit, offset = -1, 0
		projectID     = uuid.New().String()
	)

	dbMock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "product" WHERE (project_id = $1)`)).WithArgs(projectID).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{"conut"}).AddRow(0))
	dbMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "product" WHERE (project_id = $1) OFFSET 0`)).
		WithArgs(sqlmock.AnyArg()).
		WillReturnError(gorm.ErrInvalidSQL)

	tx := test_tool.MockedGORMDBForTest(t, db)
	rows := make([]Product, 0)
	page, err := Pagination(tx, limit, offset, &rows)
	assert.Error(t, err)
	assert.Empty(t, page)
}

func TestPagination_LimitOffsetSuccess(t *testing.T) {
	db, dbMock, err := sqlmock.New()
	assert.NoError(t, err)
	var (
		limit, offset = -1, -1
		projectID     = uuid.New().String()
	)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "type", "cloud_protocol", "gateway_protocol", "project_id"}).
		AddRow(uuid.New().String(), "name1", "description1", 1, "mqtt", "mqtt", projectID).
		AddRow(uuid.New().String(), "name2", "description2", 2, "http", nil, projectID).
		AddRow(uuid.New().String(), "name3", "description3", 1, "http", "http", projectID)

	dbMock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "product"`)).WillReturnRows(sqlmock.NewRows(
		[]string{"conut"}).AddRow(3))
	dbMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "product" LIMIT 10 OFFSET 0`)).WillReturnRows(rows)

	tx := test_tool.MockedGORMDBForTest(t, db)
	result := make([]Product, 0)
	page, err := Pagination(tx, limit, offset, &result)
	assert.NoError(t, err)
	assert.NotEmpty(t, page)
	assert.Equal(t, 3, page.TotalRecord)
	assert.Equal(t, page.NextCursor, 0)
}

func TestPagination_Success(t *testing.T) {
	db, dbMock, err := sqlmock.New()
	assert.NoError(t, err)

	var (
		limit, offset = 1, 0
		projectID     = uuid.New().String()
	)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "type", "cloud_protocol", "gateway_protocol", "project_id"}).
		AddRow(uuid.New().String(), "name1", "description1", 1, "mqtt", "mqtt", projectID).
		AddRow(uuid.New().String(), "name2", "description2", 2, "http", nil, projectID).
		AddRow(uuid.New().String(), "name3", "description3", 1, "http", "http", projectID)

	dbMock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "product"`)).WillReturnRows(sqlmock.NewRows(
		[]string{"conut"}).AddRow(3))
	dbMock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "product" LIMIT 1 OFFSET 0`)).WillReturnRows(rows)

	tx := test_tool.MockedGORMDBForTest(t, db)
	result := make([]Product, 0)
	page, err := Pagination(tx, limit, offset, &result)
	assert.NoError(t, err)
	assert.NotEmpty(t, page)
	assert.Equal(t, 3, page.TotalRecord)
	assert.Equal(t, page.NextCursor, 1)
}
