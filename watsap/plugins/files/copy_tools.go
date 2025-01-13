package files

import (
    "io"
    "os"
    "path/filepath"
    "strings"
)

func CopyToolsDir() error {
    srcDir := `D:\tools`
    destDir := `/desired/destination/path`

    return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        relPath := strings.TrimPrefix(path, srcDir)
        destPath := filepath.Join(destDir, relPath)

        if info.IsDir() {
            return os.MkdirAll(destPath, info.Mode())
        }

        srcFile, err := os.Open(path)
        if err != nil {
            return err
        }
        defer srcFile.Close()

        destFile, err := os.Create(destPath)
        if err != nil {
            return err
        }
        defer destFile.Close()

        _, err = io.Copy(destFile, srcFile)
        return err
    })
}
