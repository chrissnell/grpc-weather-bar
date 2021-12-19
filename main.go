package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"regexp"
	"time"

	weather "github.com/chrissnell/grpc-weather-bar/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

func main() {
	uid, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	cfgFile := flag.String("config", uid.HomeDir+"/.config/grpc-weather-bar/config", "Path to weather-bar config file (default: $HOME/.config/grpc-weather-bar/config)")
	flag.Parse()

	// Read our server configuration
	filename, _ := filepath.Abs(*cfgFile)
	cfg, err := NewConfig(filename)
	if err != nil {
		log.Fatalln("Error reading config file.  Did you pass the -config flag?  Run with -h for help.\n", err)
	}

	if cfg.Server.Hostname == "" {
		log.Fatalln("Error: must provide a server hostname in config file.")
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = "7500"
	}

	// Default timeout is 30s
	if cfg.Server.Timeout == 0 {
		cfg.Server.Timeout = 30 * time.Second
	}

	var conn *grpc.ClientConn

	errCh := make(chan error)

	for {
		// Use a client cert if we are doing mutual TLS authentication
		if cfg.Server.MTLSCert != "" {
			creds, err := credentials.NewClientTLSFromFile(cfg.Server.MTLSCert, "")
			if err != nil {
				log.Fatalln("Could not load TLS certificate:", err)
			}

			fmt.Printf("Dialing %v:%v ...\n", cfg.Server.Hostname, cfg.Server.Port)

			conn, err = grpc.Dial(cfg.Server.Hostname+":"+cfg.Server.Port, grpc.WithTransportCredentials(creds), grpc.WithBlock(), grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:    time.Second,
				Timeout: cfg.Server.Timeout,
			}))
			if err != nil {
				log.Println("Failed to connect:", err)
			} else {
				go getLiveWeather(cfg, conn, errCh)
			}
		} else {
			fmt.Printf("Dialing %v:%v ...\n", cfg.Server.Hostname, cfg.Server.Port)

			if !cfg.Server.UseTLS {
				conn, err = grpc.Dial(cfg.Server.Hostname+":"+cfg.Server.Port,
					grpc.WithInsecure(),
					grpc.WithBlock(),
					grpc.WithTimeout(cfg.Server.Timeout),
					grpc.WithKeepaliveParams(
						keepalive.ClientParameters{
							Time:    15 * time.Second,
							Timeout: cfg.Server.Timeout,
						},
					))
			} else {
				tlsConfig := tls.Config{}

				conn, err = grpc.Dial(cfg.Server.Hostname+":"+cfg.Server.Port,
					grpc.WithBlock(),
					grpc.WithTransportCredentials(credentials.NewTLS(&tlsConfig)),
					grpc.WithTimeout(cfg.Server.Timeout),
					grpc.WithKeepaliveParams(
						keepalive.ClientParameters{
							Time:    15 * time.Second,
							Timeout: cfg.Server.Timeout,
						},
					))
			}
			if err != nil {
				log.Println("Failed to connect:", err)
			} else {
				fmt.Println("Getting live weather...")
				go getLiveWeather(cfg, conn, errCh)
			}
		}

		err = <-errCh
		if err != nil {
			fmt.Println("Connection to weather server failed.  Retrying in 1s.")
			time.Sleep(time.Second)
		}

	}

}

func getLiveWeather(cfg *Config, conn *grpc.ClientConn, errCh chan error) {
	var rdg *weather.WeatherReading

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer conn.Close()

	c := weather.NewWeatherClient(conn)

	lwc, err := c.GetLiveWeather(ctx, &weather.Empty{})
	if err != nil {
		log.Fatalln("Could not create GetLiveWeather client:", err)
	}

	for {
		rdg, err = lwc.Recv()
		if err != nil {
			log.Println("Error receiving from server:", err)
			errCh <- err
			return
		}

		fmt.Println(formatOutput(cfg, rdg))

	}

}

func formatOutput(c *Config, rdg *weather.WeatherReading) string {
	var output string
	var cardIndex int

	output = c.Format.WxFormat

	tempC := (rdg.OutsideTemperature - 32) * (5.0 / 9.0)

	cardDirections := []string{"  N", "NNE", " NE", "ENE",
		"  E", "ESE", " SE", "SSE",
		"  S", "SSW", " SW", "WSW",
		"  W", "WNW", " NW", "NNW"}

	cardIndex = int((float32(rdg.WindDirection) + float32(11.25)) / float32(22.5))
	cardDirection := cardDirections[cardIndex%16]

	regTempF := regexp.MustCompile("%temperature-fahrenheit%")
	regTempC := regexp.MustCompile("%temperature-celcius%")
	regHum := regexp.MustCompile("%humidity%")
	regWindS := regexp.MustCompile("%windspeed%")
	regWindD := regexp.MustCompile("%winddirection%")
	regWindC := regexp.MustCompile("%windcardinal%")
	regRain := regexp.MustCompile("%rainfall%")

	output = regTempF.ReplaceAllLiteralString(output, fmt.Sprintf("%.1f", rdg.OutsideTemperature))
	output = regTempC.ReplaceAllLiteralString(output, fmt.Sprintf("%.1f", tempC))

	output = regHum.ReplaceAllLiteralString(output, fmt.Sprintf("%v", rdg.OutsideHumidity))

	output = regWindS.ReplaceAllLiteralString(output, fmt.Sprintf("%v", rdg.WindSpeed))
	output = regWindD.ReplaceAllLiteralString(output, fmt.Sprintf("%v", rdg.WindDirection))
	output = regWindC.ReplaceAllLiteralString(output, cardDirection)

	output = regRain.ReplaceAllLiteralString(output, fmt.Sprintf("%.2f", rdg.RainfallDay))

	return output
}
