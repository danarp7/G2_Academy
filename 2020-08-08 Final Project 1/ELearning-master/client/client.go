package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	pb "E-Learning/proto"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

//Soal dengan struct
type Soal struct {
	Pertanyaan string
	Jawaban    string
	Pilihan    []string
}

//PaketSoal dengan struct
type PaketSoal struct {
	Tingkat string
	Kelas   string
	Matpel  string
	Tema    string
	Soal    []Soal
}

type userData struct {
	Username string
	Password string
	Nama     string
	Umur     int32
	Status   string
}

var tempJawaban []string
var jawaban []string
var data userData

func main() {
	homeLogin()
}

//=================================METHOD UNTUK CONNECT GRPC==================

//ConnectServer dial ke service siswa
func ConnectServer() pb.ConnectClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return pb.NewConnectClient(conn)
}

//ConnectServerGuru dial ke service guru
func ConnectServerGuru() pb.GuruClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// defer conn.Close()
	return pb.NewGuruClient(conn)
}

func ConnectServerData() pb.DataClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// defer conn.Close()
	return pb.NewDataClient(conn)
}

//============================================================================

func CheckLogin(user string, pass string) *pb.UserData {
	conn := ConnectServerData()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	inputLogin := pb.InputLogin{Username: user, Password: pass}

	login, err := conn.CheckLogin(ctx, &inputLogin)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return login
}

func NewSignUp(ud pb.User) *pb.UserData {
	conn := ConnectServerData()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	register, err := conn.NewSignUp(ctx, &ud)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return register
}

func homeLogin() {
	data = userData{}
	var menuHome int = 0

	for menuHome != 3 {
		fmt.Println("Home")
		fmt.Println("1. Login")
		fmt.Println("2. SignUp")
		fmt.Println("3. Exit")
		fmt.Scanf("%d", &menuHome)
		switch menuHome {
		case 1:
			login()
		case 2:
			signup()
		case 3:
			fmt.Println("Terimakasih")
		default:
			fmt.Println("Mohon pilih menu yang benar")
		}
	}
}

func login() {
	var menuHome int = 0

	for menuHome != 99 {
		var user string = ""
		var pass string = ""
		fmt.Print("Username :")
		fmt.Scanf("%s", &user)
		fmt.Println()
		fmt.Print("Password :")
		fmt.Scanf("%s", &pass)
		fmt.Println()

		userLogin(user, pass)
	}
}

func signup() {
	reader := bufio.NewReader(os.Stdin)
	var user string = ""
	var pass string = ""
	var nama, status string
	var umur, intstatus int32

	fmt.Print("Username :")
	fmt.Scanf("%s", &user)

	//regex check di user n pass
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	fmt.Println(re.MatchString(user))
	if !re.MatchString(user) {
		fmt.Println("Mohon input username dengan format email yang benar")
		signup()
	}

	fmt.Println()
	fmt.Print("Password :")
	fmt.Scanf("%s", &pass)
	fmt.Println()

	err := verifyPassword(pass)
	if err != "" {
		fmt.Println(err)
		signup()
	}

	fmt.Println("Input nama Anda")
	nama, _ = reader.ReadString('\n')
	nama = strings.Replace(nama, "\n", "", -1)

	fmt.Println("Input umur Anda")
	fmt.Scanf("%d", &umur)

INPUTSTATUS:
	fmt.Println("Input status Anda")
	fmt.Println("1. Siswa")
	fmt.Println("2. Guru")
	fmt.Scanf("%d", &intstatus)

	if intstatus == 1 {
		status = "Siswa"
	} else if intstatus == 2 {
		status = "Guru"
	} else {
		fmt.Println("Mohon input status yang benar")
		goto INPUTSTATUS
	}

	newUser := pb.User{Username: user, Password: pass, Nama: nama, Umur: umur, Status: status}

	regreply := NewSignUp(newUser)
	fmt.Println(regreply.Notice)
	main()
}

func verifyPassword(password string) string {
	var uppercasePresent bool
	var numberPresent bool
	var specialCharPresent bool
	const minPassLength = 8
	const maxPassLength = 64
	var passLen int
	var errorString = ""

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch != ' ':
			passLen++
		}
	}
	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !uppercasePresent {
		appendError("minimal 1 huruf kapital")
	}
	if !numberPresent {
		appendError("minimal 1 angka")
	}
	if !specialCharPresent {
		appendError("minimal 1 special character")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(fmt.Sprintf("password harus antara %d - %d karakter", minPassLength, maxPassLength))
	}

	if len(errorString) != 0 {
		return errorString
	}
	return ""
}

