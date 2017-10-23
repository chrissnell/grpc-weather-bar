![](http://island.nu/github/weather-bar/weather-bar.gif)

# weather-bar
Pull live weather reports from [gopherwx](https://github.com/chrissnell/gopherwx) and display them in Polybar.

# How to use
1. Compile weather-bar and put the binary in your `$PATH`:  `go get -u github.com/chrissnell/weather-bar`
2. Configure a [gopherwx](https://github.com/chrissnell/gopherwx) server, or use someone else's.  Make sure you specify a gRPC server in the config.yaml.  If this is for a laptop that will leave your network, expose gopherwx's gRPC port to the Internet via your router.
3. Create your weather-bar configuration file in `${HOME}/.config/weather-bar/config` using the [example](https://github.com/chrissnell/weather-bar/blob/master/example/config) from this repo.
4. Add a new module to your Polybar config that uses Polybar's script module to run weather-bar in tail fashion.  See the config [snippet](https://github.com/chrissnell/weather-bar/blob/master/example/polybar-config) in this repo for an example.
