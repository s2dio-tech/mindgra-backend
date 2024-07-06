package common

type Nullable struct {
	Value any
}

func (n Nullable) ToStringPtr() *string {
	if (n.Value) == nil {
		return nil
	}
	str := n.Value.(string)
	return &str
}

func (n Nullable) ToStringArrayPtr() *[]string {
	if (n.Value) == nil {
		return nil
	}
	ifs := n.Value.([]interface{})
	strs := []string{}
	for _, i := range ifs {
		strs = append(strs, i.(string))
	}
	return &strs
}

func (n Nullable) ToInt64Ptr() *int64 {
	if (n.Value) == nil {
		return nil
	}
	strs := n.Value.(int64)
	return &strs
}

func ToPointer[T any](x T) *T {
	return &x
}
