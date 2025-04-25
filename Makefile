.PHONY: build

build:
	sam build

local:
	sam build
	sam local start-api

build-AttachIotPolicyFunction:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./bootstrap cmd/attach-iot-policy/main.go
	cp ./bootstrap $(ARTIFACTS_DIR)/.

build-GetDeviceSensorDataFunction:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./bootstrap cmd/get-device-sensor-data/main.go
	cp ./bootstrap $(ARTIFACTS_DIR)/.

build-ListUserDevicesFunction:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./bootstrap cmd/list-user-devices/main.go
	cp ./bootstrap $(ARTIFACTS_DIR)/.

build-AddDeviceFunction:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./bootstrap cmd/add-device/main.go
	cp ./bootstrap $(ARTIFACTS_DIR)/.
