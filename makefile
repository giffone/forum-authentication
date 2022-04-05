run :
	go run cmd/forumsqlite/main.go

docker :
	docker image build -f dockerfile . -t forumsqlite
	docker container run -p 9000:3306 -d --name forum forumsqlite

clean :
	docker system prune -a