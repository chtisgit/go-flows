package staging

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	ipfix "github.com/CN-TU/go-ipfix"

	"github.com/chtisgit/go-flows/flows"
	"github.com/chtisgit/go-flows/modules/features"
	"github.com/chtisgit/go-flows/packet"
)

/*

Features in here are subject to change. Use them with caution.

*/

////////////////////////////////////////////////////////////////////////////////

type _characters struct {
	flows.BaseFeature
	time flows.DateTimeNanoseconds
	src  []byte
}

func (f *_characters) Event(new interface{}, context *flows.EventContext, src interface{}) {
	var time flows.DateTimeNanoseconds
	if f.time != 0 {
		time = context.When() - f.time
	}
	if len(f.src) == 0 {
		tmpSrc, _ := new.(packet.Buffer).NetworkLayer().NetworkFlow().Endpoints()
		f.src = tmpSrc.Raw()
	}
	f.time = context.When()
	newTime := int(time / (100 * flows.MillisecondsInNanoseconds)) // time is now in deciseconds

	srcIP, _ := new.(packet.Buffer).NetworkLayer().NetworkFlow().Endpoints()

	var buffer bytes.Buffer
	if bytes.Equal(f.src, srcIP.Raw()) {
		buffer.WriteString("A")
	} else {
		buffer.WriteString("B")
	}

	buffer.WriteString(strings.Repeat("_", newTime))

	f.SetValue(buffer.String(), context, f)
}

func init() {
	flows.RegisterTemporaryFeature("_characters", "returns a textual representation of a packet", ipfix.StringType, 0, flows.PacketFeature, func() flows.Feature { return &_characters{} }, flows.RawPacket)
}

////////////////////////////////////////////////////////////////////////////////

type _characters2 struct {
	flows.BaseFeature
	time flows.DateTimeNanoseconds
	src  []byte
}

func (f *_characters2) Event(new interface{}, context *flows.EventContext, src interface{}) {
	var time flows.DateTimeNanoseconds
	if f.time != 0 {
		time = context.When() - f.time
	}
	if len(f.src) == 0 {
		tmpSrc, _ := new.(packet.Buffer).NetworkLayer().NetworkFlow().Endpoints()
		f.src = tmpSrc.Raw()
	}
	f.time = context.When()
	newTime := int(time / (100 * flows.MillisecondsInNanoseconds)) // time is now in deciseconds

	srcIP, _ := new.(packet.Buffer).NetworkLayer().NetworkFlow().Endpoints()
	tcp := features.GetTCP(new)
	if tcp == nil {
		return
	}

	var buffer bytes.Buffer
	if bytes.Equal(f.src, srcIP.Raw()) {
		// A->B
		if tcp.FIN && tcp.ACK {
			buffer.WriteString("n")
		} else if tcp.FIN {
			buffer.WriteString("f")
		} else if tcp.SYN && tcp.ACK {
			buffer.WriteString("k")
		} else if tcp.SYN {
			buffer.WriteString("s")
		} else if tcp.RST && tcp.ACK {
			buffer.WriteString("t")
		} else if tcp.RST {
			buffer.WriteString("r")
		} else if tcp.PSH && tcp.ACK {
			buffer.WriteString("h")
		} else if tcp.PSH {
			buffer.WriteString("p")
		} else if tcp.ACK {
			buffer.WriteString("a")
		} else if tcp.URG {
			buffer.WriteString("u")
		} else {
			buffer.WriteString("o")
		}
	} else {
		// B->A
		if tcp.FIN && tcp.ACK {
			buffer.WriteString("N")
		} else if tcp.FIN {
			buffer.WriteString("F")
		} else if tcp.SYN && tcp.ACK {
			buffer.WriteString("K")
		} else if tcp.SYN {
			buffer.WriteString("S")
		} else if tcp.RST && tcp.ACK {
			buffer.WriteString("T")
		} else if tcp.RST {
			buffer.WriteString("R")
		} else if tcp.PSH && tcp.ACK {
			buffer.WriteString("H")
		} else if tcp.PSH {
			buffer.WriteString("P")
		} else if tcp.ACK {
			buffer.WriteString("A")
		} else if tcp.URG {
			buffer.WriteString("U")
		} else {
			buffer.WriteString("O")
		}
	}

	buffer.WriteString(strings.Repeat("-", newTime))

	f.SetValue(buffer.String(), context, f)
}

func init() {
	flows.RegisterTemporaryFeature("_characters2", "returns a textual representation of a packet", ipfix.StringType, 0, flows.PacketFeature, func() flows.Feature { return &_characters2{} }, flows.RawPacket)
}

