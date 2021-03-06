package packet

import (
	"fmt"
	"io"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/chtisgit/go-flows/flows"
)

type decodeStats struct {
	decodeError uint64
	keyError    uint64
}

// EventTable represents a flow table that can handle multiple events in one go
type EventTable interface {
	// EOF expires all the flows in the table at the given point in time with EOF as end reason
	EOF(flows.DateTimeNanoseconds)
	// Print table statistics to the given writer
	PrintStats(io.Writer)
	usage() []bufferUsage
	event(buffer *shallowMultiPacketBuffer)
	flush()
	getDecodeStats() *decodeStats
	getSelector() DynamicKeySelector
}

type baseTable struct {
	selector DynamicKeySelector
	autoGC   bool
}

func (bt baseTable) getSelector() DynamicKeySelector {
	return bt.selector
}

type parallelFlowTable struct {
	baseTable
	tables      []*flows.FlowTable
	expirewg    sync.WaitGroup
	buffers     []*shallowMultiPacketBufferRing
	tmp         []*shallowMultiPacketBuffer
	wg          sync.WaitGroup
	usageBuffer []bufferUsage
	decodeStats decodeStats
	expireTime  flows.DateTimeNanoseconds
	nextExpire  flows.DateTimeNanoseconds
}

type singleFlowTable struct {
	baseTable
	table       *flows.FlowTable
	buffer      *shallowMultiPacketBufferRing
	done        chan struct{}
	usageBuffer [1]bufferUsage
	decodeStats decodeStats
	expireTime  flows.DateTimeNanoseconds
	nextExpire  flows.DateTimeNanoseconds
}

func (sft *singleFlowTable) PrintStats(w io.Writer) {
	fmt.Fprintf(w,
		`Decode statistics:
	decode errors: %d
	key function rejects: %d
`, sft.decodeStats.decodeError, sft.decodeStats.keyError)
	fmt.Fprintf(w,
		`Table statistics:
	flows: %d
	peak flows: %d
`, sft.table.Stats.Flows, sft.table.Stats.Maxflows)
}

func (sft *singleFlowTable) getDecodeStats() *decodeStats {
	return &sft.decodeStats
}

func (sft *singleFlowTable) event(buffer *shallowMultiPacketBuffer) {
	current := buffer.Timestamp()
	expire := false
	if current > sft.nextExpire {
		expire = true
		sft.nextExpire = current + sft.expireTime
	}
	if buffer.empty() && !expire {
		return
	}
	b, _ := sft.buffer.popEmpty()
	buffer.Copy(b)
	b.expire = expire
	b.timestamp = current
	b.finalize()
}

func (sft *singleFlowTable) flush() {
	close(sft.buffer.full)
	<-sft.done
}

func (sft *singleFlowTable) EOF(now flows.DateTimeNanoseconds) {
	sft.table.EOF(now)
}

func (sft *singleFlowTable) usage() []bufferUsage {
	sft.usageBuffer[0] = sft.buffer.usage()
	return sft.usageBuffer[:]
}

