
This guide is just a fast *quick start* document about how to use **bosh-init** to deploy a new Bosh Director.

HOWTO
--------------
#Prerequisites：
1.  Create vm with Ubuntu 14.04 x86_64 LTS.
2.  Do not install any Ruby run-time under 2.0 version. If installed, pls uninstall firstly.


#1. Install Ruby 2.0+
```bash
    apt-get install -y software-properties-common wget curl git vim tmux expect 
    apt-add-repository -y ppa:brightbox/ruby-ng
    apt-get update
    apt-get install -y build-essential libxml2-dev libsqlite3-dev libxslt1-dev libpq-dev libmysqlclient-dev zlib1g-dev ruby2.2 ruby2.2-dev
```

#2.Install bosh_cli gem
```bash
    gem install bosh_cli --no-ri --no-rdoc 
```

#3.Install bosh-init 
```bash
    wget -O /usr/local/bin/bosh-init  http://10.113.109.244/stable/bosh-init-enable-os-reload
```

#4. Download deployment sample file
```bash
    mkdir -p /root/myDeploy/`date -d "today" +"%Y%m%d_%H"`/director
    cd /root/myDeploy/`date -d "today" +"%Y%m%d_%H"`/director
    wget http://10.113.109.244/stable/sl-bosh.yml
```

#5. Download Bosh tarball file
```bash
    cd /root/myDeploy/`date -d "today" +"%Y%m%d_%H"`/director
    wget http://10.113.109.244/stable/bosh-230+dev.2.tgz
```

#6. Download eCPI tarball file
```bash
    cd /root/myDeploy/`date -d "today" +"%Y%m%d_%H"`/director
    wget http://10.113.109.244/stable/bosh-softlayer-cpi-0+dev.1.tgz
```

#7. Download stemcell tarball file
```bash
    cd /root/myDeploy/`date -d "today" +"%Y%m%d_%H"`/director
    wget http://10.113.109.244/stable/light-bosh-stemcell-3031-softlayer-esxi-ubuntu-trusty-go_agent.tgz
```


#8.Edit deployment sample file
`sl-bosh.yml`
```yaml
---
name: bosh

releases:
- name: bosh
  url: file://./bosh-230+dev.2.tgz                                                      # <--- Replace with bosh tarball url
- name: bosh-softlayer-cpi
  url: file://./bosh-softlayer-cpi-0+dev.1.tgz                                          # <--- Replace with eCPI tarball url
  
resource_pools:
- name: vms
  network: default
  stemcell:
    url: file://./light-bosh-stemcell-3031-softlayer-esxi-ubuntu-trusty-go_agent.tgz    # <--- Replace with light stemcell tarball url
  cloud_properties:
    Domain: softlayer.com                                                               # <--- Replace with your your domian name
    VmNamePrefix: bosh-matt-osreload-test                                               # <--- Replace with your VmNamePrefix
    EphemeralDiskSize: 100
    StartCpus: 4
    MaxMemory: 8192
    Datacenter:
       Name: lon02
    HourlyBillingFlag: true
    PrimaryNetworkComponent:
       NetworkVlan:
          Id: 524956
    PrimaryBackendNetworkComponent:
       NetworkVlan:
          Id: 524954
    NetworkComponents:
    - MaxSpeed: 1000
disk_pools:
- name: disks
  disk_size: 100_000
  cloud_properties:
    consistent_performance_iscsi: true


networks:
- name: default
  type: dynamic
  dns: 
  - 8.8.8.8                                                                                           # <--- Google DNS 


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
  - {name: powerdns, release: bosh}
  - {name: cpi, release: bosh-softlayer-cpi}

  resource_pool: vms
  persistent_disk_pool: disks

  networks:
  - name: default

  properties:
    nats:
      user: nats
      password: nats
      auth_timeout: 3
      address: 127.0.0.1
      listen_address: 0.0.0.0
      port: 4222
      no_epoll: false
      no_kqueue: true
      ping_interval: 5
      ping_max_outstanding: 2
      http:
        port: 9222
    redis:
      address: 127.0.0.1
      password: redis
      port: 25255
      loglevel: info
    postgres: &20585760
      user: postgres
      password: postges
      host: 127.0.0.1
      database: bosh
      adapter: postgres
    blobstore:
      address: 127.0.0.1
      director:
        user: director
        password: director
      agent:
        user: agent
        password: agent
      port: 25250
      provider: dav
    director:
      cpi_job: cpi
      address: 127.0.0.1
      name: bosh
      db:
        adapter: postgres
        database: bosh
        host: 127.0.0.1
        password: postges
        user: postgres
    hm:
      http:
        user: hm
        password: hm
        port: 25923
      director_account:
        user: admin
        password: Cl0udyWeather
      intervals:
        log_stats: 300
        agent_timeout: 180
        rogue_agent_alert: 180
        prune_events: 30
        poll_director: 60
        poll_grace_period: 30
        analyze_agents: 60
      pagerduty_enabled: false
      resurrector_enabled: false
    dns:
      address: 127.0.0.1
      domain_name: microbosh
      db: *20585760
      webserver:
        port: 8081
        address: 0.0.0.0
    softlayer: &softlayer
      username: fake_username                                                                      # <--- Replace with username
      apiKey: fake_password                                                                        # <--- Replace with password

    cpi:
      agent: 
        mbus: nats://nats:nats@127.0.0.1:4222  ＃<-- double check user/password for nats                                          
        ntp: []
        vcappassword: $6$EWSrWZcD$TW7FcMConryu3h436UbSKNw3SWpSUpmx14lcO7r2aKt/UhITdVdrfnXbQOgnuT7.0OLpJGqMHUdjJlFUfBaUH0
        blobstore: 
          provider: dav
          options:
            endpoint: http://127.0.0.1:25250 
            user: agent
            password: agent
    ntp: &ntp []
 
cloud_provider:
  template: {name: cpi, release: bosh-softlayer-cpi}
  mbus: "https://admin:admin@bosh-matt-osreload-test.softlayer.com:6868"            # <--- Replace with your VmNamePrefix + domian name

  properties:
    softlayer: *softlayer
    cpi:
      agent:
        mbus: https://admin:admin@127.0.0.1:6868
        ntp: *ntp
        vcappassword: $6$EWSrWZcD$TW7FcMConryu3h436UbSKNw3SWpSUpmx14lcO7r2aKt/UhITdVdrfnXbQOgnuT7.0OLpJGqMHUdjJlFUfBaUH0
        blobstore:
          provider: local
          options:
            blobstore_path: /var/vcap/micro_bosh/data/cache
```

#9.Run bosh-init to deploy bosh director
```bash
   bosh-init deploy sl-bosh.yml
```

#10.Check bosh director status
If success, the bosh Director should be ready. 
You can use bosh cli to target the bosh director, and run some commands to check the status of the director.

```bash
   cat /etc/hosts | sed -n '2p'  | awk '{print $1}' | xargs -n1 -t bosh target
   bosh status
```
