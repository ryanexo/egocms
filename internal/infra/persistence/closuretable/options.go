package closuretable

import `fmt`

type Option func(*options)

type options struct {
    TableName  string
    ModuleName string
}

func (s options) tableName() string {
    if s.TableName != "" {
        return s.TableName
    }
    return fmt.Sprintf("%s_tree_closure", s.ModuleName)
}

func WithTableName(name string) Option {
    return func(o *options) {
        o.TableName = name
    }
}

func WithModuleName(name string) Option {
    return func(o *options) {
        o.ModuleName = name
    }
}
