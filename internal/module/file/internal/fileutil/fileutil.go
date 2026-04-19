package fileutil

import (
    `time`
    
    `github.com/google/uuid`
)

func GeneratePath() (*Path, error) {
    date := time.Now().Format("2006/0102")
    id, err := uuid.NewV7()
    if err != nil {
        return nil, err
    }
    
    return &Path{
        Filename: id.String(),
        Path:     "./" + date + "/",
    }, nil
}
