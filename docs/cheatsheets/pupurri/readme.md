---
id: pupurri
title: Pupurri
sidebar_label: Pupurri
sidebar_position: 1
---

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