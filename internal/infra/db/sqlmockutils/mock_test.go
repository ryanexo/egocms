package sqlmockutils

import (
    "testing"
    
    `cms/internal/infra/db`
    
    "github.com/stretchr/testify/assert"
)

type T struct {
    A string
    B int
    C bool
}

func Test_Row(t *testing.T) {
    db, sqlmock := db.NewDBMock()
    rows := NewRows(T{}).Add(
        []T{
            {"test1", 1, true},
            {"test2", 2, true},
        },
        &T{"test3", 3, false},
        (*int)(nil),
    ).Result()
    ptrTest := NewRows(&T{})
    assert.Equal(t, 3, len(ptrTest.fields))
    
    var result []T
    sqlmock.ExpectQuery("^SELECT").WillReturnRows(rows)
    err := db.Table("test").Select([]string{"A", "B", "C"}).Find(&result).Error
    assert.Nil(t, err)
    assert.Equal(t, 3, len(result))
    assert.Equal(t, "test1", result[0].A)
    assert.Equal(t, "test2", result[1].A)
    assert.Equal(t, "test3", result[2].A)
}
