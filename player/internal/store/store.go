package store

type StoreManager interface {
	Save(interface{}) error
}

func GetStoreManager(name string) (StoreManager, error) {
	if f, ok := mapStoreManagerFactories[name]; ok {
		if sm, err := f(); err != nil {
			return sm, nil
		}

		return nil, ErrFailedToCreateStoreManager
	}

	return nil, ErrInvalidStoreManagerName
}

//////////////////////

type storeManagerFactory func() (StoreManager, error)

var mapStoreManagerFactories = map[string]storeManagerFactory{
	"test": createOneGameSM,
}
