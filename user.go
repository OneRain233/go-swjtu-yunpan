package main

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"strconv"
	"strings"
)

var client = resty.New()

func (user *User) doLogin() error {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"account":"` + user.Username + `","deviceinfo":{"ostype":1},"password":"` + user.Password + `","vcodeinfo":{"uuid":"","vcode":"","ismodify":false}}`).
		Post("http://yunpan.swjtu.edu.cn:9998/v1/auth1?method=getnew")
	if err != nil {
		return err
	}
	//fmt.Println(resp)
	// decode json
	var result map[string]interface{}
	if json.Unmarshal(resp.Body(), &result) != nil {
		return err
	}
	//fmt.Println(result)

	// check login status (has errcode key)
	if _, ok := result["errcode"]; ok {
		return fmt.Errorf("login failed")
	}

	// set user info
	user.TokenId = result["tokenid"].(string)
	user.UserId = result["userid"].(string)

	// get user info
	resp, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{}`).
		SetQueryParam("tokenid", user.TokenId).
		SetQueryParam("userid", user.UserId).
		Post("http://yunpan.swjtu.edu.cn:9998/v1/user?method=get")
	if err != nil {
		return err
	}
	//fmt.Println(resp)
	// decode json
	if json.Unmarshal(resp.Body(), &result) != nil {
		return err
	}

	// set user info
	user.Name = result["name"].(string)

	return nil
}

func (user *User) printUser() {
	s := fmt.Sprintf("username: %s, password: %s, name: %s, tokenid: %s, userid: %s", user.Username, user.Password, user.Name, user.TokenId, user.UserId)
	fmt.Println(s)

}

func (user *User) getDocsEntries() (error error, n []Entry) {
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // the certificate is invalid on this site
	resp, err := client.R().
		SetQueryParam("tokenid", user.TokenId).
		SetQueryParam("userid", user.UserId).
		SetQueryParam("method", "get").
		SetHeader("User-Agent", "Android").
		SetBody(``).
		Post("https://yunpan.swjtu.edu.cn:9999/v1/entrydoc")

	if err != nil {
		return err, []Entry{}
	}
	fmt.Println("===============================")
	fmt.Println(resp)
	fmt.Println("===============================")
	// decode json
	//var entries []Entry
	var temp tempEntries
	//var temp map[string][]byte
	err = json.Unmarshal(resp.Body(), &temp)
	if err != nil {
		return err, []Entry{}
	}

	// set user info
	user.DocEntries = temp.Docinfos
	return nil, user.DocEntries
}

func (user *User) getDir(node FileNode) (error error, n []FileNode) {
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // the certificate is invalid on this site
	DocId := node.Docid
	resp, err := client.R().
		SetQueryParam("tokenid", user.TokenId).
		SetQueryParam("userid", user.UserId).
		SetQueryParam("method", "list").
		SetHeader("User-Agent", "Android").
		// {"docid":"","by":"time","sort":"desc","attr":true}
		SetBody(`{"docid":"` + DocId + `","by":"time","sort":"desc","attr":true}`).
		Post("https://yunpan.swjtu.edu.cn:9124/v1/dir")

	newDir := Dir{}
	err = json.Unmarshal(resp.Body(), &newDir)

	// new FileNode array
	var nodes []FileNode
	// to toFileNode
	for _, entry := range newDir.Dirs {
		err = toFileNode(&entry, &node)
		if err != nil {
			return err, []FileNode{}
		}
		nodes = append(nodes, entry)
	}

	for _, entry := range newDir.Files {
		err = toFileNode(&entry, &node)
		if err != nil {
			return err, []FileNode{}
		}
		nodes = append(nodes, entry)
	}

	if err != nil {
		return err, []FileNode{}
	}

	return nil, nodes
}

func toFileNode(node *FileNode, parentFileNode *FileNode) error {
	if node.Size == -1 {
		node.isDir = true
	} else {
		node.isDir = false
	}
	node.parentNode = parentFileNode
	node.Path = node.Docid
	return nil

}

