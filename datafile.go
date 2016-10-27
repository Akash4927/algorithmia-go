package algorithmia

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type datafileClient interface {
	getHelper(url string, params url.Values) (*http.Response, error)
	headHelper(url string) (*http.Response, error)
	putHelper(url string, data []byte) (*http.Response, error)
	deleteHelper(url string) (*http.Response, error)
}

type FileAttributes struct {
	FileName     string `json:"filename"`
	LastModified string `json:"last_modified"`
	Size         int64  `json:"size"`
}

type DataFile struct {
	DataObjectType

	Path         string
	Url          string
	LastModified time.Time
	Size         int64

	client datafileClient
}

func NewDataFile(client datafileClient, dataUrl string) *DataFile {
	p := strings.TrimSpace(dataUrl)
	if strings.HasPrefix(p, "data://") {
		p = p[len("data://"):]
	} else if strings.HasPrefix(p, "/") {
		p = p[1:]
	}
	return &DataFile{
		DataObjectType: File,
		client:         client,
		Path:           p,
		Url:            getUrl(p),
	}
}

func (f *DataFile) SetAttributes(attr *FileAttributes) error {
	//%Y-%m-%dT%H:%M:%S.000Z
	t, err := time.Parse("2006-01-02T15:04:05.000Z", attr.LastModified)
	if err != nil {
		return err
	}
	f.LastModified = t
	f.Size = attr.Size
	return nil
}

func (f *DataFile) File() (*os.File, error) {
	if exists, err := f.Exists(); err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.New(fmt.Sprint("file does not exist -", f.Path))
	}

	resp, err := f.client.getHelper(f.Url, url.Values{})
	if err != nil {
		return nil, err
	}

	rf, err := ioutil.TempFile(os.TempDir(), "algorithmia")
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(rf, resp.Body)
	if err != nil {
		rf.Close()
		return nil, err
	}
	return rf, nil
}

func (f *DataFile) Exists() (bool, error) {
	resp, err := f.client.headHelper(f.Url)
	return resp.StatusCode == http.StatusOK, err
}

func (f *DataFile) Name() (string, error) {
	_, name, err := getParentAndBase(f.Path)
	return name, err
}

func (f *DataFile) Bytes() ([]byte, error) {
	if exists, err := f.Exists(); err != nil {
		return nil, err
	} else if !exists {
		return nil, errors.New(fmt.Sprint("file does not exist -", f.Path))
	}

	resp, err := f.client.getHelper(f.Url, url.Values{})
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func (f *DataFile) StringContents() (string, error) {
	if b, err := f.Bytes(); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func (f *DataFile) Json(x interface{}) error {
	if exists, err := f.Exists(); err != nil {
		return err
	} else if !exists {
		return errors.New(fmt.Sprint("file does not exist -", f.Path))
	}

	resp, err := f.client.getHelper(f.Url, url.Values{})
	if err != nil {
		return err
	}

	return getJson(resp, x)
}

func (f *DataFile) Put(data []byte) error {
	resp, err := f.client.putHelper(f.Url, data)
	if err != nil {
		return err
	}

	b, err := getRaw(resp)
	if err != nil {
		return err
	}

	err = ErrorFromJsonData(b)
	if err != nil {
		return err
	}

	return nil
}

func (f *DataFile) PutJson(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return f.Put(b)
}

func (f *DataFile) PutFile(fpath string) error {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return err
	}
	return f.Put(b)
}

func (f *DataFile) Delete() error {
	resp, err := f.client.deleteHelper(f.Url)
	if err != nil {
		return err
	}

	if err := ErrorFromResponse(resp); err != nil {
		return err
	}

	return nil
}
