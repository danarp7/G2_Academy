package main

import (
	"fmt"
	"math"
)

func main() {
	var menu int = 1
	for menu > 0 && menu <= 3 {
		menuKabataku := 0
		menuLuas := 0
		menuVolume := 0
		var angka int
		var angka1 int
		var angka2 int
		var jariJari int
		var pangkat int
		var sisi1 int
		var sisi2 int
		var sisi3 int
		var tinggi int
		var sisiAlas1 int
		var sisiAlas2 int

		fmt.Println("\n====================")
		fmt.Println("\tPilih Menu")
		fmt.Println("1. Kabataku")
		fmt.Println("2. Hitung Luas")
		fmt.Println("3. Hitung Volume")
		fmt.Println("4. Exit")
		fmt.Print("Pilih menu :\t")
		fmt.Scanf("%d", &menu)
		fmt.Println("\nMenu yang dipilih", menu)
		fmt.Println()

		switch menu {
		case 1:
			fmt.Println("\n==== Kabataku ====")
			fmt.Println("1. Perkalian")
			fmt.Println("2. Pembagian")
			fmt.Println("3. Pertambahan")
			fmt.Println("4. Pengurangan")
			fmt.Println("5. Akar")
			fmt.Println("6. Pangkat")
			fmt.Print("Pilih menu Kabataku :\t")
			fmt.Scanf("%d", &menuKabataku)
			fmt.Println("\nMenu kabataku yang dipilih", menuKabataku)
			fmt.Println()

			if menuKabataku == 1 {
				fmt.Println("\nPerkalian")
				fmt.Print("Masukkan angka pertama :\t")
				fmt.Scanf("%d", &angka1)
				fmt.Println()
				fmt.Print("Masukkan angka pertama :\t")
				fmt.Scanf("%d", &angka2)
				fmt.Println()

				kabataku := Kabataku{float64(angka1), float64(angka2)}
				kabataku.perkalian()
			} else if menuKabataku == 2 {
				fmt.Println("\nPembagian")
				fmt.Print("Masukkan angka pertama :\t")
				fmt.Scanf("%d", &angka1)
				fmt.Println()
				fmt.Print("Masukkan angka pertama :\t")
				fmt.Scanf("%d", &angka2)
				fmt.Println()

				kabataku := Kabataku{float64(angka1), float64(angka2)}
				kabataku.pembagian()
			} else if menuKabataku == 3 {
				fmt.Print("Masukkan angka pertama :\t")
				fmt.Scanf("%d", &angka1)
				fmt.Println()
				fmt.Print("Masukkan angka pertama :\t")
				fmt.Scanf("%d", &angka2)
				fmt.Println()

				kabataku := Kabataku{float64(angka1), float64(angka2)}
				kabataku.pertambahan()
			} else if menuKabataku == 4 {
				fmt.Println("\nPengurangan")
				fmt.Print("Masukkan angka pertama :\t")
				fmt.Scanf("%d", &angka1)
				fmt.Println()
				fmt.Print("Masukkan angka pertama :\t")
				fmt.Scanf("%d", &angka2)
				fmt.Println()

				kabataku := Kabataku{float64(angka1), float64(angka2)}
				kabataku.pengurangan()
			} else if menuKabataku == 5 {
				fmt.Println("\nAkar")
				fmt.Print("Masukkan angka :\t")
				fmt.Scanf("%d", &angka)
				fmt.Println()

				akar(float64(angka))
			} else if menuKabataku == 6 {
				fmt.Println("\nPangkat")
				fmt.Print("Masukkan angka :\t")
				fmt.Scanf("%d", &angka)
				fmt.Println()
				fmt.Print("Masukkan pangkat :\t")
				fmt.Scanf("%d", &pangkat)
				fmt.Println()

				Pangkat(float64(angka), float64(pangkat))
			} else {
				fmt.Println("Anda Salah Memasukan Pilihan !")
			}
		case 2:
			fmt.Println("\n====  Hitung Luas ====")
			fmt.Println("1. Persegi")
			fmt.Println("2. Lingkaran")
			fmt.Print("Pilih menu Hitung Luas :\t")
			menuLuas,_ = fmt.Scanf("%d", &menuLuas)
			fmt.Println("\nMenu hitung luas yang dipilih", menuLuas)
			fmt.Println()

			if menuLuas == 1 {
				fmt.Println("\nLuas Persegi")
				fmt.Print("Masukkan sisi1 :\t")
				fmt.Scanf("%d", &sisi1)
				fmt.Println()
				fmt.Print("Masukkan sisi2 :\t")
				fmt.Scanf("%d", &sisi2)
				fmt.Println()

				luas := LuasPersegi{float64(sisi1), float64(sisi2)}
				luas.persegi()
			} else if menuLuas == 2 {
				fmt.Println("\nLuas Lingkaran")
				fmt.Print("Masukkan jari-jari :\t")
				fmt.Scanf("%d", &jariJari)
				fmt.Println()

				lingkaran(float64(jariJari))
			} else {
				fmt.Println("Anda Salah Memasukan Pilihan !")
			}
		case 3:
			fmt.Println("\n==== Hitung Volume ====")
			fmt.Println("1. Tabung")
			fmt.Println("2. Balok")
			fmt.Println("3. Prisma")
			fmt.Print("Pilih menu Hitung Volume :\t")
			fmt.Scanf("%d", &menuVolume)
			fmt.Println("\nMenu hitung volume yang dipilih", menuVolume)
			fmt.Println()

			if menuVolume == 1 {
				fmt.Println("\nVolume Tabung")
				fmt.Print("Masukkan jari-jari :\t")
				fmt.Scanf("%d", &jariJari)
				fmt.Println()
				fmt.Print("Masukkan tinggi :\t")
				fmt.Scanf("%d", &tinggi)
				fmt.Println()

				volume := VolumeTabung{float64(jariJari), float64(tinggi)}
				volume.tabung()
			} else if menuVolume == 2 {
				fmt.Println("\nVolume Balok")
				fmt.Print("Masukkan sisi1 :\t")
				fmt.Scanf("%d", &sisi1)
				fmt.Println()
				fmt.Print("Masukkan sisi2 :\t")
				fmt.Scanf("%d", &sisi2)
				fmt.Println()
				fmt.Print("Masukkan sisi3 :\t")
				fmt.Scanf("%d", &sisi3)
				fmt.Println()

				volume := VolumeBalok{float64(sisi1), float64(sisi2), float64(sisi3)}
				volume.balok()
			} else if menuVolume == 3 {
				fmt.Println("\nVolume Prisma")
				fmt.Print("Masukkan sisi alas 1 :\t")
				fmt.Scanf("%d", &sisiAlas1)
				fmt.Println()
				fmt.Print("Masukkan sisi alas 2 :\t")
				fmt.Scanf("%d", &sisiAlas2)
				fmt.Println()
				fmt.Print("Masukkan tinggi :\t")
				fmt.Scanf("%d", &tinggi)
				fmt.Println()

				volume := VolumePrisma{float64(sisiAlas1), float64(sisiAlas2), float64(tinggi)}
				volume.prisma()
			} else {
				fmt.Println("Anda Salah Memasukan Pilihan !")

			}
		default:
			fmt.Print("Terimakasih")
		}
	}
}