func userLogin(user string, pass string) {
	//var regexUser, _ = regexp.Compile(`^(.+)@(.+)$`)
	//var regexPass, _ = regexp.Compile(`^(?=.*[A-Z])(?=.*[a-z])(?=.*[0-9]).{8,}$`)
	//
	//userMatch := regexUser.MatchString(user)
	//passMatch := regexPass.MatchString(pass)

	//if userMatch == true && passMatch == true {
	check := CheckLogin(user, pass)
	notice := check.Notice
	data = userData{Username: check.User.Username, Password: check.User.Password,
		Nama: check.User.Nama, Umur: check.User.Umur, Status: check.User.Status}

	if notice == "Sukses" {
		menu()
	} else {
		fmt.Println(notice)
	}
	//} else {
	//	fmt.Println("Maaf, regex salah.")
	//}
}

func menu() {
	var menu int
	fmt.Println("WELCOME TO RUANG GULU")
	fmt.Println("Pilih service apa yang anda ingin gunakan: ")
	fmt.Println("1. Siswa")
	fmt.Println("2. Buat Soal")
	fmt.Scan(&menu)
	switch menu {
	case 1:
		menuSiswa()
	case 2:
		buatSoal()
	default:
		fmt.Println("Mohon pilih menu yang benar")
	}
}

//=================================METHOD UNTUK SISWA==================

//Hasil untuk mendapatkan nilai
func Hasil(tingkat string, kelas string, mapel string, temaSoal string, kumpulanJawaban []string) *pb.Result {
	conn := ConnectServer()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	hasilJawaban := pb.KumpulanJawaban{Tingkat: tingkat, Kelas: kelas, Mapel: mapel,
		TemaSoal: temaSoal, JawabanList: kumpulanJawaban}

	hasil, err := conn.Hasil(ctx, &hasilJawaban)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return hasil
}

func getSoal(tingkat string, kelas string, mapel string, temaSoal string) *pb.SoalList {
	conn := ConnectServer()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pilihanSoal := pb.PilihanSoal{Tingkat: tingkat, Kelas: kelas, Mapel: mapel, TemaSoal: temaSoal}

	soalList, err := conn.Soal(ctx, &pilihanSoal)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return soalList
}

func getTemaList(tingkat string, kelas string, matpel string) *pb.TemaList {
	conn := ConnectServer()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pilihTema := pb.PilihTema{Tingkat: tingkat, Kelas: kelas, Matpel: matpel}

	TemaList, err := conn.AmbilTema(ctx, &pilihTema)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return TemaList
}

func getMapelList(tingkat string, kelas string) *pb.MapelList {
	conn := ConnectServer()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pilihMapel := pb.PilihMapel{Tingkat: tingkat, Kelas: kelas}

	MapelList, err := conn.AmbilMapel(ctx, &pilihMapel)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return MapelList
}

func getKelasList(tingkat string) *pb.KelasList {
	conn := ConnectServer()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pilihKelas := pb.PilihKelas{Tingkat: tingkat}

	KelasList, err := conn.AmbilKelas(ctx, &pilihKelas)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return KelasList
}

