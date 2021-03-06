package operations

import (
	"errors"
	"math"
	"net"

	"github.com/chtisgit/go-flows/flows"
	"github.com/chtisgit/go-flows/modules/features"
	ipfix "github.com/CN-TU/go-ipfix"
)

type get struct {
	flows.BaseFeature
	index, current int64
}

func (f *get) SetArguments(arguments []int, features []flows.Feature) {
	f.index = flows.ToInt(features[arguments[0]].Value())
}
func (f *get) Start(ec *flows.EventContext) {
	f.BaseFeature.Start(ec)
	f.current = 0
}

func (f *get) Event(new interface{}, context *flows.EventContext, src interface{}) {
	if f.current == f.index {
		f.SetValue(new, context, f)
	}
	f.current++
}

func resolveGet(args []ipfix.InformationElement) (ipfix.InformationElement, error) {
	if len(args) != 2 {
		return ipfix.InformationElement{}, errors.New("get must have exactly 2 arguments")
	}
	return args[1], nil
}

func init() {
	flows.RegisterCustomFunction("get", "gets the <value>-th element of the second argument; indexing is like in Python", resolveGet, flows.FlowFeature, func() flows.Feature { return &get{} }, flows.Const, flows.PacketFeature)
}

////////////////////////////////////////////////////////////////////////////////

type count struct {
	flows.BaseFeature
	count uint64
}

func (f *count) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.count = 0
}

func (f *count) Event(new interface{}, context *flows.EventContext, src interface{}) {
	f.count++
}

func (f *count) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	f.SetValue(f.count, context, f)
}

func init() {
	flows.RegisterTypedFunction("count", "returns number of selected objects", ipfix.Unsigned64Type, 0, flows.FlowFeature, func() flows.Feature { return &count{} }, flows.Selection)
	flows.RegisterTypedFunction("count", "returns number of selected objects", ipfix.Unsigned64Type, 0, flows.FlowFeature, func() flows.Feature { return &count{} }, flows.PacketFeature)
}

////////////////////////////////////////////////////////////////////////////////

type mean struct {
	flows.BaseFeature
	total float64
	count uint64
}

func (f *mean) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.total = 0
	f.count = 0
}

func (f *mean) Event(new interface{}, context *flows.EventContext, src interface{}) {
	f.total += flows.ToFloat(new)
	f.count++
}

func (f *mean) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	if f.count > 0 {
		f.SetValue(f.total/float64(f.count), context, f)
	}
}

func init() {
	flows.RegisterTypedFunction("mean", "returns mean of input", ipfix.Float64Type, 0, flows.FlowFeature, func() flows.Feature { return &mean{} }, flows.PacketFeature)
}

////////////////////////////////////////////////////////////////////////////////

type min struct {
	flows.BaseFeature
	current interface{}
}

func (f *min) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.current = nil
}

func (f *min) Event(new interface{}, context *flows.EventContext, src interface{}) {
	if f.current == nil {
		f.current = new
	} else {
		_, fl, a, b := flows.UpConvert(new, f.current)
		switch fl {
		case flows.UIntType:
			if a.(uint64) < b.(uint64) {
				f.current = new
			}
		case flows.IntType:
			if a.(int64) < b.(int64) {
				f.current = new
			}
		case flows.FloatType:
			if a.(float64) < b.(float64) {
				f.current = new
			}
		}
	}
}

func (f *min) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	f.SetValue(f.current, context, f)
}

func init() {
	flows.RegisterFunction("min", "returns min of input", flows.FlowFeature, func() flows.Feature { return &min{} }, flows.PacketFeature)
	flows.RegisterFunction("minimum", "returns min of input", flows.FlowFeature, func() flows.Feature { return &min{} }, flows.PacketFeature)
}

////////////////////////////////////////////////////////////////////////////////

type max struct {
	flows.BaseFeature
	current interface{}
}

func (f *max) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.current = nil
}

func (f *max) Event(new interface{}, context *flows.EventContext, src interface{}) {
	if f.current == nil {
		f.current = new
	} else {
		_, fl, a, b := flows.UpConvert(new, f.current)
		switch fl {
		case flows.UIntType:
			if a.(uint64) > b.(uint64) {
				f.current = new
			}
		case flows.IntType:
			if a.(int64) > b.(int64) {
				f.current = new
			}
		case flows.FloatType:
			if a.(float64) > b.(float64) {
				f.current = new
			}
		}
	}
}

