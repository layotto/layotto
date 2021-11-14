package sequencer

import (
	"context"
	"errors"
	"mosn.io/layotto/components/sequencer"
	"mosn.io/pkg/log"
	"mosn.io/pkg/utils"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const defaultSize = 10000
const defaultLimit = 1000
const defaultRetry = 5
const waitTime = time.Millisecond * 10

// DoubleBuffer is double segment id buffer.
// There are two buffers in DoubleBuffer: InUseBuffer is in use, BackUpBuffer is a backup buffer.
// Their default capacity is 1000. When the InUseBuffer usage exceeds 30%, the BackUpBuffer will be initialized.
// When InUseBuffer is used up, swap them.
type DoubleBuffer struct {
	Key          string
	Size         int
	InUseBuffer  *Buffer
	BackUpBuffer *Buffer
	// 1 means getting BackUpBuffer，0 means not
	Processing uint32
	lock       sync.Mutex
	Store      sequencer.Store
}

type Buffer struct {
	from int64
	to   int64
}

func NewDoubleBuffer(key string, store sequencer.Store) *DoubleBuffer {

	d := &DoubleBuffer{
		Key:   key,
		Size:  defaultSize,
		Store: store,
	}

	return d
}

//init double buffer
func (d *DoubleBuffer) init() error {

	buffer, err := d.getNewBuffer()
	if err != nil {
		return err
	}

	d.InUseBuffer = buffer

	return nil
}

//getId next id
func (d *DoubleBuffer) getId() (int64, error) {

	d.lock.Lock()
	defer d.lock.Unlock()

	if d.InUseBuffer == nil {
		return 0, errors.New("[DoubleBuffer] Get error: InUseBuffer nil ")
	}
	//check swap
	if d.InUseBuffer.from > d.InUseBuffer.to {
		err := d.swap()
		if err != nil {
			return 0, err
		}
	}
	next := d.InUseBuffer.from
	d.InUseBuffer.from++

	//when InUseBuffer id more than half used, initialize BackUpBuffer.
	//equal make sure only one thread enter
	if d.InUseBuffer.to-d.InUseBuffer.from == defaultLimit {
		utils.GoWithRecover(func() {
			//one case: here and swap are performed simultaneously,
			//swap fast,no check will cover old BackUpBuffer

			//remove lock,add  atomic for visibility
			//check nil and not on processing
			if d.BackUpBuffer == nil && atomic.CompareAndSwapUint32(&d.Processing, 0, 1) {
				defer atomic.StoreUint32(&d.Processing, 0)
				buffer, err := d.getNewBuffer()
				//for visibility
				atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&d.BackUpBuffer)), unsafe.Pointer(buffer))
				if err != nil {
					log.DefaultLogger.Errorf("[DoubleBuffer] [getNewBuffer] error: %v", err)
					return
				}
			}
		}, nil)
	}

	return next, nil
}

//swap InUseBuffer and BackUpBuffer, must be locked
func (d *DoubleBuffer) swap() error {

	//check again and BackUpBuffer == nil
	if atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&d.BackUpBuffer))) == nil {
		//retry do processing
		for i := 0; i < defaultRetry; i++ {
			if atomic.CompareAndSwapUint32(&d.Processing, 0, 1) {
				defer atomic.StoreUint32(&d.Processing, 0)
				//double check
				if atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&d.BackUpBuffer))) != nil {
					break
				}
				buffer, err := d.getNewBuffer()
				if err != nil {
					return err
				}
				d.BackUpBuffer = &Buffer{}
				atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&d.BackUpBuffer)), unsafe.Pointer(buffer))
			} else {
				// sleep to avoid competition
				log.DefaultLogger.Infof("[DoubleBuffer] wait swap")
				time.Sleep(waitTime)
			}
		}
	}
	//if still nil, give up swap and return error
	if atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&d.BackUpBuffer))) == nil {
		return errors.New("swap error")
	}

	d.InUseBuffer = (*Buffer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&d.BackUpBuffer))))
	d.BackUpBuffer = nil
	return nil
}

//getNewBuffer return a new segment
func (d *DoubleBuffer) getNewBuffer() (*Buffer, error) {
	support, result, err := d.Store.GetSegment(&sequencer.GetSegmentRequest{
		Key:  d.Key,
		Size: d.Size,
	})
	if err != nil {
		return nil, err
	}
	if !support {
		return nil, errors.New("[DoubleBuffer] unSupport Segment id")
	}
	return &Buffer{
		from: result.From,
		to:   result.To,
	}, nil
}

var BufferCatch = map[string]*DoubleBuffer{}

//common lock is enough ？
var rwLock sync.RWMutex

func GetNextIdFromCache(ctx context.Context, store sequencer.Store, req *sequencer.GetNextIdRequest) (bool, int64, error) {

	// 1. check support
	support, _, _ := store.GetSegment(&sequencer.GetSegmentRequest{
		Key:  req.Key,
		Size: 0,
	})

	// return if not support
	if !support {
		return false, 0, nil
	}

	// 2. find the DoubleBuffer for this store and key
	var d *DoubleBuffer
	var err error
	if _, ok := BufferCatch[req.Key]; !ok {
		d, err = getDoubleBufferInWL(req.Key, store)
	} else {
		d, err = getDoubleBufferInRL(req.Key)
	}

	if err != nil {
		return true, 0, err
	}

	// 3. get the next id.
	// The buffer should automatically load segment into cache if the cache is (nearly) empty
	id, err := d.getId()

	if err != nil {
		return true, 0, err
	}

	return true, id, nil
}

//DoubleBuffer for this key not exist
func getDoubleBufferInWL(key string, store sequencer.Store) (*DoubleBuffer, error) {
	d := NewDoubleBuffer(key, store)
	rwLock.Lock()
	defer rwLock.Unlock()
	//double check
	if _, ok := BufferCatch[key]; ok {
		return BufferCatch[key], nil
	}
	err := d.init()
	if err != nil {
		return nil, err
	}
	BufferCatch[key] = d
	return d, nil
}

// DoubleBuffer for this key  exist
func getDoubleBufferInRL(key string) (*DoubleBuffer, error) {
	/*	rwLock.RLock()
		defer rwLock.RUnlock()*/
	return BufferCatch[key], nil
}
