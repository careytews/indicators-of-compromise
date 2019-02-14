
TOPDIR=$(shell git rev-parse --show-toplevel)

all: init ioc.json

upload: all
	gsutil cp ./output/ioc.json gs://trust-networks-ioc/
	gsutil cp ./output/ioc.json gs://trust-networks-ioc/ioc-tn.json
	gsutil cp ./output/ioc.json gs://trust-networks-ioc/ioc-rh.json

ioc.json: ioc-sources output/tt.json output/shalla.json output/tors.json merge-ioc
	./merge-ioc ioc-sources/trust-networks/trust-networks.json \
	output/tt.json output/shalla-anonvpn.json output/shalla-dynamic.json \
	output/shalla-hacking.json output/shalla-redirector.json output/shalla-spyware.json \
	output/shalla-warez.json output/tors.json > output/ioc.json
	cp ./output/ioc.json .

target-source: target-source.go
	GOPATH=${TOPDIR} go build target-source.go common.go

output/tt.json: target-source
	./target-source

merge-ioc: merge-ioc.go
	GOPATH=${TOPDIR} go build merge-ioc.go

tor-source: tor-source.go
	GOPATH=${TOPDIR} go build tor-source.go common.go

output/tors.json: tors.html tor-source
	./tor-source

tors.html: init
	mkdir ioc-sources/tor
	wget https://www.dan.me.uk/tornodes -O ioc-sources/tor/tors.html

shalla-source: ioc-shallalist shalla-source.go
	GOPATH=${TOPDIR} go build shalla-source.go common.go

output/shalla.json: shalla-source
	./shalla-source
	touch output/shalla.json

bin/dep:
	GOPATH=${TOPDIR} go get -u github.com/golang/dep/cmd/dep

godeps: bin/dep
	GOPATH=${TOPDIR} ${TOPDIR}/bin/dep ensure -v

init: godeps
	git clone git@github.com:botherder/targetedthreats.git
	git clone git@github.com:trustnetworks/ioc-sources.git
	git clone git@github.com:trustnetworks/ioc-shallalist.git
	mkdir output

# rebuild: should rebuild the whole thing
