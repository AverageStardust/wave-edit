production:
	go build -ldflags="-s -w"
	xz -1vf wave-edit