package store

type StoreManager interface {
	Save(interface{}) error
}

func GetStoreManager(name string) (StoreManager, error) {
	f, ok := mapStoreManagerFactories[name]
	if !ok {
		return nil, ErrInvalidStoreManagerName
	}

	sm, err := f()
	if err != nil {
		return nil, ErrFailedToCreateStoreManager
	}

	return sm, nil
}

//////////////////////

type storeManagerFactory func() (StoreManager, error)

const ONEBOT_NAME = "one"

var mapStoreManagerFactories = map[string]storeManagerFactory{
	// "test": createOneGameSM,
	"one": createOneGameSM,
}
