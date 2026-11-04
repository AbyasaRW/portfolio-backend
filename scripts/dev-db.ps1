<#
 scripts/dev-db.ps1
 Menjalankan PostgreSQL KHUSUS PENGEMBANGAN di port 5433 (terpisah dari PostgreSQL utama
 di port 5432, jadi tidak mengganggu database lain dan tidak butuh password admin).

 Pemakaian (dari folder backend):
   .\scripts\dev-db.ps1 start     # buat (jika belum ada) dan nyalakan database
   .\scripts\dev-db.ps1 stop      # matikan
   .\scripts\dev-db.ps1 status    # cek berjalan atau tidak
   .\scripts\dev-db.ps1 reset     # HAPUS semua data dev lalu buat ulang

 Data disimpan di %LOCALAPPDATA%\portfolio-devdb (di luar folder proyek).
 Hanya menerima koneksi dari komputer ini (127.0.0.1), tanpa password: JANGAN dipakai untuk produksi.
#>
param([ValidateSet('start','stop','status','reset')][string]$Action = 'start')

$ErrorActionPreference = 'Stop'
$Port    = 5433
$DbName  = 'portfolio_db'
$DbUser  = 'portfolio_dev'
$DataDir = Join-Path $env:LOCALAPPDATA 'portfolio-devdb'
$LogFile = Join-Path $DataDir 'server.log'

# Cari instalasi PostgreSQL (versi tertinggi).
$pgRoot = Get-ChildItem 'C:\Program Files\PostgreSQL' -Directory -ErrorAction SilentlyContinue | Sort-Object Name -Descending | Select-Object -First 1
if (-not $pgRoot) { throw 'PostgreSQL tidak ditemukan di C:\Program Files\PostgreSQL. Install dulu, atau pakai Docker.' }
$Bin = Join-Path $pgRoot.FullName 'bin'

function Test-Running { & "$Bin\pg_isready.exe" -h 127.0.0.1 -p $Port -q; return ($LASTEXITCODE -eq 0) }

function Initialize-Cluster {
    if (Test-Path (Join-Path $DataDir 'PG_VERSION')) { return }
    Write-Host "Membuat cluster baru di $DataDir ..."
    New-Item -ItemType Directory -Force -Path $DataDir | Out-Null
    & "$Bin\initdb.exe" -D $DataDir -U $DbUser -A trust -E UTF8 --locale=C | Out-Null
    if ($LASTEXITCODE -ne 0) { throw 'initdb gagal.' }
    Add-Content (Join-Path $DataDir 'postgresql.conf') "`nport = $Port`nlisten_addresses = '127.0.0.1'`nunix_socket_directories = ''`n"
}

switch ($Action) {
    'status' { if (Test-Running) { "Berjalan di 127.0.0.1:$Port" } else { 'Tidak berjalan' } }
    'stop'   { if (Test-Running) { & "$Bin\pg_ctl.exe" -D $DataDir stop -m fast -W | Out-Null; Start-Sleep -Seconds 2; 'Dimatikan.' } else { 'Sudah mati.' } }
    'reset'  {
        if (Test-Running) { & "$Bin\pg_ctl.exe" -D $DataDir stop -m fast -W | Out-Null; Start-Sleep -Seconds 2 }
        if (Test-Path $DataDir) { Remove-Item $DataDir -Recurse -Force }
        & $PSCommandPath start
    }
    'start'  {
        Initialize-Cluster
        if (-not (Test-Running)) {
            # Jalankan postgres.exe langsung, terlepas dari pipa/terminal ini.
            # (pg_ctl start | ... bisa macet karena server mewarisi pipa keluaran.)
            Remove-Item (Join-Path $DataDir 'postmaster.pid') -Force -ErrorAction SilentlyContinue
            Start-Process -FilePath "$Bin\postgres.exe" -ArgumentList '-D', "`"$DataDir`"" -WindowStyle Hidden `
                -RedirectStandardOutput (Join-Path $DataDir 'stdout.log') -RedirectStandardError $LogFile | Out-Null
            $ready = $false
            for ($i = 0; $i -lt 30; $i++) { Start-Sleep -Milliseconds 500; if (Test-Running) { $ready = $true; break } }
            if (-not $ready) { throw "Gagal start. Lihat log: $LogFile" }
        }
        $exists = & "$Bin\psql.exe" -h 127.0.0.1 -p $Port -U $DbUser -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$DbName'"
        if (-not $exists) { & "$Bin\createdb.exe" -h 127.0.0.1 -p $Port -U $DbUser $DbName }
        Write-Host "Database siap: host=127.0.0.1 port=$Port user=$DbUser db=$DbName (tanpa password)"
    }
}