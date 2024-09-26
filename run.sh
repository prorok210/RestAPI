echo "Running tests..."
go test ./... -v

if [ $? -ne 0 ]; then
    echo "Tests failed. Aborting server start."
    exit 1
fi

echo "Building server..."
sudo go build main.go

if [ $? -ne 0 ]; then
    echo "Build failed. Aborting server start."
    exit 1
fi

sudo ./main