func menuSiswa() {
	var menuTingkat int = 0

	for menuTingkat != 5 {
		var menuKelasSD int = 0
		var menuKelasSMP int = 0
		var menuKelasSMA int = 0
		var menuKelasUmum int = 0

		fmt.Println("\n====================")
		fmt.Println("\tE-Learning")
		fmt.Println("1. SD")
		fmt.Println("2. SMP")
		fmt.Println("3. SMA")
		fmt.Println("4. Umum")
		fmt.Println("5. Keluar")
		fmt.Print("Pilih tingkat :\t")
		fmt.Scanf("%d", &menuTingkat)
		fmt.Println("\nTingkat yang dipilih", menuTingkat)
		fmt.Println()

		switch menuTingkat {
		case 1:
			nomorterakhirKelas := 1
			for menuKelasSD != nomorterakhirKelas {
				nomorterakhirKelas = 1
				var menuMapelSD int = 0
				tingkat := "SD"
				kelasList := getKelasList(tingkat)
				arrayKelas := kelasList.KelasList
				sort.Strings(arrayKelas)
				fmt.Println("\nTINGKAT > KELAS ")
				fmt.Println("\tSD")
				for j := 0; j < len(arrayKelas); j++ {
					stringnomor := strconv.Itoa(j + 1)
					fmt.Println(stringnomor + ". " + arrayKelas[j])
					nomorterakhirKelas++
				}
				stringnomor := strconv.Itoa(nomorterakhirKelas)
				fmt.Println(stringnomor + ". Kembali")
				fmt.Print("Pilih kelas :\t")
				fmt.Scanf("%d", &menuKelasSD)

				if menuKelasSD <= len(arrayKelas) {

					fmt.Println("\nKelas yang dipilih", menuKelasSD)
					fmt.Println()

					kelas := arrayKelas[menuKelasSD-1]
					nomorterakhirMapel := 1

					for menuMapelSD != nomorterakhirMapel {

						nomorterakhirMapel = 1
						var menuTemaSoalSD int = 0
						mapelList := getMapelList(tingkat, kelas)
						arrayMapel := mapelList.MapelList
						sort.Strings(arrayMapel)
						fmt.Println("\nTINGKAT > KELAS > MAPEL ")
						fmt.Println("\tMata Pelajaran")
						for j := 0; j < len(arrayMapel); j++ {
							stringnomor := strconv.Itoa(j + 1)
							fmt.Println(stringnomor + ". " + arrayMapel[j])
							nomorterakhirMapel++
						}
						stringnomor := strconv.Itoa(nomorterakhirMapel)
						fmt.Println(stringnomor + ". Kembali")
						fmt.Print("Pilih mapel :\t")
						fmt.Scanf("%d", &menuMapelSD)

						if menuMapelSD <= len(arrayMapel) {
							fmt.Println("\nMapel yang dipilih", menuMapelSD)
							fmt.Println()
							mapel := arrayMapel[menuMapelSD-1]
							nomorterakhirTema := 1
							for menuTemaSoalSD != nomorterakhirTema {

								nomorterakhirTema = 1
								temaList := getTemaList(tingkat, kelas, mapel)
								arrayTema := temaList.TemaList
								sort.Strings(arrayTema)
								fmt.Println("\nTINGKAT > KELAS > MAPEL > TEMA")
								fmt.Println("\tTema Soal")

								for j := 0; j < len(arrayTema); j++ {
									stringnomor := strconv.Itoa(j + 1)
									fmt.Println(stringnomor + ". " + arrayTema[j])
									nomorterakhirTema++
								}
								stringnomor := strconv.Itoa(nomorterakhirTema)
								fmt.Println(stringnomor + ". Kembali")
								fmt.Print("Pilih mapel :\t")
								fmt.Scanf("%d", &menuTemaSoalSD)

								if menuTemaSoalSD <= len(arrayTema) {
									fmt.Println("\nMapel yang dipilih", menuTemaSoalSD)
									fmt.Println()
									temaSoal := arrayTema[menuTemaSoalSD-1]

									soal := getSoal(tingkat, kelas, mapel, temaSoal)

									kumpulanJawaban := showAndGet(soal)

									hasil := Hasil(tingkat, kelas, mapel, temaSoal, kumpulanJawaban)

									fmt.Println(hasil.Result)

									jawaban = []string{}

								} else if menuTemaSoalSD > (len(arrayTema) + 1) {
									fmt.Println("Mohon input nomor tema soal yang benar")
								}
							}
						} else if menuMapelSD > (len(arrayMapel) + 1) {
							fmt.Println("Mohon input nomor mata pelajaran yang benar")
						}
					}
				} else if menuKelasSD > (len(arrayKelas) + 1) {
					fmt.Println("Mohon input kelas yang benar")

				}
			}
		case 2:
			nomorterakhirKelas := 1
			for menuKelasSMP != nomorterakhirKelas {
				nomorterakhirKelas = 1
				var menuMapelSMP int = 0
				tingkat := "SMP"
				kelasList := getKelasList(tingkat)
				arrayKelas := kelasList.KelasList
				sort.Strings(arrayKelas)
				fmt.Println("\n\tSMP")
				for j := 0; j < len(arrayKelas); j++ {
					stringnomor := strconv.Itoa(j + 1)
					fmt.Println(stringnomor + ". " + arrayKelas[j])
					nomorterakhirKelas++
				}
				stringnomor := strconv.Itoa(nomorterakhirKelas)
				fmt.Println(stringnomor + ". Kembali")
				fmt.Print("Pilih kelas :\t")
				fmt.Scanf("%d", &menuKelasSMP)

				if menuKelasSMP <= len(arrayKelas) {

					fmt.Println("\nKelas yang dipilih", menuKelasSMP)
					fmt.Println()

					kelas := arrayKelas[menuKelasSMP-1]
					nomorterakhirMapel := 1

					for menuMapelSMP != nomorterakhirMapel {

						nomorterakhirMapel = 1
						var menuTemaSoalSMP int = 0
						mapelList := getMapelList(tingkat, kelas)
						arrayMapel := mapelList.MapelList
						sort.Strings(arrayMapel)
						fmt.Println("\n\tMata Pelajaran")
						for j := 0; j < len(arrayMapel); j++ {
							stringnomor := strconv.Itoa(j + 1)
							fmt.Println(stringnomor + ". " + arrayMapel[j])
							nomorterakhirMapel++
						}
						stringnomor := strconv.Itoa(nomorterakhirMapel)
						fmt.Println(stringnomor + ". Kembali")
						fmt.Print("Pilih mapel :\t")
						fmt.Scanf("%d", &menuMapelSMP)

						if menuMapelSMP <= len(arrayMapel) {
							fmt.Println("\nMapel yang dipilih", menuMapelSMP)
							fmt.Println()
							mapel := arrayMapel[menuMapelSMP-1]
							nomorterakhirTema := 1
							for menuTemaSoalSMP != nomorterakhirTema {

								nomorterakhirTema = 1
								temaList := getTemaList(tingkat, kelas, mapel)
								arrayTema := temaList.TemaList
								sort.Strings(arrayTema)
								fmt.Println("\n\tTema Soal")

								for j := 0; j < len(arrayTema); j++ {
									stringnomor := strconv.Itoa(j + 1)
									fmt.Println(stringnomor + ". " + arrayTema[j])
									nomorterakhirTema++
								}
								stringnomor := strconv.Itoa(nomorterakhirTema)
								fmt.Println(stringnomor + ". Kembali")
								fmt.Print("Pilih mapel :\t")
								fmt.Scanf("%d", &menuTemaSoalSMP)

								if menuTemaSoalSMP <= len(arrayTema) {
									fmt.Println("\nMapel yang dipilih", menuTemaSoalSMP)
									fmt.Println()
									temaSoal := arrayTema[menuTemaSoalSMP-1]

									soal := getSoal(tingkat, kelas, mapel, temaSoal)

									kumpulanJawaban := showAndGet(soal)

									hasil := Hasil(tingkat, kelas, mapel, temaSoal, kumpulanJawaban)

									fmt.Println(hasil.Result)

									jawaban = []string{}

								} else if menuTemaSoalSMP > (len(arrayTema) + 1) {
									fmt.Println("Mohon input nomor tema soal yang benar")
								}
							}
						} else if menuMapelSMP > (len(arrayMapel) + 1) {
							fmt.Println("Mohon input nomor mata pelajaran yang benar")
						}
					}
				} else if menuKelasSMP > (len(arrayKelas) + 1) {
					fmt.Println("Mohon input kelas yang benar")

				}
			}
		case 3:
			nomorterakhirKelas := 1
			for menuKelasSMA != nomorterakhirKelas {
				nomorterakhirKelas = 1
				var menuMapelSMA int = 0
				tingkat := "SMA"
				kelasList := getKelasList(tingkat)
				arrayKelas := kelasList.KelasList
				sort.Strings(arrayKelas)
				fmt.Println("\n\tSMA")
				for j := 0; j < len(arrayKelas); j++ {
					stringnomor := strconv.Itoa(j + 1)
					fmt.Println(stringnomor + ". " + arrayKelas[j])
					nomorterakhirKelas++
				}
				stringnomor := strconv.Itoa(nomorterakhirKelas)
				fmt.Println(stringnomor + ". Kembali")
				fmt.Print("Pilih kelas :\t")
				fmt.Scanf("%d", &menuKelasSMA)

				if menuKelasSMA <= len(arrayKelas) {

					fmt.Println("\nKelas yang dipilih", menuKelasSMA)
					fmt.Println()

					kelas := arrayKelas[menuKelasSMA-1]
					nomorterakhirMapel := 1

					for menuMapelSMA != nomorterakhirMapel {

						nomorterakhirMapel = 1
						var menuTemaSoalSMA int = 0
						mapelList := getMapelList(tingkat, kelas)
						arrayMapel := mapelList.MapelList
						sort.Strings(arrayMapel)
						fmt.Println("\n\tMata Pelajaran")
						for j := 0; j < len(arrayMapel); j++ {
							stringnomor := strconv.Itoa(j + 1)
							fmt.Println(stringnomor + ". " + arrayMapel[j])
							nomorterakhirMapel++
						}
						stringnomor := strconv.Itoa(nomorterakhirMapel)
						fmt.Println(stringnomor + ". Kembali")
						fmt.Print("Pilih mapel :\t")
						fmt.Scanf("%d", &menuMapelSMA)

						if menuMapelSMA <= len(arrayMapel) {
							fmt.Println("\nMapel yang dipilih", menuMapelSMA)
							fmt.Println()
							mapel := arrayMapel[menuMapelSMA-1]
							nomorterakhirTema := 1
							for menuTemaSoalSMA != nomorterakhirTema {

								nomorterakhirTema = 1
								temaList := getTemaList(tingkat, kelas, mapel)
								arrayTema := temaList.TemaList
								sort.Strings(arrayTema)
								fmt.Println("\n\tTema Soal")

								for j := 0; j < len(arrayTema); j++ {
									stringnomor := strconv.Itoa(j + 1)
									fmt.Println(stringnomor + ". " + arrayTema[j])
									nomorterakhirTema++
								}
								stringnomor := strconv.Itoa(nomorterakhirTema)
								fmt.Println(stringnomor + ". Kembali")
								fmt.Print("Pilih mapel :\t")
								fmt.Scanf("%d", &menuTemaSoalSMA)

								if menuTemaSoalSMA <= len(arrayTema) {
									fmt.Println("\nMapel yang dipilih", menuTemaSoalSMA)
									fmt.Println()
									temaSoal := arrayTema[menuTemaSoalSMA-1]

									soal := getSoal(tingkat, kelas, mapel, temaSoal)

									kumpulanJawaban := showAndGet(soal)

									hasil := Hasil(tingkat, kelas, mapel, temaSoal, kumpulanJawaban)

									fmt.Println(hasil.Result)

									jawaban = []string{}

								} else if menuTemaSoalSMA > (len(arrayTema) + 1) {
									fmt.Println("Mohon input nomor tema soal yang benar")
								}
							}
						} else if menuMapelSMA > (len(arrayMapel) + 1) {
							fmt.Println("Mohon input nomor mata pelajaran yang benar")
						}
					}
				} else if menuKelasSMA > (len(arrayKelas) + 1) {
					fmt.Println("Mohon input kelas yang benar")

				}
			}
		case 4:
			nomorterakhirKelas := 1
			for menuKelasUmum != nomorterakhirKelas {
				nomorterakhirKelas = 1
				var menuMapelUmum int = 0
				tingkat := "Umum"
				kelasList := getKelasList(tingkat)
				arrayKelas := kelasList.KelasList
				sort.Strings(arrayKelas)
				fmt.Println("\n\tUmum")
				for j := 0; j < len(arrayKelas); j++ {
					stringnomor := strconv.Itoa(j + 1)
					fmt.Println(stringnomor + ". " + arrayKelas[j])
					nomorterakhirKelas++
				}
				stringnomor := strconv.Itoa(nomorterakhirKelas)
				fmt.Println(stringnomor + ". Kembali")
				fmt.Print("Pilih kelas :\t")
				fmt.Scanf("%d", &menuKelasUmum)

				if menuKelasUmum <= len(arrayKelas) {

					fmt.Println("\nKelas yang dipilih", menuKelasUmum)
					fmt.Println()

					kelas := arrayKelas[menuKelasUmum-1]
					nomorterakhirMapel := 1

					for menuMapelUmum != nomorterakhirMapel {

						nomorterakhirMapel = 1
						var menuTemaSoalUmum int = 0
						mapelList := getMapelList(tingkat, kelas)
						arrayMapel := mapelList.MapelList
						sort.Strings(arrayMapel)
						fmt.Println("\n\tMata Pelajaran")
						for j := 0; j < len(arrayMapel); j++ {
							stringnomor := strconv.Itoa(j + 1)
							fmt.Println(stringnomor + ". " + arrayMapel[j])
							nomorterakhirMapel++
						}
						stringnomor := strconv.Itoa(nomorterakhirMapel)
						fmt.Println(stringnomor + ". Kembali")
						fmt.Print("Pilih mapel :\t")
						fmt.Scanf("%d", &menuMapelUmum)

						if menuMapelUmum <= len(arrayMapel) {
							fmt.Println("\nMapel yang dipilih", menuMapelUmum)
							fmt.Println()
							mapel := arrayMapel[menuMapelUmum-1]
							nomorterakhirTema := 1
							for menuTemaSoalUmum != nomorterakhirTema {

								nomorterakhirTema = 1
								temaList := getTemaList(tingkat, kelas, mapel)
								arrayTema := temaList.TemaList
								sort.Strings(arrayTema)
								fmt.Println("\n\tTema Soal")

								for j := 0; j < len(arrayTema); j++ {
									stringnomor := strconv.Itoa(j + 1)
									fmt.Println(stringnomor + ". " + arrayTema[j])
									nomorterakhirTema++
								}
								stringnomor := strconv.Itoa(nomorterakhirTema)
								fmt.Println(stringnomor + ". Kembali")
								fmt.Print("Pilih mapel :\t")
								fmt.Scanf("%d", &menuTemaSoalUmum)

								if menuTemaSoalUmum <= len(arrayTema) {
									fmt.Println("\nMapel yang dipilih", menuTemaSoalUmum)
									fmt.Println()
									temaSoal := arrayTema[menuTemaSoalUmum-1]

									soal := getSoal(tingkat, kelas, mapel, temaSoal)

									kumpulanJawaban := showAndGet(soal)

									hasil := Hasil(tingkat, kelas, mapel, temaSoal, kumpulanJawaban)

									fmt.Println(hasil.Result)

									jawaban = []string{}

								} else if menuTemaSoalUmum > (len(arrayTema) + 1) {
									fmt.Println("Mohon input nomor tema soal yang benar")
								}
							}
						} else if menuMapelUmum > (len(arrayMapel) + 1) {
							fmt.Println("Mohon input nomor mata pelajaran yang benar")
						}
					}
				} else if menuKelasUmum > (len(arrayKelas) + 1) {
					fmt.Println("Mohon input kelas yang benar")

				}
			}
		case 5:
			fmt.Println("Terimakasih")
		default:
			fmt.Print("Maaf, no pilihan yang anda masukkan salah.")
		}
	}
}

