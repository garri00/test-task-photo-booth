# Get the build date in UTC format
$BUILD_DATE = (Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")

# Get the latest Git commit hash
$GIT_COMMIT = (git rev-parse HEAD).Trim()

# Get the latest Git tag (version)
$VERSION = (git describe --tags --abbrev=0).Trim()

# Build the Go project
go build -o bin/main.exe -ldflags="-X 'github.com/garri00/test-task-photo-booth/version.buildDate=$BUILD_DATE' -X 'github.com/garri00/test-task-photo-booth/version/version.gitCommit=$GIT_COMMIT' -X 'github.com/garri00/basic-go-project/version.gitVersion=$VERSION'" bin/main.go
