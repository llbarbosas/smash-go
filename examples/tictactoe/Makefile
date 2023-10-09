modules: modules/simulator/simulator.go modules/console_renderer/console_renderer.go modules/input_manager/input_manager.go
	for module in "simulator" "console_renderer" "input_manager" ; do \
    	go build -buildmode=plugin -o bin/modules/$$module.so modules/$$module/$$module.go  ; \
	done
clean:
	rm -f bin/modules/*