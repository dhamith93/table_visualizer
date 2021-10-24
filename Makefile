clean:
	echo "Removing dist/ dir"
	rm -rf dist

build:
	mkdir dist
	echo "Building for Linux/macOS"
	mkdir dist/table_viz-linux-amd64
	cd cmd/table_visualizer && \
	GOOS=linux GOARCH=amd64 go build -o ../../dist/table_viz-linux-amd64/table_viz-linux-amd64 .
	cd dist && zip -r table_viz-linux-amd64.zip table_viz-linux-amd64

run:
	go run cmd/table_visualizer/*.go