package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	guuid "github.com/google/uuid"
	"log"
	"net"
	"time"

	pb "Tugas_1/ProtoParkir"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type Karcis struct {
	IdParkir  string 			`json:"id_parkir"`
	StartTime time.Time  		`json:"waktu_masuk"`
}

type Struk struct {
	IdParkir      string 		`json:"id_parkir"`
	StartTime     time.Time  	`json:"waktu_masuk"`
	EndTime       time.Time  	`json:"waktu_keluar"`
	Durasi        int32      	`json:"durasi"`
	TipeKendaraan string     	`json:"tipe_kendaraan"`
	PlatNo        string     	`json:"plat_no"`
	BiayaParkir   int32      	`json:"biaya_parkir"`
}

type Parkir struct {
	IdParkir      string 		`json:"id_parkir"`
	TipeKendaraan string     	`json:"tipe_kendaraan"`
	PlatNo        string     	`json:"plat_no"`
}

var dataParkirMasuk []Karcis
var dataParkirKeluar []Struk

type server struct {
	pb.UnimplementedConnectServer
}

func simpanDataParkirMasuk(idParkir string, startTime time.Time) (string, time.Time) {
	parkirMasuk := Karcis{IdParkir: idParkir, StartTime: startTime}
	dataParkirMasuk = append(dataParkirMasuk, parkirMasuk)
	return idParkir, startTime
}

func simpanDataParkirKeluar(idParkir string, startTime time.Time, endTime time.Time, tipeKendaraan string,
	platNo string, durasi int32, biayaParkir int32) {
	parkirKeluar := Struk{IdParkir: idParkir, StartTime: startTime, EndTime: endTime,
							Durasi: durasi, TipeKendaraan: tipeKendaraan, PlatNo: platNo,
							BiayaParkir: biayaParkir}
	dataParkirKeluar = append(dataParkirKeluar,parkirKeluar)
}

func hitungBiayaParkir(idParkir string, startTime time.Time, tipeKendaraan string, platNo string) (time.Time, int32, int32) {
	endTime := time.Now()
	if tipeKendaraan == "Motor" {
		durasi, biayaParkir := parkirMotor(startTime, endTime)
		simpanDataParkirKeluar(idParkir, startTime, endTime, tipeKendaraan, platNo, durasi, biayaParkir)
		return endTime, durasi, biayaParkir
	} else if tipeKendaraan == "Mobil" {
		durasi, biayaParkir := parkirMobil(startTime, endTime)
		simpanDataParkirKeluar(idParkir, startTime, endTime, tipeKendaraan, platNo, durasi, biayaParkir)
		return endTime, durasi, biayaParkir
	} else {
		log.Println("Tipe kendaraan yang anda masukkan salah.")
		return time.Time{}, 0, 0
	}
}

func parkirMotor(startTime time.Time, endTime time.Time) (int32, int32) {
	diff := endTime.Sub(startTime)
	durasi := int32(diff.Seconds())
	satuHari := int32(12 * 60 * 60)
	if durasi == 1 {
		biayaParkir := int32(3000)
		return durasi, biayaParkir
	} else if durasi >= satuHari{
		biayaParkir := ((satuHari - 1) * 2000) + 3000
		return durasi, biayaParkir
	} else {
		biayaParkir := ((durasi - 1) * 2000) + 3000
		return durasi, biayaParkir
	}
}

func parkirMobil(startTime time.Time, endTime time.Time) (int32, int32) {
	diff := endTime.Sub(startTime)
	durasi := int32(diff.Seconds())
	satuHari := int32(12 * 60 * 60)
	if durasi == 1 {
		biayaParkir := int32(5000)
		return durasi, biayaParkir
	} else if durasi >= satuHari{
		biayaParkir := ((satuHari - 1) * 3000) + 5000
		return durasi, biayaParkir
	} else {
		biayaParkir := ((durasi - 1) * 3000) + 5000
		return durasi, biayaParkir
	}
}

func (s *server) ConnectToServer(ctx context.Context, void *empty.Empty) (*pb.ServerReply, error) {
	return &pb.ServerReply{Message: "Client is connected to Server"}, nil
}

func (s *server) KarcisMasuk(ctx context.Context, void *empty.Empty) (*pb.Karcis, error) {
	idParkir := guuid.New().String()
	startTime := time.Now()
	id_parkir, waktu_masuk := simpanDataParkirMasuk(idParkir, startTime)

	log.Println("\n\tParkir Masuk")
	log.Println("Id Parkir : ", id_parkir)
	log.Println("Waktu Masuk : ", waktu_masuk.Format("2006-Jan-02 Monday 03:04:05"))
	fmt.Println()


	var result = pb.Karcis{IdParkir: id_parkir, WaktuMasuk: waktu_masuk.Format("2006-Jan-02 Monday 03:04:05")}
	return &result, nil
}
func (s *server) StrukPembayaran(ctx context.Context, parkir *pb.Parkir) (*pb.Struk, error) {
	var waktuMasuk string = ""
	var waktuKeluar string = ""
	var durasi int32 = 0
	var biayaParkir int32 = 0

	for i := 0; i < len(dataParkirMasuk); i++ {
		if dataParkirMasuk[i].IdParkir == parkir.GetIdParkir() {
			waktuKeluarval, durasiVal, biayaParkirval := hitungBiayaParkir(dataParkirMasuk[i].IdParkir,
				dataParkirMasuk[i].StartTime,
				parkir.GetTipeKendaraan(),
				parkir.GetPlatNo())

			waktuMasuk = dataParkirMasuk[i].StartTime.Format("2006-Jan-02 Monday 03:04:05")
			waktuKeluar = waktuKeluarval.Format("2006-Jan-02 Monday 03:04:05")
			durasi = durasiVal
			biayaParkir = biayaParkirval

			log.Println("\n\tParkir Keluar")
			log.Println("Id Parkir : ", dataParkirMasuk[i].IdParkir)
			log.Println("Waktu Masuk : ", waktuMasuk)
			log.Println("Waktu Keluar : ", waktuKeluar)
			log.Println("Durasi : ", durasi)
			log.Println("Tipe Kendaraan : ", parkir.GetTipeKendaraan())
			log.Println("Plat No : ", parkir.GetPlatNo())
			log.Println("Biaya Parkir : ", biayaParkir)
			fmt.Println()

			dataParkirMasuk = append(dataParkirMasuk[:i], dataParkirMasuk[i+1:]...)
		}
	}
	
	if biayaParkir == 0 {
		return nil, nil
	}

	var result = pb.Struk{IdParkir: parkir.GetIdParkir(),
							WaktuMasuk: waktuMasuk, WaktuKeluar: waktuKeluar,
							Durasi: durasi, TipeKendaraan: parkir.GetTipeKendaraan(),
							PlatNo: parkir.GetPlatNo(), BiayaParkir: biayaParkir}

	return &result, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterConnectServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
