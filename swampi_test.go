package swampi

import (
	"bytes"
	"fmt"
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
		{"Text File", args{"README.md"}, false, "b4dd34b82fc6d9518b73676430221595f72e6c1c105adff19e1f23e7468b8565"},
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
			if resp.call != SingleFileUpload {
				t.Fatal("bad call")
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
				fmt.Println(string(contents))
				t.Fatal("bad hash returned")
			}
			resp, err = swampi.Send(SingleFileDownload, nil, nil, string(contents))
			if err != nil {
				t.Fatal(err)
			}
			if resp.call != SingleFileDownload {
				t.Fatal("bad call")
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
			if resp.call != SingleFileUpload {
				t.Fatal("bad call")
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

func TestBZZList(t *testing.T) {
	swampi := New("http://localhost:8500")
	type args struct {
		hash string
	}
	tests := []struct {
		name     string
		args     args
		wantHash string
	}{
		{"Readme", args{"b4dd34b82fc6d9518b73676430221595f72e6c1c105adff19e1f23e7468b8565"}, "49c5b6e9dd8531a05a7a1c6f91f261c2214bc93d9f1c157fe2dc68c8006c8b63"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := swampi.Send(ListFiles, nil, nil, tt.args.hash)
			if err != nil {
				t.Fatal(err)
			}
			if resp.call != ListFiles {
				t.Fatal("bad call")
			}
			su, err := resp.SwarmUnmarshal()
			if err != nil {
				t.Fatal(err)
			}
			bzzList, ok := su.(*BZZList)
			if !ok {
				t.Fatal("failed to properly unmarshal")
			}
			if len(bzzList.Entries) == 0 {
				t.Fatal("bad number of entries")
			}
			if bzzList.Entries[0].Hash != tt.wantHash {
				t.Fatal("bad hash returned")
			}
			fmt.Printf("%+v\n", bzzList)
		})
	}
}
