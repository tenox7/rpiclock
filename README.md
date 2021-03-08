# Raspberry PI Zero WiFi-NTP/RTC Desktop Clock

Simple WiFi NTP desktop clock with a large 7-segment display. Built on Raspberry PI Zero.

## Hardware BOM
* [Raspberry PI Zero with WiFi and GPIO Headers](https://www.raspberrypi.org/products/raspberry-pi-zero/)
* [Adafruit 1.2" 7-Segment Display with I2C Backpack](https://www.adafruit.com/product/1270)
* Breadboard Wires
* Optional [RTC Pi Hat](https://www.abelectronics.co.uk/p/70/rtc-pi)
* USB power supply and cable
* Case TBD

## Building the hardware

### GPIO to Adafruit Segment Display I2C
* Raspberry Pi 3.3V to 7-Segment Display IO
* Raspberry Pi 5V to 7-Segment Display VIN
* Raspberry Pi GND to 7-Segment Display GND
* Raspberry Pi SCL to 7-Segment Display SCL
* Raspberry Pi SDA to 7-Segment Display SDA

### Optional RTC Hat

RTC hat is not required and rarely used while you have NTP. However in case you lose power and your internet/wifi doesn't get back on before RPI Zero you will have no time at all or the time will be completely wrong (last time PI was up).

Install RTC Hat between PI GPIO and wires going to Adafruit. Make sure you install coin battery before powering it on.

### Case

TBD - 3D printed?

## Software configuration

### OS

Any RPI OS will do, I used [Raspberry Pi OS](https://www.raspberrypi.org/software/operating-systems/), formerly Raspbian. For this project the Lite version should be used.

### WiFi, Locale, Timezone, etc.

You can use `raspi-config` to configure WiFi, Locale, Timezone etc.

### NTP config

```shell
$ apt install ntp ntpstat
```

### RTC config (optional)

```shell
$ sudo apt install i2c-tools
$ sudo apt remove fake-hwclock
$ sudo echo dtoverlay=i2c-rtc,ds1307 >> /boot/config.txt
$ sudo echo rtc-ds1307 >> /etc/modules
$ sudo echo '5 *  *  * * *    root   /sbin/hwclock -w' >> /etc/crontab
```

Edit `/lib/udev/hwclock-set`, remove following lines:

```
if [ -e /run/systemd/system ] ; then
  exit 0
fi
```

Reboot, check if hwclock works:

```shell
$ sudo i2cdetect -y
```

should show `UU` on position `68`.

```shell
$ sudo hwclock -r
$ sudo hwclock -w
```

### Python Libs

```
$ sudo apt install python3-pip
$ sudo pip3 install adafruit-circuitpython-ht16k33 apscheduler
```

### Clock Service

Move .py in to location you want, it can run from /home/pi or /usr/local/bin, etc.

Move .service in to `/etc/systemd/system`

```shell
$ sudo systemctl daemon-reload
$ sudo systemctl enable BigSeg7x4_clock.service 
$ sudo systemctl start BigSeg7x4_clock.service 
```

## References
* [Adafruit Wiring and Setup](https://learn.adafruit.com/adafruit-led-backpack/python-wiring-and-setup-d74df15e-c55c-487a-acce-a905497ef9db)
* [RTC Pi setup on RPI OS](https://www.abelectronics.co.uk/kb/article/30/rtc-pi-on-a-raspberry-pi-raspbian-jessie)

## Legal

This is not an officially supported Google product.

```
Copyright 2021 Google LLC
```
