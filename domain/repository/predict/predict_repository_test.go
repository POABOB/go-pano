package predict_repository

// import (
// 	"database/sql"
// 	"fmt"
// 	"reflect"
// 	"regexp"
// 	"sync"
// 	"testing"
// 	"time"

// 	"go-pano/config"
// 	"go-pano/domain/model"
// 	"go-pano/utils"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func init() {
// 	gin.SetMode(gin.TestMode)
// 	config.LoadConfigTest()
// 	utils.Reset()
// }

// type PredictRepositorySuite struct {
// 	DB         *gorm.DB
// 	mock       sqlmock.Sqlmock
// 	repository IPredictRepository
// }

// // https://github.com/go-gorm/gorm/issues/3565

// func TestPredictRepository(test *testing.T) {
// 	s := &PredictRepositorySuite{}
// 	var (
// 		db  *sql.DB
// 		err error
// 	)
// 	db, s.mock, err = sqlmock.New()
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	dialector := mysql.New(mysql.Config{
// 		DSN:                       "sqlmock_db",
// 		DriverName:                "mysql",
// 		Conn:                      db,
// 		SkipInitializeWithVersion: true,
// 	})
// 	s.DB, err = gorm.Open(dialector, &gorm.Config{})

// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	var o sync.Once
// 	mockInstance := utils.NewMockInstance(s.DB, o)

// 	// 依賴注入
// 	s.repository = NewPredictRepository(mockInstance)

// 	now := time.Now()
// 	var getPredict *model.Predict = &model.Predict{
// 		ID:            1,
// 		ClinicId:      1,
// 		Filename:      "a.jpg",
// 		PredictString: `{"test":1}`,
// 		CreatedAt:     now,
// 		UpdatedAt:     now,
// 	}

// 	test.Run("成功：GetFirstByIDAndFileName，有找到資料。", func(test *testing.T) {
// 		s.mock.ExpectQuery(regexp.QuoteMeta(
// 			"SELECT * FROM `Predict` WHERE clinic_id=? AND filename=? ORDER BY `Predict`.`id` LIMIT 1")).
// 			WithArgs(getPredict.ClinicId, getPredict.Filename).
// 			WillReturnRows(sqlmock.NewRows([]string{"id", "clinic_id", "filename", "predict", "created_at", "updated_at"}).
// 				AddRow(getPredict.ID,
// 					getPredict.ClinicId,
// 					getPredict.Filename,
// 					getPredict.Predict,
// 					getPredict.CreatedAt,
// 					getPredict.UpdatedAt,
// 				))

// 		res, err := s.repository.GetFirstByIDAndFileName(getPredict.ClinicId, getPredict.Filename)

// 		assert.NoError(test, err)
// 		assert.True(test, reflect.DeepEqual(getPredict, res))
// 	})

// 	test.Run("成功：GetFirstByIDAndFileName，沒有找到資料。", func(test *testing.T) {
// 		s.mock.ExpectQuery(regexp.QuoteMeta(
// 			"SELECT * FROM `Predict` WHERE clinic_id=? AND filename=? ORDER BY `Predict`.`id` LIMIT 1")).
// 			WithArgs(getPredict.ClinicId, getPredict.Filename).
// 			WillReturnRows(sqlmock.NewRows([]string{"id", "clinic_id", "filename", "predict", "created_at", "updated_at"}))

// 		res, err := s.repository.GetFirstByIDAndFileName(getPredict.ClinicId, getPredict.Filename)

// 		assert.NoError(test, err)
// 		assert.True(test, reflect.DeepEqual(&model.Predict{}, res))
// 	})

// 	test.Run("成功：Create成功插入。", func(test *testing.T) {
// 		s.mock.ExpectBegin()
// 		s.mock.ExpectExec(
// 			regexp.QuoteMeta("INSERT INTO `Predict` (`clinic_id`,`filename`,`predict`,`created_at`,`updated_at`) VALUES (?,?,?,?,?)")).
// 			WithArgs(getPredict.ClinicId, getPredict.Filename, getPredict.Predict, sqlmock.AnyArg(), sqlmock.AnyArg()).
// 			WillReturnResult(sqlmock.NewResult(1, 1))
// 		s.mock.ExpectCommit()
// 		err := s.repository.Create(getPredict.ClinicId, getPredict.Filename, getPredict.PredictString)

// 		assert.NoError(test, err)
// 	})

// 	test.Run("失敗：Create主鍵重複插入。", func(test *testing.T) {
// 		s.mock.ExpectBegin()
// 		s.mock.ExpectExec(
// 			regexp.QuoteMeta("INSERT INTO `Predict` (`clinic_id`,`filename`,`predict`,`created_at`,`updated_at`) VALUES (?,?,?,?,?)")).
// 			WithArgs(getPredict.ClinicId, getPredict.Filename, getPredict.Predict, sqlmock.AnyArg(), sqlmock.AnyArg()).
// 			WillReturnResult(sqlmock.NewResult(0, 0))
// 		s.mock.ExpectCommit()
// 		err := s.repository.Create(getPredict.ClinicId, getPredict.Filename, getPredict.PredictString)

// 		assert.EqualError(test, err, "Key Conflict")
// 	})
// }
