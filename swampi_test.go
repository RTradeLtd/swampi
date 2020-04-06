package swampi

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestSingleFile(t *testing.T) {
	swampi := New("http://localhost:8500")
	type args struct {
		filePath string
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		wantHash string
	}{
		{"Text File", args{"README.md"}, false, "17dbcabf5b829eed6214dfeada274cfb4cde67b2c3801ffecac8ccc6e8f11c91"},
		{"Binary File", args{"./test_data/swarm"}, false, "4bd4fb200ec385993708e2ff2163b3935178deb4920c85ff184a1f5d8bd9318d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileBytes, err := ioutil.ReadFile(tt.args.filePath)
			if err != nil {
				t.Fatal(err)
			}
			resp, err := swampi.Send(SingleFileUpload, bytes.NewReader(fileBytes), map[string][]string{
				"content-type": []string{SingleFileUpload.ContentType()},
			})
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				t.Fatalf("bad status code, got %v, want %v", resp.StatusCode, http.StatusOK)
			}
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(contents) != tt.wantHash {
				t.Fatal("bad hash returned")
			}
			resp, err = swampi.Send(SingleFileDownload, nil, nil, string(contents))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if err := ioutil.WriteFile(tt.args.filePath+"1", data, os.FileMode(0640)); err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tt.args.filePath + "1")
			data, err = ioutil.ReadFile(tt.args.filePath + "1")
			if err != nil {
				t.Fatal(err)
			}
			resp, err = swampi.Send(SingleFileUpload, bytes.NewReader(fileBytes), map[string][]string{
				"content-type": []string{SingleFileUpload.ContentType()},
			})
			if err != nil {
				t.Fatal(err)
			}
			contents, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(contents) != tt.wantHash {
				t.Fatal("bad hash returned")
			}
		})
	}
}
