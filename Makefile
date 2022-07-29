TARGET_DIR = target
BUILD_DIR = $(TARGET_DIR)/build
BIN_DIR = usr/bin

INSTALL_ROOT = /

casaos-gateway: clean
	mkdir -pv $(BUILD_DIR)/$(BIN_DIR)
	cp -rv build $(TARGET_DIR)
	go build -v -o $(BUILD_DIR)/$(BIN_DIR)/casaos-gateway

clean:
	rm -rfv $(TARGET_DIR)

install:
	cp -rv $(BUILD_DIR)/* $(INSTALL_ROOT)

uninstall:
	rm -v /etc/casaos/gateway.ini
	rm -v /usr/bin/casaos-gateway
	rm -v /usr/lib/systemd/system/casaos-gateway.service