////////////////////////////////////////////////////////////////////////////////

// outputs number of consecutive seconds in which there was at least one packet
// seconds are counted from the first packet
type _consecutiveSeconds struct {
	flows.BaseFeature
	count    uint64
	lastTime flows.DateTimeNanoseconds
}

func (f *_consecutiveSeconds) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.lastTime = 0
	f.count = 0
}

func (f *_consecutiveSeconds) Event(new interface{}, context *flows.EventContext, src interface{}) {
	time := context.When()
	if f.lastTime == 0 {
		f.lastTime = time
		f.count++
	} else {
		if time-f.lastTime > 1*flows.SecondsInNanoseconds { // if time difference to f.lastTime is more than one second
			f.lastTime = time
			if time-f.lastTime < 2*flows.SecondsInNanoseconds { // if there is less than 2s between this and last packet, count
				f.count++
			} else { // otherwise, there was a break in seconds between packets
				f.SetValue(f.count, context, f)
				f.count = 1
			}
		}
	}
}

func (f *_consecutiveSeconds) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	f.SetValue(f.count, context, f)
}

func init() {
	flows.RegisterTemporaryFeature("__consecutiveSeconds", "outputs number of consecutive seconds in which there was at least one packet", ipfix.Unsigned64Type, 0, flows.PacketFeature, func() flows.Feature { return &_consecutiveSeconds{} }, flows.RawPacket)
	flows.RegisterTemporaryCompositeFeature("__maximumConsecutiveSeconds", "", ipfix.Unsigned64Type, 0, "maximum", "_consecutiveSeconds")
	flows.RegisterTemporaryCompositeFeature("__minimumConsecutiveSeconds", "", ipfix.Unsigned64Type, 0, "minimum", "_consecutiveSeconds")
}

////////////////////////////////////////////////////////////////////////////////

// outputs number of seconds in which there was at least one packet
// seconds are counted from the first packet
type _activeForSeconds struct {
	flows.BaseFeature
	count    uint64
	lastTime flows.DateTimeNanoseconds
}

func (f *_activeForSeconds) Start(context *flows.EventContext) {
	f.lastTime = 0
	f.count = 0
}

func (f *_activeForSeconds) Event(new interface{}, context *flows.EventContext, src interface{}) {
	time := context.When()
	if f.lastTime == 0 {
		f.lastTime = time
		f.count++
	} else {
		if time-f.lastTime > 1*flows.SecondsInNanoseconds { // if time difference to f.lastTime is more than one second
			f.lastTime = time
			f.count++
		}
	}
}

func (f *_activeForSeconds) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	f.SetValue(f.count, context, f)
}

func init() {
	flows.RegisterTemporaryFeature("_activeForSeconds", "outputs number of seconds in which there was at least one packet", ipfix.Unsigned64Type, 0, flows.FlowFeature, func() flows.Feature { return &_activeForSeconds{} }, flows.RawPacket)
}

////////////////////////////////////////////////////////////////////////////////

// outputs tcp options of the 1st packet
type _tcpOptionsFirstPacket struct {
	flows.BaseFeature
	done bool
}

func (f *_tcpOptionsFirstPacket) Start(context *flows.EventContext) {
	f.done = false
}

func (f *_tcpOptionsFirstPacket) Event(new interface{}, context *flows.EventContext, src interface{}) {
	var buffer bytes.Buffer
	if !f.done {
		tcp := features.GetTCP(new)
		if tcp != nil {
			opts := tcp.Options
			for _, o := range opts {
				buffer.WriteString(fmt.Sprintf("[%d %d %x]", o.OptionType, o.OptionLength, o.OptionData))
			}
			f.SetValue(buffer.String(), context, f)
		}
	}
	f.done = true
}

func init() {
	flows.RegisterTemporaryFeature("_tcpOptionsFirstPacket", "textual representation of TCP options in first packet", ipfix.StringType, 0, flows.FlowFeature, func() flows.Feature { return &_tcpOptionsFirstPacket{} }, flows.RawPacket)
}

////////////////////////////////////////////////////////////////////////////////

// outputs first 2 bytes of tcp timestamp of the 1st packet
type _tcpTimestampFirstPacket struct {
	flows.BaseFeature
	done bool
}

func (f *_tcpTimestampFirstPacket) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.done = false
}

