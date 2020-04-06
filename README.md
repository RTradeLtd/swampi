#  swampi

`swampi` is a golang client for the swarm http api, and enables programmatic access to the swarm API via your golang programs. It's primary purpose is for use with Temporal, but other users may find it useful. 

# format

Individual API calls are represented by a typed string call `APICall`. This allows reducing code surface for common task such as request handling, request formatting, etc..

The value of `APICall` instances is the actual HTTP path for the API call. For example single file upload is stored in code as `APICall("/bzz:/")`, so `APICall::String` will return the http path to be fed into `http.NewRequest`. Additionally methods like `APICall::ContentType` mean you don't have to deal with guessing the content type for requests as its done for you. Additionally certain API calls require url arguments, such as the single file download which takes the swarm hash. These are handled with `APICall::ParseArgs`.

# usage

For more detailed usage instructions see the tests, but in essence it consists of:

```Golang
// this creates our swampi client
swampi := New("http://localhost:8500")
fileBytes, err := ioutil.ReadFile(tt.args.filePath)
if err != nil {
	log.Fatal(err)
}
// handles parsing the api calls and sending it, returning the response
resp, err := swampi.Send(SingleFileUpload, bytes.NewReader(fileBytes), map[string][]string{
	"content-type": []string{SingleFileUpload.ContentType()},
})
if err != nil {
    log.Fatal(err)
}
// do stuff with response
```