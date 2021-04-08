package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strings"
)

var (
	splitor = "/"
)

func GetGOPATH() string {
	cmd := exec.Command("go", "env")
	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut
	e := cmd.Run()
	if e != nil {
		fmt.Printf("cmd.Run() | err: %v\n", e)
		return ""
	}

	fmt.Println("cmd package")

	reg := regexp.MustCompile(`set GOPATH=(.+)`)
	l := reg.FindStringSubmatch(cmdOut.String())
	if len(l) < 2 {
		fmt.Println("没有找到gopath路径")
		return ""
	}
	return l[1]
}

func DownPkg(URL string) error {
	cmd := exec.Command("go", "get", URL)
	var cmdOut bytes.Buffer
	cmd.Stdout = &cmdOut
	e := cmd.Run()
	if e != nil {
		return e
	}

	fmt.Println(cmdOut.String())
	return nil
}

func GetPkgParentPath(URL string) string {
	modPath := "/pkg/mod/"
	goPath := GetGOPATH()
	p1, _ := ParseURL(URL)
	out := fmt.Sprintf("%s%s%s", goPath, modPath, p1)
	return out
}

func GetPkgVendorPath(pkgPath string, proPath string) string {
	p1, p2 := ParseURL(pkgPath)
	return fmt.Sprintf("%s/vendor/%s/%s", proPath, p1, p2)
}

func GetPkgDir(URL string) (string, error) {
	_, pkgName := ParseURL(URL)
	parentDir := GetPkgParentPath(URL)

	t := make([]string, 0)
	files, _ := ioutil.ReadDir(parentDir)
	for _, f := range files {
		find, err := regexp.MatchString(fmt.Sprintf("%s.*", pkgName), f.Name())
		if err != nil {
			fmt.Printf("regexp err: %v\n", err)
		}
		if find {
			fmt.Println(f.Name())
			t = append(t, f.Name())
		}
	}

	sort.Strings(t)
	l := len(t)
	outDir := fmt.Sprintf("%s%s%s", parentDir, splitor, t[l-1])
	fmt.Println("outPath: ", outDir)
	return outDir, nil
}

func CopyDir(src string, dst string) error {
	// src = fmt.Sprintf(`'%s'`, src)
	// dst = fmt.Sprintf(`'%s'`, dst)
	src = strings.Replace(src, splitor, `\`, -1) // splitor = `\`
	dst = strings.Replace(dst, splitor, `\`, -1)
	// src = strings.Replace(src, `\`, splitor, -1)
	// dst = strings.Replace(dst, `\`, splitor, -1)
	fmt.Println("源文件夹：", src)
	fmt.Println("目标文件夹", dst)

	_, err := os.Stat(dst)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dst, os.ModePerm)
		}
		fmt.Printf("os.Stat | err: %v\n", err)
	}

	var cmd *exec.Cmd
	var cmdOut bytes.Buffer
	switch runtime.GOOS {
	case "windows":
		fmt.Println("windows")
		cmd = exec.Command("xcopy", src, dst, "/I", "/E", "/y")
		cmd.Stdout = &cmdOut
		e := cmd.Run()
		fmt.Println("cmdOut: ", cmdOut.String())
		if e != nil {
			fmt.Printf("cmd.Run() | err: %v\n", e)
			return e
		}
	case "darwin", "linux":
		fmt.Println("darwin", "linux")
		cmd = exec.Command("cp", "-R", src, dst)
	}

	// outPut, e := cmd.Output()
	// if e != nil {
	// 	return e
	// }
	// fmt.Printf("output: %s\n", string(outPut))
	return nil
}

func FormatUrl(URL string) string {
	return strings.Split(URL, "@")[0]
}

func ParseURL(URL string) (string, string) {
	URL = FormatUrl(URL)
	paths := strings.Split(URL, splitor)
	l := len(paths)
	return strings.Join(paths[:l-1], splitor), paths[l-1]
}

func CopyPkg(pkgURL string, proPath string) error {
	// 下载包
	err := DownPkg(pkgURL)
	if err != nil {
		return err
	}

	// 拷贝包
	srcPath, err := GetPkgDir(pkgURL)
	if err != nil {
		return err
	}
	dstPath := GetPkgVendorPath(pkgURL, proPath)
	err = CopyDir(srcPath, dstPath)
	if err != nil {
		return err
	}
	fmt.Printf("%s 拷贝到 %s 成功！\n", pkgURL, proPath)
	return nil
}
