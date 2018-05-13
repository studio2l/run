package main

// run은 환경변수를 설정한 후 다른 프로그램 실행한다.

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// die는 받아들인 에러를 출력하고 프로그램을 종료한다.
func die(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

// getEnv는 환경변수 리스트에서 해당 키를 찾아 그 값을 반환한다.
// 만약 키가 없다면 빈 문자열을 반환한다.
func getEnv(key string, env []string) string {
	for _, e := range env {
		kv := strings.SplitN(e, "=", -1)
		if len(kv) != 2 {
			continue
		}
		k := strings.TrimSpace(kv[0])
		v := strings.TrimSpace(kv[1])
		if k == key {
			return v
		}
	}
	return ""
}

// parseEnv는 env 환경변수 리스트를 참조해서 e 환경변수의 값을 해석한다.
// 예를들어 e가 "TEST=$OTHER", env가 []string{"OTHER=test"} 라면
// "TEST=test", nil을 반환한다.
// e가 환경변수 문자열로 변경 불가능 하다면 빈문자열과 에러를 반환한다.
// 반환되는 문자열 키와 값 앞 뒤의 공백은 제거된다.
func parseEnv(e string, env []string) (string, error) {
	kv := strings.SplitN(e, "=", -1)
	if len(kv) != 2 {
		return "", fmt.Errorf("invalid env string: %s", e)
	}
	k := strings.TrimSpace(kv[0])
	if k == "" {
		return "", fmt.Errorf("env key empty")
	}
	if strings.Contains(k, "$") {
		return "", fmt.Errorf("env key should not have '$' char: %s", k)
	}
	v := strings.TrimSpace(kv[1])
	re := regexp.MustCompile(`[$]\w+`)
	for {
		idxs := re.FindStringIndex(v)
		if idxs == nil {
			break
		}
		s := idxs[0]
		e := idxs[1]
		pre := v[:s]
		post := v[e:]
		envk := v[s+1 : e]
		envv := getEnv(envk, env)
		v = pre + envv + post
	}
	return k + "=" + v, nil
}

// parseEnvFile은 파일을 읽어 그 안의 환경변수 문자열을 리스트 형태로 반환한다.
// 파일을 읽는 도중 에러가 나거나 환경변수 파싱이 불가능하다면 빈문자열과 에러를 반환한다.
func parseEnvFile(f string, env []string) ([]string, error) {
	b, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	s := string(b)
	penv := []string{}
	for _, l := range strings.Split(s, "\n") {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if strings.HasPrefix(l, "#") {
			continue
		}
		e, err := parseEnv(l, env)
		if err != nil {
			return nil, err
		}
		penv = append(penv, e)

		// 기존 환경변수에도 추가하는데 그 이유는
		// 앞줄의 환경변수가 뒷줄의 환경변수를 완성하는데 쓰일 수 있기 때문이다.
		// 이렇게 한다고 해도 이 함수를 부른 함수의 env에는 적용되지 않는다.
		env = append(env, e)
	}
	return penv, nil
}

// Config는 커맨드라인 옵션 값을 담는다.
type Config struct {
	env     string
	envfile string
	dir     string
}

func main() {
	cfg := Config{}
	flag.StringVar(&cfg.env, "env", "", "미리 선언할 환경변수. envfile에 앞서 설정됩니다. 콤마(,)를 이용해 여러 환경변수를 설정할 수 있습니다.")
	flag.StringVar(&cfg.envfile, "envfile", "", "환경변수들이 설정되어있는 파일을 불러옵니다. 콤마(,)를 이용해 여러 파일을 불러 올 수 있습니다. ?(물음표) 로 시작하는 파일은 없어도 에러가 나지 않습니다.")
	flag.StringVar(&cfg.dir, "dir", "", "명령을 실행할 디렉토리를 설정합니다. 설정하지 않으면 현재 디렉토리에서 실행합니다.")
	flag.Parse()

	// OS 환경변수에 env와 envfile을 파싱해 환경변수를 추가/대체한다.
	env := os.Environ()
	for _, e := range strings.Split(cfg.env, ",") {
		e = strings.TrimSpace(e)
		if e == "" {
			continue
		}
		var err error
		e, err = parseEnv(e, env)
		if err != nil {
			die(err)
		}
		env = append(env, e)
	}
	for _, envf := range strings.Split(cfg.envfile, ",") {
		envf = strings.TrimSpace(envf)
		dieNoFile := true
		if len(envf) > 0 && envf[0] == '?' {
			dieNoFile = false
			envf = envf[1:]
		}
		if envf == "" {
			continue
		}
		envs, err := parseEnvFile(envf, env)
		if err != nil {
			if os.IsNotExist(err) && dieNoFile {
				die(err)
			}
		}
		for _, e := range envs {
			env = append(env, e)
		}
	}

	// 설정된 환경으로 명령을 실행한다.
	cmds := flag.Args()
	if len(cmds) == 0 {
		die(errors.New("need command to run"))
	}
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Env = env
	cmd.Dir = cfg.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		die(err)
	}
}
