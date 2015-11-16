.PHONY: rpm deb
rpm deb:
	fpm -f -s dir -t $@ -n drydock -v 0.0.4 \
		--architecture native --description "Docker Image Cleaner" \
		--maintainer "GoCardless Engineering <engineering@gocardless.com>" \
		drydock=/usr/local/bin/drydock
