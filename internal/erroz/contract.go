package erroz

type Option func(*businessError)

type BusinessError interface {
    WithOption(...Option) BusinessError
    Format(...any) BusinessError
    Wrap(error) BusinessError
    Write(Responsible)
    WriteWithAbort(Responsible)
    ToError() error
    prototype() businessError
}

type Responsible interface {
    JSON(status int, body any)
    Abort()
}
