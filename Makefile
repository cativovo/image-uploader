dev:
	fd -e html -e go | entr -rcs "npm run tailwind:dev; go run main.go"
