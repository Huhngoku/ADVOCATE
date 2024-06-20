package explanation

import (
	"fmt"
	"io"
	"os"
)

func copyTrace(path string, index int) error {
	// copy the folder "path/advocateTrace"
	// to "path/bugs/bug_index/advocateTrace"
	err := copyDir(path+"advocateTrace", path+"bugs/bug_"+fmt.Sprint(index)+"/advocateTrace")
	if err != nil {
		return err
	}

	return nil
}

func copyRewrite(path string, index int) error {
	// copy the folder "path/rewritten_trace_index"
	// to "path/bugs/bug_index/rewritten_trace"
	err := copyDir(path+"rewritten_trace_"+fmt.Sprint(index), path+"bugs/bug_"+fmt.Sprint(index)+"/rewritten_trace")
	if err != nil {
		return err
	}

	return nil
}

func copyDir(src string, dst string) error {
	// copy the folder "src" to "dst"
	// if the folder "dst" does not exist, create it
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		err := os.Mkdir(dst, 0755)
		if err != nil {
			return err
		}
	}

	// get the content of the folder "src"
	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// copy each file from "src" to "dst"
	for _, file := range files {
		if file.Name() == "times.log" {
			continue
		}

		srcFile := src + "/" + file.Name()
		dstFile := dst + "/" + file.Name()

		err := copyFile(srcFile, dstFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src string, dst string) error {
	// copy the file "src" to "dst"
	// open the file "src"
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// create the file "dst"
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	// copy the content of "src" to "dst"
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
