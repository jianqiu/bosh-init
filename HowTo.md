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
`Example.yml`
```yaml
---
name: bosh

releases:
- name: bosh
  url: file://./bosh-206+dev.2.tgz
- name: bosh-softlayer-cpi
  url: file://./bosh-softlayer-cpi-0+dev.1.tgz

resource_pools:
- name: vms
  network: default
  stemcell:
    url: file://./light-bosh-stemcell-3030-softlayer-esxi-ubuntu-trusty-go_agent.tgz
  cloud_properties:
    Domain: softlayer.com
    VmNamePrefix: bosh-wjq
    StartCpus: 1
    MaxMemory: 1024
    Datacenter:
       Name: lon02
    HourlyBillingFlag: true
disk_pools:
- name: disks
  disk_size: 40_000
  cloud_properties:
    consistent_performance_iscsi: true


networks:
- name: default
  type: dynamic
  dns: [8.8.8.8] # <--- Replace with your DNS


jobs:
- name: bosh
  instances: 1

  templates:
  - {name: nats, release: bosh}
  - {name: redis, release: bosh}
  - {name: postgres, release: bosh}
  - {name: blobstore, release: bosh}
  - {name: director, release: bosh}
  - {name: health_monitor, release: bosh}
  - {name: cpi, release: bosh-softlayer-cpi}

  resource_pool: vms
  persistent_disk_pool: disks

  networks:
  - name: default

  properties:
    nats:
      address: 127.0.0.1
      user: nats
      password: nats-password

    redis:
      listen_addresss: 127.0.0.1
      address: 127.0.0.1
      password: redis-password

    postgres: &db
      host: 127.0.0.1
      user: postgres
      password: postgres-password
      database: bosh
      adapter: postgres

    blobstore:
      address: 127.0.0.1   # <--- Replace with a private IP
      port: 25250
      provider: dav
      director: {user: director, password: director-password}
      agent: {user: agent, password: agent-password}

    director:
      address: 127.0.0.1
      name: my-bosh
      db: *db
      cpi_job: cpi
      max_threads: 3

    hm:
      director_account: {user: admin, password: admin}
      resurrector_enabled: true

    softlayer: &softlayer
      username: fake_user_name	# <--- Replace with username
      apiKey: fake_api_key 		# <--- Replace with password
      public_vlan_id: fake_public_vlan_id
      private_vlan_id: fake_private_vlan_id
      data_center: lon02

    cpi:
      agent: {mbus: "nats://nats:nats-password@127.0.0.1:4222"} # <--- Replace with a private IP

    ntp: &ntp []

cloud_provider:
  template: {name: cpi, release: bosh-softlayer-cpi}
  mbus: "https://admin:admin@bosh-wjq.softlayer.com:6868" # <--- Replace with a floating IP

  properties:
    softlayer: *softlayer
    cpi:
      agent:
        mbus: https://admin:admin@127.0.0.1:6868
        ntp: *ntp
        blobstore:
          provider: local
          options:
            blobstore_path: /var/vcap/micro_bosh/data/cache
```

#9.Run bosh-init
```bash
   ./bosh-init deploy sl.yml
```

#10.Make bosh cli ready
In step#9, Bosh Director should be created. So, you can use bosh cli to interactive with Bosh Director.
*Softlayer is not supported officailly in current phrase, so you MUST not install bosh cli from the CF community.*
*Instead, you MUST get the bosh cli vm from the Softlayer image template, or use the pre-build tarball.*

```bash
#install Softlayer Python client 
pip install slcli
#CID can be get from the json file created by Bosh-init
slcli vs list | grep <CID> //Get Bosh director's IP created by Bosh-init
slcli vs details <CID>      
bosh target <BOSH_DIRECTOR_IP>
bosh status
```
#11.Upload releases to bosh director
1. Copy the releases from Repo to your bosh cli (somewhere like ~/releases)
2. bosh upload release <~/releases/cf-release_cf-xxx.tgz>


#12.Upload the fake stemcell
1. Download & transfer to your bosh cli
2. bosh upload stemcell STEMCELL_FILE

#13.Deploy the env
1. bosh deployment xxx.yml
2. bosh deploy
