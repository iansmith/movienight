Dev Setup
---------


### Code setup
* Assume you want to develop in ~/movienight.src
** `$ mkdir ~/movienight.src`
* Copy enable-movienight.sample from github to ~/movienight.src/enable-movienight
** Adjust enable script as needed, but the GOPATH and DATABASE_URL are the critical bits
** DATABASE_URL is at a reasonable default for local dev
* `$ source enable-movienight`
* Make the standard go directory structure
** `$ mkdir -p ~/movienight/{bin,pkg,src/github.com/iansmith/}`
* Create the movienight directory as a "standard" go source package
** `$ pushd ~/movienight/src/github.com/iansmith; git clone git@github.com:iansmith/movienight.git; popd`
* Note that the .git is in the dir `src/github.com/iansmith/movienight`

### Database setup
* Get postgres 9.X running, but bind to port 5433 (not 5432 b/c it would conflict with mesa)
** on OS with homebrew: `$  postgres --port=5433 -D /usr/local/var/postgres` 
* Make the DB
** `$ createdb -p 5433 movienight`
* Install the command line deps, assuming you have your GOPATH set (see above)
** `$ go get github.com/tools/godep`
* Install gopherjs compiler
** `$ go get github.com/gopherjs/gopherjs`
