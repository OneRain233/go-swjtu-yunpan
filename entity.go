package main

// User define user struct
type User struct {
	Username   string  `json:"username"`
	Password   string  `json:"password"`
	Name       string  `json:"name"`
	TokenId    string  `json:"tokenid"`
	UserId     string  `json:"userid"`
	DocEntries []Entry `json:"doc_entries"`
}

//type Entry struct {
//	Docinfo []Entry `json:"docinfos"`
//}

type Entry struct {
	Attr              int    `json:"attr"`
	ClientMtime       int64  `json:"client_mtime"`
	CreaterId         string `json:"createrId"`
	Docid             string `json:"docid"`
	Docname           string `json:"docname"`
	Doctype           string `json:"doctype"`
	Downloadwatermark bool   `json:"downloadwatermark"`
	Duedate           int    `json:"duedate"`
	Modified          int64  `json:"modified"`
	Otag              string `json:"otag"`
	Siteinfo          struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"siteinfo"`
	Size            int    `json:"size"`
	Typename        string `json:"typename"`
	ViewDoctype     int    `json:"view_doctype"`
	ViewDoctypename string `json:"view_doctypename"`
	ViewName        string `json:"view_name"`
	ViewType        int    `json:"view_type"`
}

type Dir struct {
	Dirs  []FileNode `json:"dirs"`
	Files []FileNode `json:"files"`
}
type tempEntries struct {
	Docinfos []Entry `json:"docinfos"`
}

type FileNode struct {
	Attr        int    `json:"attr"`
	ClientMtime int64  `json:"client_mtime"`
	CreateTime  int64  `json:"create_time"`
	Creator     string `json:"creator"`
	Csflevel    int    `json:"csflevel"`
	Docid       string `json:"docid"`
	Duedate     int    `json:"duedate"`
	Editor      string `json:"editor"`
	Modified    int64  `json:"modified"`
	Name        string `json:"name"`
	Rev         string `json:"rev"`
	Size        int64  `json:"size"`
	isDir       bool
	parentNode  *FileNode
	Path        string
}

type DownloadInfo struct {
	Authrequest   []string `json:"authrequest"`
	ClientMtime   int64    `json:"client_mtime"`
	Editor        string   `json:"editor"`
	Modified      int64    `json:"modified"`
	Name          string   `json:"name"`
	NeedWatermark bool     `json:"need_watermark"`
	Rev           string   `json:"rev"`
	Siteid        string   `json:"siteid"`
	Size          int      `json:"size"`
}

type AuthRequest struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"header"`
}
