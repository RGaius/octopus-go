package datatype

type Decimal struct {
}

func (d *Decimal) convert(value interface{}) (interface{}, error) {
	return value, nil
}

func init() {
	Register("DECIMAL", &Decimal{})
}
