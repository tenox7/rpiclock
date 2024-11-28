all:
	GOOS=linux GOARCH=arm GOARM=6 go build -a .

clean:
	rm -f rpiclock
