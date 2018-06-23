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

You can simply run it by: `sudo ./dnsgo`. You can define the address to bind, like `sudo ./dnsgo -host 0.0.0.0`.

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
     
- nameserver.conf

    It's the file of nameserver. Config it by `nameserver 8.8.8.8`.
    
- resolver.conf

    A file like `/etc/resolv.conf`, don't know it much. But it's about dns forward. The format like `server=/*.abc.local/192.168.10.1`
    
### Example

Basic use: `./dnsgo -host 0.0.0.0`, you can test it by `dig @127.0.0.1 xxxx.xx`

To use it on the internet, first, you have to add a nameserver record for it. Config it at the domain name resolution page.

![set_ns](http://view.grt1st.cn/img/dnsgo1.png)

My config is not effective, so to be continued.

## To-do 

1. Dynamic update hosts
2. More stability test
3. Code rebuild
4. Fix ns lookup read timeout bug

## Thanks

[miekg/dns](https://github.com/miekg/dns)

[kenshinx/godns](https://github.com/kenshinx/godns)