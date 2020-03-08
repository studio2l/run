package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseEnv(t *testing.T) {
	cases := []struct {
		e    string
		env  []string
		want string
	}{
		{"TEST=abc", []string{}, "TEST=abc"},
		{"TEST=$A", []string{"A=abc"}, "TEST=abc"},
		{"TEST=$B", []string{"A=abc"}, "TEST="},
		{"TEST=$A/$B", []string{"A=abc"}, "TEST=" + filepath.FromSlash("abc/")},
		{"TEST=$A/$B", []string{"A=abc", "B=def"}, "TEST=" + filepath.FromSlash("abc/def")},
		{"TEST=$A/$B", []string{"  A = abc ", " B =def  "}, "TEST=" + filepath.FromSlash("abc/def")},
	}
	for _, c := range cases {
		got, err := parseEnv(c.e, c.env)
		if err != nil {
			t.Fatalf("parseEnv(%q, %q): unexpected error: %v", c.e, c.env, err)
		}
		if got != c.want {
			t.Fatalf("parseEnv(%q, %q): got %s, want %s", c.e, c.env, got, c.want)
		}
	}
}

func TestParseEnvfileValid(t *testing.T) {
	want := []string{
		"TL_GLOBAL_PATH=" + filepath.FromSlash("Z:/VFX/global"),
		"TL_MAYA_PATH=" + filepath.FromSlash("Z:/VFX/global/maya"),
		"TL_HOUDINI_PATH=" + filepath.FromSlash("Z:/VFX/global/houdini"),
		"A=a",
		"B=b",
		"C=c",
		"D=d",
		"E=e",
		"ABC=abc",
		"ABCDE=abcde",
	}
	f := "testdata/valid.env"
	env := []string{}
	got, err := parseEnvFile(f, env)
	if err != nil {
		t.Fatalf("parseEnvFile(%q, %v): unexpected error: %v", f, env, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseEnvfile(%q, %v):\ngot %v\nwant %v", f, env, got, want)
	}
}

func TestParseEnvsetPaths(t *testing.T) {
	want := []string{
		"testdata/site/env/all.env",
		"testdata/site/env/maya/all.env?",
		"testdata/site/env/maya/lit.env?",
		"testdata/site/show/test/show.env?",
	}
	f := "testdata/site.envs"
	env := []string{"SITE_ROOT=testdata/site", "PROGRAM=maya", "TEAM=lit", "SHOW=test"}
	got, err := parseEnvsetFile(f, env)
	if err != nil {
		t.Fatalf("parseEnvsetPaths(%q, %v): unexpected error: %v", f, env, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseEnvsetPaths(%q, %v):\ngot %v\nwant %v", f, env, got, want)
	}
}

func TestParseEnvsetEnvs(t *testing.T) {
	want := []string{
		"SITE_ROOT=testdata/site",
		"PROGRAM=maya",
		"TEAM=lit",
		"SHOW=test",
		"ALL=all",
		"MAYA_ALL=maya_all",
		"MAYA_LIT=maya_all/lit",
	}
	f := "testdata/site.envs"
	env := []string{"SITE_ROOT=testdata/site", "PROGRAM=maya", "TEAM=lit", "SHOW=test"}
	got := env
	envfiles, err := parseEnvsetFile(f, got)
	if err != nil {
		t.Fatalf("parseEnvsetEnvs(%q, %v): unexpected error: %v", f, env, err)
	}
	for _, envf := range envfiles {
		envs, err := parseEnvFile(envf, got)
		if err != nil {
			t.Fatalf("parseEnvsetEnvs(%q, %v): unexpected error: %v", f, env, err)
		}
		for _, e := range envs {
			got = append(got, e)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseEnvsetEnvs(%q, %v):\ngot %v\nwant %v", f, env, got, want)
	}
}
