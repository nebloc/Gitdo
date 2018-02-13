install: 
	go build ./
	mv Gitdo .git/gitdo/

rm_commit:
	git reset --soft HEAD~
	git status
