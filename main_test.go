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
	got, err := parseEnvFile(f, []string{})
	if err != nil {
		t.Fatalf("parseEnvFile(%q, []): unexpected error: %v", f, err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseEnvfile(%q, []): got %s, want %s", f, got, want)
	}
}
