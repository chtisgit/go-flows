package operations

import (
	"github.com/chtisgit/go-flows/flows"
	ipfix "github.com/CN-TU/go-ipfix"
)

type andPacket struct {
	flows.MultiBasePacketFeature
}

func (f *andPacket) Event(new interface{}, context *flows.EventContext, src interface{}) {
	values := f.EventResult(new, src)
	if values == nil {
		return
	}

	val := values[0].(bool)
	for _, v := range values[1:] {
		if val == false {
			f.SetValue(val, context, f)
			return
		}
		val = val && v.(bool)
	}
	f.SetValue(val, context, f)
}

type andFlow struct {
	flows.MultiBaseFlowFeature
}

func (f *andFlow) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	values := f.GetValues(context)

	val := values[0].(bool)
	for _, v := range values[1:] {
		if val == false {
			f.SetValue(val, context, f)
			return
		}
		val = val && v.(bool)
	}
	f.SetValue(val, context, f)
}

func init() {
	flows.RegisterTypedFunction("and", "returns logical conjunction of arguments", ipfix.BooleanType, 0, flows.PacketFeature, func() flows.Feature { return &andPacket{} }, flows.PacketFeature, flows.PacketFeature, flows.Ellipsis)
	flows.RegisterTypedFunction("and", "returns logical conjunction of arguments", ipfix.BooleanType, 0, flows.FlowFeature, func() flows.Feature { return &andFlow{} }, flows.FlowFeature, flows.FlowFeature, flows.Ellipsis)
}

////////////////////////////////////////////////////////////////////////////////

type orPacket struct {
	flows.MultiBasePacketFeature
}

func (f *orPacket) Event(new interface{}, context *flows.EventContext, src interface{}) {
	values := f.EventResult(new, src)
	if values == nil {
		return
	}

	val := values[0].(bool)
	for _, v := range values[1:] {
		if val == true {
			f.SetValue(val, context, f)
			return
		}
		val = val || v.(bool)
	}
	f.SetValue(val, context, f)
}

type orFlow struct {
	flows.MultiBaseFlowFeature
}

func (f *orFlow) Stop(reason flows.FlowEndReason, context *flows.EventContext) {
	values := f.GetValues(context)

	val := values[0].(bool)
	for _, v := range values[1:] {
		if val == true {
			f.SetValue(val, context, f)
			return
		}
		val = val || v.(bool)
	}
	f.SetValue(val, context, f)
}

func init() {
	flows.RegisterTypedFunction("or", "returns logical disjunction of arguments", ipfix.BooleanType, 0, flows.PacketFeature, func() flows.Feature { return &orPacket{} }, flows.PacketFeature, flows.PacketFeature, flows.Ellipsis)
	flows.RegisterTypedFunction("or", "returns logical disjunction of arguments", ipfix.BooleanType, 0, flows.FlowFeature, func() flows.Feature { return &orFlow{} }, flows.FlowFeature, flows.FlowFeature, flows.Ellipsis)
}
