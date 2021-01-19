package utils

import (
	"archive/zip"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Eval(expStr string) int {
	regx, err := regexp.Compile(`[A-F0-9]{2}`)
	if err != nil {
		panic(err)
	}
	result := regx.FindAllString(expStr, -1)

	for _, v := range result {
		nv, err := strconv.ParseUint(v, 16, 16)
		if err != nil {
			panic(err)
		}
		expStr = strings.Replace(expStr, v, strconv.Itoa(int(nv)), -1)
	}

	exp, err := parser.ParseExpr(expStr)
	if err != nil {
		// fmt.Println("expStr:", expStr)
		panic(err)
	}
	return evalExp(exp)
}

func evalExp(exp ast.Expr) int {
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return evalExpBinaryExpr(exp)
	case *ast.BasicLit:
		switch exp.Kind {
		case token.INT:
			i, _ := strconv.Atoi(exp.Value)
			return i
		}
	}

	return 0
}

func evalExpBinaryExpr(exp *ast.BinaryExpr) int {
	left := evalExp(exp.X)
	right := evalExp(exp.Y)

	switch exp.Op {
	case token.ADD:
		return left + right
	case token.SUB:
		return left - right
	case token.MUL:
		return left * right
	case token.QUO:
		return left / right
	}

	return 0
}

func ReadFile(path string) []byte {
	fi, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	file := make([]byte, fi.Size()+1)
	l, err := f.Read(file)
	if err != nil {
		panic(err)
	}
	return file[:l]
}

func ReadZipFile(file *zip.File) []byte {
	fc, err := file.Open()
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadAll(fc)
	if err != nil {
		panic(err)
	}

	err = fc.Close()
	if err != nil {
		panic(err)
	}

	return content
}
