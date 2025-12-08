Test ONLINE Mode (2 terminal)

   Terminal 1 - Cloud API:

   bash
     cd /home/fad/Documents/myProject/shosha/shosha_desktop
     make dev-cloud

   Terminal 2 - Frontend:

   bash
     cd /home/fad/Documents/myProject/shosha/shosha_desktop
     make dev-frontend

   Expected:
   •  Sidebar menampilkan "Online" (hijau)
   •  Data sync ke cloud setiap 30 detik
   •  "belum sync" count = 0 setelah sync berhasil

   ──────────────────────────────────────────

   Test OFFLINE Mode (1 terminal)

   bash
     cd /home/fad/Documents/myProject/shosha/shosha_desktop
     make dev-offline

   Expected:
   •  Sidebar menampilkan "Offline" (kuning)
   •  Data tersimpan lokal saja
   •  "belum sync" count bertambah saat input transaksi

   ──────────────────────────────────────────

   Test Reconnect

   1. Jalankan make dev-offline → input beberapa transaksi → lihat "belum sync" count naik
   2. Stop aplikasi (Ctrl+C)
   3. Jalankan make dev-cloud di terminal 1
   4. Jalankan make dev-frontend di terminal 2
   5. Lihat status berubah ke "Online" dan data auto sync

   ──────────────────────────────────────────

   Commands Summary

   Command             │ Mode    │ Keterangan                
   --------------------+---------+---------------------------
   `make dev-cloud`    │ -       │ Cloud API saja (port 3000)
   `make dev-frontend` │ Online  │ Frontend + sync ke cloud
   `make dev-offline`  │ Offline │ Frontend tanpa sync
   `make clean`        │ -       │ Hapus semua database

   Mau langsung test sekarang?

      Login tersedia:

   Username         │ Password       │ Role   
   -----------------+----------------+--------
   `admin`          │ `admin123`     │ Admin
   `adminShosha`    │ `password123*` │ Admin
   `adminCabang`    │ `password123*` │ Admin
   `officialShosha` │ `password123*` │ Manager
   `officialCabang` │ `password123*` │ Manager