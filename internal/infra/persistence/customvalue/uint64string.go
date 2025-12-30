package customvalue

import (
    `database/sql/driver`
    `errors`
    `fmt`
    `strconv`
)

type Uint64String string

func (u *Uint64String) Scan(value any) error {
    if value == nil {
        *u = ""
        return nil
    }
    
    switch v := value.(type) {
    case int64:
        if v < 0 {
            return fmt.Errorf("negative value for Uint64String")
        }
        *u = Uint64String(strconv.FormatUint(uint64(v), 10))
    case uint64:
        *u = Uint64String(strconv.FormatUint(v, 10))
    case []byte:
        s := string(v)
        if _, err := strconv.ParseUint(s, 10, 64); err != nil {
            return err
        }
        *u = Uint64String(s)
    default:
        *u = ""
        return fmt.Errorf("unsupported Scan type %T", v)
    }
    
    return nil
}

func (u Uint64String) Value() (driver.Value, error) {
    if u == "" {
        return nil, nil
    }
    
    n, err := strconv.ParseUint(string(u), 10, 64)
    if err != nil {
        return nil, err
    }
    
    return n, nil
}

func (u Uint64String) Uint64() (uint64, error) {
    if u == "" {
        return 0, errors.New("invalid Uint64String")
    }
    return strconv.ParseUint(string(u), 10, 64)
}

func (u Uint64String) String() string {
    return string(u)
}

func NewUint64String(n uint64) Uint64String {
    return Uint64String(strconv.FormatUint(n, 10))
}
