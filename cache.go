package cache

type node struct {
	key   string
	value int
	freq  int
	prev  *node
	next  *node
}

type dll struct {
	head *node
	tail *node
}

type store struct {
	cap     int
	keyMap  map[string]*node
	freqMap map[int]*dll
	minFreq int
}

func New(cap int) *store {
	return &store{
		cap:     cap,
		keyMap:  make(map[string]*node),
		freqMap: make(map[int]*dll),
	}
}

func (s *store) Get(key string) (int, bool) {
	n, ok := s.keyMap[key]
	if !ok {
		return 0, false
	}

	n.freq++
	s.removeFromOldFreqList(n)
	s.addToFreqList(n)
	s.updateMinFreqIfNeeded(n)

	return n.value, true
}

func (s *store) updateMinFreqIfNeeded(n *node) {
	oldFreq := n.freq - 1
	oldFreqList := s.freqMap[oldFreq]
	// list becomes empty after remove
	if oldFreqList.head.next == oldFreqList.tail {
		delete(s.freqMap, oldFreq)
		if oldFreq == s.minFreq {
			s.minFreq = n.freq
		}
	}
}

func (s *store) removeFromOldFreqList(n *node) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

func (s *store) addToFreqList(n *node) {
	ll, ok := s.freqMap[n.freq]
	if !ok {
		// use dummy head & tail nodes has 2 advantages:
		// - no edge cases: because all nodes now have a prev/next node, this
		// simplify insertion/removal
		// - head and tail are fixed, they don't move when nodes are added/removed
		s.freqMap[n.freq] = &dll{
			head: &node{},
			tail: &node{},
		}
		ll = s.freqMap[n.freq]
		ll.head.next = ll.tail
		ll.tail.prev = ll.head
	}
	n.next = ll.head.next
	ll.head.next.prev = n
	n.prev = ll.head
	ll.head.next = n
}

func (s *store) Set(key string, value int) {
	n, ok := s.keyMap[key]
	if ok {
		n.value = value
		n.freq++
		s.removeFromOldFreqList(n)
		s.addToFreqList(n)
		s.updateMinFreqIfNeeded(n)
		return
	}

	if len(s.keyMap) == s.cap {
		s.evict()
		// no need to update s.minFreq because the will-be inserted node below
		// will result in s.minFreq = 1
	}

	newNode := &node{
		key:   key,
		value: value,
		freq:  1,
	}
	s.keyMap[key] = newNode
	s.addToFreqList(newNode)
	s.minFreq = 1
}

func (s *store) evict() {
	evictedList := s.freqMap[s.minFreq]
	toEvict := evictedList.tail.prev

	s.removeFromOldFreqList(toEvict)
	delete(s.keyMap, toEvict.key)

	if evictedList.head.next == evictedList.tail {
		delete(s.freqMap, s.minFreq)
	}
}
