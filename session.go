package httpsession

import (
	"net/http"
	"time"
)

type Id string

type Store interface {
	Get(id Id, key string) interface{}
	Set(id Id, key string, value interface{})
	Del(id Id, key string) bool
	DelAll(id Id) bool
	Run() error
}

type Session struct {
	id    Id
	store Store
}

func (session *Session) Get(key string) interface{} {
	return session.store.Get(session.id, key)
}

func (session *Session) Set(key string, value interface{}) {
	session.store.Set(session.id, key, value)
}

func (session *Session) Del(key string) bool {
	return session.store.Del(session.id, key)
}

const (
	DefaultExpireTime = 30 * time.Minute
)

type Manager struct {
	store                  Store
	MaxAge                 int
	Path                   string
	generator              IdGenerator
	transfer               Transfer
	beforeReleaseListeners map[BeforeReleaseListener]bool
	afterCreatedListeners  map[AfterCreatedListener]bool
}

func Default() *Manager {
	store := NewMemoryStore(DefaultExpireTime)
	key := string(GenRandKey(16))
	return NewManager(store,
		NewSha1Generator(key),
		NewCookieTransfer("test"))
}

func NewManager(store Store, gen IdGenerator, transfer Transfer) *Manager {
	return &Manager{
		store:     store,
		generator: gen,
		transfer:  transfer,
	}
}

func (manager *Manager) Session(req *http.Request, rw http.ResponseWriter) *Session {
	id, err := manager.transfer.Get(req)
	if err != nil {
		// TODO:
		println("error:", err.Error())
		return nil
	}

	if !manager.generator.IsValid(id) {
		id = manager.generator.Gen(req)
		manager.transfer.Set(rw, id)
	}

	session := &Session{id: id, store: manager.store}
	// is exist?
	manager.afterCreated(session)
	return session
}

func (manager *Manager) Invalidate(rw http.ResponseWriter, session *Session) {
	manager.beforeReleased(session)
	manager.store.DelAll(session.id)
	manager.transfer.Clear(rw)
}

func (manager *Manager) afterCreated(session *Session) {
	for listener, _ := range manager.afterCreatedListeners {
		listener.OnAfterCreated(session)
	}
}

func (manager *Manager) beforeReleased(session *Session) {
	for listener, _ := range manager.beforeReleaseListeners {
		listener.OnBeforeRelease(session)
	}
}

func (manager *Manager) Run() error {
	return manager.store.Run()
}
