all: docker

vwes-backend:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

docker: vwes-backend
	docker image build -t nankeen/vwes-backend .

deploy: docker
	docker push vwesbackend.azurecr.io/vwes
