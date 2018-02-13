package diffparser

import "testing"

var exampleDiff = `diff --git a/Main.go b/Main.go
index 984e788..09fba66 100644
--- a/Main.go
+++ b/Main.go
@@ -17,6 +17,7 @@ type PageInfo struct {
	 User *user.User
	 LoginOut string
 }
-//TODO: Fix init handler
 
 func init() {
	 http.HandleFunc("/", HomeHandler)
@@ -43,6 +44,7 @@ func NewRequestHandler(w http.ResponseWriter, r *http.Request) {
	 }
 }
 
+//TODO: Talk to matt about handle here
 /**
 Handles the index page for header, footer and the likes
  */
diff --git a/Users.go b/Users.go
index 34e4c3b..7a81965 100644
--- a/Users.go
+++ b/Users.go
@@ -4,7 +4,7 @@ import (
	 "golang.org/x/net/context"
	 "google.golang.org/appengine/user"
 )
-
-//TODO: Testing
 func userControl(ctx context.Context, currentPage string) (*user.User, string){
 
	 user := user.Current(ctx)
@@ -33,6 +33,7 @@ func getLoginURL(ctx context.Context, currentPage string) (string, error) {
	 return loginURL, nil
 }
 
+   //TODO: Get out
 func getLogoutURL(ctx context.Context, currentPage string) (string, error) {
	 //LogoutURL
logoutURL, err := user.LogoutURL(ctx, currentPage)`

func TestSplitToFiles(t *testing.T) {
	result := SplitToFiles(exampleDiff)
	if len(result) != 2 {
		t.Errorf("Expected 2 file sections, got %d", len(result))
	}
}
