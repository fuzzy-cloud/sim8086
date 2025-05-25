asm:
	./scripts/gen_asm.sh

test: asm
	go test ./...
