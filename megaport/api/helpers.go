package api

func String(v interface{}) *string {
	r := v.(string)
	return &r
}

func Uint64(v interface{}) *uint64 {
	r := v.(uint64)
	return &r
}

func Uint64FromInt(v interface{}) *uint64 {
	r := uint64(v.(int))
	return &r
}

func Bool(v interface{}) *bool {
	r := v.(bool)
	return &r
}
