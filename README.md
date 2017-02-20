# garage-go

A simple garage door controller and monitor

I far too frequently leave my garage door open.  I wanted a way to know if it is open or not.
And a way to close it from my phone.

I wired up a Raspberry PI zero to a spare wireless door opener remote, and a simple magnetic
reed switch for doors/windows.

This code implements a simple REST api to open/close/check the door status.

It also implements a Monitor to check if the door has been open for longer than 1 hour,
and closes it if that time is exceeded.

## TODO
* Connect to IFTTT so I can control using Alexa or Google Home
* Create a Telegram Bot so I can control via secure text message

