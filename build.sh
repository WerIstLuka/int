#!/bin/bash

pkgversion=$(cat int.go | grep "^var Version string" | grep -Po '(?<=").+(?=")')

command -v go >/dev/null 2>&1 || { echo >&2 "Go is not installed, exiting"; exit 1; }


echo "building version" $pkgversion

go build

#check if build failed
if [ "$?" != "0" ]; then
	echo "error building, exiting"
	exit 1
fi


failed_tests=0
echo "----------testing----------"
test_int () {
	if [ "$(./int $1)" == "$2" ]; then
		echo "passed"
	else
		echo "----------failed-----------"
		echo "input: $1"
		echo "expected output: $2"
		echo "actual output: $(./int $1)"
		let failed_tests+=1
	fi
}

test_int_pipe () {
	if [ $(echo $1 | ./int $2) == "$3" ]; then
		echo "passed"
	else
		echo "----------failed-----------"
		echo "input: $1 | $2"
		echo "expected output: $3"
		echo "actual output: $(echo $1 | ./int $2)"
		let failed_tests+=1
	fi
}

test_int -v $pkgversion
test_int 0b101 5
test_int 0x20 32
test_int "a -x" 10
test_int "2 -b" "Error: Number 2 is not valid for base 2"
test_int "10 -B7" 7
test_int "10 -Ob" 1010
test_int "-O7 10" 13
test_int "-B63 10" "Error: invalid base: 63"
test_int "-x A" 10
test_int "a" "Error: Number a is not valid for base 10"
test_int "0xA -l" "Error: Number A is not valid for base 16"
test_int "12 -Ox --long" "c"
test_int "1.5" "Error: Illegal Character: ."
test_int_pipe 0b101 ""  5
test_int_pipe "-x 10" "" 16
test_int_pipe "-x -Ob" "20" 100000


if [ "$failed_tests" != "0" ]; then
	echo "$failed_tests tests failed, exiting"
	exit 1
fi

if [ "$1" == "-dev" ]; then
	exit 0
fi


#build releases
if [ -d "build" ]; then
	rm -rf "build"
fi

mkdir "build"

echo "----------building---------"
build_release () {
	echo "building for $1 $2"
	env GOOS=$1 GOARCH=$2 go build -o build/$1-$2
}
build_release linux amd64
build_release linux arm64
build_release linux riscv64
build_release freebsd amd64
build_release freebsd arm64
build_release openbsd amd64
build_release openbsd arm64
build_release windows amd64
build_release windows arm64
build_release darwin amd64
build_release darwin arm64


command -v dpkg-deb >/dev/null 2>&1 || { echo >&2 "dpkg-deb is not installed, skipping debian packages"; exit 0; }

echo "------------dpkg-----------"

echo "building debian packages"
pushd build >> /dev/null

make_debian_package () {
	installsize=$(du -ks "linux-$1"|cut -f 1)
	mkdir -p int_$1/{DEBIAN,usr/bin}
	cp ../control "int_$1/DEBIAN"
	sed -i -e "s/pkgver/${pkgversion}/g" "int_$1/DEBIAN/control"
	sed -i -e "s/pkgarch/$1/g" "int_$1/DEBIAN/control"
	sed -i -e "s/pkgsize/$installsize/g" "int_$1/DEBIAN/control"
	cp linux-$1 "int_$1/usr/bin/int"
	chmod -w "int_$1/usr/bin/int"
	chmod +x "int_$1/usr/bin/int"
	dpkg-deb --root-owner-group --build "int_$1"
	rm -rf "int_$1"
}

make_debian_package amd64
make_debian_package arm64
make_debian_package riscv64

popd >> /dev/null
