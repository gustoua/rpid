# RPId
Log CPU temperature, fan RPM, environmental data (Ambient temperature, Relative humidity, Atmospheric pressure) from external sensors connected to Raspberry Pi GPIO. Control fan state without a built-in gpio-fan overlay 

# ToDo's
- fan activation temps - on, off boundaries, wait before turn off... basically minimize state changes, but do it so it doesn't work all night cooling what's already cold
- configure fan acticavion/deactivation temps. No config file for now
- web server that serves /status with current state to gatus

# Setup
- pi user must be added to the same group as /dev/gpiomem in ([source](https://raspberrypi.stackexchange.com/questions/40105/access-gpio-pins-without-root-no-access-to-dev-mem-try-running-as-root)).

- pi user must be added to syslog group to write into /var/log/rpid.log
