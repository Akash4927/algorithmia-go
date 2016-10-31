Algorithmia Common Library (Golang)
===================================

Golang client library for accessing the Algorithmia API.

For API documentation, see the [Godoc](https://godoc.org/github.com/algebraic-brain/algorithmia-go).

## Install

```bash
go get github.com/algebraic-brain/algorithmia-go
```


## Authentication

First, create an Algorithmia client and authenticate with your API key:

```Go
import (
	algorithmia "github.com/algebraic-brain/algorithmia-go"
)

var apiKey = "{{Your API key here}}"
var client = algorithmia.NewClient(apiKey, "")
```

Now you're ready to call algorithms.

## Calling algorithms

The following examples of calling algorithms are organized by type of input/output which vary between algorithms.

Note: a single algorithm may have different input and output types, or accept multiple types of input,
so consult the algorithm's description for usage examples specific to that algorithm.

### Text input/output

Call an algorithm with text input by simply passing a string into its `Pipe` method.
If the algorithm output is text, then the `Result` field of the response will be a string.

```Go
algo, _ := client.Algo("demo/Hello/0.1.1")
resp, _ := algo.Pipe("Author")
response := resp.(*algorithmia.AlgoResponse)
fmt.Println(response.Result)            //Hello Author
fmt.Println(response.Metadata)          //Metadata(content_type='text',duration=0.0002127)
fmt.Println(response.Metadata.Duration) //0.0002127
```

### JSON input/output

Call an algorithm with JSON input by simply passing in a type that can be serialized to JSON.
For algorithms that return JSON, the `Result` field of the response will be the appropriate
deserialized type.

```Go
algo, _ := client.Algo("WebPredict/ListAnagrams/0.1.0")
resp, _ := algo.Pipe([]string{"transformer", "terraforms", "retransform"})
response := resp.(*algorithmia.AlgoResponse)
fmt.Println(response.Result) //[transformer retransform]
```

### Binary input/output

Call an algorithm with binary input by passing a byte array into the `Pipe` method.
Similarly, if the algorithm response is binary data, then the `Result` field of the response
will be a byte array.

```Go
input, _ := ioutil.ReadFile("/path/to/bender.png")
algo, _ := client.Algo("opencv/SmartThumbnail/0.1")
resp, _ := algo.Pipe(input)
response := resp.(*algorithmia.AlgoResponse)
ioutil.WriteFile("thumbnail.png", response.Result.([]byte), 0666)
fmt.Println(response.Result) //[binary byte sequence]
```

### Error handling

API errors and Algorithm exceptions will result in calls to `Pipe` returning an error:

```Go
algo, _ := client.Algo("util/whoopsWrongAlgo")
_, err := algo.Pipe("Hello, World!")
fmt.Println(err) //algorithm algo://util/whoopsWrongAlgo not found
```

### Request options

The client exposes options that can configure algorithm requests.
This includes support for changing the timeout or indicating that the API should include stdout in the response.

```Go
algo, _ = client.Algo("util/echo")
algo.SetOptions(algorithmia.AlgoOptions{Timeout: 60, Stdout: false})
```

## Working with data
The Algorithmia client also provides a way to manage both Algorithmia hosted data
and data from Dropbox or S3 accounts that you've connected to you Algorithmia account.

TODO
