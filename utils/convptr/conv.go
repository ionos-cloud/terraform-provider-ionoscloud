package convptr

func Int64ToInt32(i *int64) *int32 {
	if i == nil {
		return nil
	}

	v := int32(*i)
	return &v
}

func Int32ToInt64(i *int32) *int64 {
	if i == nil {
		return nil
	}

	v := int64(*i)
	return &v
}