//============================================================================

//=================================MEMBUAT SOAL============================

//KirimSoalBuatan kirim Paket Soal yang dibuat untuk disimpan di server
func KirimSoalBuatan(paketSoal PaketSoal) string {
	guru := ConnectServerGuru()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	convertedPS := &pb.PaketSoal{
		Tingkat: "",
		Kelas:   "",
		Matpel:  "",
		Tema:    "",
		Soal:    []*pb.Soal{},
	}
	convertedPS.Tingkat = paketSoal.Tingkat
	convertedPS.Kelas = paketSoal.Kelas
	convertedPS.Matpel = paketSoal.Matpel
	convertedPS.Tema = paketSoal.Tema

	var result2 []*pb.Soal
	for i := 0; i < len(paketSoal.Soal); i++ {
		tempsoal := pb.Soal{Pertanyaan: paketSoal.Soal[i].Pertanyaan, Jawaban: paketSoal.Soal[i].Jawaban, Pilihan: paketSoal.Soal[i].Pilihan}
		result2 = append(result2, &tempsoal)
	}
	convertedPS.Soal = result2

	msg, err := guru.KirimSoalKeServer(ctx, convertedPS)
	if err != nil {
		log.Fatalf("Tidak bisa terhubung: %v", err)
	}
	fmt.Println("Laporan", msg.GetMessage())
	return msg.GetMessage()
}

