package gotest

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jstemmer/go-junit-report/v2/pkg/gtr"

	"github.com/google/go-cmp/cmp"
)

const testdataRoot = "../../../testdata/"

var tests = []struct {
	in       string
	expected []gtr.Event
}{
	{"01-pass",
		[]gtr.Event{
			{Type: "run_test", Name: "TestZ"},
			{Type: "end_test", Name: "TestZ", Result: "PASS", Duration: 60 * time.Millisecond},
			{Type: "run_test", Name: "TestA"},
			{Type: "end_test", Name: "TestA", Result: "PASS", Duration: 100 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "summary", Result: "ok", Name: "package/name", Duration: 160 * time.Millisecond},
		}},
	{"02-fail",
		[]gtr.Event{
			{Type: "run_test", Name: "TestOne"},
			{Type: "end_test", Name: "TestOne", Result: "FAIL", Duration: 20 * time.Millisecond},
			{Type: "output", Data: "\tfile_test.go:11: Error message"},
			{Type: "output", Data: "\tfile_test.go:11: Longer"},
			{Type: "output", Data: "\t\terror"},
			{Type: "output", Data: "\t\tmessage."},
			{Type: "run_test", Name: "TestTwo"},
			{Type: "end_test", Name: "TestTwo", Result: "PASS", Duration: 130 * time.Millisecond},
			{Type: "status", Result: "FAIL"},
			{Type: "output", Data: "exit status 1"},
			{Type: "summary", Result: "FAIL", Name: "package/name", Duration: 151 * time.Millisecond},
		}},
	{"03-skip",
		[]gtr.Event{
			{Type: "run_test", Name: "TestOne"},
			{Type: "end_test", Name: "TestOne", Result: "SKIP", Duration: 20 * time.Millisecond},
			{Type: "output", Data: "\tfile_test.go:11: Skip message"},
			{Type: "run_test", Name: "TestTwo"},
			{Type: "end_test", Name: "TestTwo", Result: "PASS", Duration: 130 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "summary", Result: "ok", Name: "package/name", Duration: 150 * time.Millisecond},
		}},
	{"04-go_1_4",
		[]gtr.Event{
			{Type: "run_test", Name: "TestOne"},
			{Type: "end_test", Name: "TestOne", Result: "PASS", Duration: 60 * time.Millisecond},
			{Type: "run_test", Name: "TestTwo"},
			{Type: "end_test", Name: "TestTwo", Result: "PASS", Duration: 100 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "summary", Result: "ok", Name: "package/name", Duration: 160 * time.Millisecond},
		}},
	// Test 05 is skipped, because it was actually testing XML output
	{"06-mixed",
		[]gtr.Event{
			{Type: "run_test", Name: "TestOne"},
			{Type: "end_test", Name: "TestOne", Result: "PASS", Duration: 60 * time.Millisecond},
			{Type: "run_test", Name: "TestTwo"},
			{Type: "end_test", Name: "TestTwo", Result: "PASS", Duration: 100 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "summary", Result: "ok", Name: "package/name1", Duration: 160 * time.Millisecond},
			{Type: "run_test", Name: "TestOne"},
			{Type: "end_test", Name: "TestOne", Result: "FAIL", Duration: 20 * time.Millisecond},
			{Type: "output", Data: "\tfile_test.go:11: Error message"},
			{Type: "output", Data: "\tfile_test.go:11: Longer"},
			{Type: "output", Data: "\t\terror"},
			{Type: "output", Data: "\t\tmessage."},
			{Type: "run_test", Name: "TestTwo"},
			{Type: "end_test", Name: "TestTwo", Result: "PASS", Duration: 130 * time.Millisecond},
			{Type: "status", Result: "FAIL"},
			{Type: "output", Data: "exit status 1"},
			{Type: "summary", Result: "FAIL", Name: "package/name2", Duration: 151 * time.Millisecond},
		}},
	{"07-compiled_test",
		[]gtr.Event{
			{Type: "run_test", Name: "TestOne"},
			{Type: "end_test", Name: "TestOne", Result: "PASS", Duration: 60 * time.Millisecond},
			{Type: "run_test", Name: "TestTwo"},
			{Type: "end_test", Name: "TestTwo", Result: "PASS", Duration: 100 * time.Millisecond},
			{Type: "status", Result: "PASS"},
		}},
	{"08-parallel",
		[]gtr.Event{
			{Type: "run_test", Name: "TestDoFoo"},
			{Type: "run_test", Name: "TestDoFoo2"},
			{Type: "end_test", Name: "TestDoFoo", Result: "PASS", Duration: 270 * time.Millisecond},
			{Type: "output", Data: "\tcov_test.go:10: DoFoo log 1"},
			{Type: "output", Data: "\tcov_test.go:10: DoFoo log 2"},
			{Type: "end_test", Name: "TestDoFoo2", Result: "PASS", Duration: 160 * time.Millisecond},
			{Type: "output", Data: "\tcov_test.go:21: DoFoo2 log 1"},
			{Type: "output", Data: "\tcov_test.go:21: DoFoo2 log 2"},
			{Type: "status", Result: "PASS"},
			{Type: "summary", Result: "ok", Name: "package/name", Duration: 440 * time.Millisecond},
		}},
	{"09-coverage",
		[]gtr.Event{
			{Type: "run_test", Name: "TestZ"},
			{Type: "end_test", Name: "TestZ", Result: "PASS", Duration: 60 * time.Millisecond},
			{Type: "run_test", Name: "TestA"},
			{Type: "end_test", Name: "TestA", Result: "PASS", Duration: 100 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "coverage", CovPct: 13.37},
			{Type: "summary", Result: "ok", Name: "package/name", Duration: 160 * time.Millisecond},
		}},
	{"10-multipkg-coverage",
		[]gtr.Event{
			{Type: "run_test", Name: "TestA"},
			{Type: "end_test", Name: "TestA", Result: "PASS", Duration: 100 * time.Millisecond},
			{Type: "run_test", Name: "TestB"},
			{Type: "end_test", Name: "TestB", Result: "PASS", Duration: 300 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "coverage", CovPct: 10},
			{Type: "summary", Result: "ok", Name: "package1/foo", Duration: 400 * time.Millisecond, CovPct: 10},
			{Type: "run_test", Name: "TestC"},
			{Type: "end_test", Name: "TestC", Result: "PASS", Duration: 4200 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "coverage", CovPct: 99.8},
			{Type: "summary", Result: "ok", Name: "package2/bar", Duration: 4200 * time.Millisecond, CovPct: 99.8},
		}},
	{"11-go_1_5",
		[]gtr.Event{
			{Type: "run_test", Name: "TestOne"},
			{Type: "end_test", Name: "TestOne", Result: "PASS", Duration: 20 * time.Millisecond},
			{Type: "run_test", Name: "TestTwo"},
			{Type: "end_test", Name: "TestTwo", Result: "PASS", Duration: 30 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "summary", Result: "ok", Name: "package/name", Duration: 50 * time.Millisecond},
		}},
	{"12-go_1_7",
		[]gtr.Event{
			{Type: "run_test", Name: "TestOne"},
			{Type: "run_test", Name: "TestOne/Child"},
			{Type: "run_test", Name: "TestOne/Child#01"},
			{Type: "run_test", Name: "TestOne/Child=02"},
			{Type: "end_test", Name: "TestOne", Result: "PASS", Duration: 10 * time.Millisecond},
			{Type: "end_test", Name: "TestOne/Child", Result: "PASS", Indent: 1, Duration: 20 * time.Millisecond},
			{Type: "end_test", Name: "TestOne/Child#01", Result: "PASS", Indent: 1, Duration: 30 * time.Millisecond},
			{Type: "end_test", Name: "TestOne/Child=02", Result: "PASS", Indent: 1, Duration: 40 * time.Millisecond},
			{Type: "run_test", Name: "TestTwo"},
			{Type: "run_test", Name: "TestTwo/Child"},
			{Type: "run_test", Name: "TestTwo/Child#01"},
			{Type: "run_test", Name: "TestTwo/Child=02"},
			{Type: "end_test", Name: "TestTwo", Result: "PASS", Duration: 10 * time.Millisecond},
			{Type: "end_test", Name: "TestTwo/Child", Result: "PASS", Indent: 1, Duration: 20 * time.Millisecond},
			{Type: "end_test", Name: "TestTwo/Child#01", Result: "PASS", Indent: 1, Duration: 30 * time.Millisecond},
			{Type: "end_test", Name: "TestTwo/Child=02", Result: "PASS", Indent: 1, Duration: 40 * time.Millisecond},
			{Type: "run_test", Name: "TestThree"},
			{Type: "run_test", Name: "TestThree/a#1"},
			{Type: "run_test", Name: "TestThree/a#1/b#1"},
			{Type: "run_test", Name: "TestThree/a#1/b#1/c#1"},
			{Type: "end_test", Name: "TestThree", Result: "PASS", Duration: 10 * time.Millisecond},
			{Type: "end_test", Name: "TestThree/a#1", Result: "PASS", Indent: 1, Duration: 20 * time.Millisecond},
			{Type: "end_test", Name: "TestThree/a#1/b#1", Result: "PASS", Indent: 2, Duration: 30 * time.Millisecond},
			{Type: "end_test", Name: "TestThree/a#1/b#1/c#1", Result: "PASS", Indent: 3, Duration: 40 * time.Millisecond},
			{Type: "run_test", Name: "TestFour"},
			{Type: "run_test", Name: "TestFour/#00"},
			{Type: "run_test", Name: "TestFour/#01"},
			{Type: "run_test", Name: "TestFour/#02"},
			{Type: "end_test", Name: "TestFour", Result: "FAIL", Duration: 20 * time.Millisecond},
			{Type: "end_test", Name: "TestFour/#00", Result: "FAIL", Indent: 1, Duration: 0},
			{Type: "output", Data: "    \texample.go:12: Expected abc  OBTAINED:"},
			{Type: "output", Data: "    \t\txyz"},
			{Type: "output", Data: "    \texample.go:123: Expected and obtained are different."},
			{Type: "end_test", Name: "TestFour/#01", Result: "SKIP", Indent: 1, Duration: 0},
			{Type: "output", Data: "    \texample.go:1234: Not supported yet."},
			{Type: "end_test", Name: "TestFour/#02", Result: "PASS", Indent: 1, Duration: 0},
			{Type: "run_test", Name: "TestFive"},
			{Type: "end_test", Name: "TestFive", Result: "SKIP", Duration: 0},
			{Type: "output", Data: "\texample.go:1392: Not supported yet."},
			{Type: "run_test", Name: "TestSix"},
			{Type: "end_test", Name: "TestSix", Result: "FAIL", Duration: 0},
			{Type: "output", Data: "\texample.go:371: This should not fail!"},
			{Type: "status", Result: "FAIL"},
			{Type: "summary", Result: "FAIL", Name: "package/name", Duration: 50 * time.Millisecond},
		}},
	{"13-syntax-error",
		[]gtr.Event{
			{Type: "output", Data: "# package/name/failing1"},
			{Type: "output", Data: "failing1/failing_test.go:15: undefined: x"},
			{Type: "output", Data: "# package/name/failing2"},
			{Type: "output", Data: "failing2/another_failing_test.go:20: undefined: y"},
			{Type: "output", Data: "# package/name/setupfailing1"},
			{Type: "output", Data: "setupfailing1/failing_test.go:4: cannot find package \"other/package\" in any of:"},
			{Type: "output", Data: "	/path/vendor (vendor tree)"},
			{Type: "output", Data: "	/path/go/root (from $GOROOT)"},
			{Type: "output", Data: "	/path/go/path (from $GOPATH)"},
			{Type: "run_test", Name: "TestA"},
			{Type: "end_test", Name: "TestA", Result: "PASS", Duration: 100 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "summary", Result: "ok", Name: "package/name/passing1", Duration: 100 * time.Millisecond},
			{Type: "run_test", Name: "TestB"},
			{Type: "end_test", Name: "TestB", Result: "PASS", Duration: 100 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "summary", Result: "ok", Name: "package/name/passing2", Duration: 100 * time.Millisecond},
			{Type: "summary", Result: "FAIL", Name: "package/name/failing1", Data: "[build failed]"},
			{Type: "summary", Result: "FAIL", Name: "package/name/failing2", Data: "[build failed]"},
			{Type: "summary", Result: "FAIL", Name: "package/name/setupfailing1", Data: "[setup failed]"},
		}},
	{"14-panic",
		[]gtr.Event{
			{Type: "output", Data: "panic: init"},
			{Type: "output", Data: "stacktrace"},
			{Type: "summary", Result: "FAIL", Name: "package/panic", Duration: 3 * time.Millisecond},
			{Type: "output", Data: "panic: init"},
			{Type: "output", Data: "stacktrace"},
			{Type: "summary", Result: "FAIL", Name: "package/panic2", Duration: 3 * time.Millisecond},
		}},
	{"15-empty",
		[]gtr.Event{
			{Type: "output", Data: "testing: warning: no tests to run"},
			{Type: "status", Result: "PASS"},
			{Type: "summary", Result: "ok", Name: "package/empty", Duration: 1 * time.Millisecond},
		}},
	{"16-repeated-names",
		[]gtr.Event{
			{Type: "run_test", Name: "TestRepeat"},
			{Type: "end_test", Name: "TestRepeat", Result: "PASS"},
			{Type: "run_test", Name: "TestRepeat"},
			{Type: "end_test", Name: "TestRepeat", Result: "PASS"},
			{Type: "run_test", Name: "TestRepeat"},
			{Type: "end_test", Name: "TestRepeat", Result: "PASS"},
			{Type: "status", Result: "PASS"},
			{Type: "summary", Result: "ok", Name: "package/repeated-names", Duration: 1 * time.Millisecond},
		}},
	{"17-race",
		[]gtr.Event{
			{Type: "run_test", Name: "TestRace"},
			{Type: "output", Data: "test output"},
			{Type: "output", Data: "2 0xc4200153d0"},
			{Type: "output", Data: "=================="},
			{Type: "output", Data: "WARNING: DATA RACE"},
			{Type: "output", Data: "Write at 0x00c4200153d0 by goroutine 7:"},
			{Type: "output", Data: "  race_test.TestRace.func1()"},
			{Type: "output", Data: "      race_test.go:13 +0x3b"},
			{Type: "output", Data: ""},
			{Type: "output", Data: "Previous write at 0x00c4200153d0 by goroutine 6:"},
			{Type: "output", Data: "  race_test.TestRace()"},
			{Type: "output", Data: "      race_test.go:15 +0x136"},
			{Type: "output", Data: "  testing.tRunner()"},
			{Type: "output", Data: "      /usr/local/Cellar/go/1.8.3/libexec/src/testing/testing.go:657 +0x107"},
			{Type: "output", Data: ""},
			{Type: "output", Data: "Goroutine 7 (running) created at:"},
			{Type: "output", Data: "  race_test.TestRace()"},
			{Type: "output", Data: "      race_test.go:14 +0x125"},
			{Type: "output", Data: "  testing.tRunner()"},
			{Type: "output", Data: "      /usr/local/Cellar/go/1.8.3/libexec/src/testing/testing.go:657 +0x107"},
			{Type: "output", Data: ""},
			{Type: "output", Data: "Goroutine 6 (running) created at:"},
			{Type: "output", Data: "  testing.(*T).Run()"},
			{Type: "output", Data: "      /usr/local/Cellar/go/1.8.3/libexec/src/testing/testing.go:697 +0x543"},
			{Type: "output", Data: "  testing.runTests.func1()"},
			{Type: "output", Data: "      /usr/local/Cellar/go/1.8.3/libexec/src/testing/testing.go:882 +0xaa"},
			{Type: "output", Data: "  testing.tRunner()"},
			{Type: "output", Data: "      /usr/local/Cellar/go/1.8.3/libexec/src/testing/testing.go:657 +0x107"},
			{Type: "output", Data: "  testing.runTests()"},
			{Type: "output", Data: "      /usr/local/Cellar/go/1.8.3/libexec/src/testing/testing.go:888 +0x4e0"},
			{Type: "output", Data: "  testing.(*M).Run()"},
			{Type: "output", Data: "      /usr/local/Cellar/go/1.8.3/libexec/src/testing/testing.go:822 +0x1c3"},
			{Type: "output", Data: "  main.main()"},
			{Type: "output", Data: "      _test/_testmain.go:52 +0x20f"},
			{Type: "output", Data: "=================="},
			{Type: "end_test", Name: "TestRace", Result: "FAIL"},
			{Type: "output", Data: "\ttesting.go:610: race detected during execution of test"},
			{Type: "status", Result: "FAIL"},
			{Type: "output", Data: "exit status 1"},
			{Type: "summary", Result: "FAIL", Name: "race_test", Duration: 15 * time.Millisecond},
		}},
	{"18-coverpkg",
		[]gtr.Event{
			{Type: "run_test", Name: "TestA"},
			{Type: "end_test", Name: "TestA", Result: "PASS", Duration: 100 * time.Millisecond},
			{Type: "run_test", Name: "TestB"},
			{Type: "end_test", Name: "TestB", Result: "PASS", Duration: 300 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "coverage", CovPct: 10, CovPackages: []string{"fmt", "encoding/xml"}},
			{Type: "summary", Result: "ok", Name: "package1/foo", Duration: 400 * time.Millisecond, CovPct: 10, CovPackages: []string{"fmt", "encoding/xml"}},
			{Type: "run_test", Name: "TestC"},
			{Type: "end_test", Name: "TestC", Result: "PASS", Duration: 4200 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "coverage", CovPct: 99.8, CovPackages: []string{"fmt", "encoding/xml"}},
			{Type: "summary", Result: "ok", Name: "package2/bar", Duration: 4200 * time.Millisecond, CovPct: 99.8, CovPackages: []string{"fmt", "encoding/xml"}},
		}},
	{"19-pass",
		[]gtr.Event{
			{Type: "run_test", Name: "TestZ"},
			{Type: "output", Data: "some inline text"},
			{Type: "end_test", Name: "TestZ", Result: "PASS", Duration: 60 * time.Millisecond},
			{Type: "run_test", Name: "TestA"},
			{Type: "end_test", Name: "TestA", Result: "PASS", Duration: 100 * time.Millisecond},
			{Type: "status", Result: "PASS"},
			{Type: "summary", Result: "ok", Name: "package/name", Duration: 160 * time.Millisecond},
		}},
	{"20-parallel",
		[]gtr.Event{
			{Type: "run_test", Name: "FirstTest"},
			{Type: "output", Data: "Message from first"},
			{Type: "pause_test", Name: "FirstTest"},
			{Type: "run_test", Name: "SecondTest"},
			{Type: "output", Data: "Message from second"},
			{Type: "pause_test", Name: "SecondTest"},
			{Type: "cont_test", Name: "FirstTest"},
			{Type: "output", Data: "Supplemental from first"},
			{Type: "run_test", Name: "ThirdTest"},
			{Type: "output", Data: "Message from third"},
			{Type: "end_test", Name: "ThirdTest", Result: "FAIL", Duration: 10 * time.Millisecond},
			{Type: "output", Data: "\tparallel_test.go:32: ThirdTest error"},
			{Type: "end_test", Name: "FirstTest", Result: "FAIL", Duration: 2 * time.Second},
			{Type: "output", Data: "\tparallel_test.go:14: FirstTest error"},
			{Type: "end_test", Name: "SecondTest", Result: "FAIL", Duration: 1 * time.Second},
			{Type: "output", Data: "\tparallel_test.go:23: SecondTest error"},
			{Type: "status", Result: "FAIL"},
			{Type: "output", Data: "exit status 1"},
			{Type: "summary", Result: "FAIL", Name: "pkg/parallel", Duration: 3010 * time.Millisecond},
		}},
	{"21-cached",
		[]gtr.Event{
			{Type: "run_test", Name: "TestOne"},
			{Type: "end_test", Name: "TestOne", Result: "PASS"},
			{Type: "status", Result: "PASS"},
			{Type: "summary", Result: "ok", Name: "package/one", Data: "(cached)"},
		}},
	{"22-whitespace",
		[]gtr.Event{}},
}

func TestParse(t *testing.T) {
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			testParse(t, test.in, test.expected)
		})
	}
}

func testParse(t *testing.T, name string, expected []gtr.Event) {
	if len(expected) == 0 {
		t.SkipNow()
		return
	}
	f, err := os.Open(filepath.Join(testdataRoot, name+".txt"))
	if err != nil {
		t.Errorf("error reading %s: %v", name, err)
		return
	}
	defer f.Close()

	actual, err := Parse(f)
	if err != nil {
		t.Errorf("Parse(%s) error: %v", name, err)
		return
	}

	if diff := cmp.Diff(actual, expected); diff != "" {
		t.Errorf("Parse %s returned unexpected events, diff (-got, +want):\n%v", name, diff)
	}
}
