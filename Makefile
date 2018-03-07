install:
	rm -rf ./hooks/
	mkdir ./hooks
	go build -o Gitdo ./app
	mv Gitdo ./hooks/
	cp pre-commit ./hooks/
	cp ./app/config.json ./hooks/
	cp -R ./plugins/trello_js ./hooks/

rm_commit:
	git reset --soft HEAD~
	git status

test_trello:
	node plugins/trello_js/trello.js "[{\"FileName\":\"app/cli.go\",\"TaskName\":\"Better handling of plugins\",\"FileLine\":103}]"
