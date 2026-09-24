<p align="center">
  <img src="https://i.imgur.com/J8l3YP3.png" alt="Logo SundaLang" width="180"/>
</p>

<h1 align="center">SundaLang</h1>
<p align="center">
  <strong>Bahasa Pemrograman Sunda Pandeglang</strong><br>
  <em>Sederhana, Modern, Nyunda Pisan</em>
</p>

<p align="center">
  <a href="https://sundalang.emandev.xyz">
    <img src="https://img.shields.io/badge/web-official-blue" alt="web" />
  </a>
  <a href="https://raw.githubusercontent.com/bromanprjkt/sundalang/refs/heads/main/LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License: MIT" />
  </a>
  <a href="https://github.com/bromanprjkt/sundalang">
    <img src="https://img.shields.io/github/stars/bromanprjkt/sundalang?style=social" alt="GitHub stars" />
  </a>
  <a href="https://github.com/bromanprjkt/sundalang/releases">
    <img src="https://img.shields.io/github/v/release/bromanprjkt/sundalang?label=release" alt="Latest Release" />
  </a>
</p>

<p align="center">
  <a href="#-fitur-utama">Fitur</a> •
  <a href="#-kamus-syntax">Kamus</a> •
  <a href="#-cara-install">Install</a> •
  <a href="#-cara-pakai">Cara Pakai</a>
</p>

> **"Diajar koding bari ngamumule bahasa indung."**

**SundaLang** (Sunda Pandeglang Programming Language) adalah bahasa pemrograman esoteric (esolang) berbasis interpreter yang menggunakan kosa kata asli **Bahasa Sunda Pandeglang**, Banten. Dibuat untuk edukasi, melatih logika berpikir, sekaligus ngamumule budaya lokal lewat coding.

## ✨ Fitur Utama

- **Sederhana** → Sintaks mirip Python/Go, mudah dipelajari pemula.
- **Nyunda 100%** → Keyword pakai dialek Pandeglang (`kedap`, `lamun`, `cetakkeun`, dsb).
- **Ringan & Cepat** → Dibangun dengan Go, tanpa dependency berat.
- **Lengkap** → Variabel, I/O, File Handling, Error Handling (`cobaan-sanya`), Array, Map (`wadah`).
## 📚 Kamus Syntax

Berikut adalah daftar padanan kata kunci dalam SundaLang:

| SundaLang | Konsep Umum | Keterangan |
| :--- | :--- | :--- |
| `tanda` | `var` / `let` | Deklarasi variabel |
| `tetep` | `const` | Konstanta |
| `cetakkeun` | `print` | Mencetak output |
| `tanyakeun` | `input` | Meminta input user |
| `lamun` | `if` | Percabangan |
| `lamunteu` | `else` | Percabangan alternatif |
| `milih` | `switch` | Seleksi kondisi |
| `kasus` | `case` | Kasus dalam switch |
| `baku` | `default` | Default dalam switch |
| `pikeun` | `for` | Perulangan |
| `kedap` | `while` | Perulangan kondisi |
| `eureun` | `break` | Berhenti dari loop |
| `tuluy` | `continue` | Lanjut loop |
| `pungsi` | `function` | Definisi fungsi |
| `balik` / `balikkeun` | `return` | Mengembalikan nilai |
| `buka` | `import` | Mengimpor file lain |
| `wadah` | `map` / `struct` | Struktur data key-value |
| `cobaan` | `try` | Mencoba blok kode |
| `sanya` | `catch` | Menangkap error |
| `bener` | `true` | Boolean benar |
| `salah` | `false` | Boolean salah |
| `ewehan` | `null` | Nilai kosong |

### Fungsi Bawaan (Built-in)

- `panjang(obj)`: Menghitung panjang string/array.
- `mimiti(array)`: Mengambil elemen pertama array.
- `tungtung(array)`: Mengambil elemen terakhir array.
- `asupkeun(array, item)`: Menambahkan item ke array.
- `garede(str)`: Uppercase string.
- `laleutik(str)`: Lowercase string.
- `kana_angka(str)`: Convert ke integer.
- `kana_tulisan(obj)`: Convert ke string.
- `tipe(obj)`: Cek tipe data.
- `waktu()`: Cek waktu sekarang (jam:menit:detik).
- `acak(max)`: Generate angka acak 0 s.d max.
- `sare(ms)` / `reureuh(ms)`: Pause eksekusi (sleep).
- `maca(path)`: Baca isi file.
- `nyerat(path, content)`: Tulis ke file.

## 📦 Cara Install

### ⚡ Quick Install (Recommended)

Download dan jalankan installer otomatis dengan one-liner:

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/bromanprjkt/sundalang/main/install.ps1 | iex
```

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/bromanprjkt/sundalang/main/install.sh | bash
```

Script akan otomatis:
- Download binary terbaru (v1.0.4+)
- Install ke `~/.sundalang/bin/`
- Tambahkan ke PATH

### 🎯 Manual Install via Binary

1. Download binary `sundalang.exe` (Windows) ataw `sundalang` (Linux/macOS) dari [GitHub Releases](https://github.com/bromanprjkt/sundalang/releases).
2. Jalankan installer built-in:
   ```bash
   ./sundalang install
   ```

## 🚀 Cara Pakai

**Jalankan file .sl:**
```bash
sundalang 01_dasar.sl
```

**Interactive Launcher:**
Cukup ketik `sundalang` di terminal untuk membuka menu interaktif (Install, Uninstall, REPL, Run).

**Cek Versi:**
```bash
sundalang --version
```

## 💻 Contoh Kode

```javascript
// Halo Dunya di SundaLang
cetakkeun("Sampurasun Dunya!")

// Variabel
tanda ngaran = "Kabayan"
tanda umur = 25

// Percabangan
lamun (umur >= 17) {
    cetakkeun(ngaran + " geus dewasa")
} lamunteu {
    cetakkeun(ngaran + " budak keneh")
}

// Perulangan
pikeun tanda i = 0; i < 5; i = i + 1 {
    cetakkeun("Looping ka-" + i)
}
```

## 🤝 Kontribusi

Kami sangat senang kalau baraya rek nyumbangkeun!

1. Fork repo ini
2. Buat branch baru (`git checkout -b fitur-anyar`)
3. Commit (`git commit -m 'Nambihan fitur mantap'`)
4. Push & buat Pull Request

## 📄 Lisensi

Dirilis di bawah **MIT License**. Mangga dianggo, dimodifikasi, disebarluaskan saginana.
WARNING: Gunakan dengan bijak, jangan dipakai nyieun malware (kualat).
