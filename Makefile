TARGET_DIR = target
BUILD_DIR = $(TARGET_DIR)/build

INSTALL_ROOT = /

casaos-gateway: clean
	mkdir -pv $(BUILD_DIR)/usr/bin
	go build -v -o $(BUILD_DIR)/usr/bin/casaos-gateway
	cp -rv build $(TARGET_DIR)

clean:
	rm -rfv $(TARGET_DIR)

install:
	cp -rv $(BUILD_DIR)/* $(INSTALL_ROOT)
	systemctl enable --now casaos-gateway.service

uninstall:
	systemctl disable --now casaos-gateway.service
	rm -v $(INSTALL_ROOT)/etc/casaos/gateway.ini
	rm -v $(INSTALL_ROOT)/usr/bin/casaos-gateway
	rm -v $(INSTALL_ROOT)/usr/lib/systemd/system/casaos-gateway.service