func (f *_tcpTimestampFirstPacket) Event(new interface{}, context *flows.EventContext, src interface{}) {
	if !f.done {
		tcp := features.GetTCP(new)
		if tcp != nil {
			opts := tcp.Options
			for _, o := range opts {
				if o.OptionType.String() == "Timestamps" {
					ts := binary.BigEndian.Uint32(o.OptionData[0:4])
					f.SetValue(ts, context, f)
				}
			}
		}
	}
	f.done = true
}

func init() {
	flows.RegisterTemporaryFeature("_tcpTimestampFirstPacket", "TCP timestamp of first packet", ipfix.Unsigned32Type, 0, flows.FlowFeature, func() flows.Feature { return &_tcpTimestampFirstPacket{} }, flows.RawPacket)
}

////////////////////////////////////////////////////////////////////////////////

// outputs option data of tcp options before timestamp of the 1st packet
type _tcpOptionDataFirstPacket struct {
	flows.BaseFeature
	done bool
}

func (f *_tcpOptionDataFirstPacket) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.done = false
}

func (f *_tcpOptionDataFirstPacket) Event(new interface{}, context *flows.EventContext, src interface{}) {
	var buffer bytes.Buffer
	if !f.done {
		tcp := features.GetTCP(new)
		if tcp != nil {
			opts := tcp.Options
			for _, o := range opts {
				if o.OptionType.String() == "Timestamps" {
					break
				}
				buffer.WriteString(fmt.Sprintf("[%x]", o.OptionData))
			}
			f.SetValue(buffer.String(), context, f)
		}
	}
	f.done = true
}

func init() {
	flows.RegisterTemporaryFeature("_tcpOptionDataFirstPacket", "textual representation of TCP options in first packet", ipfix.StringType, 0, flows.FlowFeature, func() flows.Feature { return &_tcpOptionDataFirstPacket{} }, flows.RawPacket)
}

////////////////////////////////////////////////////////////////////////////////

// outputs list of the difference of tcp timestamp divided by actual time in the packets in the flow
type _tcpTimestampsPerSeconds struct {
	flows.BaseFeature
	timestamps []uint32
	times      []uint32
}

func (f *_tcpTimestampsPerSeconds) Event(new interface{}, context *flows.EventContext, src interface{}) {
	tcp := features.GetTCP(new)
	if tcp != nil {
		opts := tcp.Options
		for _, o := range opts {
			if o.OptionType.String() == "Timestamps" {
				timestamp := binary.BigEndian.Uint32(o.OptionData[0:4])
				f.timestamps = append(f.timestamps, timestamp)
				time := context.When()
				newTime := uint32(time / flows.SecondsInNanoseconds)
				f.times = append(f.times, newTime)
			}
		}
	}
}

func (f *_tcpTimestampsPerSeconds) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	var buffer []float64
	if len(f.timestamps) > 1 {
		for i := range f.timestamps[0 : len(f.timestamps)-1] {
			tcpStampDif := f.timestamps[i+1] - f.timestamps[i]
			stampDif := f.times[i+1] - f.times[i]

			if stampDif > 0 {
				res := float64(tcpStampDif) / float64(stampDif)
				buffer = append(buffer, res)
			}
		}

		if len(buffer) > 0 {
			for _, val := range buffer {
				f.SetValue(val, context, f)
			}
		}
	}
}

func init() {
	flows.RegisterTemporaryFeature("_tcpTimestampsPerSeconds", "list of the difference of tcp timestamp divided by actual time in the packets in the flow", ipfix.Float64Type, 0, flows.PacketFeature, func() flows.Feature { return &_tcpTimestampsPerSeconds{} }, flows.RawPacket)
}

////////////////////////////////////////////////////////////////////////////////

type _label struct {
	flows.BaseFeature
}

func (f *_label) Event(new interface{}, context *flows.EventContext, src interface{}) {
	label := new.(packet.Buffer).Label()
	if label != nil {
		f.SetValue(label, context, f)
	}
}

func init() {
	flows.RegisterTemporaryFeature("__label", "label of the packet", ipfix.OctetArrayType, 0, flows.PacketFeature, func() flows.Feature { return &_label{} }, flows.RawPacket)
}

////////////////////////////////////////////////////////////////////////////////

type _flowKey struct {
	flows.BaseFeature
}

func (f *_flowKey) Event(new interface{}, context *flows.EventContext, src interface{}) {
	flowkey := context.Flow().Key()
	f.SetValue(flowkey, context, f)
}

func init() {
	flows.RegisterTemporaryFeature("__flowKey", "string go-flows uses as a key", ipfix.StringType, 0, flows.PacketFeature, func() flows.Feature { return &_flowKey{} }, flows.RawPacket)
}
