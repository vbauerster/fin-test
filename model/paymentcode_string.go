// Code generated by "stringer -type=PaymentCode -trimprefix=Payment"; DO NOT EDIT.

package model

import "strconv"

const _PaymentCode_name = "DEPOSITWITHDRAWTRANSFER"

var _PaymentCode_index = [...]uint8{0, 7, 15, 23}

func (i PaymentCode) String() string {
	i -= 1
	if i < 0 || i >= PaymentCode(len(_PaymentCode_index)-1) {
		return "PaymentCode(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _PaymentCode_name[_PaymentCode_index[i]:_PaymentCode_index[i+1]]
}
