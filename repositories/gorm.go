package corerp

import (
	"log"
	"net"
	"os"
	"syscall"

	corecontainer "github.com/rayyone/go-core/container"
	"github.com/rayyone/go-core/helpers/method"
	"github.com/rayyone/go-core/ryerr"
	"gorm.io/gorm"
)

// CoreGormRepository Base Repo
type CoreGormRepository struct {
	BaseQuery   func(r corecontainer.RequestInf) *gorm.DB
	IsDebugging bool
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

// Create records by a given condition. out *interface
func (br *CoreGormRepository) Create(r corecontainer.RequestInf, out interface{}) (*gorm.DB, error) {
	tx := DefaultBaseQuery(r).Create(out)
	err := tx.Error
	if err != nil {
		err = ryerr.Newf("Base Repo [Create] Error: %s", err)
		if br.IsDebugging {
			return tx, err
		} else {
			return tx, ryerr.Msg(err, "Something went wrong. Please try again later")
		}
	}
	return tx, nil
}

// FindBy Find one record by a given condition
func (br *CoreGormRepository) FindBy(r corecontainer.RequestInf, out interface{}, where string, args ...interface{}) (*gorm.DB, error) {
	tx := br.BaseQuery(r).Where(where, args...).Find(out)
	err := GetFindByErrorType(tx.Error, br.IsDebugging)

	return tx, err
}

func (br *CoreGormRepository) FirstBy(r corecontainer.RequestInf, out interface{}, where string, args ...interface{}) (*gorm.DB, error) {
	tx := br.BaseQuery(r).Where(where, args...).First(out)
	return tx, tx.Error
}

// FindByID Find one record by ID
func (br *CoreGormRepository) FindByID(r corecontainer.RequestInf, out interface{}, id interface{}) (*gorm.DB, error) {
	tx := br.BaseQuery(r).Where("id = ?", id).First(out)
	return tx, tx.Error
}

// Update Update model. model *interface, field: not pointer
func (br *CoreGormRepository) Update(r corecontainer.RequestInf, model interface{}, fields interface{}) (*gorm.DB, error) {
	tx := DefaultBaseQuery(r).Model(model).Updates(fields)
	err := tx.Error
	if err != nil {
		err = ryerr.Newf("Base Repo [Update] Error: %s", err)
		if br.IsDebugging {
			return tx, err
		} else {
			return tx, ryerr.Msg(err, "Something went wrong. Please try again later")
		}
	}
	return tx, nil
}

// UpdateWhere Update by a given condition
func (br *CoreGormRepository) UpdateWhere(r corecontainer.RequestInf, model interface{}, fields interface{}, where string, args ...interface{}) (*gorm.DB, error) {
	tx := DefaultBaseQuery(r).Model(model).Where(where, args...).Updates(fields)
	err := tx.Error
	if err != nil {
		err = ryerr.Newf("Base Repo [UpdateWhere] Error: %s", err)
		if br.IsDebugging {
			return tx, err
		} else {
			return tx, ryerr.Msg(err, "Something went wrong. Please try again later")
		}
	}
	return tx, nil
}

// Save Update model if ID is present / Create if not. model *interface
func (br *CoreGormRepository) Save(r corecontainer.RequestInf, model interface{}) (*gorm.DB, error) {
	tx := DefaultBaseQuery(r).Save(model)
	err := tx.Error
	if err != nil {
		err = ryerr.Newf("Base Repo [Save] Error: %s", err)
		if br.IsDebugging {
			return tx, err
		} else {
			return tx, ryerr.Msg(err, "Something went wrong. Please try again later")
		}
	}
	return tx, nil
}

// DeleteWhere Delete by a given condition
func (br *CoreGormRepository) DeleteWhere(r corecontainer.RequestInf, model interface{}, where string, args ...interface{}) (*gorm.DB, error) {
	tx := DefaultBaseQuery(r).Where(where, args...).Delete(model)
	err := tx.Error
	if err != nil {
		err = ryerr.Newf("Base Repo [DeleteWhere] Error: %s", err)
		if br.IsDebugging {
			return tx, err
		} else {
			return tx, ryerr.Msg(err, "Something went wrong. Please try again later")
		}
	}
	return tx, nil
}

// ForceDeleteWhere Delete by a given condition & ignore soft deletes
func (br *CoreGormRepository) ForceDeleteWhere(r corecontainer.RequestInf, model interface{}, where string, args ...interface{}) (*gorm.DB, error) {
	tx := DefaultBaseQuery(r).Unscoped().Where(where, args...).Delete(model)
	err := tx.Error
	if err != nil {
		err = ryerr.Newf("Base Repo [ForceDeleteWhere] Error: %s", err)
		if br.IsDebugging {
			return tx, err
		} else {
			return tx, ryerr.Msg(err, "Something went wrong. Please try again later")
		}
	}
	return tx, nil
}

// Pluck model, out *[]interface
func (br *CoreGormRepository) Pluck(r corecontainer.RequestInf, model interface{}, out interface{}, col string, where string, args ...interface{}) (*gorm.DB, error) {
	tx := br.BaseQuery(r).Model(model).Where(where, args...).Pluck(col, out)
	err := tx.Error
	if err != nil {
		err = ryerr.Newf("Base Repo [Pluck] Error: %s", err)
		if br.IsDebugging {
			return tx, err
		} else {
			return tx, ryerr.Msg(err, "Something went wrong. Please try again later")
		}
	}
	return tx, nil
}

// Load Load relation
func (br *CoreGormRepository) Load(r corecontainer.RequestInf, model interface{}, out interface{}, rel string) error {
	err := br.BaseQuery(r).Model(model).Select("*").Association(rel).Find(out)
	if err != nil {
		err = ryerr.Newf("Base Repo [Load] Error: %s", err)
		if br.IsDebugging {
			return err
		} else {
			return ryerr.Msg(err, "Something went wrong. Please try again later")
		}
	}
	return err
}

func (br *CoreGormRepository) GetORM(r corecontainer.RequestInf) *gorm.DB {
	return r.GetDBM().GetTx()
}

// GetFindByErrorType get error type from error thrown from gorm find() method
func GetFindByErrorType(err error, isDebugging bool) error {
	if err != nil {
		opError, isOpError := err.(*net.OpError)
		if isOpError {
			if se, ok := opError.Err.(*os.SyscallError); ok {
				if se.Err == syscall.EPIPE {
					log.Printf("Error: Broken Pipe | %+v", err)
					if isDebugging {
						err = ryerr.Wrap(err, "Error: Broken Pipe")
					} else {
						err = ryerr.NewAndDontReport("Something went wrong. Please try again later")
					}
				} else if se.Err == syscall.ECONNRESET {
					log.Printf("Error: Connection Reset | %+v", err)
					if isDebugging {
						err = ryerr.Wrap(err, "Error: Connection Reset")
					} else {
						err = ryerr.NewAndDontReport("Something went wrong. Please try again later")
					}
				}
			}
		} else {
			_, _, fnName := method.TraceCaller(3)
			err = ryerr.Newf("Repo [%s] Error: %s", fnName, err)
			if isDebugging {
				return err
			} else {
				return ryerr.Msg(err, "Something went wrong. Please try again later")
			}
		}
	}
	return err
}
