# File API

## grpc API definition

```
Put(*PutFileStu) error
Get(*GetFileStu) (io.ReadCloser, error)
List(*ListRequest) (*ListResp, error)
Del(*DelRequest) error
```

## Research

Refer：

```
https://github.com/mosn/layotto/issues/98
```

## Explanation

### Put

#### Entry type
The put interface is used to upload files. The input types are as follows：

```
type PutFileStu struct {
    DataStream io.Reader //used to read the file data transmitted by grpc stream
    FileName string //File name
    Metadata mapping [string]string //Find the field
}

```
#### return type

error

----

### Get

#### Entry type

get interface used download file：

```
    type GetFileStu struct {
    ObjectName string  //FileName
    Metadata   map[string]string //extended fields， eg.bucketName，endpoint
    }
```
#### return type

The return type is io.ReadCloser, error. io.ReadCloser implements the read and write interfaces and can be implemented by yourself, as long as it supports streaming, such as net.Pipe() type

---

### List

#### Entry type

The List interface is used to query files in a certain directory (bucket). The input types are as follows:

```
     type ListRequest struct {
         DirectoryName string //Directory name
         Metadata map[string]string //Extension field
     }
```
#### Return value type

```
     type ListResp struct {
     FilesName []string //List of all files in the directory
     }
```
---

### Del

#### Entry type

The Del interface is used to delete a file. The input types are as follows:

```
     type DelRequest struct {
         FileName string //File name to delete
         Metadata map[string]string //Extension field
     }
```

#### Return value type

Return error type

---
