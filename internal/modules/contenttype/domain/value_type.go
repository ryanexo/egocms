package domain

var (
    ValueString = valueType{0}
    ValueNumber = valueType{1}
    ValueBool   = valueType{2}
    ValueTime   = valueType{3}
)

type valueType struct {
    v int16
}

func (t valueType) Is(typ int16) bool {
    return t.v == typ
}
