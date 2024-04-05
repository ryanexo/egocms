package erroz

func WithData(data any) Option {
    return func(j *json) {
        j.Data = data
    }
}

func WithStatus(status int) Option {
    return func(j *json) {
        j.status = status
    }
}
