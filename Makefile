run-with-flags:
	go run ./cmd/app --tick_interval=2s --message="Hello, Raymond!" --health_port=8081

run-with-config:
	go run ./cmd/app --config .config/preset-1.yaml

run-with-env:
	RX9PN_HEALTH_PORT=8081 go run ./cmd/app/main.go

run-with-config-override:
	# Flag > Env > Config file
	RX9PN_HEALTH_PORT=8083 go run ./cmd/app/main.go --config .config/preset-1.yaml --metrics_port 2111
