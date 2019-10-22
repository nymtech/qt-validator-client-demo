build_gui:
	qtdeploy build desktop qt-demo

build_release_gui:
	make build_gui
	mkdir -p build
	# move client configs to include them in the release zip
	cp -a internal/demo_configs/. qt-demo/deploy/linux/
	(cd qt-demo/deploy/linux; zip -r ../../../build/linux_amd64.zip *)

run_gui:
	qt-demo/deploy/linux/qt-demo
