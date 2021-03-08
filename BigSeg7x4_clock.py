#!/usr/bin/python3
# apt install ntp ntpstat python3-pip
# pip3 install adafruit-circuitpython-ht16k33 apscheduler
import time, datetime, atexit, signal, subprocess, board, socket
from apscheduler.schedulers.blocking import BlockingScheduler
from adafruit_ht16k33.segments import BigSeg7x4

sched = BlockingScheduler()
display = BigSeg7x4(board.I2C())
ti = "    " # time
bl = False  # bottom left dot
tl = False  # top left dot
co = False  # seconds colon
ap = False  # am/pm indicator
br = 0.50   # brightness
ip = ""

def bye(a=0,b=0):
        dispip()
        quit()

def getip():
        try:
                s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
                s.connect(("1.1.1.1", 1))
                ip = s.getsockname()[0].split('.')[3]
                s.close()
                return ip
        except:
                return 0

def dispip():
        display.bottom_left_dot = True
        display.top_left_dot = False
        display.ampm = False
        display.brightness = 1
        display.colon = False
        display.print('{:4}'.format(getip()))

def update():
        global ti,bl,tl,co,ap,br

        display.print(ti)
        display.top_left_dot = tl
        display.bottom_left_dot = bl
        #display.ampm = ap
        display.brightness = br
        display.colon = co

def ntp():
        global bl

        ntp = subprocess.Popen('ntpstat', stdout=subprocess.PIPE, stderr=subprocess.DEVNULL, shell=False)
        out = ntp.communicate()[0].decode()

        if ntp.returncode != 0:
                bl = False
                return

        if 'NTP server' in out:
                bl = True
        else:
                bl = False

def hourly():
        global tl,br

        if '{d:%p}'.format(d=datetime.datetime.now()) == "PM":
                tl = True
        else:
                tl = False

        if time.localtime().tm_hour > 7 and time.localtime().tm_hour < 20:
                br = 0.50
        else:
                br = 0.35

def tick():
        global ti,co

        if (time.localtime().tm_sec % 2) == 0:
                co = True
        else:
                co = False

        ti = '{d:%l}:{d.minute:02}'.format(d=datetime.datetime.now())

        update()


dispip()
time.sleep(5)
atexit.register(bye)
signal.signal(signal.SIGTERM, bye)
signal.signal(signal.SIGINT, bye)
sched.add_job(tick, 'interval', seconds=1)
sched.add_job(ntp, 'interval', seconds=10)
sched.add_job(hourly, 'cron', minute=0, second=0)
hourly() # run once to set right am/pm and brightness from start
sched.start()
