package diffparse

import (
	"testing"
)

func TestGetGitDiff(t *testing.T) {
	diff, err := GetGitDiff()
	if err != nil {
		t.Errorf("GetGitDiff returned error: %v", err)
	}
	t.Logf("GetGitDiff returned: %s", diff)
}

func TestParseGitDiff(t *testing.T) {
	lines, err := ParseGitDiff(example_diff)
	if err != nil {
		t.Errorf("parse diff returned error: %v", err)
	}
	t.Logf("Output:\n%v", lines)
}

var example_diff string = `diff --git a/cli.go b/cli.go
index 30d4ebc..ca12716 100644
--- a/cli.go
+++ b/cli.go
@@ -1,8 +1,9 @@
 package main
 
 import (
+	"encoding/json"
 	"fmt"
-	"github.com/waigani/diffparser"
+	"github.com/appscode/diffparser"
 	"log"
 	"os"
 	"os/exec"
@@ -24,12 +25,15 @@ func main() {
 	}
 
 	// Save output as string
-	cmdOutput := fmt.Sprintf("\n%s", resp)
+	cmdOutput := fmt.Sprintf("%s", resp)
+
+	fmt.Println(cmdOutput + "\n\n\n\n")
 
 	// Parse diff output
 	diff, err := diffparser.Parse(cmdOutput)
 	if err != nil {
 		log.Fatalf("Error processing diff: %v", err)
+		os.Exit(1)
 	}
 
 	// Create waitgroup to sync handling of all files
@@ -50,31 +54,29 @@ func ProcessFileDiff(file *diffparser.DiffFile, wg *sync.WaitGroup) {
 
 	re := regexp.MustCompile((?:[[:space:]]|)//(?:[[:space:]]|)TODO:[[:space:]](.*))
 
-	output := fmt.Sprintf("%s\n", file.NewName)
+	stagedTasks := make([]Task, 0)
+
+	// TODO: Clean up this spaghetti code
 	for _, hunk := range file.Hunks {
 		for _, line := range hunk.NewRange.Lines {
-			if line.Mode == 0 {
+			if line.Mode == 0 { // if line was added
 				match := re.FindStringSubmatch(line.Content)
-				if len(match) > 0 {
+				if len(match) > 0 { // if match was found
 					t := Task{
 						file.NewName,
 						match[1],
-						line.Position,
+						line.Number,
 					}
-					output += t.ToString() + "\n"
+					stagedTasks = append(stagedTasks, t)
 				}
 			}
 		}
 	}
-	fmt.Println(output)
-}
 
-type Task struct {
-	FileName string
-	TaskName string
-	Position int
-}
+	b, err := json.Marshal(stagedTasks)
+	if err != nil {
+		log.Fatal(err)
+	}
 
-func (t *Task) ToString() string {
-	return fmt.Sprintf("File: %s, Task: %s, Pos: %d", t.FileName, t.TaskName, t.Position)
+	fmt.Printf("%s\n", b)
 }
diff --git a/config.go b/config.go
index ea10e7d..eda3313 100644
--- a/config.go
+++ b/config.go
@@ -5,5 +5,8 @@ package main
 // TODO: with tab then space
 
 func test() {
-	// TODO: Create test suite
+}
+
+func hello() {
+	// TODO: Cleanup
 }
diff --git a/diffparser/diffparser.go b/diffparser/diffparser.go
new file mode 100644
index 0000000..afed857
--- /dev/null
+++ b/diffparser/diffparser.go
@@ -0,0 +1,43 @@
+package diffparser
+
+import (
+	"strings"
+)
+
+// SplitToFiles takes the given diff and splits it in to an array of strings, each a different file section of the diff
+func SplitToFiles(diff string) []string {
+	diffLines := strings.Split(diff, "\n")
+	var files []string
+
+	currentFile := ""
+	for _, line := range diffLines {
+		if strings.HasPrefix(line, "diff") && currentFile != "" {
+			files = append(files, currentFile)
+			currentFile = ""
+		}
+		currentFile += line + "\n"
+	}
+	files = append(files, currentFile)
+
+	return files
+}
+
+// ParseDiff takes a given string diff and converts it in to useable structs
+func ParseDiff(diff string) {
+
+}
+
+type File struct {
+	OrigFile string
+	NewFile  string
+	Hunks    []hunk
+}
+type hunk struct {
+	origStartInd int
+	origRange    int
+	newStartInd  int
+	newRange     int
+	content      string
+}
+
+func SplitToChunk(file string) {}
diff --git a/diffparser/diffparser_test.go b/diffparser/diffparser_test.go
new file mode 100644
index 0000000..2fde0d4
--- /dev/null
+++ b/diffparser/diffparser_test.go
@@ -0,0 +1,52 @@
+package diffparser
+
+import "testing"
+
+var exampleDiff = diff --git a/Main.go b/Main.go
+index 984e788..09fba66 100644
+--- a/Main.go
++++ b/Main.go
+@@ -17,6 +17,7 @@ type PageInfo struct {
+	 User *user.User
+	 LoginOut string
+ }
+-//TODO: Fix init handler
+ 
+ func init() {
+	 http.HandleFunc("/", HomeHandler)
+@@ -43,6 +44,7 @@ func NewRequestHandler(w http.ResponseWriter, r *http.Request) {
+	 }
+ }
+ 
++//TODO: Talk to matt about handle here
+ /**
+ Handles the index page for header, footer and the likes
+  */
+diff --git a/Users.go b/Users.go
+index 34e4c3b..7a81965 100644
+--- a/Users.go
++++ b/Users.go
+@@ -4,7 +4,7 @@ import (
+	 "golang.org/x/net/context"
+	 "google.golang.org/appengine/user"
+ )
+-
+-//TODO: Testing
+ func userControl(ctx context.Context, currentPage string) (*user.User, string){
+ 
+	 user := user.Current(ctx)
+@@ -33,6 +33,7 @@ func getLoginURL(ctx context.Context, currentPage string) (string, error) {
+	 return loginURL, nil
+ }
+ 
++   //TODO: Get out
+ func getLogoutURL(ctx context.Context, currentPage string) (string, error) {
+	 //LogoutURL
+logoutURL, err := user.LogoutURL(ctx, currentPage)
+
+func TestSplitToFiles(t *testing.T) {
+	result := SplitToFiles(exampleDiff)
+	if len(result) != 2 {
+		t.Errorf("Expected 2 file sections, got %d", len(result))
+	}
+}
diff --git a/structs.go b/structs.go
new file mode 100644
index 0000000..9c1ba55
--- /dev/null
+++ b/structs.go
@@ -0,0 +1,7 @@
+package main
+
+type Task struct {
+	FileName string json: file_name
+	TaskName string json: task_name
+	FileLine int    json: FileLine
+}
`
