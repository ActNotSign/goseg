APP = segment
SEG = github.com/huichen/sego
CLI = github.com/codegangsta/cli
export GOPATH = ${PWD}
develop:
	go get ${SEG}
	go get ${CLI}
	go build -o ${APP} -ldflags '-s -w'
build:
	go get ${SEG}
	go get ${CLI}
	go build -o ${APP} -ldflags '-s -w'
	goupx -s=true -u ${APP}

run:
	@go run *.go

clean:
	@rm ${APP}
