package map_options

// Path: map-options/option.go
import "sort"

type DataType interface {
	int | int8 | int16 | int32 | int64 | string | float32 | float64
}

type Slice[T DataType] []T

type DataMap[KEY DataType, VALUE DataType] map[KEY]VALUE

type DataMapOption[KEYS DataType, VALUE DataType] struct {
	keys    Slice[KEYS]
	options DataMap[KEYS, VALUE]
}

func NewDataMapOption[KEYS, VALUE DataType](options DataMap[KEYS, VALUE]) *DataMapOption[KEYS, VALUE] {
	keys := make(Slice[KEYS], 0, len(options))
	for k := range options {
		keys = append(keys, k)
	}
	sortKeys(keys)
	return &DataMapOption[KEYS, VALUE]{
		keys:    keys,
		options: options,
	}
}

func (o *DataMapOption[KEYS, VALUE]) Keys() []KEYS {
	return o.keys
}

func (o *DataMapOption[KEYS, VALUE]) Option(key KEYS) VALUE {
	return o.options[key]
}

func (o *DataMapOption[KEYS, VALUE]) Options() DataMap[KEYS, VALUE] {
	return o.options
}

func sortKeys[T DataType](keys []T) {
	sort.Slice(keys, func(i, j int) bool {
		return keys[j] > keys[i]
	})
}

// GetMapKey type DataMap[KEYS DataType, VALUE DataType] map[KEYS]VALUE
// eg: GetMapKey[uint32, string](map[uint32]string{1: "1"}, 1)
func GetMapKey[KEYS, VALUE DataType](m DataMap[KEYS, VALUE], key KEYS) VALUE {
	var value VALUE
	if v, ok := m[key]; ok {
		return v
	}
	return value
}
