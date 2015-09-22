This README is just a fast *quick start* document about how to try **bosh-init** function.

HOWTO
--------------
#1. Install Dependences
1.  install redis-server
2.  install mySQL-server
3.  install PostgreSQL Server

#2. Install Ruby runtime
1.  install rvm
2.  install Ruby 1.9.3
3.  install bundler gem
4.  install bosh-cli gem

#3.Install Go runtime
1. install Go1.3
2. setup GOPATH

#4.Build bosh-init project
1. clone bosh-init source code
2. build bosh-init
```bash
    cd bosh-init
    bin/build
```

#5.Build Bosh release
1. clone bosh source code
2. build bosh
```bash
    cd bosh
    bundle install
    cd bosh/bosh-dev
    bundle exec rake release:create_dev_release
    cd bosh/release/
    bosh create release --force --with-tarball
```

#6.Clone softlayer eCPI source code
```bash
    git clone bosh-softlayer-cpi-release 
    git submodule update --init —-recursive
```


#7. Build softlayer cpi release
```bash
    cd bosh;
    bundle install;
    cd bosh/bosh-dev
    bundle exec rake release:create_dev_release
    cd bosh/release/
    bosh create release --force --with-tarball
```


#8.Edit yaml file 

#9.Run bosh-init
```bash
   ./bosh deploy sl.yml
```