func (user *User) downloadFile(node FileNode) error {
	if node.isDir {
		return fmt.Errorf("cannot download a directory")
	}
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // the certificate is invalid on this site

	//{
	//"docid":"",
	//"reqhost":"yunpan.swjtu.edu.cn",
	//"rev":""
	//}
	reqhost := "yunpan.swjtu.edu.cn"
	rev := node.Rev
	docid := node.Path
	body := `{"docid":"` + docid + `","reqhost":"` + reqhost + `","rev":"` + rev + `"}`
	resp, err := client.R().
		SetQueryParam("tokenid", user.TokenId).
		SetQueryParam("userid", user.UserId).
		SetQueryParam("method", "osdownload").
		SetHeader("User-Agent", "Android").
		SetBody(body).
		Post("https://yunpan.swjtu.edu.cn:9124/v1/file")

	if err != nil {
		return err
	}

	err = saveFile(resp.Body())
	if err != nil {
		return err
	}

	return nil
}

func saveFile(message json.RawMessage) error {
	var file DownloadInfo
	err := json.Unmarshal(message, &file)
	if err != nil {
		return err
	}
	fmt.Println("=============Downloading================")
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // the certificate is invalid on this site
	// TODO: clean these code
	resp := client.R()
	aR, err := parseAuthRequest(file.Authrequest)
	url := aR.Url
	headers := aR.Headers
	for k, v := range headers {
		resp.SetHeader(k, v)
	}
	resp.SetOutput(file.Name)
	get, err := resp.Get(url)
	if err != nil {
		return err
	}
	fmt.Println(get)

	return nil
}

func parseHeader(rawHeader string) (header map[string]string) {
	header = make(map[string]string)
	lines := strings.Split(rawHeader, "\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			s := strings.Split(line, ":")
			header[s[0]] = strings.Join(s[1:], ":")
		}
	}
	return header
}

