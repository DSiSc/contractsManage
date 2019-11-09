package contracts

import (
	"errors"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG/common"
	"github.com/DSiSc/evm-NG/common/hexutil"
	"github.com/DSiSc/evm-NG/common/math"
	"github.com/DSiSc/evm-NG/constant"
	"github.com/DSiSc/statedb-NG/util"
	"math/big"
	"reflect"
	"sync"
)

var (
	ContractsAddr sync.Map
)

var (
	UnSupportedTypeError  = errors.New("unsupported arg type")
	InvalidUnmarshalError = errors.New("invalid unmarshal error")
)

// AbiEncode encode contract call param
func AbiEncode(args ...interface{}) ([]byte, error) {
	retPre := make([]byte, 0)
	retData := make([]byte, 0)
	preOffsetPadding := len(args) * constant.EvmWordSize
	for _, arg := range args {
		valType := reflect.TypeOf(arg)
		switch valType.Kind() {
		case reflect.String:
			offset := preOffsetPadding + len(retData)
			retPre = append(retPre, math.PaddedBigBytes(big.NewInt(int64(offset)), constant.EvmWordSize)...)
			retData = append(retData, encodeString(arg.(string))...)
		case reflect.Slice:
			if reflect.Array == valType.Elem().Kind() && valType.AssignableTo(reflect.TypeOf([]types.Address{})) {
				offset := preOffsetPadding + len(retData)
				retPre = append(retPre, math.PaddedBigBytes(big.NewInt(int64(offset)), constant.EvmWordSize)...)
				addrs := arg.([]types.Address)
				retData = append(retData, math.PaddedBigBytes(big.NewInt(0).SetUint64(uint64(len(addrs))), constant.EvmWordSize)...)
				for _, addr := range addrs {
					retData = append(retData, common.LeftPadBytes(addr[:], constant.EvmWordSize)...)
				}
				continue
			}

			if reflect.Uint8 != valType.Elem().Kind() {
				return nil, UnSupportedTypeError
			}
			offset := preOffsetPadding + len(retData)
			retPre = append(retPre, math.PaddedBigBytes(big.NewInt(int64(offset)), constant.EvmWordSize)...)
			retData = append(retData, encodeBytes(arg.([]byte))...)
		case reflect.Uint64:
			retPre = append(retPre, math.PaddedBigBytes(big.NewInt(0).SetUint64(arg.(uint64)), constant.EvmWordSize)...)
		case reflect.Array:
			if valType.AssignableTo(reflect.TypeOf(types.Address{})) {
				addr := arg.(types.Address)
				retPre = append(retPre, common.LeftPadBytes(addr[:], constant.EvmWordSize)...)
			}
		default:
			return nil, errors.New("unsupported return type")
		}
	}
	return append(retPre, retData...), nil
}

// AbiDecode decode contract return values
func AbiDecode(input []byte, vals ...interface{}) error {
	for i := 0; i < len(vals); i++ {
		rv := reflect.ValueOf(vals[i])
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return InvalidUnmarshalError
		}

		switch rv.Elem().Type().Kind() {
		case reflect.String:
			arg := string(extractDynamicTypeData(input, i))
			rv.Elem().SetString(arg)
		case reflect.Slice:
			if reflect.Array == rv.Elem().Type().Elem().Kind() && reflect.TypeOf(vals[i]).AssignableTo(reflect.TypeOf(&[]types.Address{})) {
				addrNum, addrBytes := extractAddressTypeArray(input, i)
				if addrNum <= 0 {
					continue
				}
				for i := uint64(0); i < addrNum; i++ {
					addr := util.BytesToAddress(addrBytes[i*constant.EvmWordSize : (i+1)*constant.EvmWordSize])
					rv.Elem().Set(reflect.Append(rv.Elem(), reflect.ValueOf(addr)))
				}
				continue
			}
			if reflect.Uint8 != rv.Elem().Type().Elem().Kind() {
				return UnSupportedTypeError
			}
			arg := extractDynamicTypeData(input, i)
			rv.Elem().SetBytes(arg)
		case reflect.Uint64:
			arg, _ := math.ParseUint64(hexutil.Encode(input[i*constant.EvmWordSize : (i+1)*constant.EvmWordSize]))
			rv.Elem().SetUint(arg)
		case reflect.Array:
			arg := arrayByte20(input, i)
			rv.Elem().SetBytes(arg)
		case reflect.Bool:
			arg, _ := math.ParseUint64(hexutil.Encode(input[i*constant.EvmWordSize : (i+1)*constant.EvmWordSize]))
			rv.Elem().SetBool(arg != 0)
		default:
			return UnSupportedTypeError
		}
	}
	return nil
}

// extract dynamic type data
func extractAddressTypeArray(data []byte, index int) (uint64, []byte) {
	offset, _ := math.ParseUint64(hexutil.Encode(data[index*constant.EvmWordSize : (index+1)*constant.EvmWordSize]))
	if offset >= uint64(len(data)) {
		return 0, nil
	}
	dataLen, _ := math.ParseUint64(hexutil.Encode(data[offset : offset+constant.EvmWordSize]))
	argStart := offset + constant.EvmWordSize
	argEnd := argStart + dataLen*constant.EvmWordSize
	return dataLen, data[argStart:argEnd]
}

// extract dynamic type data
func extractDynamicTypeData(data []byte, index int) []byte {
	offset, _ := math.ParseUint64(hexutil.Encode(data[index*constant.EvmWordSize : (index+1)*constant.EvmWordSize]))
	if offset >= uint64(len(data)) {
		// address type
		addr := util.BytesToAddress(arrayByte20(data, index))
		addrStr := util.AddressToHex(addr)
		return []byte(addrStr)
	}
	dataLen, _ := math.ParseUint64(hexutil.Encode(data[offset : offset+constant.EvmWordSize]))
	argStart := offset + constant.EvmWordSize
	argEnd := argStart + dataLen
	return data[argStart:argEnd]
}

// encode the string to the format needed by evm
func encodeString(val string) []byte {
	return encodeBytes([]byte(val))
}

// encode the byte array to the format needed by evm
func encodeBytes(val []byte) []byte {
	ret := make([]byte, 0)
	ret = append(ret, math.PaddedBigBytes(big.NewInt(int64(len(val))), constant.EvmWordSize)...)
	for i := 0; i < len(val); {
		if (len(val) - i) > constant.EvmWordSize {
			ret = append(ret, val[i:i+constant.EvmWordSize]...)
			i += constant.EvmWordSize
		} else {
			ret = append(ret, common.RightPadBytes(val[i:], constant.EvmWordSize)...)
			i += len(val)
		}
	}
	return ret
}

func arrayByte20(totalInput []byte, varIndex int) []byte {
	offset := varIndex*constant.EvmWordSize + constant.AddressOffset
	addrByte := totalInput[offset : (varIndex+1)*constant.EvmWordSize]
	return []byte(addrByte)
}
