package main

import (
"encoding/json"
"fmt"
guuid "github.com/google/uuid"
"log"
"net/http"
"time"
)

type ParkirMasuk struct {
	IdParkir  string `json:"id_parkir"`
	StartTime time.Time  `json:"waktu_masuk"`
}

type ParkirKeluar struct {
	IdParkir      string `json:"id_parkir"`
	StartTime     time.Time  `json:"waktu_masuk"`
	EndTime       time.Time  `json:"waktu_keluar"`
	Durasi        int        `json:"durasi"`
	TipeKendaraan string     `json:"tipe_kendaraan"`
	PlatNo        string     `json:"plat_no"`
	BiayaParkir   int        `json:"biaya_parkir"`
}

var dataParkirMasuk []ParkirMasuk
var dataParkirKeluar []ParkirKeluar

func simpanDataParkirMasuk(idParkir string, startTime time.Time) {
	parkirMasuk := ParkirMasuk{IdParkir: idParkir, StartTime: startTime}
	dataParkirMasuk = append(dataParkirMasuk, parkirMasuk)
}

func simpanDataParkirKeluar(idParkir string, startTime time.Time, endTime time.Time, tipeKendaraan string,
	platNo string, durasi int, biayaParkir int) {
	parkirKeluar := ParkirKeluar{IdParkir: idParkir, StartTime: startTime, EndTime: endTime,
		Durasi: durasi, TipeKendaraan: tipeKendaraan, PlatNo: platNo,
		BiayaParkir: biayaParkir}
	dataParkirKeluar = append(dataParkirKeluar,parkirKeluar)
}

func hitungBiayaParkir(idParkir string, startTime time.Time, endTime time.Time, inpTipe string,
	inpPlatNo string) {
	if inpTipe == "Motor" {
		parkirMotor(idParkir, startTime, endTime, inpTipe, inpPlatNo)
	} else if inpTipe == "Mobil" {
		parkirMobil(idParkir, startTime, endTime, inpTipe, inpPlatNo)
	} else {
		fmt.Println("Tipe kendaraan yang anda masukkan salah.")
	}
}

func parkirMotor(idParkir string, startTime time.Time, endTime time.Time, tipeKendaraan string, platNo string) {
	diff := endTime.Sub(startTime)
	durasi := int(diff.Seconds())
	satuHari := (12 * 60 * 60)
	if durasi == 1 {
		biayaParkir := 3000
		simpanDataParkirKeluar(idParkir, startTime, endTime, tipeKendaraan, platNo,
			durasi, biayaParkir)
	} else if durasi >= satuHari{
		biayaParkir := ((satuHari - 1) * 2000) + 3000
		simpanDataParkirKeluar(idParkir, startTime, endTime, tipeKendaraan, platNo,
			durasi, biayaParkir)
	} else {
		biayaParkir := ((durasi - 1) * 2000) + 3000
		simpanDataParkirKeluar(idParkir, startTime, endTime, tipeKendaraan, platNo,
			durasi, biayaParkir)
	}
}

func parkirMobil(idParkir string, startTime time.Time, endTime time.Time, tipeKendaraan string, platNo string) {
	diff := endTime.Sub(startTime)
	durasi := int(diff.Seconds())
	satuHari := (12 * 60 * 60)
	if durasi == 1 {
		biayaParkir := 5000
		simpanDataParkirKeluar(idParkir, startTime, endTime, tipeKendaraan, platNo,
			durasi, biayaParkir)
	} else if durasi >= satuHari{
		biayaParkir := ((satuHari - 1) * 3000) + 5000
		simpanDataParkirKeluar(idParkir, startTime, endTime, tipeKendaraan, platNo,
			durasi, biayaParkir)
	} else {
		biayaParkir := ((durasi - 1) * 3000) + 5000
		simpanDataParkirKeluar(idParkir, startTime, endTime, tipeKendaraan, platNo,
			durasi, biayaParkir)
	}
}

func customerMasuk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")


	if r.Method == "POST" {
		idParkir := guuid.New().String()
		startTime := time.Now()
		karcis := ParkirMasuk{idParkir,  startTime}

		var result, err = json.Marshal(karcis)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(result)
		if err != nil {
			log.Fatal(err)
		}

		simpanDataParkirMasuk(idParkir, startTime)

		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func customerKeluar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var inpIdP = r.FormValue("id_parkir")
		var inpTipe = r.FormValue("tipe_kendaraan")
		var inpPlatNo = r.FormValue("plat_no")
		endTime := time.Now()

		var result []byte
		var err error

		for i := 0; i < len(dataParkirMasuk); i++ {
			if dataParkirMasuk[i].IdParkir == inpIdP {
				hitungBiayaParkir(dataParkirMasuk[i].IdParkir, dataParkirMasuk[i].StartTime, endTime, inpTipe, inpPlatNo)
				for _, dpk := range dataParkirKeluar {
					if dpk.IdParkir == dataParkirMasuk[i].IdParkir {
						result, err = json.Marshal(dpk)

						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}

						dataParkirMasuk = append(dataParkirMasuk[:i], dataParkirMasuk[i+1:]...)

						w.Write(result)
						return
					}
				}
			}
		}

		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func parkirMasuk(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var result, err = json.Marshal(dataParkirMasuk)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func parkirKeluar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var result, err = json.Marshal(dataParkirKeluar)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func main() {
	http.HandleFunc("/server/customer/masuk", customerMasuk)
	http.HandleFunc("/server/customer/keluar", customerKeluar)
	http.HandleFunc("/server/parkir/masuk", parkirMasuk)
	http.HandleFunc("/server/parkir/keluar", parkirKeluar)

	fmt.Println("starting web server at http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}