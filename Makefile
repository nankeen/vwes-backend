all: docker

vwes-backend:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

docker: vwes-backend
	docker image build -t vwesbackend.azurecr.io/vwes .

deploy: docker
	docker push vwesbackend.azurecr.io/vwes

clean:
	rm vwes-backend
