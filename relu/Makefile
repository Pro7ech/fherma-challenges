# Flags default values
sk = temps/sk.bin
cc = ./artifact/context
key_eval = temps/mult.bin
input = temps/in.bin
output = temps/out.bin

all-0: 
	go run setup.go --sk=$(sk) --cc=$(cc)_0 --key_eval=$(key_eval) --input=$(input)
	go run main.go --cc=$(cc)_0 --key_eval=$(key_eval) --input=$(input) --output=$(output)
	go run verify.go --sk=$(sk) --cc=$(cc)_0 --input=$(input) --output=$(output)
	go run clean.go
	go clean

all-1: 
	go run setup.go --sk=$(sk) --cc=$(cc)_1 --key_eval=$(key_eval) --input=$(input)
	go run main.go --cc=$(cc)_1 --key_eval=$(key_eval) --input=$(input) --output=$(output)
	go run verify.go --sk=$(sk) --cc=$(cc)_1 --input=$(input) --output=$(output)
	go run clean.go
	go clean

clean:
	go run clean.go
	go clean
