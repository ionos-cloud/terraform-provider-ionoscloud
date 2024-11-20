package convptr

// Int64ToInt32 converts an int64 to an int32 pointer.
func Int64ToInt32(i *int64) *int32 {
	if i == nil {
		return nil
	}

	v := int32(*i)
	return &v
}

// Int32ToInt64 converts an int32 to an int64 pointer.
func Int32ToInt64(i *int32) *int64 {
	if i == nil {
		return nil
	}

	v := int64(*i)
	return &v
}
