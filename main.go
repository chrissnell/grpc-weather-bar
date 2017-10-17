package main

import (
	"flag"
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"regexp"

	weather "github.com/chrissnell/weather-bar/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	var rdg *weather.WeatherReading

	uid, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}

	cfgFile := flag.String("config", uid.HomeDir+"/.config/weather-bar/config", "Path to weather-bar config file (default: $HOME/.config/weather-bar/config)")
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var conn *grpc.ClientConn

	if cfg.Server.Cert != "" {
		creds, err := credentials.NewClientTLSFromFile(cfg.Server.Cert, "")
		if err != nil {
			log.Fatalln("Could not load TLS certificate:", err)
		}

		conn, err = grpc.Dial(cfg.Server.Hostname+":"+cfg.Server.Port, grpc.WithTransportCredentials(creds))
	} else {
		conn, err = grpc.Dial(cfg.Server.Hostname+":"+cfg.Server.Port, grpc.WithInsecure())
		if err != nil {
			log.Fatalln("Failed to connect:", err)
		}
	}

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
			break
		}

		fmt.Println(formatOutput(cfg, rdg))

	}

}

func formatOutput(c *Config, rdg *weather.WeatherReading) string {
	var output string
	var cardIndex int

	output = c.Format.WxFormat

	tempC := (rdg.OutsideTemp - 32) * (5 / 9)

	cardDirections := []string{"  N", "NNE", " NE", "ENE",
		"  E", "ESE", " SE", "SSE",
		"  S", "SSW", " SW", "WSW",
		"  W", "WNW", " NW", "NNW"}

	cardIndex = int((float32(rdg.WindDir) + float32(11.25)) / float32(22.5))
	cardDirection := cardDirections[cardIndex%16]

	regTempF := regexp.MustCompile("%temperature-fahrenheit%")
	regTempC := regexp.MustCompile("%temperature-celcius%")
	regHum := regexp.MustCompile("%humidity%")
	regWindS := regexp.MustCompile("%windspeed%")
	regWindD := regexp.MustCompile("%winddirection%")
	regWindC := regexp.MustCompile("%windcardinal%")
	regRain := regexp.MustCompile("%rainfall%")

	output = regTempF.ReplaceAllLiteralString(output, fmt.Sprintf("%.1f", rdg.OutsideTemp))
	output = regTempC.ReplaceAllLiteralString(output, fmt.Sprintf("%.1f", tempC))

	output = regHum.ReplaceAllLiteralString(output, fmt.Sprintf("%v", rdg.OutsideHumidity))

	output = regWindS.ReplaceAllLiteralString(output, fmt.Sprintf("%v", rdg.WindSpeed))
	output = regWindD.ReplaceAllLiteralString(output, fmt.Sprintf("%v", rdg.WindDir))
	output = regWindC.ReplaceAllLiteralString(output, cardDirection)

	output = regRain.ReplaceAllLiteralString(output, fmt.Sprintf("%.1f", rdg.RainfallDay))

	return output
}
