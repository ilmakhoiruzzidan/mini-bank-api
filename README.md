# Mini Bank API

## Deskripsi

Mini Bank API adalah aplikasi backend sederhana untuk sistem perbankan yang memungkinkan pengguna untuk melakukan login,
logout, dan mengakses informasi akun mereka.

## Teknologi yang Digunakan

- **Go** (Golang) sebagai bahasa pemrograman utama.
- **JWT** untuk otentikasi.
- **Gin Framework** untuk routing dan middleware.
- **Godotenv** untuk memuat variabel lingkungan dari file `.env`.
- **Testify** untuk framework testing.

## Langkah-langkah Menjalankan Aplikasi

### 1. Persiapan Lingkungan Pengembangan

Pastikan Anda sudah menginstal perangkat berikut di sistem Anda:

- **Go** (Golang): [Install Go](https://golang.org/doc/install)
- **Gin Framework** untuk routing: `go get -u github.com/gin-gonic/gin`

### 2. Menyiapkan File `.env`

Buat file `.env` di root direktori proyek Anda dan pastikan berisi variabel lingkungan yang diperlukan untuk aplikasi
Anda. Contoh isi file `.env`:

```env
JWT_SECRET_KEY=testsecretkey
```

Penjelasan:

- `JWT_SECRET_KEY`: Kunci untuk menandatangani JWT yang digunakan dalam otentikasi.

### 3. Menjalankan Aplikasi Secara Lokal

1. Pastikan Anda berada di direktori root proyek.
2. Install dependensi yang diperlukan dengan perintah:

```bash
go mod tidy
```

3. Setelah itu, Anda dapat menjalankan aplikasi dengan perintah berikut:

```bash
go run main.go
```

Aplikasi akan berjalan pada port default, misalnya `localhost:8080`.

### 4. Menjalankan Tes

Untuk menjalankan tes unit yang telah ditulis dengan menggunakan framework `Testify`, Anda bisa menggunakan perintah
berikut:

```bash
go test ./...
```

Ini akan menjalankan semua tes yang ada di dalam direktori proyek.

### 5. Dokumentasi API

Berikut adalah daftar endpoint API yang tersedia:

#### **1. Authentication (Login & Logout)**

- **POST /auth/login**
    - Deskripsi: Mengautentikasi pengguna dan memberikan JWT untuk akses lebih lanjut.
    - Input:
      ```json
      {
        "username": "username_example",
        "password": "password_example"
      }
      ```
    - Output:
      ```json
      {
        "access_token": "JWT_ACCESS_TOKEN"
      }
      ```

- **POST /auth/logout** _(protected)_
    - Deskripsi: Menghapus token JWT dari sesi pengguna (logout).
    - Header:
        - Authorization: `Bearer <access_token>`

- **GET /auth/me** _(protected)_
    - Deskripsi: Mendapatkan informasi pengguna yang sedang login menggunakan JWT.
    - Header:
        - Authorization: `Bearer <access_token>`
    - Output:
      ```json
      {
        "id": 1,
        "username": "username_example"
      }
      ```

#### **2. Transaction**

- **POST /transactions** _(protected)_
    - Deskripsi: Membuat transaksi baru.
    - Header:
        - Authorization: `Bearer <access_token>`
    - Input:
      ```json
      {
        "amount": 100,
        "description": "Transfer to savings"
      }
      ```
    - Output:
      ```json
      {
        "transaction_id": 1,
        "amount": 100,
        "description": "Transfer to savings"
      }
      ```

### 6. Struktur Direktori

Berikut adalah struktur direktori utama proyek ini:

```
/mini-bank-api
│
├── /middlewares         # Middleware untuk otentikasi JWT
├── /handlers            # Pengelola untuk routing dan logic API
├── /models              # Model untuk data (misalnya Customer, Transaction)
├── /services            # Service bisnis logic
├── /router              # Routing enpoints
├── /utils               # Utilitas umum (misalnya untuk JWT)
├── /tests               # Testing
├── main.go              # File utama aplikasi
├── .env                 # File untuk variabel lingkungan
└── README.md            # Dokumentasi proyek
```