type Kabataku struct {
	Angka1 float64
	Angka2 float64
}

func (k Kabataku) perkalian() {
	hasil := k.Angka1 * k.Angka2
	fmt.Println("Hasil Perkalian :", hasil)
}

func (k Kabataku) pembagian() {
	hasil := k.Angka1 / k.Angka2
	fmt.Println("Hasil Pembagian :", hasil)
}

func (k Kabataku) pertambahan() {
	hasil := k.Angka1 + k.Angka2
	fmt.Println("Hasil Pertambahan :", hasil)
}

func (k Kabataku) pengurangan() {
	hasil := k.Angka1 - k.Angka2
	fmt.Println("Hasil Pengurangan :", hasil)
}

func akar(angka float64) {
	hasil := math.Sqrt(angka)
	fmt.Println("Hasil akar :", hasil)
}

func Pangkat(angka float64, pangkat float64) {
	hasil := math.Pow(angka, pangkat)
	fmt.Println("Hasil pangkat :", hasil)
}

type LuasPersegi struct {
	Sisi1 float64
	Sisi2 float64
}

func (l LuasPersegi) persegi() {
	hasil := l.Sisi1 * l.Sisi2
	fmt.Println("Hasil luas persegi :", hasil)
}

func lingkaran(jariJari float64) {
	phi := 3.14
	hasil := phi * math.Pow(jariJari, 2)
	fmt.Println("Hasil luas lingkaran :", hasil)
}

type VolumeTabung struct {
	JariJari float64
	Tinggi float64
}

type VolumeBalok struct {
	Sisi1 float64
	Sisi2 float64
	Sisi3 float64
}

type VolumePrisma struct {
	SisiAlas1  float64
	SisiAlas2  float64
	Tinggi float64
}

func (vt VolumeTabung) tabung() {
	phi := 3.14
	hasil := phi * math.Pow(vt.JariJari, 2) * vt.Tinggi
	fmt.Println("Hasil volume tabung :", hasil)
}

func (vb VolumeBalok) balok() {
	hasil := vb.Sisi1 * vb.Sisi2 * vb.Sisi3
	fmt.Println("Hasil volume balok :", hasil)
}

func (vp VolumePrisma) prisma() {
	hasil := 0.5 * vp.SisiAlas1 * vp.SisiAlas2 * vp.Tinggi
	fmt.Println("Hasil volume prisma :", hasil)
}