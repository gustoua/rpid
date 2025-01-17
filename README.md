# RPId
Log CPU temperature, fan RPM, environmental data (Ambient temperature, Relative humidity, Atmospheric pressure) from external sensors connected to Raspberry Pi GPIO. Control fan state without a built-in gpio-fan overlay 

# ToDo's

- extract env sensors polls to external plugins. I'll need some kind of main service with plugins structure

- throttling detection :
    - https://chewett.co.uk/blog/258/vchi-initialization-failed-raspberry-pi-fixed/
    - https://jamesachambers.com/measure-raspberry-pi-undervoltage-true-clock-speeds/

# Setup
- pi user must be added to the same group as /dev/gpiomem in ([source](https://raspberrypi.stackexchange.com/questions/40105/access-gpio-pins-without-root-no-access-to-dev-mem-try-running-as-root)).

# Real life usage example
This is example of revision [e07bbed](https://github.com/parMaster/rpid/commit/e07bbed66f5384c41d595c599d575dca676a7c38) working on a Pi4 4Gb with in a radiator-case with a 50mm 12v fan installed on top, connected to 5v power through a npn-transistor, like this:
![IMG_3351](https://user-images.githubusercontent.com/1956191/222020060-eb204c20-2573-484a-a245-0be3da81abb7.jpeg)
It would be much better if I could PWM this fan, but cheap fans like this can't be easily pwm-ed, so there are several ways people usually connect these:
- connect permanently to 5v - is a pretty good option. It's completely silent this way, but the downside is that it always works, obviously.
- connect to 3.3v permanently - the same as previous, but under high loads it's not enough, it always spins and not every fan starts from 3.3v
- connect to 5v through a transistor and use stock gpio_fan device tree overlay. I guess it's ok, but not much configurable.

So, this package, apart of being a great exercice, is configurable, it provides an endpoint (/status) accessible by monitoring software and it provides nice charts (/charts endpoint), like this one:
![newplot](https://user-images.githubusercontent.com/1956191/222021495-85ca3665-fb5d-47d2-8218-2eb4e2c78d2b.png)

The chart shows, that under normal-high load, like youtube 1080p videoplay, it cycles from 40 to 45˚C, turning the fan on roughly 50% of time. In the extreme scenario of 4-core stress-test, the temperature can spike to as high as 50˚C, which still is considered "great" for the Raspberry Pi 4. The deepest cooldown to around 36˚C I could achieve by cutting stress-test right after it peaked, so the fan works for the next 3 minutes under no load.
