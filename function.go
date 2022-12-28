package eval

import (
	"bytes"
	"go/ast"
	"strings"
	"time"
)

func init() {
	funcMap["left_pad"] = leftPad
	funcMap["right_pad"] = rightPad
	funcMap["date_parse"] = dateParse
	funcMap["sub_str"] = subStr
	funcMap["contact"] = contact
	funcMap["index_of"] = indexOf
}

func leftPad(args []ast.Expr, data map[string]any) (any, error) {
	str, padStr, toLength, err := _padParam(args, data)
	if err != nil {
		return nil, err
	}

	if len(str) >= toLength {
		return str, nil
	}

	needPad := toLength - len(str)
	var buffer bytes.Buffer
	for i := 0; i < needPad; i++ {
		buffer.WriteString(padStr)
	}
	buffer.WriteString(str)

	return buffer.String(), nil
}

func rightPad(args []ast.Expr, data map[string]any) (any, error) {
	str, padStr, toLength, err := _padParam(args, data)
	if err != nil {
		return nil, err
	}

	if len(str) >= toLength {
		return str, nil
	}

	needPad := toLength - len(str)
	var buffer bytes.Buffer
	buffer.WriteString(str)
	for i := 0; i < needPad; i++ {
		buffer.WriteString(padStr)
	}

	return buffer.String(), nil
}

func subStr(args []ast.Expr, data map[string]any) (any, error) {
	if len(args) != 3 {
		return nil, illegalParams
	}

	strArg, err := eval(args[0], data)
	if err != nil {
		return nil, err
	}
	startArg, err := eval(args[1], data)
	if err != nil {
		return nil, err
	}
	endArg, err := eval(args[2], data)
	if err != nil {
		return nil, err
	}

	str := String(strArg)
	start, err := Int(startArg)
	if err != nil {
		return nil, err
	}
	end, err := Int(endArg)
	if err != nil {
		return nil, err
	}

	if start > end || start > len(str) || end > len(str) {
		return nil, illegalParams
	}

	return str[start:end], nil
}

func contact(args []ast.Expr, data map[string]any) (any, error) {
	ss := make([]string, 0, len(args))

	for _, arg := range args {
		strArg, err := eval(arg, data)
		if err != nil {
			return nil, err
		}
		ss = append(ss, String(strArg))
	}

	return strings.Join(ss, ""), nil
}

func indexOf(args []ast.Expr, data map[string]any) (any, error) {
	if len(args) != 2 {
		return nil, illegalParams
	}

	strArg, err := eval(args[0], data)
	if err != nil {
		return nil, err
	}
	subStrArg, err := eval(args[1], data)
	if err != nil {
		return nil, err
	}

	str := String(strArg)
	sub := String(subStrArg)

	return strings.Index(str, sub), nil
}

func dateParse(args []ast.Expr, data map[string]any) (any, error) {
	if len(args) != 2 {
		return nil, illegalParams
	}

	dateArg, err := eval(args[0], data)
	if err != nil {
		return nil, err
	}
	formatArg, err := eval(args[1], data)
	if err != nil {
		return nil, err
	}

	date := String(dateArg)
	format := String(formatArg) // YYYY MM DD 格式

	return _parseDate(format, date)
}

func _padParam(args []ast.Expr, data map[string]any) (string, string, int, error) {
	if len(args) != 3 {
		return "", "", 0, illegalParams
	}
	strArg, err := eval(args[0], data)
	if err != nil {
		return "", "", 0, err
	}
	padStrArg, err := eval(args[1], data)
	if err != nil {
		return "", "", 0, err
	}
	toLengthArg, err := eval(args[2], data)
	if err != nil {
		return "", "", 0, err
	}

	str := String(strArg)
	padStr := String(padStrArg)
	toLength, err := Int(toLengthArg)
	return str, padStr, toLength, err
}

func _parseDate(format, date string) (time.Time, error) {
	if strings.Contains(format, "YYYY") {
		format = strings.ReplaceAll(format, "YYYY", "2006")
	}
	if strings.Contains(format, "yyyy") {
		format = strings.ReplaceAll(format, "yyyy", "2006")
	}
	if strings.Contains(format, "MM") {
		format = strings.ReplaceAll(format, "MM", "01")
	}
	if strings.Contains(format, "DD") {
		format = strings.ReplaceAll(format, "DD", "02")
	}
	if strings.Contains(format, "dd") {
		format = strings.ReplaceAll(format, "dd", "02")
	}
	if strings.Contains(format, "HH") {
		format = strings.ReplaceAll(format, "HH", "15")
	}
	if strings.Contains(format, "hh") {
		format = strings.ReplaceAll(format, "hh", "15")
	}
	if strings.Contains(format, "mm") {
		format = strings.ReplaceAll(format, "mm", "04")
	}
	if strings.Contains(format, "ss") {
		format = strings.ReplaceAll(format, "ss", "05")
	}
	// 纳秒
	format, date = _nsFormatStyle(format, date, "SSSSSSSSS", "000000000")

	// 微秒
	format, date = _nsFormatStyle(format, date, "SSSSSS", "000000")

	// 毫秒
	format, date = _nsFormatStyle(format, date, "SSS", "000")

	return time.Parse(format, date)
}

func _nsFormatStyle(format, date string, old, new string) (string, string) {
	if index := strings.Index(format, old); index > 0 {
		format = strings.ReplaceAll(format, old, new)
		if format[index-1:index] != "." {
			format = format[:index] + "." + format[index:]
			date = date[:index] + "." + date[index:]
		}
	}
	return format, date
}
