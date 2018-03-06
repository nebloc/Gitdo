install:
	rm -rf ./hooks/
	mkdir ./hooks
	go build -o Gitdo ./app
	mv Gitdo ./hooks/
	cp pre-commit ./hooks/
	cp ./app/config.json ./hooks/
	cp ./plugins/gitdo_trello.py ./hooks/

rm_commit:
	git reset --soft HEAD~
	git status
