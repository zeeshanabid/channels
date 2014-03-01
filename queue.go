package channels

const minQueueLen = 8

// A fast, ring-buffer queue based on the version suggested by Dariusz Górecki.
// Using this instead of a simple slice+append provides substantial memory and time
// benefits, and fewer GC pauses.
type queue struct {
	buf               []interface{}
	head, tail, count int
}

func newQueue() *queue {
	return &queue{buf: make([]interface{}, minQueueLen)}
}

func (q *queue) length() int {
	return q.count
}

func (q *queue) resize() {
	newBuf := make([]interface{}, q.count*2)

	if q.tail > q.head {
		copy(newBuf, q.buf[q.head:q.tail])
	} else {
		copy(newBuf, q.buf[q.head:len(q.buf)])
		copy(newBuf[len(q.buf)-q.head:], q.buf[:q.tail])
	}

	q.head = 0
	q.tail = q.count
	q.buf = newBuf
}

func (q *queue) enqueue(elem interface{}) {
	if q.count == len(q.buf) {
		q.resize()
	}

	q.buf[q.tail] = elem
	q.tail = (q.tail + 1) % len(q.buf)
	q.count++
}

func (q *queue) peek() interface{} {
	return q.buf[q.head]
}

func (q *queue) dequeue() {
	q.head = (q.head + 1) % len(q.buf)
	q.count--
	if len(q.buf) > minQueueLen && q.count*4 <= len(q.buf) {
		q.resize()
	}
}
