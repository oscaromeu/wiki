---
id: pupurri
title: Pupurri
sidebar_label: Pupurri
sidebar_position: 1
---

## TCPDump

### Traffic Isolation Examples

Let’s start with a basic command that will get us HTTPS traffic:

```
tcpdump -nnSX port 443
```

```bash
04:45:40.573686 IP 78.149.209.110.27782 > 172.30.0.144.443: Flags [.], ack 
278239097, win 28, options [nop,nop,TS val 939752277 ecr 1208058112], length 0
    0x0000:  4500 0034 0014 0000 2e06 c005 4e8e d16e  E..4........N..n
    0x0010:  ac1e 0090 6c86 01bb 8e0a b73e 1095 9779  ....l......>...y
    0x0020:  8010 001c d20
```


This showed some HTTPS traffic, with a hex display visible on the right portion of the output (alas, it’s encrypted). Just remember—when in doubt, run the command above with the port you’re interested in, and you should be on your way.


#### Everything on an interface

Just see what’s going on, by looking at what’s hitting your interface.

```bash
tcpdump -i eth0
```

Or get all interfaces with `-i any`.

#### Find traffic by IP

One of the most common queries, using host, you can see traffic that’s going to or from `1.1.1.1`.



```bash
tcpdump host 1.1.1.1
```

```bash
06:20:25.593207 IP 172.30.0.144.39270 > one.one.one.one.domain: 
12790+ A? google.com. 
(28) 06:20:25.594510 IP one.one.one.one.domain > 172.30.0.144.39270: 
12790 1/0/0 A 172.217.15.78 (44)
```

#### Filtering by Source and/or Destination

If you only want to see traffic in one direction or the other, you can use `src` and `dst`

```
tcpdump src 1.1.1.1
tcpdump dst 1.0.0.1
```


## Bluetooth

### Initializing and stopping the bluetooth

First we should check if the bluetooth service is enabled with

```bash
systemctl status bluetooth.service
```

the output of the previous command has to be something like 

```
● bluetooth.service - Bluetooth service
     Loaded: loaded (/usr/lib/systemd/system/bluetooth.service; enabled; preset: disabled)
     Active: active (running) since Fri 2022-09-02 20:37:36 CEST; 1 day 21h ago
       Docs: man:bluetoothd(8)
   Main PID: 878 (bluetoothd)
     Status: "Running"
      Tasks: 1 (limit: 19041)
     Memory: 2.8M
        CPU: 704ms
     CGroup: /system.slice/bluetooth.service
```

Once we know that the bluetooth service is running we can check how many bluetooth devices are present on our computer

```bash
$ bluetoothctl list
```

the command has the following output

```
$ bluetoothctl list
Controller 34:E1:2D:4F:D8:D6 doraemon [default]
``` 

we can execute `bluetoothctl show <ctrl>` to get more details, where `<ctrl>` is the bluetooth controller device.  

In order to initialize the bluetooth execute `bluetooth power on` and to stop it `bluetooth power off`. 

### Managing devices

The following step is searching and pairing devices. 

1. `bluetooth select <ctrl>`
1. `bluetooth pairable on`
1. `bluetooth discoverable on`


After our bluetooth device is discoverable and pairable we need to scan for bluetooth devices

```bash
bluetoothctl scan on
```

we execute `bluetoothctl pair <device-id>` when the desired device is found. Finally we can bring the device back to hidden state with

```bash
bluetoothctl connect <device>
```

we can automate the connection with the following bash one liner

```bash
bluetoothctl connect  $(bluetoothctl devices Paired | grep Device | cut -d ' ' -f 2)
```

and similarly to disconnect the device we can use

```bash
bluetoothctl disconnect  $(bluetoothctl devices Paired | grep Device | cut -d ' ' -f 2)
```


## SSH

### Proxyjump

```
ssh -J internal-proxy last-host -f -N

```
## Enlaces

[1] - [Tuneles proxyjump](https://sysarmy.com/blog/posts/proxyjump-tuneles-ssh/)
