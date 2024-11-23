### **_http.ResponseWriter_** (net/http package)

- A ResponseWriter interface is used by an HTTP handler to construct an HTTP response.
- A ResponseWriter may not be used after [Handler.ServeHTTP] has returned.

### **_http.MaxBytesReader_**

- It's intended for limiting the size of incoming request bodies.
- MaxBytesReader prevents clients from accidentally or maliciously sending a large request and wasting server resources. If possible, it tells the [ResponseWriter] to close the connection after the limit has been reached.

### **_json.Marshal()_**

- The json.Marshal() function is a powerful tool for converting Go data structures into JSON format. It provides flexibility through interfaces and struct tags while ensuring that only relevant data is included in the output.
- example :

  ```go
  package main
  import (
  "encoding/json"
  "fmt"
  )

  type User struct {
    Name string `json:"name"`
    Age int `json:"age"`
    Email string `json:"email,omitempty"` // omitempty means this field won't be included if empty
  }

    func main() {
    user := User{Name: "Alice", Age: 30}
    jsonData, err := json.Marshal(user)
    if err != nil {
    fmt.Println("Error marshaling to JSON:", err)
    return
    }
    fmt.Println(string(jsonData)) // Output: {"name":"Alice","age":30}
    }
  ```
