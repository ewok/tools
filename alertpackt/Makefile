all: alertpackt alertpackt_l32 alertpackt_l64

alertpackt: main.go
	go build

alertpackt_l64: main.go
	GOOS=linux GOARCH=amd64 go build -o alertpackt_l64

alertpackt_l32: main.go
	GOOS=linux GOARCH=386 go build -o alertpackt_l32