func (f *max) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	f.SetValue(f.current, context, f)
}

func init() {
	flows.RegisterFunction("max", "returns max of input", flows.FlowFeature, func() flows.Feature { return &max{} }, flows.PacketFeature)
	flows.RegisterFunction("maximum", "returns max of input", flows.FlowFeature, func() flows.Feature { return &max{} }, flows.PacketFeature)
}

////////////////////////////////////////////////////////////////////////////////

// Calculate online stdev according to Welford "Note on a method for calculating corrected sums of squares and products."
type stdev struct {
	flows.BaseFeature
	count    uint64
	mean, m2 float64
}

func (f *stdev) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.count = 0
	f.mean = 0
	f.m2 = 0
}

func (f *stdev) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	if f.count != 0 {
		f.SetValue(math.Sqrt(f.m2/(float64(f.count)-1)), context, f)
	}
}

func (f *stdev) Event(new interface{}, context *flows.EventContext, src interface{}) {
	val := flows.ToFloat(new)
	f.count++
	delta := val - f.mean
	f.mean = f.mean + delta/float64(f.count)
	delta2 := val - f.mean
	f.m2 = f.m2 + delta*delta2
}

func init() {
	flows.RegisterTypedFunction("stdev", "", ipfix.Float64Type, 0, flows.FlowFeature, func() flows.Feature { return &stdev{} }, flows.PacketFeature)
}

////////////////////////////////////////////////////////////////////////////////

// Calculate online variance according to Welford "Note on a method for calculating corrected sums of squares and products."
type variance struct {
	flows.BaseFeature
	count    uint64
	mean, m2 float64
}

func (f *variance) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.count = 0
	f.mean = 0
	f.m2 = 0
}

func (f *variance) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	if f.count > 1 {
		f.SetValue(f.m2/(float64(f.count)-1), context, f)
	}
}

func (f *variance) Event(new interface{}, context *flows.EventContext, src interface{}) {
	val := flows.ToFloat(new)
	f.count++
	delta := val - f.mean
	f.mean = f.mean + delta/float64(f.count)
	delta2 := val - f.mean
	f.m2 = f.m2 + delta*delta2
}

func init() {
	flows.RegisterTypedFunction("variance", "", ipfix.Float64Type, 0, flows.FlowFeature, func() flows.Feature { return &variance{} }, flows.PacketFeature)
}

////////////////////////////////////////////////////////////////////////////////

type median struct {
	flows.BaseFeature
	vector features.TypedSlice
}

func (f *median) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.vector = nil
}

func (f *median) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	if f.vector == nil {
		return // No median
	}
	k := f.vector.Len()
	// Start with trivial cases
	switch k {
	case 0:
		return // No median
	case 1:
		f.SetValue(f.vector.Get(0), context, f) // take the only value
		return
	case 2:
		if f.vector.IsNumeric() {
			// for numeric value between item 0 and 1
			f.SetValue((f.vector.GetFloat(0)+f.vector.GetFloat(1))/2, context, f)
			return
		}
		// for non-numeric take the lower value
		if f.vector.Less(0, 1) {
			f.SetValue(f.vector.Get(0), context, f)
			return
		}
		f.SetValue(f.vector.Get(1), context, f)
		return
	}
	if k%2 == 1 {
		median := k / 2 // middle element
		f.vector.Select(0, k-1, median)
		f.SetValue(f.vector.Get(median), context, f)
		return
	}
	if !f.vector.IsNumeric() {
		// non-numeric elements -> take the smaller one
		median := k/2 - 1
		f.vector.Select(0, k-1, median)
		f.SetValue(f.vector.Get(median), context, f)
	}
	median1 := k / 2
	f.vector.Select(0, k-1, median1)
	median2 := f.vector.Max(0, median1-1)
	f.SetValue((f.vector.GetFloat(median1)+f.vector.GetFloat(median2))/2, context, f)
}

func (f *median) Event(new interface{}, context *flows.EventContext, src interface{}) {
	if f.vector == nil {
		f.vector = features.NewTypedSlice(new)
	} else {
		f.vector.Append(new)
	}
}

