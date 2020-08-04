package main

import (
	"fmt"   
)

func main() {
	tahun := 2020
	hariPertama := 3
	spasi := hariPertama

	var bulan = [13]string {
			"",
			"Januari", "Februari", "Maret",
			"April", "Mei", "Juni",
			"Juli", "Agustus", "September",
			"Oktober", "November", "Desember"}

	var jumlahHari = [13]int {0, 31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	for m := 1; m <= 12; m++ {

		// Cek tahun kabisat
		if  ((((tahun % 4 == 0) && (tahun % 100 != 0)) ||  (tahun % 400 == 0)) && m == 2) {
			jumlahHari[m] = 29
		}
			
		// Print calendar
		fmt.Println("          "+ bulan[m] + " ", tahun)
		fmt.Println("_____________________________________")
		fmt.Printf(" %3s %4s %4s %4s %4s %4s %4s ", "S", "S", "R", "K", "J", "S", "M")
		fmt.Println()

		spasi = (jumlahHari[m-1] + spasi) % 7

		for i := 0; i < spasi; i++ {
			fmt.Print("     ")
		}
			
		for i := 1; i <= jumlahHari[m]; i++ {
			fmt.Printf(" %3d ", i)
			if (((i + spasi) % 7 == 0) || (i == jumlahHari[m])) {
				fmt.Println()
			}
		}
		fmt.Println()
	}
}