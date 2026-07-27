package fileutil

import `slices`

func IsImage(ext string) bool {
    imageExt := []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp"}
    return slices.Contains(imageExt, ext)
}
