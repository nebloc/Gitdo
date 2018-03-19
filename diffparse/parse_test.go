package diffparse

import (
	"testing"
)

func TestParseGitDiff(t *testing.T) {
	ExpectedLines := []SourceLine{
		{"test.txt", "test.txt", "", 3, REMOVED},
		{"test.txt", "test.txt", "//TODO: hello", 5, ADDED},
		{"config.go", "", "package main", 0, REMOVED},
		{"config.go", "", "//TODO: load config from file test", 0, REMOVED}, // TODO: Do we need to know the file line of removed? currently the counter is off
		{"", "diffparser/diffparser.go", "+package diffparser", 1, ADDED},
	}
	lines, err := ParseGitDiff(example_diff)
	if err != nil {
		t.Errorf("parse diff returned error: %v", err)
	}
	for i, expLine := range ExpectedLines {
		if expLine != lines[i] {
			t.Errorf("Line %d failed: expected %v, got %v", i, expLine, lines[i])
		}
	}
}

var example_diff string = `diff --git a/test.txt b/test.txt
index d90eea3..849126d 100644
--- a/test.txt
+++ b/test.txt
@@ -1,5 +1,5 @@
 Hello Ben
 How Are you
-
 Testing git pos
 To see if it works
+//TODO: hello
diff --git a/config.go b/config.go
deleted file mode 100644
index aa891e2..0000000
--- a/config.go
+++ /dev/null
@@ -1,3 +0,0 @@
-package main
-//TODO: load config from file test
diff --git a/diffparser/diffparser.go b/diffparser/diffparser.go
new file mode 100644
index 0000000..afed857
--- /dev/null
+++ b/diffparser/diffparser.go
@@ -0,0 +1,43 @@
++package diffparser`
