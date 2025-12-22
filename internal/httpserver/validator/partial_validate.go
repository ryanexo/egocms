package validator

type PartialFields struct {
    fields []string
}

var _ PartialValidation = new(PartialFields)

func (p *PartialFields) ValidationFields() []string {
    return p.fields
}

func (p *PartialFields) SetValidationFields(fields []string) {
    p.fields = fields
}
