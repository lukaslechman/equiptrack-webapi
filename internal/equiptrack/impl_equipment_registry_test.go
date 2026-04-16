package equiptrack

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lukaslechman/equiptrack-webapi/internal/db_service"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DbServiceMock[DocType interface{}] struct {
	mock.Mock
}

func (m *DbServiceMock[DocType]) CreateDocument(ctx context.Context, id string, document *DocType) error {
	args := m.Called(ctx, id, document)
	return args.Error(0)
}

func (m *DbServiceMock[DocType]) FindDocument(ctx context.Context, id string) (*DocType, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*DocType), args.Error(1)
}

func (m *DbServiceMock[DocType]) FindDocuments(ctx context.Context, filter interface{}) ([]*DocType, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*DocType), args.Error(1)
}

func (m *DbServiceMock[DocType]) UpdateDocument(ctx context.Context, id string, document *DocType) error {
	args := m.Called(ctx, id, document)
	return args.Error(0)
}

func (m *DbServiceMock[DocType]) DeleteDocument(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *DbServiceMock[DocType]) Disconnect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type EquipmentSuite struct {
	suite.Suite
	dbServiceMock *DbServiceMock[Equipment]
}

func TestEquipmentSuite(t *testing.T) {
	suite.Run(t, new(EquipmentSuite))
}

func (suite *EquipmentSuite) SetupTest() {
	suite.dbServiceMock = &DbServiceMock[Equipment]{}

	// compile time check
	var _ db_service.DbService[Equipment] = suite.dbServiceMock

	suite.dbServiceMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return(
			&Equipment{
				Id:           "eq-001",
				Name:         "Röntgen EVO-3000",
				Category:     "imaging",
				Manufacturer: "Siemens",
				SerialNumber: "RX-2021-00123",
				PurchaseDate: "2021-03-15",
				Status:       "active",
				Location: Location{
					Building:   "Budova A",
					Department: "Rádiológia",
					Room:       "M-12",
				},
			},
			nil,
		)

	// FindDocuments – default odpoveď
	suite.dbServiceMock.
		On("FindDocuments", mock.Anything, mock.Anything).
		Return(
			[]*Equipment{
				{
					Id:     "eq-001",
					Name:   "Röntgen EVO-3000",
					Status: "active",
				},
				{
					Id:     "eq-002",
					Name:   "EKG Monitor",
					Status: "damaged",
				},
			},
			nil,
		)
}

func (suite *EquipmentSuite) Test_ListEquipment_FindDocumentsCalled() {
	// ARRANGE
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Request = httptest.NewRequest("GET", "/api/equipment", nil)

	sut := implEquipmentRegistryAPI{}

	// ACT
	sut.ListEquipment(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "FindDocuments", mock.Anything, mock.Anything)
	suite.Equal(200, recorder.Code)
}

func (suite *EquipmentSuite) Test_UpdateStatus_DbServiceUpdateCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("UpdateDocument", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	body := `{"status": "decommissioned"}`

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "equipmentId", Value: "eq-001"},
	}
	ctx.Request = httptest.NewRequest("PATCH", "/api/equipment/eq-001", strings.NewReader(body))

	sut := implEquipmentRegistryAPI{}

	// ACT
	sut.UpdateEquipmentStatus(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "UpdateDocument", mock.Anything, "eq-001", mock.Anything)
}

func (suite *EquipmentSuite) Test_GetEquipment_NotFound() {
	// ARRANGE
	// nový mock – nepotrebujeme default FindDocument zo SetupTest
	freshMock := &DbServiceMock[Equipment]{}
	freshMock.
		On("FindDocument", mock.Anything, mock.Anything).
		Return((*Equipment)(nil), db_service.ErrNotFound)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", freshMock)
	ctx.Params = []gin.Param{
		{Key: "equipmentId", Value: "nonexistent"},
	}
	ctx.Request = httptest.NewRequest("GET", "/api/equipment/nonexistent", nil)

	sut := implEquipmentRegistryAPI{}

	// ACT
	sut.GetEquipment(ctx)

	// ASSERT
	suite.Equal(404, recorder.Code)
}

// Test: DeleteEquipment zavolá DeleteDocument
func (suite *EquipmentSuite) Test_DeleteEquipment_DbServiceDeleteCalled() {
	// ARRANGE
	suite.dbServiceMock.
		On("DeleteDocument", mock.Anything, "eq-001").
		Return(nil)

	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set("db_service", suite.dbServiceMock)
	ctx.Params = []gin.Param{
		{Key: "equipmentId", Value: "eq-001"},
	}
	ctx.Request = httptest.NewRequest("DELETE", "/api/equipment/eq-001", nil)

	sut := implEquipmentRegistryAPI{}

	// ACT
	sut.DeleteEquipment(ctx)

	// ASSERT
	suite.dbServiceMock.AssertCalled(suite.T(), "DeleteDocument", mock.Anything, "eq-001")
	suite.Equal(204, recorder.Code)
}
