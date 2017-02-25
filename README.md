# garage-go

A simple garage door controller and monitor

I far too frequently leave my garage door open.  I wanted a way to know if it is open or not.
And a way to close it from my phone.

I wired up a Raspberry PI zero to a spare wireless door opener remote, and a simple magnetic
reed switch for doors/windows.

This code implements a simple REST api to open/close/check the door status.
It also implements a simple Telegram chat bot that has the same basic features.

There is also a monitor to check if the door has been open for longer than 1 hour,
and closes it if that time is exceeded.

## Configuration
Using golang's 'flag' module for command line configuration.  Use -h or --help for a listing.

If you only want Telegram support use:  `./garage -telegram -token "your token"`
See Telegram's Botfather for how to create your own bot and get a token.

If you want the REST API use: `./garage -rest -user "username" -pass "password"`
If you use the -ssl flag you will need to create the server.crt/server,key files.

### systemd
The service can be controlled by systemd with a service file such as the following
```
/lib/systemd/system/garage.service
[Unit]
Description=Garage Door Control/Monitor service
After=network.target auditd.service

[Service]
EnvironmentFile=-/etc/default/garage
ExecStart=/usr/bin/garage $OPTS
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

Then put your preferred options in /etc/default/garage like this:

```
OPTS=-telegram -token your-token-here
```

## Raspberry Pi GPIO
I'm using https://github.com/stianeikeland/go-rpio for accessing the GPIO pins on the pi.

One GPIO pin is set to Output and is connected to an NPN transistor with a built-in bias resistor.  The emitter/collector are connected across the remote's switch that opens/closes the door.  When the output pin is set Low the circuit is closed.  I hold it closed for a half second before setting it High again.

The remote is normally powered by a 2032 button cell at 3V.  Now it is powered by the 3.3V ping on the pi, so a single power source is used for everything.

Another GPIO pin is set as Input with a pull-down resistor and is connected to the magnetic reed switch mounted to the door.  This detects if the door is open by more than a few inches.

The monitor goroutine checks the door status every few seconds and records the time it is first detected as open.  When the door is detected as closed the time is reset.  If the open time exceeds some value (1 hour) the door is closed.

## Telegram Bot
I was looking for a good/secure method to control this thing and saw multiple projects using Telegram's bot API.  It was really simple to get it working.

I'm using https://github.com/go-telegram-bot-api/telegram-bot-api for the client API.

It currently supports /start, /status, /open and /close actions.  I can envision hooking up a pi camera to this and getting a picture of the state of the door as well as sending a notification whenever the door state changes or the automatic close feature activates.  I'm also planning on learning how to set up custom a ResponseKeyboard.

## Pictures and Screenshots
Put some of those here...

## TODO
* Connect to IFTTT so I can control using Alexa or Google Home