func buatArraySoal() []Soal {
	reader := bufio.NewReader(os.Stdin)
	//loop input soal. stop ketika guru input "selesai"
	var kumpulansoal []Soal
	var soal, jawaban, pilihan string
	var arraypilihan []string
	var jumlahsoal int

	for jumlahsoal >= 0 {
		fmt.Println("Jumlah soal yang sudah anda buat :", jumlahsoal, "Ketik \"selesai\" jika sudah")
		fmt.Println("Input soal")
		soal, _ = reader.ReadString('\n')
		soal = strings.Replace(soal, "\n", "", -1)
		if soal == "selesai" {
			if jumlahsoal == 0 {
				fmt.Println("Anda belum membuat soal sama sekali... kembali ke menu...")
				main()
			}
			fmt.Println("Input selesai")
			return kumpulansoal
		}
		fmt.Println("Input jawaban")
		jawaban, _ = reader.ReadString('\n')
		jawaban = strings.Replace(jawaban, "\n", "", -1)
		arraypilihan = append(arraypilihan, jawaban)

		for i := 0; i < 3; i++ {

			fmt.Println("Input pilihan ", i+1)
			pilihan, _ = reader.ReadString('\n')
			pilihan = strings.Replace(pilihan, "\n", "", -1)

			//handle kalau pilihan yang diinput sama dengan yang sebelum"nya
		CHECK:
			for _, a := range arraypilihan {
				if a == pilihan {
					fmt.Println("Input pilihan ada yang sama. Mohon ganti pilihan anda.")
					pilihan, _ = reader.ReadString('\n')
					pilihan = strings.Replace(pilihan, "\n", "", -1)
					goto CHECK
				}
			}
			arraypilihan = append(arraypilihan, pilihan)
		}
		tempsoal := Soal{soal, jawaban, arraypilihan}
		arraypilihan = []string{}
		kumpulansoal = append(kumpulansoal, tempsoal)
		jumlahsoal++
	}
	return kumpulansoal
}

