package main

import (
	"flag"
	"fmt"
	"vendor_helper/utils"
)

func main() {
	var proPath string // 项目路径
	var pkgURL string  // 包路径
	flag.StringVar(&proPath, "pro", "", "项目路径，会把包文件考本到项目下的vendor中")
	flag.StringVar(&pkgURL, "pkg", "", "包名，需要安装包的地址")
	flag.Parse()

	fmt.Println("hello world")
	fmt.Println(proPath)
	fmt.Println(pkgURL)

	gopath := utils.GetGOPATH()
	fmt.Println(gopath)

	pkgURL = "github.com/golang/snappy"
	proPath = `C:\工作文件\TsepVoteServer`
	err := utils.CopyPkg(pkgURL, proPath)
	if err != nil {
		fmt.Printf("CopyPkg failed, err: %v\n", err)
	}

	ori := `C:\workspace_go\test_pb`
	target := `C:\workspace_go\fff`
	err = utils.CopyDir(ori, target)
	if err != nil {
		fmt.Println("文件拷贝失败", err)
	}
}
