package erroz

func WithData(data any) Option {
    return func(j *businessError) {
        j.Data = data
    }
}

func WithStatus(status int) Option {
    return func(j *businessError) {
        j.status = status
    }
}

func WithDebug(data any) Option {
    return func(j *businessError) {
        j.Debug = append(j.Debug, data)
    }
}
