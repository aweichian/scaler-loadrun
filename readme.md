go test -bench=. .
go test -bench=. . -cpu=1,2 .
go test -bench=. . -benchtime=10s -cpu=1,2 .
go test -bench=Fib -count=3 . 