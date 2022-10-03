package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
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

func (user *User) getDocsEntries() error {
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // the certificate is invalid on this site
	resp, err := client.R().
		SetQueryParam("tokenid", user.TokenId).
		SetQueryParam("userid", user.UserId).
		SetQueryParam("method", "get").
		SetHeader("User-Agent", "Android").
		SetBody(``).
		Post("https://yunpan.swjtu.edu.cn:9999/v1/entrydoc")

	if err != nil {
		return err
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
		return err
	}

	// set user info
	user.DocEntries = temp.Docinfos
	return nil
}

func (user *User) getDir(node FileNode) (error error, n []FileNode) {
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) // the certificate is invalid on this site
	DocId := node.Docid
	resp, err := client.R().
		SetQueryParam("tokenid", user.TokenId).
		SetQueryParam("userid", user.UserId).
		SetQueryParam("method", "list").
		SetHeader("User-Agent", "Android").
		// {"docid":"gns:\/\/8A7A69A73B8D428C90320800E6AD6783","by":"time","sort":"desc","attr":true}
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
	return nil

}
