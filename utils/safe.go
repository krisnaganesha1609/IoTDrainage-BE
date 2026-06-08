package utils

// Safe type-assertion helpers to prevent runtime panics when InfluxDB
// returns nil for fields that are absent in a particular record.

func SafeFloat64(v interface{}) float64 {
	if f, ok := v.(float64); ok {
		return f
	}
	return 0
}

func SafeBool(v interface{}) bool {
	if b, ok := v.(bool); ok {
		return b
	}
	return false
}

func SafeString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func SafeInt64(v interface{}) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case float64:
		return int64(val)
	case int:
		return int64(val)
	}
	return 0
}
