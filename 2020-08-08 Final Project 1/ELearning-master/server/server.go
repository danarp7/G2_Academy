package main

import (
	pb "E-Learning/proto"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

var kumpulanSoal = map[string]map[string]map[string][]Soal{
	"Kelas 1":  map[string]map[string][]Soal{},
	"Kelas 2":  map[string]map[string][]Soal{},
	"Kelas 3":  map[string]map[string][]Soal{},
	"Kelas 4":  map[string]map[string][]Soal{},
	"Kelas 5":  map[string]map[string][]Soal{},
	"Kelas 6":  map[string]map[string][]Soal{},
	"Kelas 7":  map[string]map[string][]Soal{},
	"Kelas 8":  map[string]map[string][]Soal{},
	"Kelas 9":  map[string]map[string][]Soal{},
	"Kelas 10": map[string]map[string][]Soal{},
	"Kelas 11": map[string]map[string][]Soal{},
	"Kelas 12": map[string]map[string][]Soal{},
	"Umum":     map[string]map[string][]Soal{},
}

var userInSession userData

//Soal dengan struct
type Soal struct {
	Pertanyaan string
	Jawaban    string
	Pilihan    []string
}

type userData struct {
	Username string
	Password string
	Nama     string
	Umur     int32
	Status   string
}

var userDataList = map[string]userData{}

type server struct {
	pb.UnimplementedConnectServer
	pb.UnimplementedGuruServer
	pb.UnimplementedDataServer
}

// type serverGuru struct {
// 	pb.UnimplementedGuruServer
// }

func certified(kelas string, mapel string, tema string, nilai string) {
	f, err := os.Create("Certified.txt")
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	d := []string{"SERTIFIKAT",
		("KELAS : " + kelas), ("MAPEL : " + mapel), ("TEMA : " + tema),
		("NILAI : " + nilai),
		("NAMA : " + userInSession.Nama),
	}

	for _, v := range d {
		fmt.Fprintln(f, v)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("file written successfully")
}

func (s *server) Soal(ctx context.Context, ps *pb.PilihanSoal) (*pb.SoalList, error) {

	var result []Soal

	if v, found := kumpulanSoal[ps.GetKelas()][ps.GetMapel()][ps.GetTemaSoal()]; found {
		result = v
	}

	//fmt.Println(result)

	var result2 []*pb.Soal
	for i := 0; i < len(result); i++ {
		tempsoal := pb.Soal{Pertanyaan: result[i].Pertanyaan, Jawaban: result[i].Jawaban, Pilihan: result[i].Pilihan}
		result2 = append(result2, &tempsoal)
	}
	//fmt.Println(result2)

	result3 := pb.SoalList{SoalList: result2}

	return &result3, nil
}

func (s *server) Hasil(ctx context.Context, kj *pb.KumpulanJawaban) (*pb.Result, error) {
	var result []Soal
	var nilai int = 0
	var correct int = 0
	var notice string
	kumpulanJawaban := kj.GetJawabanList()

	if v, found := kumpulanSoal[kj.GetKelas()][kj.GetMapel()][kj.GetTemaSoal()]; found {
		result = v
	}

	//koreksi jawaban
	for i := 0; i < len(kumpulanJawaban); i++ {
		fmt.Println(result[i].Jawaban)
		fmt.Println(kumpulanJawaban[i])
		if kumpulanJawaban[i] == result[i].Jawaban {
			correct++
		}
	}

	nilai = (100 / len(kumpulanJawaban)) * correct

	if nilai >= 70 {
		notice = "Nilai " + strconv.Itoa(nilai) + "\nSELAMAT\nKamu berhasil lulus tema soal ini.\nSilahkan cek sertifikatmu ya."
		certified(kj.GetKelas(), kj.GetMapel(), kj.GetTemaSoal(), strconv.Itoa(nilai))
	} else {
		notice = "Nilai " + strconv.Itoa(nilai) + "\nJANGAN PATAH SEMANGAT\nTetaplah berusaha."
	}
	fmt.Println(notice)

	notice2 := pb.Result{Result: notice}

	return &notice2, nil
}

func (s *server) AmbilTema(ctx context.Context, pt *pb.PilihTema) (*pb.TemaList, error) {

	var result []string

	if v, found := kumpulanSoal[pt.GetKelas()][pt.GetMatpel()]; found {

		for k, _ := range v {
			result = append(result, k)
		}

		fmt.Println(result)
	}
	temaList := pb.TemaList{TemaList: result}
	return &temaList, nil
}

func (s *server) AmbilMapel(ctx context.Context, pt *pb.PilihMapel) (*pb.MapelList, error) {

	var result []string

	if v, found := kumpulanSoal[pt.GetKelas()]; found {

		for k, _ := range v {
			result = append(result, k)
		}

		fmt.Println(result)
	}
	mapelList := pb.MapelList{MapelList: result}
	return &mapelList, nil
}

func (s *server) AmbilKelas(ctx context.Context, pt *pb.PilihKelas) (*pb.KelasList, error) {

	var result []string
	if pt.GetTingkat() == "SD" {
		result = []string{"Kelas 1", "Kelas 2", "Kelas 3", "Kelas 4", "Kelas 5", "Kelas 6"}
	} else if pt.GetTingkat() == "SMP" {
		result = []string{"Kelas 7", "Kelas 8", "Kelas 9"}
	} else if pt.GetTingkat() == "SMA" {
		result = []string{"Kelas 10", "Kelas 11", "Kelas 12"}
	} else if pt.GetTingkat() == "Umum" {
		result = []string{"Umum"}
	}

	fmt.Println(result)
	kelasList := pb.KelasList{KelasList: result}
	return &kelasList, nil
}

func (s *server) KirimSoalKeServer(ctx context.Context, req *pb.PaketSoal) (*pb.ServerGuruReply, error) {
	// fmt.Println(req.GetTingkat())
	// fmt.Println(req.GetKelas())
	// fmt.Println(req.GetMatpel())
	// fmt.Println(req.GetTema())
	// fmt.Println(req.GetSoal()[0])

	var soal []Soal
	for i := 0; i < len(req.GetSoal()); i++ {
		tempsoal := Soal{Pertanyaan: req.GetSoal()[i].Pertanyaan, Jawaban: req.GetSoal()[i].Jawaban, Pilihan: req.GetSoal()[i].Pilihan}
		soal = append(soal, tempsoal)
	}

	//check if soal exists in tingkat-kelas-matpel-tema
	//if not exists, mapnya dibuat dulu
CHECKEXISTS:
	if _, found := kumpulanSoal[req.GetKelas()]; found {
		if _, found := kumpulanSoal[req.GetKelas()][req.GetMatpel()]; found {
			if _, found := kumpulanSoal[req.GetKelas()][req.GetMatpel()]; found {
				if _, found := kumpulanSoal[req.GetKelas()][req.GetMatpel()][req.GetTema()]; found {
					//beri nama lain
					newname := req.GetTema() + " (new)"
					//save
					saveSoal(req.GetTingkat(), req.GetKelas(), req.GetMatpel(), newname, soal)
					return &pb.ServerGuruReply{Message: ("Paket Soal dengan Tema ini sudah tersedia\nSoal disimpan dengan judul " + newname)}, nil
				} else {
					fmt.Println("check 4", kumpulanSoal)
					saveSoal(req.GetTingkat(), req.GetKelas(), req.GetMatpel(), req.GetTema(), soal)
					return &pb.ServerGuruReply{Message: ("Paket Soal dengan Tema ini disimpan dengan judul " + req.GetTema())}, nil
				}
			} else {
				kumpulanSoal[req.GetKelas()][req.GetMatpel()][req.GetTema()] = []Soal{}
				fmt.Println("check 3", kumpulanSoal)
				goto CHECKEXISTS
			}
		} else {
			kumpulanSoal[req.GetKelas()][req.GetMatpel()] = map[string][]Soal{}
			fmt.Println("check 2", kumpulanSoal)
			goto CHECKEXISTS
		}
	} else {
		kumpulanSoal[req.GetKelas()] = map[string]map[string][]Soal{}
		fmt.Println("check 1", kumpulanSoal)
		goto CHECKEXISTS
	}
}
func main() {
	loadSoal()
	loadUser()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterConnectServer(s, &server{})
	pb.RegisterGuruServer(s, &server{})
	pb.RegisterDataServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}

func (s *server) CheckLogin(ctx context.Context, il *pb.InputLogin) (*pb.UserData, error) {
	tempUser := pb.User{}
	if v, found := userDataList[il.GetUsername()]; found {
		if il.GetPassword() == v.Password {
			tempUser = pb.User{Username: v.Username,
				Password: v.Password,
				Nama:     v.Nama,
				Umur:     v.Umur,
				Status:   v.Status}
			result := pb.UserData{Notice: "Sukses", User: &tempUser}
			userInSession = v
			return &result, nil
		} else {
			result := pb.UserData{Notice: "Password anda salah.", User: &tempUser}
			return &result, nil
		}
	}

	result := pb.UserData{Notice: "Username tidak terdaftar.", User: &tempUser}
	return &result, nil
}

func (s *server) NewSignUp(ctx context.Context, user *pb.User) (*pb.UserData, error) {
	tempUser := user
	if _, found := userDataList[user.GetUsername()]; found {
		result := pb.UserData{Notice: "Username sudah terdaftar. Tidak bisa mendaftar lagi", User: user}
		return &result, nil
	} else {
		saveUser(tempUser)
		result := pb.UserData{Notice: "Registrasi sukses! Silahkan login.", User: user}
		return &result, nil
	}
}

func saveUser(user *pb.User) {

	newUser := userData{Username: user.GetUsername(), Password: user.GetPassword(),
		Nama: user.GetNama(), Umur: user.GetUmur(), Status: user.GetStatus()}

	userDataList[user.GetUsername()] = newUser

}

func saveSoal(ti string, k string, m string, te string, s []Soal) {
	fmt.Println(ti, k, m, te, s)
	kumpulanSoal[k][m][te] = s
	fmt.Println(kumpulanSoal)
}

func loadSoal() {

	//Kelas 1 - Bahasa Inggris - Bab 1
	tempsoal1 := Soal{"Bahasa Inggris \"Selamat pagi\" adalah ...", "Good morning", []string{"Good morning", "Good afternoon", "Good night", "Good day"}}
	tempsoal2 := Soal{"How are you?", "I am fine, thank you.", []string{"I am fine, thank you.", "Hello!", "Goodbye!", "How about you?"}}
	tempsoal3 := Soal{"What is your name?", "My name is Budi.", []string{"My name is Budi.", "I am fine, thank you.", "I am 6 years old.", "Goodbye!"}}
	tempsoal4 := Soal{"Bahasa Inggris \"Selamat malam\" adalah ...", "Good night", []string{"Good morning", "Good afternoon", "Good night", "Good day"}}
	tempsoal5 := Soal{"How old are you?", "I am 6 years old.", []string{"My name is Budi.", "I am fine, thank you.", "I am 6 years old.", "Goodbye!"}}

	arraySoal := make([]Soal, 0)
	arraySoal = append(arraySoal, tempsoal1)
	arraySoal = append(arraySoal, tempsoal2)
	arraySoal = append(arraySoal, tempsoal3)
	arraySoal = append(arraySoal, tempsoal4)
	arraySoal = append(arraySoal, tempsoal5)

	kumpulanSoal["Kelas 1"] = map[string]map[string][]Soal{}
	kumpulanSoal["Kelas 1"]["Bahasa Inggris"] = map[string][]Soal{}
	kumpulanSoal["Kelas 1"]["Bahasa Inggris"]["Bab 1"] = arraySoal

	// Kelas 1 - Bahasa Inggris - Bab 2
	tempsoal1 = Soal{"What are you doing?", "I am studying English.", []string{"I am studying English.", "I eat bread for breakfast.", "I like my bicycle.", "I am a student."}}
	tempsoal2 = Soal{"Mom ... cooking dinner right now.", "is", []string{"is", "am", "are", "was"}}
	tempsoal3 = Soal{"My friends ... staying at home.", "are", []string{"are", "were", "is", "am"}}
	tempsoal4 = Soal{"Dad is (write)... a letter for his friend.", "writing", []string{"writing", "writeing", "writying", "writiing"}}
	tempsoal5 = Soal{"Budi and Bambang ... best friends", "are", []string{"am", "is", "was", "are"}}

	arraySoal = make([]Soal, 0)
	arraySoal = append(arraySoal, tempsoal1)
	arraySoal = append(arraySoal, tempsoal2)
	arraySoal = append(arraySoal, tempsoal3)
	arraySoal = append(arraySoal, tempsoal4)
	arraySoal = append(arraySoal, tempsoal5)

	kumpulanSoal["Kelas 1"]["Bahasa Inggris"]["Bab 2"] = arraySoal

	//Kelas 2 - Bahasa Inggris - Bab 1
	tempsoal1 = Soal{"I like airplanes, I want to be a ... ", "Pilot", []string{"Pilot", "Nurse", "Teacher", "Youtuber"}}
	tempsoal2 = Soal{"Mary likes animals, she wants to work in a ... ", "Zoo", []string{"Zoo", "School", "Hospital", "Office"}}
	tempsoal3 = Soal{"Bob likes to draw, he wants to be a ... ", "Painter", []string{"Programmer", "Engineer", "Painter", "Musician"}}
	tempsoal4 = Soal{"Gordon likes to cook, he wants to be a ...", "Chef", []string{"Chef", "Racer", "Athlete", "Programmer"}}
	tempsoal5 = Soal{"Patrick likes to act, he wants to be an ...", "Actor", []string{"Actor", "Actress", "Athlete", "Attorney"}}

	arraySoal = make([]Soal, 0)
	arraySoal = append(arraySoal, tempsoal1)
	arraySoal = append(arraySoal, tempsoal2)
	arraySoal = append(arraySoal, tempsoal3)
	arraySoal = append(arraySoal, tempsoal4)
	arraySoal = append(arraySoal, tempsoal5)

	kumpulanSoal["Kelas 2"] = map[string]map[string][]Soal{}
	kumpulanSoal["Kelas 2"]["Bahasa Inggris"] = map[string][]Soal{}
	kumpulanSoal["Kelas 2"]["Bahasa Inggris"]["Bab 1"] = arraySoal

	//Kelas 2 - Bahasa Inggris - Bab 2
	tempsoal1 = Soal{"Sam yawns all day, he is ...", "Sleepy", []string{"Sleepy", "Hungry", "Angry", "Sad"}}
	tempsoal2 = Soal{"I skipped the breakfast this morning, now I am very ...", "Hungry", []string{"Hungry", "Sleepy", "Bored", "Tired"}}
	tempsoal3 = Soal{"Mom is ... at me because I lost her tupperware", "Angry", []string{"Angry", "Happy", "Sad", "Bored"}}
	tempsoal4 = Soal{"Alice forget her water bottle, now she is ...", "Thirsty", []string{"Thirsty", "Hungry", "Sleepy", "Bored"}}
	tempsoal5 = Soal{"Psyduck uses Confuse, now I am ...", "Confused", []string{"Confused", "Sleepy", "Hungry", "Angry"}}

	arraySoal = make([]Soal, 0)
	arraySoal = append(arraySoal, tempsoal1)
	arraySoal = append(arraySoal, tempsoal2)
	arraySoal = append(arraySoal, tempsoal3)
	arraySoal = append(arraySoal, tempsoal4)
	arraySoal = append(arraySoal, tempsoal5)

	kumpulanSoal["Kelas 2"]["Bahasa Inggris"]["Bab 2"] = arraySoal

	// fmt.Println(kumpulanSoal)
}

func loadUser() {
	user1 := userData{Username: "mamat123@yahoo.co.id", Password: "*Mamat123",
		Nama: "Mamat", Umur: 14, Status: "Siswa"}
	user2 := userData{Username: "juminten123@yahoo.co.id", Password: "*Juminten123",
		Nama: "Juminten", Umur: 25, Status: "Guru"}
	user3 := userData{Username: "shir", Password: "leen",
		Nama: "Shirleen", Umur: 10, Status: "Siswa"}

	userDataList["mamat123@yahoo.co.id"] = user1
	userDataList["juminten123@yahoo.co.id"] = user2
	userDataList["shir"] = user3

	// fmt.Println(userDataList)
}