func init() {
	flows.RegisterFunction("median", "median; numeric even: arithmetic mean of two middle values; non-numeric even: lower of the two middle values", flows.FlowFeature, func() flows.Feature { return &median{} }, flows.PacketFeature)
}

////////////////////////////////////////////////////////////////////////////////

type mode struct {
	flows.BaseFeature
	vector map[interface{}]uint64
}

func (f *mode) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.vector = make(map[interface{}]uint64)
}

func (f *mode) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	var max uint64
	var m interface{}
	for val, num := range f.vector {
		if num > max {
			max = num
			m = val
		} else if num == max && features.Less(val, m) {
			m = val
		}
	}
	if max > 0 {
		f.SetValue(m, context, f)
	}
}

func (f *mode) Event(new interface{}, context *flows.EventContext, src interface{}) {
	switch val := new.(type) {
	case []byte:
		f.vector[string(val)]++
	case net.IP:
		f.vector[val.String()]++
	default:
		f.vector[val]++
	}
}

func init() {
	flows.RegisterFunction("mode", "mode of value; if multimodal then smallest value; no special handling for continous", flows.FlowFeature, func() flows.Feature { return &mode{} }, flows.PacketFeature)
}

////////////////////////////////////////////////////////////////////////////////

type modeCount struct {
	flows.BaseFeature
	vector map[interface{}]uint64
}

func (f *modeCount) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.vector = make(map[interface{}]uint64)
}

func (f *modeCount) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	var max uint64
	var m interface{}
	for val, num := range f.vector {
		if num > max {
			max = num
			m = val
		} else if num == max && features.Less(val, m) {
			m = val
		}
	}
	if max > 0 {
		f.SetValue(max, context, f)
	} else {
		f.SetValue(math.NaN(), context, f)
	}
}

func (f *modeCount) Event(new interface{}, context *flows.EventContext, src interface{}) {
	switch val := new.(type) {
	case []byte:
		f.vector[string(val)]++
	case net.IP:
		f.vector[val.String()]++
	default:
		f.vector[val]++
	}
}

func init() {
	flows.RegisterFunction("modeCount", "NUmber of packets for the mode of value; if multimodal then smallest value; no special handling for continous", flows.FlowFeature, func() flows.Feature { return &modeCount{} }, flows.PacketFeature)
}

////////////////////////////////////////////////////////////////////////////////

type distinct struct {
	flows.BaseFeature
	vector map[interface{}]uint64
}

func (f *distinct) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.vector = make(map[interface{}]uint64)
}

func (f *distinct) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	f.SetValue(len(f.vector), context, f)
}

func (f *distinct) Event(new interface{}, context *flows.EventContext, src interface{}) {
	switch val := new.(type) {
	case []byte:
		f.vector[string(val)]++
	case net.IP:
		f.vector[string(val)]++
	default:
		f.vector[val]++
	}
}

func init() {
	flows.RegisterFunction("distinct", "number of distinct elements in a list", flows.FlowFeature, func() flows.Feature { return &distinct{} }, flows.PacketFeature)
}

////////////////////////////////////////////////////////////////////////////////

type set struct {
	flows.BaseFeature
	vector map[interface{}]bool
	set []interface{}
}

func (f *set) Start(context *flows.EventContext) {
	f.BaseFeature.Start(context)
	f.vector = make(map[interface{}]bool)
	f.set = make([]interface{}, 0)
}

func (f *set) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	f.SetValue(f.set, context, f)
}

func (f *set) Event(new interface{}, context *flows.EventContext, src interface{}) {
	switch val := new.(type) {
	case []byte:
		if !f.vector[string(val)] {
			f.vector[string(val)] = true
			f.set = append(f.set, string(val))
		}
	case net.IP:
		if !f.vector[string(val)] {
			f.vector[string(val)] = true
			f.set = append(f.set, val.String())
		}
	default:
		if !f.vector[val] {
			f.vector[val] = true
			f.set = append(f.set, val)
		}
	}
}

func init() {
	flows.RegisterFunction("set", "distinct elements in a list", flows.FlowFeature, func() flows.Feature { return &set{} }, flows.PacketFeature)
}

////////////////////////////////////////////////////////////////////////////////
