# httpclient

## QUICK START

import it in your project

```go
import "github.com/Houjingchao/httpclient"
```

#### GET

```go
str, err := httpclient.Get("http://baidu.com").
		Execute().
		String()
	if err != nil {
		// do something
	}
```

### Post

```go
err := httpclient.Post("http://0.0.0.0:8080/json").
		Param("sex", "boy").
		Execute().
		ToJson(&response)
	if err != nil {
		fmt.Println(err)
	}
```

###  json result
```go
response := struct {
		Name string
		Sex  string
	}{}

	err := httpclient.Post("http://0.0.0.0:8080/json").
		Param("sex", "boy").
		Execute().
		ToJson(&response)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v", response)
```

