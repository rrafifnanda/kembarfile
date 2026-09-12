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

## Aplikasi desktop Linux

Jalankan aplikasi desktop dengan:

```bash
go run -tags migrated_fynedo ./cmd/kembarfile-desktop
```

Untuk membuat binary:

```bash
go build -tags migrated_fynedo -o kembarfile-desktop ./cmd/kembarfile-desktop
```

Build desktop memerlukan Go, GCC, OpenGL, X11, Wayland, dan header pengembangan terkait. Lihat [persyaratan Linux Fyne](https://docs.fyne.io/started/quick/).

Versi siap pakai dapat diunduh dari halaman [Releases](https://github.com/rrafifnanda/kembarfile/releases). Setelah mengunduh AppImage, aktifkan izin eksekusi dari Properties jika file manager Anda belum mengaktifkannya, lalu klik dua kali untuk membuka.

## Pengembangan

Jalankan test dengan:

```bash
go test ./...
```

Kode pemindai berada di `internal/finder` dan digunakan bersama oleh CLI serta aplikasi desktop.
