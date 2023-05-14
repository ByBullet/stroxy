all:build-gui

outputDir="build/"

build-gui:
	if [ -d $(outputDir)]; then
		rm -rf $(outputDir)
	fi
	mkdir -p build/resources build/script
	cp -r resources/ build/resources
	cp -r script/ build/script
	go build -o build/stroxy app/gui/main.go