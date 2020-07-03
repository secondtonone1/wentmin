DIR=$(cd $(dirname $0) && pwd )
echo  $DIR
cd $DIR
go build  -o wentmin main.go
./wentmin >output.txt 2>&1 &
