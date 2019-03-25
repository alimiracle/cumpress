package main

import (
    "archive/zip"
    "fmt"
    "io"
    "os"
    "path/filepath"
"strings"
)

func unzip(src, dest string) error {
    r, err := zip.OpenReader(src)
    if err != nil {
        return err
    }
    defer func() {
        if err := r.Close(); err != nil {
            panic(err)
        }
    }()

    os.MkdirAll(dest, 0755)

    // Closure to address file descriptors issue with all the deferred .Close() methods
    extractAndWriteFile := func(f *zip.File) error {
        rc, err := f.Open()
        if err != nil {
            return err
        }
        defer func() {
            if err := rc.Close(); err != nil {
                panic(err)
            }
        }()

        path := filepath.Join(dest, f.Name)

        if f.FileInfo().IsDir() {
            os.MkdirAll(path, 0755)
        } else {
            os.MkdirAll(filepath.Dir(path), 0755)
            f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
            if err != nil {
                return err
            }
            defer func() {
                if err := f.Close(); err != nil {
                    panic(err)
                }
            }()

            _, err = io.Copy(f, rc)
            if err != nil {
                return err
            }
        }
        return nil
    }

    for _, f := range r.File {
        err := extractAndWriteFile(f)
        if err != nil {
            return err
        }
    }

    return nil
}
func listFiles(root string) ([]string, error) {
    var files []string

    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
if !info.IsDir() {
    files = append(files, path)
}
        return nil
    })
    if err != nil {
        return nil, err
    }
    return files, nil

}

func zipMe(filepaths []string, target string) error {

    flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
    file, err := os.OpenFile(target, flags, 0644)

    if err != nil {
        return fmt.Errorf("Failed to open zip for writing: %s", err)
    }
    defer file.Close()

    zipw := zip.NewWriter(file)
    defer zipw.Close()

    for _, filename := range filepaths {
        if err := addFileToZip(filename, zipw); err != nil {
            return fmt.Errorf("Failed to add file %s to zip: %s", filename, err)
        }
    }
    return nil

}

func addFileToZip(filename string, zipw *zip.Writer) error {
    file, err := os.Open(filename)

    if err != nil {
        return fmt.Errorf("Error opening file %s: %s", filename, err)
    }
    defer file.Close()

    wr, err := zipw.Create(filename)
    if err != nil {

        return fmt.Errorf("Error adding file; '%s' to zip : %s", filename, err)
    }

    if _, err := io.Copy(wr, file); err != nil {
        return fmt.Errorf("Error writing %s to zip: %s", filename, err)
    }

    return nil
}
func zip_dir(dir string, out string) {

    files, err := listFiles(dir)
    if err != nil {
        panic(err)
    }

    zipMe(files, out)

    for _, f := range files {
        fmt.Println(f)
    }
    fmt.Println("Done!")
}
func zip_file(file string, out string) {

    zipMe(strings.Fields(file), out)

    fmt.Println("Done!")
}
func run_zip(src string, des string) {
    fi, err := os.Stat(src)
    if err != nil {
        fmt.Println(err)
        return
    }
    switch mode := fi.Mode(); {
    case mode.IsDir():
        // do directory stuff
zip_dir(src, des)

    case mode.IsRegular():
        // do file stuff
zip_file(src, des)

    }
}
func run(file_type string, filename string, output string) {
if file_type=="--zip" {

run_zip(filename, output)
} else if file_type=="--unzip" {
unzip(filename, output)
} else if file_type=="-u" {
unzip(filename, output)
} else if file_type=="-z" {
run_zip(filename, output)
} else {
fmt.Println("type error")

}
}
func main() {
env := os.Args
if len(env)==1 {
fmt.Println("no file type")

} else if len(env)==2 {
fmt.Println("no input file")
} else if len(env)==3 {
fmt.Println("no output file")

} else {
run(env[1], env[2], env[3])

}
}

