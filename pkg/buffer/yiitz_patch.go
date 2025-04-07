package buffer

import (
	"fmt"
	"io"
	"sort"
	"sync/atomic"
	"unsafe"

	"gvisor.dev/gvisor/pkg/sync"
)

func (v *View) ReadFromOnce(r io.Reader, maxSize int) (n int, err error) {
	n, err = r.Read(v.chunk.data[:maxSize])
	if err == nil {
		v.write = n
	}
	return
}

var MonitoredPoolLevel = 1

type MonitoredPool struct {
	pool         sync.Pool
	currentObj   sync.Map
	currentCount int32
	maxCount     int32
	totalCount   int32
	size         int32
}

func (p *MonitoredPool) Get() any {
	x := p.pool.Get()
	switch MonitoredPoolLevel {
	case 0:
		p.currentObj.Store(unsafe.SliceData(x.([]byte)), GetStackTrace(10))
		fallthrough
	case 1:
		atomic.AddInt32(&p.totalCount, 1)
		n := atomic.AddInt32(&p.currentCount, 1)
		if n > p.maxCount {
			p.maxCount = n
		}
	}
	return x
}

func (p *MonitoredPool) Put(x any) {
	p.pool.Put(x)
	switch MonitoredPoolLevel {
	case 0:
		p.currentObj.Delete(unsafe.SliceData(x.([]byte)))
		fallthrough
	case 1:
		atomic.AddInt32(&p.currentCount, -1)
	}
}

var GetStackTrace func(int) string

func GetChunkPool(size int32) *MonitoredPool {
	return getChunkPool(int(size))
}

type kv struct {
	Key   string
	Value int
}

func WriteBytesPoolStats(w io.Writer) {
	fmt.Fprintln(w, "Bytes pool:")
	for idx := range numPools {
		p := &chunkPools[idx]
		cur2 := 0
		var items []kv
		if MonitoredPoolLevel == 0 {
			m := map[string]int{}
			p.currentObj.Range(func(key, value any) bool {
				cur2 += 1
				k := value.(string)
				if v, ok := m[k]; ok {
					m[k] = v + 1
				} else {
					m[k] = 1
				}
				return true
			})
			for k, v := range m {
				items = append(items, kv{k, v})
			}

			sort.Slice(items, func(i, j int) bool {
				return items[i].Value < items[j].Value
			})
		}

		fmt.Fprintf(w, "pool %d: cur %d, cur2 %d, max %d, total %d\n", p.size, p.currentCount, cur2, p.maxCount, p.totalCount)

		if MonitoredPoolLevel == 0 {
			for _, item := range items {
				fmt.Fprintf(w, "%s: %d\n", item.Key, item.Value)
			}
		}
	}
}
