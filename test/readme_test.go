package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	algorithmia "github.com/algebraic-brain/algorithmia-go"
)

var apiKey = os.Getenv("ALGORITHMIA_API_KEY")
var client = algorithmia.NewClient(apiKey, "")

func Test1(t *testing.T) {
	algo, _ := client.Algo("demo/Hello/0.1.1")
	resp, _ := algo.Pipe("Author")
	response := resp.(*algorithmia.AlgoResponse)
	fmt.Println(response.Result)            //Hello Author
	fmt.Println(response.Metadata)          //Metadata(content_type='text',duration=0.0002127)
	fmt.Println(response.Metadata.Duration) //0.0002127
}

func Test2(t *testing.T) {
	algo, _ := client.Algo("WebPredict/ListAnagrams/0.1.0")
	resp, _ := algo.Pipe([]string{"transformer", "terraforms", "retransform"})
	response := resp.(*algorithmia.AlgoResponse)
	fmt.Println(response.Result) //[transformer retransform]
}

func Test3(t *testing.T) {
	input, _ := ioutil.ReadFile(`C:\Users\Osman\Documents\Тимур\Практическая №3 по информатике\Пароход.png`)
	algo, _ := client.Algo("opencv/SmartThumbnail/0.1")
	resp, _ := algo.Pipe(input)
	response := resp.(*algorithmia.AlgoResponse)
	ioutil.WriteFile("thumbnail.png", response.Result.([]byte), 0666)
	fmt.Println(response.Result) //[binary byte sequence]
}

func Test4(t *testing.T) {
	algo, _ := client.Algo("util/whoopsWrongAlgo")
	_, err := algo.Pipe("Hello, World!")
	fmt.Println(err)
}

func Test5(t *testing.T) {
	foo := client.Dir("data://.my/foo")
	foo.Create(nil)
	foo.File("sample.txt").Put("sample text contents")
	foo.File("binary_file").Put([]byte{72, 101, 108, 108, 111})
}

func Test6(t *testing.T) {
	foo := client.Dir("data://.my/foo")
	sampleText, _ := foo.File("sample.txt").StringContents() //string object
	fmt.Println(sampleText)                                  //"sample text contents"
	binaryContent, _ := foo.File("binary_file").Bytes()      //binary data
	fmt.Println(string(binaryContent))                       //"Hello"
	tempFile, _ := foo.File("binary_file").File()            //Open file descriptor for read
	defer tempFile.Close()
	binaryContent, _ = ioutil.ReadAll(tempFile)
	fmt.Println(string(binaryContent)) //"Hello"
}

func Test7(t *testing.T) {
	foo := client.Dir("data://.my/foo")
	foo.File("sample.txt").Delete()
	foo.ForceDelete() // force deleting the directory and its contents
}
