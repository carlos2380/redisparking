package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type CarLog struct {
	License string    `json:"license"`
	TimeIn  time.Time `json:"timeIn"`
}

var clientDB = redis.NewClient(&redis.Options{
	Addr:     "cache:6379",
	Password: "",
	DB:       0,
})

func main() {
	log.Println("Go Redis Tutorial")

	pong, err := clientDB.Ping().Result()
	log.Println(pong, err)

	carLogJson1, err := json.Marshal(CarLog{License: "1234ABC", TimeIn: time.Now()})
	if err != nil {
		log.Println(err)
	}

	carLogJson2, err := json.Marshal(CarLog{License: "4444BBB", TimeIn: time.Now()})
	if err != nil {
		log.Println(err)
	}

	err = clientDB.Set("1234ABC", carLogJson1, 0).Err()
	if err != nil {
		log.Println(err)
	}

	err = clientDB.Set("4444BBB", carLogJson2, 0).Err()
	if err != nil {
		log.Println(err)
	}

	val, err := clientDB.Get("4444BBB").Result()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(clientDB.Exists("1234ABC"))

	log.Info(val)
	fmt.Println(clientDB.DBSize())
	clientDB.Del("1234ABC")
	fmt.Println(clientDB.DBSize())
	val, err = clientDB.Get("1234ABC").Result()
	if err != nil {
		fmt.Println(err)
	}

	val, err = clientDB.Get("4444BBB").Result()
	if err != nil {
		fmt.Println(err)
	}
	log.WithFields(log.Fields{
		"msg": "Init Parking",
	}).Info("Init Parking")
	myRouter := mux.NewRouter()
	myRouter.HandleFunc("/hello", HelloServer).Methods(http.MethodGet, http.MethodOptions)
	myRouter.HandleFunc("/{id}", addCarLog).Methods(http.MethodPost, http.MethodOptions)
	myRouter.HandleFunc("/{id}", delCarLog).Methods(http.MethodDelete, http.MethodOptions)
	myRouter.HandleFunc("/", getall).Methods(http.MethodGet, http.MethodOptions)
	myRouter.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", myRouter)

}

func getall(w http.ResponseWriter, r *http.Request) {
	keys := clientDB.Keys("*")
	list := keys.String()
	log.WithFields(log.Fields{
		"msg": "Get",
	}).Info("Get all Car")
	log.Info(list)
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(
		struct {
			LICENSES string `json:"licenses"`
		}{LICENSES: list})
}

func addCarLog(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"msg": "Add",
	}).Info("Add Car")
	params := mux.Vars(r)
	id := params["id"]
	log.Warn(id)
	carLogJson, err := json.Marshal(CarLog{License: id, TimeIn: time.Now()})
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = clientDB.Set(id, carLogJson, 0).Err()
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func delCarLog(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	log.WithFields(log.Fields{
		"msg": "Del",
	}).Info("Delete Car")
	log.Error(id)
	clientDB.Del(id)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
