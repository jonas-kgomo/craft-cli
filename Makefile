BINARY_NAME=craft
INSTALL_PATH=/usr/local/bin

.PHONY: build install clean

build:
	go build -o $(BINARY_NAME) main.go

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "✅ Installed! You can now run 'craft' from anywhere."
	@echo "⚠️  Note: Make sure your GROQ_API_KEY is export in your .zshrc or .bash_profile"

clean:
	rm -f $(BINARY_NAME)
