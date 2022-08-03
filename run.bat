go build -o bookings.exe ./cmd/web/. || exit /b
bookings.exe
./booking -dbname=booking -dbuser=postgres -cache=false -production=false