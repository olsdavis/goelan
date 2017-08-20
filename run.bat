go build
if not exist "testServer\" mkdir testServer\
move goelan.exe testServer\goelan.exe
cd testServer\
goelan
cd ..