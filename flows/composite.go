package flows

import (
	"fmt"

	ipfix "pm.cn.tuwien.ac.at/ipfix/go-ipfix"
)

type compositeFeatureMaker struct {
	definition []interface{}
	ie         ipfix.InformationElement
}

func compositeToCall(features []interface{}) (ret []string) {
	flen := len(features) - 1
	for i, feature := range features {
		if list, ok := feature.([]interface{}); ok {
			ret = append(ret, compositeToCall(list)...)
		} else {
			ret = append(ret, fmt.Sprint(feature))
		}
		if i == 0 {
			ret = append(ret, "(")
		} else if i < flen {
			ret = append(ret, ",")
		} else {
			ret = append(ret, ")")
		}
	}
	return
}