func buatSoal() {
	reader := bufio.NewReader(os.Stdin)
INPUTTINGKAT:
	var tingkat, kelas, mapel, tema string
	var menuTingkat, menuKelas int
	fmt.Println("Pilih tingkat")
	fmt.Println("1. SD")
	fmt.Println("2. SMP")
	fmt.Println("3. SMA")
	fmt.Println("4. Umum")
	fmt.Println("5. Kembali")
	fmt.Scan(&menuTingkat)

	switch menuTingkat {
	case 1:
		tingkat = "SD"
	INPUTKELASSD:
		fmt.Println("Pilih Kelas")
		fmt.Println("1. Kelas I")
		fmt.Println("2. Kelas II")
		fmt.Println("3. Kelas III")
		fmt.Println("4. Kelas IV")
		fmt.Println("5. Kelas V")
		fmt.Println("6. Kelas VI")
		fmt.Println("7. Kembali")
		fmt.Scan(&menuKelas)

		switch menuKelas {
		case 1:
			kelas = "Kelas 1"
		case 2:
			kelas = "Kelas 2"
		case 3:
			kelas = "Kelas 3"
		case 4:
			kelas = "Kelas 4"
		case 5:
			kelas = "Kelas 5"
		case 6:
			kelas = "Kelas 6"
		case 7:
			goto INPUTTINGKAT
		default:
			fmt.Println("Input kelas salah, mohon input yang benar")
			goto INPUTKELASSD
		}
	case 2:
		tingkat = "SMP"
	INPUTKELASSMP:
		fmt.Println("Pilih Kelas")
		fmt.Println("1. Kelas VII")
		fmt.Println("2. Kelas VIII")
		fmt.Println("3. Kelas IX")
		fmt.Println("4. Kembali")
		fmt.Scan(&menuKelas)
		switch menuKelas {
		case 1:
			kelas = "Kelas 7"
		case 2:
			kelas = "Kelas 8"
		case 3:
			kelas = "Kelas 9"
		case 4:
			goto INPUTTINGKAT
		default:
			fmt.Println("Input kelas salah, mohon input yang benar")
			goto INPUTKELASSMP
		}
	case 3:
		tingkat = "SMA"
	INPUTKELASSMA:
		fmt.Println("Pilih Kelas")
		fmt.Println("1. Kelas X")
		fmt.Println("2. Kelas XI")
		fmt.Println("3. Kelas XII")
		fmt.Println("4. Kembali")
		fmt.Scan(&menuKelas)
		switch menuKelas {
		case 1:
			kelas = "Kelas 10"
		case 2:
			kelas = "Kelas 11"
		case 3:
			kelas = "Kelas 12"
		case 4:
			goto INPUTTINGKAT
		default:
			fmt.Println("Input kelas salah, mohon input yang benar")
			goto INPUTKELASSMA
		}
	case 4:
		tingkat = "Umum"
		kelas = "Umum"
	case 5:
		main()
	}

	fmt.Println("Input Matpel")
	mapel, _ = reader.ReadString('\n')
	mapel = strings.Replace(mapel, "\n", "", -1)

	fmt.Println("Input tema:")
	tema, _ = reader.ReadString('\n')
	tema = strings.Replace(tema, "\n", "", -1)

	fmt.Println(tingkat, kelas, mapel, tema)
	kumpulansoal := buatArraySoal()
	fmt.Println(kumpulansoal)
	paketSoal := PaketSoal{tingkat, kelas, mapel, tema, kumpulansoal}
	KirimSoalBuatan(paketSoal)
	main()
}

