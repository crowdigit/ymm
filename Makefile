.PHONY: golds run

golds:
	golds -theme=dark ./internal

run:
	go run ./cmd/ymm single 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'
