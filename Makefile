APP_NAME = casaos-gateway
APP_NAME_SHORT = gateway

TARGET_DIR = target
BUILD_DIR = $(TARGET_DIR)/build

INSTALL_ROOT = /

all: $(TARGET_DIR)

$(BUILD_DIR): clean
	mkdir -pv $(BUILD_DIR)/usr/bin

$(APP_NAME): $(BUILD_DIR)
	go build -v -o $(BUILD_DIR)/usr/bin/$(APP_NAME)

$(TARGET_DIR): $(APP_NAME)
	cp -rv build $(TARGET_DIR)

clean:
	rm -rfv $(TARGET_DIR)

install:
	cp -rv $(BUILD_DIR)/* $(INSTALL_ROOT)
	if [ ! -f $(INSTALL_ROOT)/etc/casaos/$(APP_NAME_SHORT).ini ]; then \
		cp -v $(INSTALL_ROOT)/etc/casaos/$(APP_NAME_SHORT).ini.sample $(INSTALL_ROOT)/etc/casaos/$(APP_NAME_SHORT).ini; \
	fi
	systemctl enable --now $(APP_NAME).service

uninstall:
	systemctl disable --now $(APP_NAME).service
	rm -v $(INSTALL_ROOT)/etc/casaos/$(APP_NAME_SHORT).ini.sample
	rm -v $(INSTALL_ROOT)/usr/bin/$(APP_NAME)
	rm -v $(INSTALL_ROOT)/usr/lib/systemd/system/$(APP_NAME).service

