##
## WARNING: This makefile does not have dependencies on libraries used by the 
## WARNING: go packages in pages such as seven5.  If you change those libraries
## WARNING: you have to do "make clean; make" to rebuild and include the new
## WARNING: library versions.
##
## This file must be run from within the client directory.
##
## This file may have other unknown dependency flaws. YMMV.

all: \
	../static/en/web/index.js \
	../static/en/web/index.html \
	../static/en/web/login.js \
	../static/en/web/login.html \
	../static/en/web/newmovie.js \
	../static/en/web/newmovie.html

clean:
	rm -rf \
	../../../../../pkg/darwin_js_min/ \
	../../../../../pkg/linux_js_min/ \
	../static/en/web/index.js* \
	../static/en/web/index.html \
	../static/en/web/login.js* \
	../static/en/web/login.html \
	../static/en/web/newmovie.js* \
	../static/en/web/newmovie.html

### extra go files throw into to the dependencies, it compiles too fast to bother
### trying to make this exact
### WARNING: for deployment reasons, we don't list the "true" dependency on 
### WARNING: on seven5/seven5/client here. That will break build on heroku.
extras=../shared/*.go ../wire/*.go uilib/*.go


GOPHERJS:=gopherjs
GOPHERJSOPTS:=-m

#
# There, sadly, is no godep <cmd> so we have jury-rig this via the godep "path" command. Sigh.
#
WITHDEPS:=$(shell cd ..; godep path)

../static/en/web/index.js: pages/index.go $(extras)
	GOPATH=$(GOPATH):$(WITHDEPS) $(GOPHERJS) build $(GOPHERJSOPTS) --output=../static/en/web/index.js pages/index.go

../static/en/web/login.js: pages/login.go $(extras)
	GOPATH=$(GOPATH):$(WITHDEPS) $(GOPHERJS) build $(GOPHERJSOPTS) --output=../static/en/web/login.js pages/login.go

../static/en/web/newmovie.js: pages/newmovie.go $(extras)
	GOPATH=$(GOPATH):$(WITHDEPS) $(GOPHERJS) build $(GOPHERJSOPTS) --output=../static/en/web/newmovie.js pages/newmovie.go


##
## HTML PAGES
## 
## These are done this way largely because it's hard to uniquely identify a page
## from only the basename.  
##
../static/en/web/index.html: template/index.html template/index.json template/support/*
	pagegen --support=support --dir=template --start=index.html --json=index.json > ../static/en/web/index.html

../static/en/web/login.html: template/login.html template/login.json template/support/*
	pagegen --support=support --dir=template --start=login.html --json=login.json > ../static/en/web/login.html

../static/en/web/newmovie.html: template/newmovie.html template/newmovie.json template/support/*
	pagegen --support=support --dir=template --start=newmovie.html --json=newmovie.json > ../static/en/web/newmovie.html
