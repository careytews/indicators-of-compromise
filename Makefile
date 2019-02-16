
all: ioc.json

ioc.json: ioc-sources output/tt.json output/shalla.json output/tors.json merge-ioc
	./merge-ioc ioc-sources/trust-networks/trust-networks.json \
	output/tt.json output/shalla-anonvpn.json output/shalla-dynamic.json \
	output/shalla-hacking.json output/shalla-redirector.json output/shalla-spyware.json \
	output/shalla-warez.json output/tors.json > output/ioc.json
	cp ./output/ioc.json .

target-source: target-source.go
	go build target-source.go common.go types.go

output/tt.json: target-source
	./target-source

merge-ioc: merge-ioc.go
	go build merge-ioc.go types.go

tor-source: tor-source.go
	go build tor-source.go common.go types.go

output/tors.json: tor-source
	./tor-source

shalla-source: ioc-shallalist shalla-source.go
	go build shalla-source.go common.go types.go

output/shalla.json: shalla-source
	./shalla-source
	touch output/shalla.json

init: 
	git clone git@github.com:botherder/targetedthreats.git
	git clone git@github.com:tnw-open-source/ioc-sources.git
	git clone git@github.com:tnw-open-source/ioc-shallalist.git
	mkdir output
	mkdir ioc-sources/tor
	wget https://www.dan.me.uk/tornodes -O ioc-sources/tor/tors.html

clean:
	rm -f target-source && \
	rm -f shalla-source && \
	rm -f tor-source && \
	rm -f merge-ioc && \
	rm -f translate && \
	rm -f debug && \
	rm -f ioc.json
	
clean-dirs: 
	rm -rf ioc-sources && \
	rm -rf output && \
	rm -rf targetedthreats && \
	rm -rf ioc-shallalist

clean-all: clean clean-dirs
