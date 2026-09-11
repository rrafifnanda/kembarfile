# KembarFile

KembarFile adalah CLI sederhana untuk mencari file duplikat di dalam sebuah folder. File dibandingkan berdasarkan ukuran dan hash SHA-256, jadi nama file boleh berbeda.

Program ini hanya membaca dan menampilkan hasil. Tidak ada file yang dihapus atau dipindahkan.

## Fitur

- Memindai folder beserta seluruh subfolder
- Mengabaikan symbolic link dan file non-reguler
- Menampilkan kelompok file duplikat
- Menghitung ruang penyimpanan yang bisa dibebaskan
- Tetap melanjutkan pemindaian ketika sebagian file tidak dapat dibaca

## Menjalankan

Pastikan Go sudah terpasang, lalu jalankan:

```bash
go run ./cmd/kembarfile /path/ke/folder
```

Atau build menjadi binary:

```bash
go build -o kembarfile ./cmd/kembarfile
./kembarfile /path/ke/folder
```

Contoh hasil:

```text
Group 1 (1024 bytes):
  /home/user/foto/copy.jpg
  /home/user/foto/original.jpg
1 duplicate group(s), 1024 bytes reclaimable
```

## Pengembangan

Jalankan test dengan:

```bash
go test ./...
```

Kode pemindai berada di `internal/finder`, terpisah dari CLI agar bisa digunakan kembali saat versi desktop dibuat.
