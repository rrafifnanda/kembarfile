# KembarFile

Aplikasi kecil buat mencari file yang isinya sama, meskipun nama atau lokasinya berbeda. KembarFile cuma membaca file untuk dibandingkan—tidak ada file yang dihapus atau dipindahkan.

## Download

Versi siap pakai ada di halaman [Releases](https://github.com/rrafifnanda/kembarfile/releases/latest).

- **Linux:** download file `.AppImage`, izinkan file untuk dieksekusi lewat Properties, lalu klik dua kali.
- **Windows:** download file `.exe`, lalu buka seperti aplikasi biasa.

## Cara pakai

1. Buka KembarFile.
2. Klik **Pilih Folder**.
3. Pilih folder yang ingin diperiksa.
4. Klik **Scan** dan tunggu hasilnya.

File dikelompokkan berdasarkan ukuran, kemudian dicocokkan dengan SHA-256. Symbolic link dilewati dan file yang tidak bisa dibaca akan ditampilkan sebagai peringatan.

## CLI

Kalau lebih nyaman lewat terminal:

```bash
go run ./cmd/kembarfile /path/ke/folder
```

Build binary CLI:

```bash
go build -o kembarfile ./cmd/kembarfile
```

## Build desktop dari source

Linux membutuhkan Go, GCC, dan library pengembangan untuk OpenGL, X11, serta Wayland. Daftar paketnya tersedia di [dokumentasi Fyne](https://docs.fyne.io/started/quick/).

```bash
go build -tags migrated_fynedo -o kembarfile-desktop ./cmd/kembarfile-desktop
```

Jalankan test dengan:

```bash
go test ./...
```
