##
## WARNING: This makefile does not have dependencies on libraries used by the 
## WARNING: go packages in pages such as seven5.  If you change those libraries
## WARNING: you have to do "make clean; make" to rebuild and include the new
## WARNING: library versions.
##
## This file must be run from within the client directory.
##

all: ../static/en/web/index.js \
	../static/en/web/index.html

clean:
	rm -rf \
	../static/en/web/index.js* \
	../static/en/web/index.html

### extra go files throw into to the dependencies, it compiles too fast to bother
### trying to make this exact
### WARNING: for deployment reasons, we don't list the "true" dependency on 
### WARNING: on seven5/seven5/client here. That will break build on heroku.
extras=../shared/*.go ../wire/*.go 


GOPHERJS:=gopherjs
GOPHERJSOPTS:=-m

../static/en/web/index.js: pages/index.go $(extras)
	$(GOPHERJS) build $(GOPHERJSOPTS) --output=../static/en/web/index.js pages/index.go


##
## HTML PAGES
## 
## These are done this way largely because it's hard to uniquely identify a page
## from only the basename.  
##
../static/en/web/index.html: template/index.html template/index.json template/support/*
	pagegen --support=support --dir=template --start=index.html --json=index.json > ../static/en/web/index.html
