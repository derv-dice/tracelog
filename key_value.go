package tracelog

type KeyValue struct {
	vtype int
	key   string
	value any
}

func (kv *KeyValue) IsValid() bool {
	return kv.vtype != vTypeNone
}

const (
	vTypeNone = iota
	vTypeString
	vTypeInt
	vTypeBool
)

func String(key string, val string) KeyValue {
	return KeyValue{
		vtype: vTypeString,
		key:   key,
		value: val,
	}
}

func Int(key string, val int) KeyValue {
	return KeyValue{
		vtype: vTypeInt,
		key:   key,
		value: val,
	}
}

func Bool(key string, val bool) KeyValue {
	return KeyValue{
		vtype: vTypeBool,
		key:   key,
		value: val,
	}
}
