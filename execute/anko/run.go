package anko

import (
	"fmt"
	"github.com/mattn/anko/core"
	"github.com/mattn/anko/env"
	_ "github.com/mattn/anko/packages"
	"github.com/mattn/anko/vm"
	"golang.design/x/clipboard"
	"log"
	"os"
	"reflect"
	"strings"
)

func Execute(sf string, args []string) error {
	scriptb, err := os.ReadFile(sf)
	if err != nil {
		return err
	}
	e, err := initEnv(args)
	if err != nil {
		log.Fatalf("init env err: %v\n", err)
		return err
	}
	script := string(scriptb)
	e, err = initClip(script, e)
	if err != nil {
		return err
	}

	_, err = vm.Execute(e, nil, script)
	if err != nil {
		return err
	}

	oclip, err := e.Get("_oclip")
	if err == nil {
		if c, ok := oclip.(string); ok {
			clipboard.Write(clipboard.FmtText, []byte(c))
		}
	}
	return nil
}

func initEnv(args []string) (*env.Env, error) {
	e := env.NewEnv()

	core.Import(e)
	// 增加fmt包, 和 go 保持一致
	e.Define("Sprintf", fmt.Sprintf)
	e.Define("Printf", fmt.Printf)
	e.Define("Println", fmt.Println)
	e.Define("Print", fmt.Print)
	ankoPackage()

	// 补充命令行参数
	if len(args) > 0 {
		for i, v := range args {
			v := v
			e.Define(fmt.Sprintf("_%d", i+1), v)
		}
	}

	return e, nil
}

func initClip(s string, e *env.Env) (*env.Env, error) {
	if !(strings.Contains(s, "_clip") || strings.Contains(s, "_oclip")) {
		return e, nil
	}
	err := initClipboard()
	if err != nil {
		return e, err
	}
	// 剪贴板内容
	if strings.Contains(s, "_clip") {
		clip := readClipboard()
		e.Define("_clip", clip)
	}
	if strings.Contains(s, "_oclip") {
		oclip := ""
		e.Define("_oclip", oclip)
	}
	return e, nil
}

func initClipboard() error {
	return clipboard.Init()
}

func readClipboard() string {
	bs := clipboard.Read(clipboard.FmtText)
	return string(bs)
}

func ankoPackage() {
	env.Packages["s_http"] = map[string]reflect.Value{
		"httpGet":  reflect.ValueOf(httpGet),
		"reqErr":   reflect.ValueOf(reqErr),
		"reqNotOk": reflect.ValueOf(reqNotOk),
	}
}
