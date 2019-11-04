package api

func String(v interface{}) *string {
	if v == nil {
		return nil
	}
	r := v.(string)
	return &r
}

func Uint64(v interface{}) *uint64 {
	if v == nil {
		return nil
	}
	r := v.(uint64)
	return &r
}

func Uint64FromInt(v interface{}) *uint64 {
	if v == nil {
		return nil
	}
	r := uint64(v.(int))
	return &r
}

func Bool(v interface{}) *bool {
	if v == nil {
		return nil
	}
	r := v.(bool)
	return &r
}
