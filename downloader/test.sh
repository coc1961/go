export GOMAXPROCS=8
go run  main.go -n 5 -url http://mirror.cc.columbia.edu/pub/software/eclipse/technology/epp/downloads/release/luna/SR2/eclipse-standard-luna-SR2-win32-x86_64.zip  -o eclipse.zip && unzip -t eclipse.zip > /dev/null && echo Test OK
rm eclipse.zip 2> /dev/null