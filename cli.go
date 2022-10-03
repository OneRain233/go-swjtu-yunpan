package main

import (
	"fmt"
)

var isLogin bool = false

func cli(user *User) error {
	showBanner()
	showOptions()
	option := 0
	for option != 9 {
		_, err := fmt.Scanln(&option)
		if err != nil {
			return err
		}
		switch option {
		case 1:
			// TODO
			//err, username, password := attemptLogin()
			//if err != nil {
			//	return err
			//}
			//user.Username = username
			//user.Password = password
			err = user.doLogin()
			if err != nil {
				return err
			}
			user.printUser()
			err = user.getDocsEntries()
			if err != nil {
				return err
			}
			isLogin = true

		case 2:
			fmt.Println("logout")

		case 3:
			if !isLogin {
				fmt.Println("Please login first")
			}
			fmt.Println("Entries")
			err = showDocEntries(user)
			if err != nil {
				return err
			}
		case 4:
			fmt.Println("List Dir")
			// attempt user input dir index
			var dirIndex int
			_, err := fmt.Scanln(&dirIndex)
			if err != nil {
				return err
			}
			// get dir id
			entry := user.DocEntries[dirIndex]
			err, filenode := entryToFileNode(entry)
			// get dir entries
			var dir []FileNode
			err, dir = user.getDir(filenode)
			if err != nil {
				return err
			}
			showDir(dir)
			//fmt.Println(dirId)
		}
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

	return nil, node
}

func showBanner() {
	fmt.Println("Welcome to SWJTU Cloud Pan CLI")
	fmt.Println("Version: 0.0.1")
	//return nil
}

func showOptions() {
	fmt.Println("1. Login")
	fmt.Println("2. Logout")
	fmt.Println("3. List entries")
	fmt.Println("4. List Dir")
	fmt.Println("5. List files in a dir")
	fmt.Println("9. Exit")
	//return nil
}

func showDocEntries(user *User) error {

	for idx, doc := range user.DocEntries {
		fmt.Println(idx)
		fmt.Println(doc.Docname)
		fmt.Println(doc.Docid)
		fmt.Println("----------")
	}
	return nil
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

func showDir(nodes []FileNode) {
	for _, node := range nodes {
		if node.isDir {
			fmt.Print(node.Name + "/")
		} else {
			fmt.Print(node.Name)
		}
		fmt.Println()
	}

}

func showUserinfo(user *User) {
	fmt.Println(user.Username)
	fmt.Println(user.Name)
	fmt.Println(user.TokenId)
	fmt.Println(user.UserId)
}
