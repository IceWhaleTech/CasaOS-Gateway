TARGET_DIR = target
BUILD_DIR = $(TARGET_DIR)/build

casaos-gateway: clean
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/casaos-gateway

clean:
	rm -rf $(TARGET_DIR)
