package validate

type partialValidate interface {
	ValidationFields() []string
}

type PartialValidate struct {
	fields []string
}

func (s *PartialValidate) SetValidationFields(fields ...string) {
	s.fields = fields
}

func (s *PartialValidate) ValidationFields() []string {
	return s.fields
}
