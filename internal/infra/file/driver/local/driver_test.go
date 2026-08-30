package local

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalStorageLifecycle(t *testing.T) {
	driver := localStorage{savePath: t.TempDir()}
	require.NoError(t, driver.Create(context.Background(), "2026/0826/report.txt", strings.NewReader("content"), 7))
	require.NoError(t, driver.Create(context.Background(), "2026/0826/report.txt", strings.NewReader("again"), 5))

	reader, err := driver.OpenReader(context.Background(), "2026/0826/report.txt")
	require.NoError(t, err)
	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.NoError(t, reader.Close())
	require.Equal(t, []byte("content"), data)
	require.NoError(t, driver.Delete(context.Background(), "2026/0826/report.txt"))
	require.NoError(t, driver.Put(context.Background(), "2026/0826/report.txt", strings.NewReader("replaced"), 8))
	reader, err = driver.OpenReader(context.Background(), "2026/0826/report.txt")
	require.NoError(t, err)
	data, err = io.ReadAll(reader)
	require.NoError(t, err)
	require.NoError(t, reader.Close())
	require.Equal(t, []byte("replaced"), data)
	require.NoError(t, driver.Delete(context.Background(), "2026/0826/report.txt"))
}

func TestLocalStorageRejectsEscapingPath(t *testing.T) {
	driver := localStorage{savePath: t.TempDir()}
	err := driver.Create(context.Background(), "../outside.txt", strings.NewReader("x"), 1)
	require.Error(t, err)
}
