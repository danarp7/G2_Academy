package main

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "Tugas_1/ProtoParkir"
)

const (
	address     = "localhost:50051"
)

func ConnectServer() pb.ConnectClient {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return pb.NewConnectClient(conn)
}

func parkirMasuk() *pb.Karcis {
	con := ConnectServer()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	karcis, err := con.KarcisMasuk(ctx, new(empty.Empty))
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return karcis
}

func parkirKeluar(inpIdP string, inpTipe string, inpPlatNo string) *pb.Struk {
	con := ConnectServer()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if inpTipe != "Motor" && inpTipe != "Mobil" {
		fmt.Println("Tipe kendaraan salah")
		main()
	}

	parkir := pb.Parkir{IdParkir: inpIdP, TipeKendaraan: inpTipe, PlatNo: inpPlatNo}

	struk, err := con.StrukPembayaran(ctx, &parkir)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return struk
}

func main() {
	var menu int = 0
	for menu != 3 {
		var inpIdP string
		var inpTipe string
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
			karcis := parkirMasuk()
			fmt.Println("\n==== Karcis ====")
			fmt.Println("Id Parkir : ", karcis.IdParkir)
			fmt.Println("Waktu Masuk : ", karcis.WaktuMasuk)
			fmt.Println()
		case 2:
			fmt.Println("\n==== Parkir Keluar ====")
			fmt.Print("\nID Parkir :\t")
			fmt.Scanf("%s", &inpIdP)
			fmt.Println("\n\"Motor\" atau \"Mobil\"")
			fmt.Print("Tipe Kendaraan :\t")
			fmt.Scanf("%s", &inpTipe)
			fmt.Print("Plat No :\t")
			fmt.Scanf("%s", &inpPlatNo)
			fmt.Println()

			struk := parkirKeluar(inpIdP, inpTipe, inpPlatNo)

			fmt.Println("\n==== Struk ====")
			fmt.Println("Id Parkir : ", struk.IdParkir)
			fmt.Println("Waktu Masuk : ", struk.WaktuMasuk)
			fmt.Println("Waktu Keluar : ", struk.WaktuKeluar)
			fmt.Println("Durasi : ", struk.Durasi)
			fmt.Println("Tipe Kendaraan : ", struk.TipeKendaraan)
			fmt.Println("Plat No : ", struk.PlatNo)
			fmt.Println("Biaya Parkir : ", struk.BiayaParkir)
			fmt.Println()
		default:
			fmt.Print("Terimakasih")
		}
	}
}