//============================================================================
//=================================METHOD PENDUKUNG==================
func getJawaban(soal []*pb.Soal, soalAcak []*pb.Soal) []string {
	for i := 0; i < len(soal); i++ {
		for j := 0; j < len(soalAcak); j++ {
			if soal[i].Pertanyaan == soalAcak[j].Pertanyaan {
				jawaban = append(jawaban, tempJawaban[j])
				break
			}
		}
	}
	return jawaban
}

func acakSoal(tempSoal []*pb.Soal) {
	rand.Seed(time.Now().UnixNano())
	for i := range tempSoal {
		j := rand.Intn(i + 1)
		tempSoal[i], tempSoal[j] = tempSoal[j], tempSoal[i]
	}
}

func acakPG(pg []string) {
	rand.Seed(time.Now().UnixNano())
	for i := range pg {
		j := rand.Intn(i + 1)
		pg[i], pg[j] = pg[j], pg[i]
	}
}

func showPG(pilihan []string) {
	for j := 0; j < len(pilihan); j++ {
		switch j {
		case 0:
			fmt.Println("a.", pilihan[j])
		case 1:
			fmt.Println("b.", pilihan[j])
		case 2:
			fmt.Println("c.", pilihan[j])
		case 3:
			fmt.Println("d.", pilihan[j])
		}
	}
}

