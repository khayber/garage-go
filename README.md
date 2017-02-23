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
Some of the configuration is currently hardcoded.  Specifically the GPIO pins and the automatic close feature.  These need to have external configuration.

The REST API requires server.crt and server.key PEM files for HTTPS.  It also uses BasicAuth authentication and requires $USER and $PASS environment variables to be set.

The Telegram bot requires a $TOKEN environment variable to be set.  Get this from Telegram when you set up your bot with Botfather.

## Raspberry Pi GPIO
I'm using the github.com/stianeikeland/go-rpio library for accessing the GPIO pins on the pi.

One GPIO pin is set to Output and is connected to an NPN transistor with a built-in bias resistor.  The emitter/collector are connected across the remote's switch that opens/closes the door.  When the output pin is set Low the circuit is closed.  I hold it closed for a half second before setting it High again.

The remote is normally powered by a 2032 button cell at 3V.  Now it is powered by the 3.3V ping on the pi, so a single power source is used for everything.

Another GPIO pin is set as Input with a pull-down resistor and is connected to the magnetic reed switch mounted to the door.  This detects if the door is open by more than a few inches.

The monitor goroutine checks the door status every few seconds and records the time it is first detected as open.  When the door is detected as closed the time is reset.  If the open time exceeds some value (1 hour) the door is closed.

## Telegram Bot
I was looking for a good/secure method to control this thing and saw multiple projects using Telegram's bot API.  It was really simple to get it working.

It currently supports /start, /status, /open and /close actions.  I can envision hooking up a pi camera to this and getting a picture of the state of the door as well as sending a notification whenever the door state changes or the automatic close feature activates.  I'm also planning on learning how to set up custom a ResponseKeyboard.

## Pictures and Screenshots
Put some of those here...

## TODO
* Connect to IFTTT so I can control using Alexa or Google Home
* More configuration settings (environment variables and/or config file)
* Enable/disable REST or Telegram Bot features
* Configure Automatic close feature (off, notify, close)
