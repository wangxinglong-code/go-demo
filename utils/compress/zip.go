package compress

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"strings"
)

// 压缩
func Compress(src []string, dest string) error {
	files := []*os.File{}
	for _, f := range src {
		file, err := os.Open(f)
		if err != nil {
			return err
		}
		files = append(files, file)
	}
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		// 增加对空目录的判断
		if len(fileInfos) <= 0 {
			header, err := zip.FileInfoHeader(info)
			header.Name = prefix
			if err != nil {
				return err
			}
			_, err = zw.CreateHeader(header)
			if err != nil {
				return err
			}
			file.Close()
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//解压
func DeCompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		filename := dest + file.Name
		path, err := getDir(filename)
		if err != nil {
			return err
		}
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

func getDir(path string) (string, error) {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) (string, error) {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		return "", errors.New("start is wrong")
	}

	if end < start || end > length {
		return "", errors.New("end is wrong")
	}

	return string(rs[start:end]), nil
}
