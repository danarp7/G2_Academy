package main

import (
	"fmt"
	guuid "github.com/google/uuid"
	"time"
)

type Parkir struct {
	idParkir guuid.UUID
	startTime time.Time
}

var dataParkir = make(map[Parkir]guuid.UUID)

func simpanDataParkir(idParkir guuid.UUID, startTime time.Time) {
	parkir := Parkir{idParkir: idParkir, startTime: startTime}
	dataParkir[parkir] = idParkir

	fmt.Println("ID Parkir : ", idParkir)
	fmt.Println("Waktu Masuk : ", startTime.Format("2006-Jan-02 Monday 03:04:05"))
}

func hitungBiayaParkir(inpIdP guuid.UUID, inpTipe int, inpPlatNo string, endTime time.Time) {
	endTime = endTime
	platNo := inpPlatNo
	if inpTipe == 1 {
		tipe := "Motor"
		for dp := range dataParkir {
			if dp.idParkir == inpIdP {
				parkirMotor(dp.idParkir, dp.startTime, endTime, tipe, platNo)
			}
		}
	} else if inpTipe == 2 {
		tipe := "Mobil"
		for dp := range dataParkir {
			if dp.idParkir == inpIdP {
				parkirMobil(dp.idParkir, dp.startTime, endTime, tipe, platNo)
			}
		}
	} else {
		fmt.Println("Tipe kendaraan yang anda masukkan salah.")
	}
}

func parkirMotor(idParkir guuid.UUID, startTime time.Time, endTime time.Time, tipe string, platNo string) {
	diff := endTime.Sub(startTime)
	durasi := int(diff.Seconds())
	satuHari := (12 * 60 * 60)
	if durasi == 1 {
		biaya := 3000
		cetakStruk(idParkir, startTime, endTime, tipe, platNo, durasi, biaya)
	} else if durasi >= satuHari{
		biaya := ((satuHari - 1) * 2000) + 3000
		cetakStruk(idParkir, startTime, endTime, tipe, platNo, durasi, biaya)
	} else {
		biaya := ((durasi - 1) * 2000) + 3000
		cetakStruk(idParkir, startTime, endTime, tipe, platNo, durasi, biaya)
	}
}

func parkirMobil(idParkir guuid.UUID, startTime time.Time, endTime time.Time, tipe string, platNo string) {
	diff := endTime.Sub(startTime)
	durasi := int(diff.Seconds())
	satuHari := (12 * 60 * 60)
	if durasi == 1 {
		biaya := 5000
		cetakStruk(idParkir, startTime, endTime, tipe, platNo, durasi, biaya)
	} else if durasi >= satuHari{
		biaya := ((satuHari - 1) * 3000) + 5000
		cetakStruk(idParkir, startTime, endTime, tipe, platNo, durasi, biaya)
	} else {
		biaya := ((durasi - 1) * 3000) + 5000
		cetakStruk(idParkir, startTime, endTime, tipe, platNo, durasi, biaya)
	}
}

func cetakStruk(idParkir guuid.UUID, startTime time.Time, endTime time.Time, tipe string, platNo string,
				durasi int, biaya int) {
	fmt.Println("\n Biaya Parkir")
	fmt.Println("Id Parkir : ", idParkir)
	fmt.Println("Waktu Masuk : ", startTime.Format("2006-Jan-02 Monday 03:04:05"))
	fmt.Println("Waktu Keluar : ", endTime.Format("2006-Jan-02 Monday 03:04:05"))
	fmt.Println("Selisih Waktu : ", durasi)
	fmt.Println("Tipe Kendaraan : ", tipe)
	fmt.Println("Plat No : ", platNo)
	fmt.Println("Biaya Parkir : ", biaya)
	for key, value := range dataParkir {
		if value == idParkir {
			delete(dataParkir, key)
		}
	}
}

func main() {
	var menu int = 0
	for menu != 3 {
		var inpIdP string
		var inpTipe int
		var inpPlatNo string

		fmt.Println("\n====================")
		fmt.Println("\tPilih Menu")
		fmt.Println("1. Parkir Masuk")
		fmt.Println("2. Parkir Keluar")
		fmt.Println("3. Exit")
		fmt.Print("Pilih menu :\t")
		fmt.Scanf("%d", &menu)
		fmt.Println("\nMenu yang dipilih", menu)
		fmt.Println()

		switch menu {
		case 1:
			fmt.Println("\n==== Parkir Masuk ====")
			idParkir := guuid.New()
			startTime := time.Now()

			simpanDataParkir(idParkir, startTime)
		case 2:
			fmt.Println("\n==== Parkir Keluar ====")
			for dp := range dataParkir {
				fmt.Println("ID Parkir : ", dp.idParkir)
				fmt.Println("Waktu Masuk : ", dp.startTime.Format("2006-Jan-02 Monday 03:04:05"))
			}
			fmt.Print("\nID Parkir :\t")
			fmt.Scanf("%s", &inpIdP)
			idParkir := guuid.MustParse(inpIdP)
			fmt.Println("\n1. Motor")
			fmt.Println("2. Mobil")
			fmt.Print("Tipe Kendaraan :\t")
			fmt.Scanf("%d", &inpTipe)
			fmt.Print("Plat No :\t")
			fmt.Scanf("%s", &inpPlatNo)
			fmt.Println()
			endTime := time.Now()

			hitungBiayaParkir(idParkir, inpTipe, inpPlatNo, endTime)
		default:
			fmt.Print("Terimakasih")
		}
	}
}