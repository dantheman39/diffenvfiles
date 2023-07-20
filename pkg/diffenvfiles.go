package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

type EnvFile struct {
	Path string
	Data []byte
}

func DiffEnvFiles(env1 EnvFile, env2 EnvFile) error {
	w := os.Stdout
	parseResult1, err := parseContents(env1.Data)
	if err != nil {
		return err
	}
	parseResult2, err := parseContents(env2.Data)
	if err != nil {
		return err
	}

	env1dups := findDups(getFirstElement(parseResult1.EnvVars))
	env2dups := findDups(getFirstElement(parseResult2.EnvVars))

	foundDups := false
	checkDups := func(dups []string, p string) {
		if len(dups) > 0 {
			foundDups = true
			fmt.Fprintf(w, "\nDuplicate variables found in %q\n", p)
			for _, d := range dups {
				fmt.Fprintf(w, "  %s\n", d)
			}
		}
	}
	checkDups(env1dups, env1.Path)
	checkDups(env2dups, env2.Path)
	if foundDups {
		return errors.New("duplicate entries found in env file")
	}

	env1Common, onlyA := extractOnlyInA(parseResult1.EnvVars, parseResult2.EnvVars)
	env2Common, onlyB := extractOnlyInA(parseResult2.EnvVars, parseResult1.EnvVars)

	if len(env1Common) != len(env2Common) {
		return errors.New("unexpected error, list of common keys between files doesn't match")
	}

	sameKeysDifferentValues := [][]string{}
	for i := 0; i < len(env1Common); i++ {
		e1val := env1Common[i]
		e2val := env2Common[i]
		k := e1val[0]
		k2 := e2val[0]
		if k != k2 {
			return fmt.Errorf("unexpected error, common keys weren't properly parsed and sorted. %q does not match %q", k, k2)
		}
		v1 := e1val[1]
		v2 := e2val[1]
		if v1 != v2 {
			sameKeysDifferentValues = append(sameKeysDifferentValues, []string{k, v1, v2})
		}
	}

	if len(sameKeysDifferentValues) == 0 {
		fmt.Fprintln(w, "\nAll shared variables had the same value.")
	} else {
		fmt.Fprintln(w, "\nThe following variables are different between the two env files")
		for _, sk := range sameKeysDifferentValues {
			fmt.Fprintf(w, "  %s=\n", sk[0])
			fmt.Fprintf(w, "    %s\n", sk[1])
			fmt.Fprintf(w, "    %s\n", sk[2])
		}
	}

	printOnly := func(o [][]string, p string) {
		if len(o) > 0 {
			fmt.Fprintf(w, "\nThe following variables were only in %q\n", p)
			for _, oa := range o {
				fmt.Fprintf(w, "  %s=%s\n", oa[0], oa[1])
			}
		}
	}
	printOnly(onlyA, env1.Path)
	printOnly(onlyB, env2.Path)

	return nil
}

func splitlines(s string) []string {
	// should work for windows
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

type LineParsingError struct {
	LineNumber int
	Content    string
	Err        error
}

type ParseContentsResult struct {
	EnvVars           [][]string
	LineParsingErrors []LineParsingError
}

func parseContents(contents []byte) (ParseContentsResult, error) {
	lines := splitlines(string(contents))

	envVars := [][]string{}
	for _, originalLine := range lines {
		line := strings.TrimSpace(originalLine)
		if line == "" {
			continue
		}
		if line[0:1] == "#" {
			continue
		}

		varNameAndVal := strings.Split(line, "=")
		varName := strings.TrimSpace(varNameAndVal[0])
		varVal := ""
		if len(varNameAndVal) > 1 {
			varVal = strings.TrimSpace(strings.Join(varNameAndVal[1:], "="))
		}
		envVars = append(envVars, []string{varName, varVal})
	}
	sort.Sort(sortableByFirst(envVars))
	return ParseContentsResult{EnvVars: envVars, LineParsingErrors: []LineParsingError{}}, nil
}

func extractOnlyInA(listA [][]string, listB [][]string) ([][]string, [][]string) {
	alsoInB := [][]string{}

	bKeys := getFirstElement(listB)
	isInB := func(s string) bool {
		for _, otherS := range bKeys {
			if otherS == s {
				return true
			}
		}
		return false
	}

	onlyInA := [][]string{}
	for _, element := range listA {
		if isInB(element[0]) {
			alsoInB = append(alsoInB, element)
		} else {
			onlyInA = append(onlyInA, element)

		}
	}
	return alsoInB, onlyInA
}

func findDups(keys []string) []string {
	m := map[string]bool{}
	dups := []string{}
	for _, k := range keys {
		if _, present := m[k]; present {
			dups = append(dups, k)
		} else {
			m[k] = true
		}
	}
	return dups
}

func getFirstElement(arr [][]string) []string {
	out := []string{}
	for _, a := range arr {
		out = append(out, a[0])
	}
	return out
}
