# 10_fitur_v105.sl - Demo Fitur Anyar SundaLang v1.0.5

cetakkeun("=== DEMO FITUR ANYAR SUNDALANG v1.0.5 ===")

# 1. Manipulasi Teks & String
tanda kalimah = "pandeglang,serang,lebak,cilegon"
tanda kota = beulah(kalimah, ",")
cetakkeun("\n1. Hasil beulah (split):")
pikeun tanda i = 0; i < panjang(kota); i = i + 1 {
    cetakkeun("   - Kota ka-" + (i + 1) + ": " + kota[i])
}

tanda gabung_kota = gabung(kota, " -> ")
cetakkeun("\n2. Hasil gabung (join):")
cetakkeun("   " + gabung_kota)

tanda basa = "abdi resep diajar python"
tanda basa_anyar = ganti(basa, "python", "sundalang")
cetakkeun("\n3. Hasil ganti (replace):")
cetakkeun("   " + basa_anyar)

cetakkeun("\n4. Pamariksaan ngandung (contains):")
cetakkeun("   Naha 'pandeglang' ngandung 'lang'? " + ngandung("pandeglang", "lang"))
cetakkeun("   Naha daftar kota ngandung 'serang'? " + ngandung(kota, "serang"))
cetakkeun("   Naha daftar kota ngandung 'jakarta'? " + ngandung(kota, "jakarta"))

# 2. Manipulasi Array
cetakkeun("\n5. Miceun elemen array (index 1 / serang):")
tanda kota_sesa = miceun(kota, 1)
cetakkeun("   Kota sanggeus miceun: " + gabung(kota_sesa, ", "))

cetakkeun("\n6. Malik tulisan jeung array (reverse):")
cetakkeun("   Malik 'sundalang' -> " + malik("sundalang"))
tanda angka = [1, 2, 3, 4, 5]
cetakkeun("   Malik [1, 2, 3, 4, 5] -> " + malik(angka))

# 3. Pungsi Matematika
cetakkeun("\n7. Matematika Anyar:")
cetakkeun("   mutlak(-99)       = " + mutlak(-99))
cetakkeun("   pangkat(2, 8)     = " + pangkat(2, 8))
cetakkeun("   panggedena(12, 85, 34)  = " + panggedena(12, 85, 34))
cetakkeun("   pangleutikna(12, 85, 34) = " + pangleutikna(12, 85, 34))

# 4. Argumen & Sistem
cetakkeun("\n8. Argumen CLI:")
tanda arg = argumen()
lamun panjang(arg) > 0 {
    cetakkeun("   Argumen nu diasupkeun: " + gabung(arg, ", "))
} lamunteu {
    cetakkeun("   (Eweuh argumen tambahan ti terminal)")
}

cetakkeun("\n=== DEMO BERES! ===")
