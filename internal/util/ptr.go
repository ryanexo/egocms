package util

func PtrOrDefault[T any](p *T, def T) T {
    if p == nil {
        return def
    }
    return *p
}

func PtrOrZero[T any](p *T) T {
    var zero T
    return PtrOrDefault(p, zero)
}

func ToPtr[T any](v T) *T {
    return &v
}
