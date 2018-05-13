// +build windows

package main

import "strings"

// envSepFromColon은 환경변수의 값 안에 있는 콜론(:)을
// 해당 OS의 환경변수 구분자로 변경한다.
func envSepFromColon(v string) string {
	return strings.Replace(v, ":", ";", -1)
}
