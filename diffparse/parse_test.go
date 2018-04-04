package diffparse

import (
	"testing"
)

func TestParseGitDiff(t *testing.T) {
	ExpectedLines := []SourceLine{
		{"config", "config", "", 3, REMOVED},
		{"config", "config", "//TODO: hello", 5, ADDED},
		{"config.go", "", "package main", 0, REMOVED},
		{"config.go", "", "//TODO: load config from file test", 0, REMOVED},
		// TODO: Find out if the removed line number ever needs to be accurate
		{"", "diffparser/diffparser.go", "++ b/package diffparser", 1, ADDED},
		{"commit_test.go", "commit_test.go", "+// TODO: Test <9ypvkCD1>", 242, REMOVED},
		{"commit_test.go", "commit_test.go", "+// TODO: Test", 242, ADDED},
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

	lines, err = ParseGitDiff(fileLineExample)
	if err != nil {
		t.Errorf("parse diff returned error: %v", err)
	}
	for _, line := range lines {
		if line.Content == "Quisque mauris in orci cursus lobortis. Sed sed faucibus tellus." {
			if line.Position != 30 {
				t.Errorf("Expected line 30, got %d", line.Position)
			}
		}
	}
}

var example_diff string = `diff --git a/config b/config
index d90eea3..849126d 100644
--- a/config
+++ b/config
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
+++ b/package diffparser

diff --git a/commit_test.go b/commit_test.go
index 0c208c7..8b81f2b 100644
--- a/commit_test.go
+++ b/commit_test.go
@@ -220,7 +239,7 @@ index 0000000..a30278c
 +
 +import "fmt"
 +
-+// TODO: Test <9ypvkCD1>
++// TODO: Test
 +func main(){
 +	fmt.Println("Hello Ben")
 +}
`

const fileLineExample string = `diff --git a/lorem.txt b/lorem.txt
index a3653c4..6f9f76d 100644
--- a/lorem.txt
+++ b/lorem.txt
@@ -1,19 +1,17 @@
 Lorem ipsum dolor sit amet, consectetur adipiscing elit. Integer vel orci tincidunt,
 sollicitudin odio vel, mollis augue. Praesent rhoncus finibus maximus. Donec pulvinar
-sodales mollis. Donec pharetra et nibh sed suscipit. Cras quis maximus risus. Duis vitae
-metus justo. Quisque vel dignissim felis. Aenean aliquam ipsum non auctor tempor.
 Phasellus condimentum velit consequat ipsum imperdiet, ac tristique lectus varius.
 Mauris vitae nunc a lectus eleifend fringilla mollis id ante. Morbi pharetra orci
 id euismod rhoncus.

 Phasellus turpis lorem, venenatis a dui non, bibendum varius lectus. Pellentesque
-scelerisque libero vestibulum mollis imperdiet. Proin nec justo quis mauris ullamcorper
 finibus a quis nisl. Proin lacinia aliquet ligula non viverra. Morbi a sodales diam.
 Vivamus et sapien a metus commodo ornare vel et metus. Aenean eget odio varius, viverra
 ipsum ut, euismod risus. Nullam in dui aliquam, tempus justo eu, facilisis ante. Proin
 malesuada semper dui nec ultricies. Pellentesque habitant morbi tristique senectus et
 netus et malesuada fames ac turpis egestas. Vestibulum ante ipsum primis in faucibus orci
-luctus et ultrices posuere cubilia Curae; Integer nibh quam, finibus in augue id,
+Etiam eget nulla vel mauris aliquam pulvinar in ut urna. Vestibulum ante ipsum primis in
+faucibus orci luctus et ultrices posuere cubilia Curae; Aliquam erat volutpat.
 fringilla fermentum urna. Sed interdum, enim sed porta faucibus, metus metus hendrerit
 magna, vitae luctus arcu neque at magna. Aliquam non est velit. Sed suscipit libero in
 purus sodales, quis maximus velit tincidunt.
@@ -23,15 +21,14 @@ ultricies, lacus ligula gravida sapien, congue interdum ipsum metus a ante. Null
 tellus augue. Duis convallis eget elit vitae sagittis. Aliquam sem leo, tempus non aliquet
 ut, placerat id enim. Pellentesque consequat nisl justo. Ut iaculis enim eget ante dictum,
 a laoreet urna sollicitudin. Donec sit amet ligula id ligula venenatis tincidunt sit amet
-et ipsum.
+et ipsum

 Vestibulum a erat vitae massa vestibulum efficitur. Maecenas a diam cursus felis feugiat
-laoreet imperdiet sit amet erat. Morbi et purus pulvinar, imperdiet sapien vitae,
 efficitur erat. Phasellus eget nisi ligula. Nulla purus lacus, posuere ac cursus vitae,
 maximus lacinia sem. Nunc cursus lectus et velit tempor, et elementum sapien sodales.
 Cras ultricies facilisis ipsum a placerat.
+Quisque mauris in orci cursus lobortis. Sed sed faucibus tellus.

 Sed porttitor purus id posuere dictum. Ut vitae leo ac felis mollis accumsan ac at eros.
 Donec at nisl a lorem mollis finibus. Fusce vel gravida nulla. Phasellus pharetra
 consectetur nulla vel porta. Pellentesque fringilla turpis eu iaculis iaculis. Proin
-consequat neque mi, quis congue purus luctus non.
\ No newline at end of file`
