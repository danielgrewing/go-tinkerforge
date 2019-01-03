package main

import(
	"fmt"
    "github.com/Tinkerforge/go-api-bindings/ipconnection"
	"github.com/Tinkerforge/go-api-bindings/outdoor_weather_bricklet"
    "log"
    "net/http"
	"github.com/gorilla/mux"	
	"strings"
)

const addr string = "192.168.0.101:4223"
const uid string = "DYn"

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/prometheus", prometheus).Methods("GET")
    log.Println(http.ListenAndServe(":7080", router))
}

func prometheus(w http.ResponseWriter, r *http.Request) {
	ipcon := ipconnection.New()
	defer ipcon.Close()
	ow, _ := outdoor_weather_bricklet.New(uid, &ipcon)

	ipcon.Connect(addr)
	defer ipcon.Disconnect()
	var sb strings.Builder

	identifiers, _ := ow.GetSensorIdentifiers()
	for _ ,element := range identifiers{
		temperature , humidity, _ , _ := ow.GetSensorData(element)
		tempvalue := float64(temperature)/10.0
		sb.WriteString("#tinkerforge_outdoor_weather_bricklet\n")
		sb.WriteString("tinkerforge_outdoor_weather_bricklet_" + fmt.Sprintf("%o", element) + "_temperatur_celcius" + " ")
		sb.WriteString(fmt.Sprintf("%f", tempvalue))
		sb.WriteString("\n")
		sb.WriteString("tinkerforge_outdoor_weather_bricklet_" + fmt.Sprintf("%o", element) + "_humidity_procent" + " ")
		sb.WriteString(fmt.Sprintf("%v", humidity))
    }  
	exportstring := sb.String()
	w.Write([]byte(exportstring))
  }