func (user *User) uploadFile(node FileNode, filepath string) error {
	// TODO: multi parts upload
	// TODO: upload a duplicate file (rename)
	if !node.isDir {
		return fmt.Errorf("cannot upload to a file")
	}
	fmt.Println("Uploading to " + node.Name)
	client.SetProxy("http://192.168.123.65:9999")
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // the certificate is invalid on this site
	// get upload info
	targetDocId := node.Docid
	fi, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	fileSize := fi.Size()
	fileName := fi.Name()
	// time stamp like 1664775766000000
	clientMtime := fi.ModTime().UnixNano()
	ondup := 0
	//{"docid":"","length":2662,"name":"","client_mtime":1664775766000000,"ondup":0}
	post, err := client.R().
		SetQueryParam("tokenid", user.TokenId).
		SetQueryParam("userid", user.UserId).
		SetQueryParam("method", "osinitmultiupload").
		SetHeader("User-Agent", "Android").
		SetBody(`{"docid":"` + targetDocId + `","length":` + strconv.FormatInt(fileSize, 10) + `,"name":"` + fileName + `","client_mtime":` + strconv.FormatInt(clientMtime, 10) + `,"ondup":` + strconv.Itoa(ondup) + `}`).
		Post("https://yunpan.swjtu.edu.cn:9124/v1/file")
	if err != nil {
		return err
	}
	//{
	//	"docid":"",
	//	"name":"",
	//	"rev":"",
	//	"uploadid":""
	//}
	uploadInfo := struct {
		Docid    string `json:"docid"`
		Name     string `json:"name"`
		Rev      string `json:"rev"`
		Uploadid string `json:"uploadid"`
	}{}
	err = json.Unmarshal(post.Body(), &uploadInfo)
	if err != nil {
		return err
	}

	// get request info
	//{"docid":"","rev":"","uploadid":"","parts":"1","reqhost":"yunpan.swjtu.edu.cn","usehttps":true}
	post, err = client.R().
		SetQueryParam("tokenid", user.TokenId).
		SetQueryParam("userid", user.UserId).
		SetQueryParam("method", "osuploadpart").
		SetHeader("User-Agent", "Android").
		SetBody(`{"docid":"` + uploadInfo.Docid + `","rev":"` + uploadInfo.Rev + `","uploadid":"` + uploadInfo.Uploadid + `","parts":"1","reqhost":"yunpan.swjtu.edu.cn","usehttps":true}`).
		Post("https://yunpan.swjtu.edu.cn:9124/v1/file")
	if err != nil {
		return err
	}
	// {"authrequests":
	//{"1":
	//["PUT",
	//"https://yunpan.swjtu.edu.cn:9029/",
	//"Content-Type: application/octet-stream",
	//"X-Eoss-Date: Wed, 05 Oct 2022 03:46:36 GMT",
	//"Authorization: "",
	//"x-as-userid: ""]}}
	requestInfo := struct {
		Authrequests map[string][]string `json:"authrequests"`
	}{}
	err = json.Unmarshal(post.Body(), &requestInfo)
	if err != nil {
		return err
	}

	// upload file (PUT)
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	fmt.Println(fileContent)
	//fmt.Println("debug")
	var Etag string
	for _, authrequest := range requestInfo.Authrequests {
		headers := parseHeader(strings.Join(authrequest[2:], "\n"))
		req := client.R()
		for k, v := range headers {
			req.SetHeader(k, v)
		}
		req.SetBody(fileContent)
		put, err := req.Put(authrequest[1])
		if err != nil {
			return err
		}
		fmt.Println(put.Result())
		Etag = put.Header().Get("Etag")

		fmt.Println(Etag)
	}

	// complete upload
	//{"docid":"",
	//"rev":"",
	//"uploadid":"",
	//"partinfo":{"1":["string",2662]},    // string int fucking data structure ????????? Stupid API
	//"reqhost":"yunpan.swjtu.edu.cn",
	//"usehttps":true}
	//partinfo["1"] = []string{Etag, strconv.FormatInt(fileSize, 10)}
	//partinfo := fmt.Sprintf(`{"1":["%s",%d]}`, Etag, fileSize)
	//body := struct {
	//	Docid    string `json:"docid"`
	//	Rev      string `json:"rev"`
	//	Uploadid string `json:"uploadid"`
	//	Partinfo string `json:"partinfo"` // TODO: fix bug
	//	Reqhost  string `json:"reqhost"`
	//	Usehttps bool   `json:"usehttps"`
	//}{}
	//body.Docid = uploadInfo.Docid
	//body.Rev = uploadInfo.Rev
	//body.Uploadid = uploadInfo.Uploadid
	//body.Partinfo = partinfo
	//body.Reqhost = "yunpan.swjtu.edu.cn"
	//body.Usehttps = true
	body := fmt.Sprintf(`{"docid":"%s","rev":"%s","uploadid":"%s","partinfo":{"1":["%s",%d]},"reqhost":"yunpan.swjtu.edu.cn","usehttps":true}`, uploadInfo.Docid, uploadInfo.Rev, uploadInfo.Uploadid, Etag, fileSize)
	post, err = client.R().
		SetQueryParam("tokenid", user.TokenId).
		SetQueryParam("userid", user.UserId).
		SetQueryParam("method", "oscompleteupload").
		SetHeader("User-Agent", "Android").
		SetBody(body).
		Post("https://yunpan.swjtu.edu.cn:9124/v1/file")
	if err != nil {
		return err
	}
	contentType := post.Header().Get("Content-Type")
	boundary := parseBoundary(contentType)
	fmt.Println(boundary)
	body = string(post.Body())
	newBody, err := parseBodyWithBoundary(body, boundary)
	if err != nil {
		return err
	}
	fmt.Println(len(newBody))
	for _, v := range newBody {
		fmt.Println(v)
	}
	if len(newBody) != 2 {
		return errors.New("upload failed")
	}
	// post complete upload
	aR := struct {
		AuthRequest []string `json:"authrequest"`
	}{}
	err = json.Unmarshal([]byte(newBody[1]), &aR)
	if err != nil {
		return err
	}
	aRR, err := parseAuthRequest(aR.AuthRequest)
	if err != nil {
		return err
	}
	headers := aRR.Headers
	url := aRR.Url

	req := client.R()
	for k, v := range headers {
		req.SetHeader(k, v)
	}
	req.SetBody([]byte(newBody[0]))
	_, err = req.Post(url)
	if err != nil {
		return err
	}

	// end upload
	//{"docid":"gns:\/\/\\/","rev":""}
	req = client.R().
		SetQueryParam("tokenid", user.TokenId).
		SetQueryParam("userid", user.UserId).
		SetQueryParam("method", "osendupload").
		SetHeader("User-Agent", "Android").
		SetBody(`{"docid":"` + uploadInfo.Docid + `","rev":"` + uploadInfo.Rev + `"}`)
	_, err = req.Post("https://yunpan.swjtu.edu.cn:9124/v1/file")
	if err != nil {
		return err
	}
	return nil

}

