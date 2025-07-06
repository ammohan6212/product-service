module gin-gcs-backend

go 1.21

require (
    cloud.google.com/go/storage v1.38.0
    github.com/gin-gonic/gin v1.10.0
    github.com/go-sql-driver/mysql v1.7.1
    github.com/google/uuid v1.3.1
    github.com/golang/protobuf v1.5.3 // indirect
    golang.org/x/net v0.38.0           // added to fix security issues
    golang.org/x/oauth2 v0.27.0        // added to fix security issues
)
