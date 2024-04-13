package datatype

var binaryTypeList = [...]string{"BINARY", "VARBINARY", "TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB"}

type Binary struct {
}

func (b *Binary) convert(v interface{}) (interface{}, error) {
	// 将v转换为[]byte
	return v, nil
}

func init() {
	b := &Binary{}
	for _, t := range binaryTypeList {
		Register(t, b)
	}
}