func scanJawaban() string {
	var pilihJawaban string = ""
	scan := false
	for scan != true {
		fmt.Print("Pilih jawaban :\t")
		fmt.Scanf("%s", &pilihJawaban)
		if pilihJawaban == "a" || pilihJawaban == "b" || pilihJawaban == "c" || pilihJawaban == "d" {
			scan = true
		} else {
			fmt.Println("Maaf, anda salah input. Silahkan ulangi jawaban anda.")
		}
	}
	return pilihJawaban
}

func getTempJawaban(pilihan []string, scanJawaban string) string {
	if scanJawaban == "a" {
		jawaban := pilihan[0]
		return jawaban
	} else if scanJawaban == "b" {
		jawaban := pilihan[1]
		return jawaban
	} else if scanJawaban == "c" {
		jawaban := pilihan[2]
		return jawaban
	} else if scanJawaban == "d" {
		jawaban := pilihan[3]
		return jawaban
	}
	return ""
}

func gantiJawaban(noSoal int, soalAcak []*pb.Soal) {
	fmt.Println("Pertanyaan no", strconv.Itoa(noSoal))
	fmt.Println(soalAcak[noSoal-1].Pertanyaan)
	showPG(soalAcak[noSoal-1].Pilihan)
	scanJawaban := scanJawaban()
	jawaban := getTempJawaban(soalAcak[noSoal-1].Pilihan, scanJawaban)
	tempJawaban[noSoal-1] = jawaban
}

func showAndGet(soal *pb.SoalList) []string {
	var soalAcak []*pb.Soal
	for i := 0; i < len(soal.SoalList); i++ {
		tempSoal := pb.Soal{Pertanyaan: soal.SoalList[i].Pertanyaan,
			Jawaban: soal.SoalList[i].Jawaban, Pilihan: soal.SoalList[i].Pilihan}
		soalAcak = append(soalAcak, &tempSoal)
	}
	acakSoal(soalAcak)
	for i := 0; i < len(soalAcak); i++ {
		fmt.Println("Pertanyaan no", i+1)
		fmt.Println(soalAcak[i].Pertanyaan)

		acakPG(soalAcak[i].Pilihan)
		showPG(soalAcak[i].Pilihan)
		scanJawaban := scanJawaban()
		jawaban := getTempJawaban(soalAcak[i].Pilihan, scanJawaban)
		tempJawaban = append(tempJawaban, jawaban)
	}
	var menuSoal int = 0

	for menuSoal != 2 {
		for i := 0; i < len(soalAcak); i++ {
			fmt.Println("Pertanyaan no", i+1)
			fmt.Println(soalAcak[i].Pertanyaan)
			fmt.Println("Jawaban")
			fmt.Println(tempJawaban[i])
		}
		fmt.Println()
		fmt.Println("1. Ganti jawaban")
		fmt.Println("2. Submit")
		fmt.Print("Pilih :\t")
		fmt.Scanf("%d", &menuSoal)

		switch menuSoal {
		case 1:
			var noSoal int = 0
			fmt.Print("Ganti jawaban no soal:\t")
			fmt.Scanf("%d", &noSoal)
			if noSoal <= len(soalAcak) {
				gantiJawaban(noSoal, soalAcak)
			} else if noSoal == 0 {
				fmt.Println("No soal tidak ada.")
			} else {
				fmt.Println("No soal tidak ada.")
			}
		case 2:
			fmt.Println("Submit")
		default:
			fmt.Print("Maaf, no pilihan yang anda masukkan salah.")
		}
	}
	kumpulanJawaban := getJawaban(soal.SoalList, soalAcak)
	// fmt.Println(tempJawaban)
	tempJawaban = []string{}
	// fmt.Println(kumpulanJawaban)
	return kumpulanJawaban
}