func parseBoundary(contentType string) string {
	// Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW
	boundary := strings.Split(contentType, ";")[1]
	boundary = strings.TrimSpace(boundary)
	boundary = strings.TrimPrefix(boundary, "boundary=")
	return boundary
}

func parseBodyWithBoundary(body string, boundary string) ([]string, error) {
	//	--eyyGHBAJBIdGxINw1LCiyh4S4cL2f8st
	//
	//[{"path":"db/-0","etag":"","size_bytes":2662}]
	//	--eyyGHBAJBIdGxINw1LCiyh4S4cL2f8st
	//
	//	{"authrequest":["POST","https://yunpan.swjtu.edu.cn:9029/anyshares3accesstestbucket//-i","Content-Type: application/x-www-form-urlencoded","X-Eoss-Date: Wed, 05 Oct 2022 03:46:38 GMT","Authorization: ","x-as-userid: "]}
	//	--eyyGHBAJBIdGxINw1LCiyh4S4cL2f8st--
	var result []string
	body = strings.TrimPrefix(body, "--"+boundary)
	body = strings.TrimSuffix(body, "--"+boundary+"--")
	body = strings.TrimSpace(body)
	result = strings.Split(body, "--"+boundary)
	for i, v := range result {
		result[i] = strings.TrimSpace(v)
	}
	return result, nil
}

func parseAuthRequest(rawMsg []string) (AuthRequest, error) {
	//["GET","https://yunpan.swjtu.edu.cn:9029///-i","X-Eoss-Date: Wed, 05 Oct 2022 03:53:46 GMT","X-Eoss-Length: ","Authorization: ","x-as-userid: "]
	method := rawMsg[0]
	url := rawMsg[1]
	headers := parseHeader(strings.Join(rawMsg[2:], "\n"))
	return AuthRequest{
		Method:  method,
		Url:     url,
		Headers: headers,
	}, nil
}

func (user *User) deleteFile(file FileNode) error {
	docID := file.Docid
	req, err := client.R().
		SetQueryParam("tokenid", user.TokenId).
		SetQueryParam("userid", user.UserId).
		SetQueryParam("method", "delete").
		SetHeader("User-Agent", "Android").
		SetBody(`{"docid":"` + docID + `"}`).
		Post("https://yunpan.swjtu.edu.cn:9124/v1/file")

	if err != nil {
		return err
	}
	if req.StatusCode() != 200 {
		return errors.New("delete failed")
	}
	return nil
}

func (user *User) renameFile(file FileNode, newName string) error {
	docId := file.Docid
	req, err := client.R().
		SetQueryParam("tokenid", user.TokenId).
		SetQueryParam("userid", user.UserId).
		SetQueryParam("method", "rename").
		SetHeader("User-Agent", "Android").
		SetBody(`{"docid":"` + docId + `","newname":"` + newName + `"}`).
		Post("https://yunpan.swjtu.edu.cn:9124/v1/file")
	if err != nil {
		return err
	}
	if req.StatusCode() != 200 {
		return errors.New("rename failed")
	}
	return nil

}
