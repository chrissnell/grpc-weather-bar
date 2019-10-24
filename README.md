![](https://island.nu/github/weather-bar/weather-bar.gif)

# grpc-weather-bar
Pull live weather reports from [gopherwx](https://github.com/chrissnell/gopherwx) and display them in Polybar or lemonbar.

# How to use
1. Compile grpc-weather-bar and put the binary in your `$PATH`:  `go get -u github.com/chrissnell/grpc-weather-bar`
2. Configure a [gopherwx](https://github.com/chrissnell/gopherwx) server, or use someone else's.  Make sure you specify a gRPC server in the config.yaml.  If this is for a laptop that will leave your network, expose gopherwx's gRPC port to the Internet via your router.https://github.com/chrissnell/grpc-weather-bar/blob/master/README.md
3. Create your grpc-weather-bar configuration file in `${HOME}/.config/grpc-weather-bar/config` using the [example](https://github.com/chrissnell/grpc-weather-bar/blob/master/example/config) from this repo.
4. Choose one of the options below, depending on your bar.
## Polybar
Add a new module to your Polybar config that uses Polybar's script module to run grpc-weather-bar in tail fashion.  See the config [snippet](https://github.com/chrissnell/grpc-weather-bar/blob/master/example/polybar-config) in this repo for an example.
## Lemonbar
Simply pipe the output of grpc-weather-bar to lemonbar:   `grpc-weather-bar | lemonbar`.  I recommend the [patched version](https://github.com/krypt-n/bar) that supports Xft fonts so that you can have some sweet icons.

# Fonts
grpc-weather-bar looks best with the Font Awesome icons from [Nerd Fonts](https://github.com/ryanoasis/nerd-fonts).  I used the "Sauce Code Pro" for the screenshot above.
