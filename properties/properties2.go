package properties

import (
	"container/list"
	"sync"
)

type hashtable2 struct {
	mutex    sync.Mutex
	mapper   map[interface{}]interface{}
	element  map[interface{}]*list.Element
	lister   *list.List
}

func NewHashtable2() Hashtable {
	return new(hashtable2)
}

func (h *hashtable2) Init() {
	h.mapper = make(map[interface{}]interface{})
	h.element = make(map[interface{}]*list.Element)
	h.lister = list.New()
}

func (h *hashtable2) Put(key, value interface{}) interface{} {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	old := h.mapper[key]
	h.mapper[key] = value
	h.element[key] = h.lister.PushBack(key)

	return old
}

func (h *hashtable2) Get(key interface{}) interface{} {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.mapper[key]
}

func (h *hashtable2) Remove(key interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.mapper, key)
	h.lister.Remove(h.element[key])
	delete(h.element, key)
}

func (h *hashtable2) Size() int {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	return h.lister.Len()
}

func (h *hashtable2) Keys() []interface{} {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var keys = make([]interface{}, 0, h.lister.Len())
	for node := h.lister.Front(); node != nil; node = node.Next() {
		keys = append(keys, node.Value)
	}

	return keys
}


type Properties2 struct {
	Properties
}

func NewProperties2() *Properties2 {
	properties := *NewProperties()
	properties.Hashtable = NewHashtable2()
	properties.Hashtable.Init()
	return &Properties2{
		Properties:properties,
	}
}