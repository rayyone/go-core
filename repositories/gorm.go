package corerp

import (
	"log"
	"net"
	"os"
	"syscall"

	corecontainer "github.com/rayyone/go-core/container"
	"github.com/rayyone/go-core/errors"
	"github.com/rayyone/go-core/helpers/method"
	"github.com/jinzhu/gorm"
)

// CoreGormRepository Base Repo
type CoreGormRepository struct {
	BaseQuery func(r corecontainer.RequestInf) *gorm.DB
}

// NewCoreGormRepository Initiates new base repo
func NewCoreGormRepository() *CoreGormRepository {
	return &CoreGormRepository{
		BaseQuery: func(r corecontainer.RequestInf) *gorm.DB {
			return DefaultBaseQuery(r)
		},
	}
}

func DefaultBaseQuery(r corecontainer.RequestInf) *gorm.DB {
	return r.GetDBM().GetTx().
		Set("gorm:auto_preload", false).
		Set("gorm:association_autocreate", false).
		Set("gorm:association_autoupdate", false)
}

func (br *CoreGormRepository) ResetBaseQuery() *CoreGormRepository {
	return NewCoreGormRepository()
}

func (br *CoreGormRepository) Preload(column string, conditions ...interface{}) *CoreGormRepository {
	newCoreGormRepository := NewCoreGormRepository()
	newCoreGormRepository.BaseQuery = func(r corecontainer.RequestInf) *gorm.DB {
		return br.BaseQuery(r).Preload(column, conditions...)
	}

	return newCoreGormRepository
}

// Create Create records by a given condition. out *interface
func (br *CoreGormRepository) Create(r corecontainer.RequestInf, out interface{}) error {
	err := DefaultBaseQuery(r).Create(out).Error
	if err != nil {
		err = errors.Newf("Base Repo [Create] Error: %s", err)
		return errors.Msg(err, "Something went wrong. Please try again")
	}
	return nil
}

// FindBy Find one record by a given condition
func (br *CoreGormRepository) FindBy(r corecontainer.RequestInf, out interface{}, where string, args ...interface{}) error {
	err := br.BaseQuery(r).Where(where, args...).Find(out).Error
	return GetFindByErrorType(err)
}

// FindByID Find one record by ID
func (br *CoreGormRepository) FindByID(r corecontainer.RequestInf, out interface{}, id interface{}) error {
	return br.FindBy(r, out, "id = ?", id)
}

// GetBy Get multiple records by a given condition. If record not found, this won't return an error!!!!!
func (br *CoreGormRepository) GetBy(r corecontainer.RequestInf, out interface{}, where string, args ...interface{}) error {
	err := br.BaseQuery(r).Where(where, args...).Find(out).Error
	if err != nil {
		err = errors.Newf("Base Repo [GetBy] Error: %s", err)
		return errors.Msg(err, "Something went wrong. Please try again")
	}
	return nil
}

// Update Update model. model *interface, field: not pointer
func (br *CoreGormRepository) Update(r corecontainer.RequestInf, model interface{}, fields interface{}) error {
	err := DefaultBaseQuery(r).Model(model).Updates(fields).Error
	if err != nil {
		err = errors.Newf("Base Repo [Update] Error: %s", err)
		return errors.Msg(err, "Something went wrong. Please try again")
	}
	return nil
}

// UpdateWhere Update by a given condition
func (br *CoreGormRepository) UpdateWhere(r corecontainer.RequestInf, model interface{}, fields interface{}, where string, args ...interface{}) error {
	err := DefaultBaseQuery(r).Model(model).Where(where, args...).Updates(fields).Error
	if err != nil {
		err = errors.Newf("Base Repo [UpdateWhere] Error: %s", err)
		return errors.Msg(err, "Something went wrong. Please try again")
	}
	return nil
}

// Save Update model if ID is present / Create if not. model *interface
func (br *CoreGormRepository) Save(r corecontainer.RequestInf, model interface{}) error {
	err := DefaultBaseQuery(r).Save(model).Error
	if err != nil {
		err = errors.Newf("Base Repo [Save] Error: %s", err)
		return errors.Msg(err, "Something went wrong. Please try again")
	}
	return nil
}

// DeleteWhere Delete by a given condition
func (br *CoreGormRepository) DeleteWhere(r corecontainer.RequestInf, model interface{}, where string, args ...interface{}) error {
	err := DefaultBaseQuery(r).Where(where, args...).Delete(model).Error
	if err != nil {
		err = errors.Newf("Base Repo [DeleteWhere] Error: %s", err)
		return errors.Msg(err, "Something went wrong. Please try again")
	}
	return nil
}

// ForceDeleteWhere Delete by a given condition & ignore soft deletes
func (br *CoreGormRepository) ForceDeleteWhere(r corecontainer.RequestInf, model interface{}, where string, args ...interface{}) error {
	err := DefaultBaseQuery(r).Unscoped().Where(where, args...).Delete(model).Error
	if err != nil {
		err = errors.Newf("Base Repo [ForceDeleteWhere] Error: %s", err)
		return errors.Msg(err, "Something went wrong. Please try again")
	}
	return nil
}

// Pluck model, out *[]interface
func (br *CoreGormRepository) Pluck(r corecontainer.RequestInf, model interface{}, out interface{}, col string, where string, args ...interface{}) error {
	err := br.BaseQuery(r).Model(model).Where(where, args...).Pluck(col, out).Error
	if err != nil {
		err = errors.Newf("Base Repo [Pluck] Error: %s", err)
		return errors.Msg(err, "Something went wrong. Please try again")
	}
	return nil
}

// Load Load relation
func (br *CoreGormRepository) Load(r corecontainer.RequestInf, model interface{}, out interface{}, rel string) error {
	err := br.BaseQuery(r).Model(model).Association(rel).Find(out).Error
	if err != nil {
		err = errors.Newf("Base Repo [Load] Error: %s", err)
		return errors.Msg(err, "Something went wrong. Please try again")
	}
	return err
}

// GetFindByErrorType get error type from error throwed from gorm find() method
func GetFindByErrorType(err error) error {
	if err != nil {
		opError, isOpError := err.(*net.OpError)
		if gorm.IsRecordNotFoundError(err) {
			err = errors.NotFound.New("Data not found")
		} else if isOpError {
			if se, ok := opError.Err.(*os.SyscallError); ok {
				if se.Err == syscall.EPIPE {
					log.Printf("Error: Broken Pipe | %+v", err)
					err = errors.NewAndDontReport("Something went wrong. Please try again")
				} else if se.Err == syscall.ECONNRESET {
					log.Printf("Error: Connection Reset | %+v", err)
					err = errors.NewAndDontReport("Something went wrong. Please try again")
				}
			}
		} else {
			_, _, fnName := method.TraceCaller(3)
			err = errors.Newf("Repo [%s] Error: %s", fnName, err)
			err = errors.Msg(err, "Something went wrong. Please try again")
		}
	}
	return err
}
