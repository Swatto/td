all:
	rm -fr build
	goxc -bc="linux,windows,darwin" -d=./build -pv=1.0.0
	rmdir debian

clean:
	rm -fr ./build
	rm -fr .DS_Store
	rm -fr td
