## Configuration

The file storage system is configured using a YAML or JSON file. Here's an example configuration:

```yaml
storage:
  default: "s3"
  disks:
    - driver: "local"
      name: "local"
      local:
        rootPath: "./storage"
        baseUrl: "http://localhost:8080/storage"

    - driver: "local"
      name: "local-2"
      local:
        rootPath: "./storage2"
        baseUrl: "http://localhost:8080/storage2"

    - driver: "s3"
      name: "s3"
      s3:
        region: "ap-southeast-1"
        bucket: "bucket-1"
        profile: "rayyone"

    - driver: "s3"
      name: "s3-2"
      s3:
        region: "ap-southeast-1"
        bucket: "bucket-2"
        accessKey: "<YOUR_ACCESS_KEY>"
        secretKey: "<YOUR_SECRET_KEY>"
```

---

# File Storage Usage Example

This document shows how to use the `fileStorage` abstraction with different disk drivers (e.g., local, S3).

## fileStorage Initialize
```go
filestorage.Initialize(config.Storage.Disks, config.Storage.Default)
```
---

## Put using Default Disk

```go
filename := "uploads/default-file.txt"
file := strings.NewReader("Hello from default disk!")

err := fileStorage.Put(filename, file)
if err != nil {
log.Fatalf("Put failed: %v", err)
}
fmt.Println("Uploaded using default disk.")
```

> ✅ Uses the default disk defined in your config: `storage.default`

---


## 📤 Put using Specific Disk (e.g., S3)

```go
filename := "uploads/s3-specific.txt"
file := strings.NewReader("Hello from S3!")

err := fileStorage.Disk("s3-2").Put(filename, file)
if err != nil {
log.Fatalf("S3 Put failed: %v", err)
}
fmt.Println("Uploaded using S3 disk.")
```

---

## Get File Content

```go
data, err := fileStorage.Get("uploads/default-file.txt")
fmt.Println(string(data))
```

---

## Delete File

```go
err := fileStorage.Delete("uploads/default-file.txt")
```

---

## Get URL

```go
url := fileStorage.URL("uploads/default-file.txt")
fmt.Println("URL:", url)
```

---

## Get SignedURL

```go
url := fileStorage.SignedURL("uploads/default-file.txt")
fmt.Println("URL:", url)
```

---

## Copy

```go
fileStorage.Copy("from/file.txt", "to/file.txt")
```

---

## Move

```go
fileStorage.Move("from/file.txt", "to/file.txt")
```

---

## ⚠️ Notes on `Disk(...)`

- `fileStorage.Disk("s3-2")` allows choosing a specific disk explicitly, regardless of the default.
- Useful when you have multiple S3 buckets or local disks:


