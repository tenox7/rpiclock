all:
	GOOS=linux GOARCH=arm GOARM=6 go build -a .