// NewFlowTable creates a new flowtable with the given record list, a flow creator, flow options,
// expire time, a key selector, if empty values in the key are allowed and if automatic gc should be used.
//
// num specifies the number of parallel flow tables.
func NewFlowTable(num int, features flows.RecordListMaker, newflow flows.FlowCreator, options flows.FlowOptions, expire flows.DateTimeNanoseconds, selector DynamicKeySelector, autoGC bool) EventTable {
	bt := baseTable{
		selector: selector,
		autoGC:   autoGC,
	}
	if num == 1 {
		ret := &singleFlowTable{
			baseTable:  bt,
			table:      flows.NewFlowTable(features, newflow, options, selector.fivetuple && options.TCPExpiry, 0),
			expireTime: expire,
		}
		ret.buffer = newShallowMultiPacketBufferRing(fullBuffers, batchSize)
		ret.done = make(chan struct{})
		go func() {
			t := ret.table
			defer close(ret.done)
			for buffer := range ret.buffer.full {
				// fixup buffer stats -> not possible with popFull due to select
				atomic.AddInt32(&ret.buffer.currentBuffers, -1)
				atomic.AddInt32(&ret.buffer.currentPackets, -int32(buffer.windex))
				for {
					b := buffer.read()
					if b == nil {
						break
					}
					t.Event(b)
				}
				if buffer.expire {
					t.Expire(buffer.timestamp)
					if !autoGC {
						go runtime.GC()
					}
				}
				buffer.recycle()
			}
		}()
		return ret
	}
	if num > 256 {
		panic("Maximum of 256 tables allowed")
	}
	ret := &parallelFlowTable{
		baseTable:   bt,
		tables:      make([]*flows.FlowTable, num),
		buffers:     make([]*shallowMultiPacketBufferRing, num),
		usageBuffer: make([]bufferUsage, num),
		tmp:         make([]*shallowMultiPacketBuffer, num),
		expireTime:  expire,
	}
	for i := 0; i < num; i++ {
		c := newShallowMultiPacketBufferRing(fullBuffers, batchSize)
		ret.buffers[i] = c
		t := flows.NewFlowTable(features, newflow, options, selector.fivetuple && options.TCPExpiry, uint8(i))
		ret.tables[i] = t
		ret.wg.Add(1)
		go func() {
			defer ret.wg.Done()
			for buffer := range c.full {
				// fixup buffer stats -> not possible with popFull due to select
				atomic.AddInt32(&c.currentBuffers, -1)
				atomic.AddInt32(&c.currentPackets, -int32(buffer.windex))
				for {
					b := buffer.read()
					if b == nil {
						break
					}
					t.Event(b)
				}
				if buffer.expire {
					t.Expire(buffer.timestamp)
					if !autoGC {
						ret.expirewg.Done()
					}
				}
				buffer.recycle()
			}
		}()
	}
	return ret
}

func (pft *parallelFlowTable) usage() []bufferUsage {
	for i, buffer := range pft.buffers {
		pft.usageBuffer[i] = buffer.usage()
	}
	return pft.usageBuffer
}

func (pft *parallelFlowTable) PrintStats(w io.Writer) {
	fmt.Fprintf(w,
		`Decode statistics:
	decode errors: %d
	key function rejects: %d
`, pft.decodeStats.decodeError, pft.decodeStats.keyError)
	fmt.Fprintln(w, "Table statistics:")
	var sumPackets, sumFlows uint64
	for _, table := range pft.tables {
		sumPackets += table.Stats.Packets
		sumFlows += table.Stats.Flows
	}
	for i, table := range pft.tables {
		fmt.Fprintf(w, `	Table #%d:
		packets: %d (%2.2f)
		flows: %d (%2.2f)
		peak flows: %d
`, i+1, table.Stats.Packets, float64(table.Stats.Packets)/float64(sumPackets)*100, table.Stats.Flows, float64(table.Stats.Flows)/float64(sumFlows)*100, table.Stats.Maxflows)
	}
}

func (pft *parallelFlowTable) getDecodeStats() *decodeStats {
	return &pft.decodeStats
}

// event partitions all packets from buffer over the tables based on fnvHash of the flow key.
// Expiry is carried out in every table at about the same time (a flag in the buffer is set - after this buffer is handle the table does expiry).
// event waits for every table to finish expiry and does GC afterwards - this hurts concurrency a bit, but improves memory usage.
func (pft *parallelFlowTable) event(buffer *shallowMultiPacketBuffer) {
	current := buffer.Timestamp()
	expire := false
	if current > pft.nextExpire {
		pft.expirewg.Add(len(pft.buffers))
		expire = true
		pft.nextExpire = current + pft.expireTime
	}
	if buffer.empty() && !expire {
		return
	}

	tmp := pft.tmp[:len(pft.buffers)]
	for i, buf := range pft.buffers {
		tmp[i], _ = buf.popEmpty()
		tmp[i].timestamp = current
		tmp[i].expire = expire
	}
	for {
		b := buffer.read()
		if b == nil {
			break
		}
		h := fnvHash(b.Key()) % uint64(len(tmp))
		tmp[h].push(b)
	}
	for _, buf := range pft.tmp {
		buf.finalize()
	}
	if expire && !pft.autoGC {
		pft.expirewg.Wait()
		go func() {
			runtime.GC()
		}()
	}
}

func (pft *parallelFlowTable) flush() {
	for _, c := range pft.buffers {
		close(c.full)
	}
	pft.wg.Wait()
}

func (pft *parallelFlowTable) EOF(now flows.DateTimeNanoseconds) {
	for _, t := range pft.tables {
		pft.wg.Add(1)
		go func(table *flows.FlowTable) {
			defer pft.wg.Done()
			table.EOF(now)
		}(t)
	}
	pft.wg.Wait()
}
