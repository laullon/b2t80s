package utils

import (
	"go/ast"
	"go/parser"
	"go/token"
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
