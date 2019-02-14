all rebuild init upload godeps:
	cd src/ioc && make $@

clean:
	rm -rf bin && \
	rm -rf pkg && \
	rm -rf src/github.com && \
	rm -rf src/ioc/vendor && \
	rm -f src/ioc/Gopkg.lock \
	rm -f src/ioc/target-source && \
	rm -f src/ioc/shalla-source && \
	rm -f src/ioc/tor-source && \
	rm -f src/ioc/merge-ioc && \
	rm -f src/ioc/translate && \
	rm -f src/ioc/debug && \
	rm -rf src/ioc/ioc-sources && \
	rm -rf src/ioc/output && \
	rm -rf src/ioc/targetedthreats && \
	rm -rf src/ioc/ioc-shallalist && \
	rm -f src/ioc/ioc.json
	