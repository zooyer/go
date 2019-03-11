package fastdfs

import "sync"

// An instance of <code>Hashtable</code> has two parameters that affect its
// performance: <i>initial capacity</i> and <i>load factor</i>.  The
// <i>capacity</i> is the number of <i>buckets</i> in the hash table, and the
// <i>initial capacity</i> is simply the capacity at the time the hash table
// is created.  Note that the hash table is <i>open</i>: in the case of a "hash
// collision", a single bucket stores multiple entries, which must be searched
// sequentially.  The <i>load factor</i> is a measure of how full the hash
// table is allowed to get before its capacity is automatically increased.
// The initial capacity and load factor parameters are merely hints to
// the implementation.  The exact details as to when and whether the rehash
// method is invoked are implementation-dependent.<p>
//
// Generally, the default load factor (.75) offers a good tradeoff between
// time and space costs.  Higher values decrease the space overhead but
// increase the time cost to look up an entry (which is reflected in most
// <tt>Hashtable</tt> operations, including <tt>get</tt> and <tt>put</tt>).<p>
//
// The initial capacity controls a tradeoff between wasted space and the
// need for <code>rehash</code> operations, which are time-consuming.
// No <code>rehash</code> operations will <i>ever</i> occur if the initial
// capacity is greater than the maximum number of entries the
// <tt>Hashtable</tt> will contain divided by its load factor.  However,
// setting the initial capacity too high can waste space.<p>
//
// If many entries are to be made into a <code>Hashtable</code>,
// creating it with a sufficiently large capacity may allow the
// entries to be inserted more efficiently than letting it perform
// automatic rehashing as needed to grow the table. <p>
type Hashtable interface {
	Put(key, value interface{}) interface{}
	Get(key interface{}) interface{}
	Remove(key interface{})
	Size() int
	Keys() []interface{}
}

type hashtable struct {
	mutex    sync.RWMutex
	mapper   map[interface{}]interface{}
}

func NewHashtable() Hashtable {
	var h = new(hashtable)
	h.mapper = make(map[interface{}]interface{})

	return h
}

func (h *hashtable) Put(key, value interface{}) interface{} {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var old = h.mapper[key]
	h.mapper[key] = value

	return old
}

func (h *hashtable) Get(key interface{}) interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.mapper[key]
}

func (h *hashtable) Remove(key interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	delete(h.mapper, key)
}

func (h *hashtable) Size() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return len(h.mapper)
}

func (h *hashtable) Keys() []interface{} {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var keys = make([]interface{}, 0, len(h.mapper))
	for key,_ := range h.mapper {
		keys = append(keys, key)
	}

	return keys
}
