package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	forecast "github.com/mlbright/forecast/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron"
)

var addr = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
var apikey = flag.String("api-key", "", "DarkSky API Key")
var lat = flag.String("latitude", "51.4416", "Latitude")
var long = flag.String("longitude", "5.4697", "Longitude")
var every = flag.String("every", "2m", "Update time")

var (
	temperatureGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_temperature_celcius",
			Help: "Temperature in degree Celcius",
		},
		[]string{"latitude", "longitude"},
	)
	precipIntensity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_precipitation_intensity",
			Help: "Precipitation intensity",
		},
		[]string{"latitude", "longitude"},
	)
	precipProbability = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_precipitation_probability",
			Help: "Precipitation probability",
		},
		[]string{"latitude", "longitude"},
	)
	apparentTemperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_apparent_temperature_celcius",
			Help: "Apparent temperature in degree Celcius",
		},
		[]string{"latitude", "longitude"},
	)
	dewPointTemperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_dew_point_celcius",
			Help: "Dew point in degree Celcius",
		},
		[]string{"latitude", "longitude"},
	)
	humidity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_humidity",
			Help: "Humidity",
		},
		[]string{"latitude", "longitude"},
	)
	pressure = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_pressure_mbar",
			Help: "Pressure in mB",
		},
		[]string{"latitude", "longitude"},
	)
	windSpeed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_wind_speed_kmh",
			Help: "Wind speed in km/h",
		},
		[]string{"latitude", "longitude"},
	)
	windGust = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_wind_gust_kmh",
			Help: "Wind gust in km/h",
		},
		[]string{"latitude", "longitude"},
	)
	windBearing = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_bearing_degree",
			Help: "Wind bearing",
		},
		[]string{"latitude", "longitude"},
	)
	cloudCover = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_cloud_cover",
			Help: "Cloud cover",
		},
		[]string{"latitude", "longitude"},
	)
	uvIndex = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_uv_index",
			Help: "UV index",
		},
		[]string{"latitude", "longitude"},
	)
	visibility = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_visibility_km",
			Help: "Visibility km",
		},
		[]string{"latitude", "longitude"},
	)
	ozone = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "darksky_ozone_dobson",
			Help: "Ozone in dobson",
		},
		[]string{"latitude", "longitude"},
	)
)

func init() {
	prometheus.MustRegister(temperatureGauge)
	prometheus.MustRegister(precipIntensity)
	prometheus.MustRegister(precipProbability)
	prometheus.MustRegister(apparentTemperature)
	prometheus.MustRegister(dewPointTemperature)
	prometheus.MustRegister(humidity)
	prometheus.MustRegister(pressure)
	prometheus.MustRegister(windSpeed)
	prometheus.MustRegister(windBearing)
	prometheus.MustRegister(cloudCover)
	prometheus.MustRegister(visibility)
	prometheus.MustRegister(ozone)
}

func f2s(f float64) string {
	return fmt.Sprintf("%f", f)
}

func collectSample() {
	log.Println("Collecting sample...")
	f, err := forecast.Get(*apikey, *lat, *long, "now", forecast.SI, forecast.English)
	if err != nil {
		log.Println(err)
		log.Println("Skipping measurement due to error.")
		return
	}

	temperatureGauge.With(prometheus.Labels{"latitude": f2s(f.Latitude), "longitude": f2s(f.Longitude)}).Set(f.Currently.Temperature)
	precipIntensity.With(prometheus.Labels{"latitude": f2s(f.Latitude), "longitude": f2s(f.Longitude)}).Set(f.Currently.PrecipIntensity)
	precipProbability.With(prometheus.Labels{"latitude": f2s(f.Latitude), "longitude": f2s(f.Longitude)}).Set(f.Currently.PrecipProbability)
	apparentTemperature.With(prometheus.Labels{"latitude": f2s(f.Latitude), "longitude": f2s(f.Longitude)}).Set(f.Currently.ApparentTemperature)
	dewPointTemperature.With(prometheus.Labels{"latitude": f2s(f.Latitude), "longitude": f2s(f.Longitude)}).Set(f.Currently.DewPoint)
	humidity.With(prometheus.Labels{"latitude": f2s(f.Latitude), "longitude": f2s(f.Longitude)}).Set(f.Currently.Humidity)
	pressure.With(prometheus.Labels{"latitude": f2s(f.Latitude), "longitude": f2s(f.Longitude)}).Set(f.Currently.Pressure)
	windSpeed.With(prometheus.Labels{"latitude": f2s(f.Latitude), "longitude": f2s(f.Longitude)}).Set(f.Currently.WindSpeed)
	windBearing.With(prometheus.Labels{"latitude": f2s(f.Latitude), "longitude": f2s(f.Longitude)}).Set(f.Currently.WindBearing)
	cloudCover.With(prometheus.Labels{"latitude": f2s(f.Latitude), "longitude": f2s(f.Longitude)}).Set(f.Currently.CloudCover)
	visibility.With(prometheus.Labels{"latitude": f2s(f.Latitude), "longitude": f2s(f.Longitude)}).Set(f.Currently.Visibility)
	ozone.With(prometheus.Labels{"latitude": f2s(f.Latitude), "longitude": f2s(f.Longitude)}).Set(f.Currently.Ozone)
}

func main() {
	flag.Parse()
	http.Handle("/metrics", prometheus.Handler())

	c := cron.New()
	c.AddFunc(fmt.Sprintf("@every %s", *every), collectSample)
	c.Start()

	log.Printf("Listening on %s!", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
