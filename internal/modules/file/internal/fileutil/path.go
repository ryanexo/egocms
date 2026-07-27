package fileutil

import `path`

type Path struct {
    Path     string
    Filename string
}

func (s Path) FullPath() string {
    return path.Join(s.Path, s.Filename)
}

func (s Path) FullFilename(ext string) string {
    return s.Filename + ext
}
