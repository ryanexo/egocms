package domain

var (
    ValueString = valueType{typ: 0}
    ValueNumber = valueType{typ: 1}
    ValueBool   = valueType{typ: 2}
    ValueTime   = valueType{typ: 3}
)

type valueType struct {
    typ int16
}

func (v valueType) Is(typ int16) bool {
    return v.typ == typ
}
