package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() *mongo.Client {
	var (
		err				error
		MongoDb			string
		clientOptions	*options.ClientOptions
		client			*mongo.Client
	)

	err = godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to fund .env file")
	}

	MongoDb = os.Getenv("MONGODB_URI")
	if MongoDb == "" {
		log.Fatal("MONGODB_URI not set!")
	}

	fmt.Println("MongoDB URI: ", MongoDb)

	clientOptions = options.Client().ApplyURI(MongoDb)

	client, err = mongo.Connect(clientOptions)
	if err != nil {
		return nil
	}

	return client
}

/*

1. clientOptions = options.Client().ApplyURI(MongoDb)

Fungsi: Mengonfigurasi parameter koneksi.

Cara kerja: options.Client() membuat objek konfigurasi kosong. Method .ApplyURI(MongoDb) membaca string URI MongoDB (yang diambil dari .env) lalu mem-parsing informasi di dalamnya, seperti:

- Alamat host server / cluster Atlas.
- Username dan password otentikasi.
- Opsi tambahan (seperti retryWrites=true, timeout, SSL/TLS, dsb.).

Variabel clientOptions ini menjadi "buku panduan" sebelum aplikasi benar-benar menghubungi MongoDB.

###

2. client, err = mongo.Connect(clientOptions)

Fungsi: Menginisialisasi client instance dan menyiapkan connection pool (jalur pipa komunikasi) ke MongoDB berdasarkan konfigurasi clientOptions.

Output:

- client bertipe *mongo.Client: Objek pegangan utama yang nanti dipakai oleh fungsi lain (seperti OpenCollection) untuk masuk ke database dan koleksi tertentu.

- err: Menampung error jika format konfigurasi salah atau inisialisasi driver internal gagal.

*/

func OpenCollection(collectionName string, client *mongo.Client) *mongo.Collection {
	var (
		err          	error
		databaseName	string
		collection   	*mongo.Collection

	)
	
	err = godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}

	databaseName = os.Getenv("DATABASE_NAME")
	
	fmt.Println("DATABASE_NAME: ", databaseName)

	collection = client.Database(databaseName).Collection(collectionName)
	if collection == nil {
		return nil
	}

	return collection
}

/*

Baris tersebut berfungsi untuk mendapatkan pointer/objek referensi ke tabel (koleksi data) tertentu di dalam database MongoDB, agar kamu bisa langsung menjalankan operasi data (seperti InsertOne, Find, UpdateOne, atau DeleteOne).

Hierarkinya dipecah menjadi dua tahap:

1. client.Database(databaseName)

  - Mengarahkan koneksi MongoDB (client) ke database tertentu yang namanya diambil dari environment variable databaseName (misalnya: "toko_online").

  - Menghasilkan objek *mongo.Database.

2. .Collection(collectionName)

  - Di dalam database yang sudah dipilih tadi, method ini menunjuk ke koleksi/tabel spesifik yang dioper via parameter collectionName (misalnya: "users", "products", atau "orders").

  - Menghasilkan objek *mongo.Collection.

*/