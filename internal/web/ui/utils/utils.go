// Package utils provides small class-name helpers used by the vendored
// heroui-go components (copied from workspace/heroui-go).
package utils

import "strings"

// Cn joins class name arguments (string or []string), skipping empty ones.
func Cn(args ...any) string {
	var sb strings.Builder

	for _, a := range args {
		switch v := a.(type) {
		case string:
			if v != "" {
				sb.WriteString(v)
				sb.WriteByte(' ')
			}
		case []string:
			for _, s := range v {
				if s != "" {
					sb.WriteString(s)
					sb.WriteByte(' ')
				}
			}
		}
	}

	return sb.String()
}

// If is the ternary operator: cond ? then : els
func If[T any](cond bool, then, els T) T {
	if cond {
		return then
	}
	return els
}

func BoolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
