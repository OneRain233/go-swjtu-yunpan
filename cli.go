package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var isLogin = false

// current dir
var curNode *FileNode
var entries []Entry

func initCli(user *User) error {
	err := loginInterface(user)
	if err != nil {
		return err
	}
	err, entries = user.getDocsEntries()
	if err != nil {
		return err
	}
	return nil
}

func cli(user *User) error {
	showBanner()
	err := initCli(user)
	if err != nil {
		return err
	}
	err = interactiveInterface(user)
	if err != nil {
		return err
	}
	return nil
}

func loginInterface(user *User) error {
	err := user.doLogin()
	if err != nil {
		return err
	}
	isLogin = true
	err, _ = user.getDocsEntries()
	if err != nil {
		return err
	}
	return nil
}

func entryToFileNode(entry Entry) (error error, node FileNode) {
	node = FileNode{}
	node.Name = entry.Docname
	node.Attr = entry.Attr
	node.Docid = entry.Docid
	node.ClientMtime = entry.ClientMtime
	node.Size = int64(entry.Size)
	node.isDir = true
	node.parentNode = nil
	node.Path = entry.Docid
	return nil, node
}

func showBanner() {
	fmt.Println("Welcome to SWJTU Cloud Pan CLI")
	fmt.Println("Version: 0.0.1")
	//return nil
}

func attemptLogin() (error error, username string, password string) {
	// TODO: encrypt password
	fmt.Println("Username: ")
	_, err := fmt.Scanln(&username)
	if err != nil {
		return err, "", ""
	}
	fmt.Println("Password: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		return err, "", ""
	}
	return nil, username, password
}

func parseCmd(cmd string) (error error, cmdType string, args []string) {
	temp := strings.Split(cmd, " ")
	cmdType = strings.TrimSpace(temp[0])
	args = temp[1:]
	for idx, arg := range args {
		args[idx] = strings.TrimSpace(arg)
	}
	return nil, cmdType, args
}

func interactiveInterface(user *User) error {
	// TODO
	inputReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf(">")
		// read command with args
		cmd, err := inputReader.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			continue
		}

		err, cmdType, args := parseCmd(cmd)
		//fmt.Println(cmdType, args)
		switch cmdType {
		case "ls":
			err = ls(user)
			if err != nil {
				return err
			}
		case "cd":
			path := strings.Join(args, "/")
			err = cd(user, path)
		case "upload":
			err = upload(user, args[0])
			if err != nil {
				return err
			}
		case "download":
			err = download(user, args[0])
			if err != nil {
				return err
			}
		case "pwd":
			err = pwd()
			if err != nil {
				return err
			}
		case "exit":
			fmt.Println("exit")
			os.Exit(0)
		default:
			fmt.Println("Unknown command")
		}
	}
}

func ls(user *User) error {
	// TODO: make cache for file listing
	if curNode == nil {
		for _, entry := range entries {
			fmt.Println(entry.Docname)
		}
		return nil
	}
	if !curNode.isDir {
		fmt.Println(curNode.Name)
	}
	err, dir := user.getDir(*curNode)
	if err != nil {
		return err
	}

	for _, node := range dir {
		fmt.Printf("%s\t%s\t%d\n", node.Name, node.Docid, node.Size)
	}
	return nil
}

func pwd() error {
	if curNode == nil {
		fmt.Println("/")
		return nil
	}
	if curNode.isDir {
		fmt.Println(curNode.Path)
		return nil
	}
	fmt.Println(curNode.Path)
	return nil
}

func cd(user *User, path string) error {
	if path == "/" {
		curNode = nil
		return nil
	}
	if curNode == nil {
		for _, entry := range entries {
			if entry.Docname == path {
				err, node := entryToFileNode(entry)
				if err != nil {
					return err
				}
				curNode = &node
				return nil
			}
		}
	}
	if !curNode.isDir {
		return errors.New("not a directory")
	}
	err, dir := user.getDir(*curNode)
	if err != nil {
		return err
	}
	for _, node := range dir {
		if node.Name == path {
			curNode = &node
			return nil
		}
	}
	return errors.New("no such file or directory")
}

func download(user *User, filename string) error {
	// search in current dir
	var target *FileNode
	if curNode == nil {
		for _, entry := range entries {
			if entry.Docname == filename {
				err, node := entryToFileNode(entry)
				if err != nil {
					return err
				}
				target = &node
				break
			}
		}
	} else {
		err, dir := user.getDir(*curNode)
		if err != nil {
			return err
		}
		for _, node := range dir {
			if node.Name == filename {
				break
			}
			target = &node
		}
	}
	if target.isDir {
		return errors.New("not a file")
	}
	err := user.downloadFile(*target)
	if err != nil {
		return err
	}
	return nil
}

func upload(user *User, filepath string) error {
	// check if file exists
	_, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	// check if file is too large
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	if fileInfo.Size() > 1024*1024*1024 {
		return errors.New("file too large")
	}
	// upload file
	err = user.uploadFile(*curNode, filepath)
	if err != nil {
		return err
	}
	return nil
}
