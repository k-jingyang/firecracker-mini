## /etc/cni/conf.d/ptp.conflist

```json
{ 
  "name": "ptp-net",
  "cniVersion": "0.4.0",
  "plugins": [
    {
      "type": "ptp",
      "ipMasq": true,
      "ipam": { 
        "type": "host-local",
        "subnet": "172.16.0.0/24",
      }
    },
    {
      "type": "tc-redirect-tap"
    }
  ]
}
```

## Repro

```bash
$ firecracker --version
Firecracker v1.10.1

# From https://github.com/firecracker-microvm/firecracker/blob/main/docs/getting-started.md
$ latest=$(wget "http://spec.ccfc.min.s3.amazonaws.com/?prefix=firecracker-ci/v1.11/x86_64/vmlinux-5.10&list-type=2" -O - 2>/dev/null | grep -oP "(?<=<Key>)(firecracker-ci/v1.11/x86_64/vmlinux-5\.10\.[0-9]{1,3})(?=</Key>)")

# Download a linux kernel binary
$ wget "https://s3.amazonaws.com/spec.ccfc.min/${latest}"

$ wget -O ubuntu-24.04.squashfs.upstream "https://s3.amazonaws.com/spec.ccfc.min/firecracker-ci/v1.11/x86_64/ubuntu-24.04.squashfs"

$ go build . && sudo ./firecracker-mini 

....

Ubuntu 24.04.1 LTS ubuntu-fc-uvm ttyS0

ubuntu-fc-uvm login: root (automatic login)

Welcome to Ubuntu 24.04.1 LTS (GNU/Linux 5.10.225 x86_64)

 * Documentation:  https://help.ubuntu.com
 * Management:     https://landscape.canonical.com
 * Support:        https://ubuntu.com/pro

This system has been minimized by removing packages and content that are
not required on a system that users do not log into.

To restore this content, you can run the 'unminimize' command.

The programs included with the Ubuntu system are free software;
the exact distribution terms for each program are described in the
individual files in /usr/share/doc/*/copyright.

Ubuntu comes with ABSOLUTELY NO WARRANTY, to the extent permitted by
applicable law.


The programs included with the Ubuntu system are free software;
the exact distribution terms for each program are described in the
individual files in /usr/share/doc/*/copyright.

Ubuntu comes with ABSOLUTELY NO WARRANTY, to the extent permitted by
applicable law.

root@ubuntu-fc-uvm:~# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host 
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 42:16:13:e9:7a:7f brd ff:ff:ff:ff:ff:ff
    inet 172.16.0.7/24 brd 172.16.0.255 scope global eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::4016:13ff:fee9:7a7f/64 scope link 
       valid_lft forever preferred_lft forever
```

In another terminal, on your host
```bash
$ ping 172.16.0.7
PING 172.16.0.7 (172.16.0.7) 56(84) bytes of data.
64 bytes from 172.16.0.7: icmp_seq=1 ttl=127 time=0.273 ms
64 bytes from 172.16.0.7: icmp_seq=2 ttl=127 time=0.207 ms
64 bytes from 172.16.0.7: icmp_seq=3 ttl=127 time=0.254 ms
```
