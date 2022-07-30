APP_NAME = casaos-gateway
APP_NAME_SHORT = gateway

TARGET_DIR = target
BUILD_DIR = $(TARGET_DIR)/build

INSTALL_ROOT = /

$(APP_NAME): clean
	mkdir -pv $(BUILD_DIR)/usr/bin
	go build -v -o $(BUILD_DIR)/usr/bin/$(APP_NAME)
	cp -rv build $(TARGET_DIR)

clean:
	rm -rfv $(TARGET_DIR)

install:
	cp -rv $(BUILD_DIR)/* $(INSTALL_ROOT)
	systemctl enable --now $(APP_NAME).service

uninstall:
	systemctl disable --now $(APP_NAME).service
	rm -v $(INSTALL_ROOT)/etc/casaos/$(APP_NAME_SHORT).ini
	rm -v $(INSTALL_ROOT)/usr/bin/$(APP_NAME)
	rm -v $(INSTALL_ROOT)/usr/lib/systemd/system/$(APP_NAME).service

