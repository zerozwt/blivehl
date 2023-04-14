package main

import (
	"flag"
	"fmt"

	"github.com/zerozwt/blivehl/server/db"
	"github.com/zerozwt/blivehl/server/service"
)

var toolCmd struct {
	Cmd  string
	User struct {
		SetUser  bool
		ListUser bool
		ID       string
		NickName string
		Password string
		IsAdmin  bool
		Ban      bool
	}
}

func main() {
	dbFile := ""
	flag.StringVar(&dbFile, "db", "blivehl.db", "specify db file")
	flag.StringVar(&toolCmd.Cmd, "cmd", "", "tool command")
	flag.BoolVar(&toolCmd.User.SetUser, "set_user", false, "edit user info (cmd=user)")
	flag.BoolVar(&toolCmd.User.ListUser, "list_user", false, "list all users (cmd=user)")
	flag.StringVar(&toolCmd.User.ID, "user_id", "", "specify user id")
	flag.StringVar(&toolCmd.User.NickName, "user_name", "", "specify user nickname")
	flag.StringVar(&toolCmd.User.Password, "user_pass", "", "specify user password")
	flag.BoolVar(&toolCmd.User.IsAdmin, "user_admin", false, "specify user is an admin")
	flag.BoolVar(&toolCmd.User.Ban, "user_ban", false, "specify whether user is banned")
	flag.Parse()

	if toolCmd.Cmd != "user" {
		fmt.Println("Invalid cmd: " + toolCmd.Cmd)
		return
	}

	if err := db.InitDB(dbFile); err != nil {
		fmt.Println("Init db failed: ", err)
		return
	}

	userHeader := []string{"ID", "Name", "Password", "IsAdmin", "Banned"}
	makeLine := func(cells ...any) []string {
		ret := []string{}
		for _, item := range cells {
			ret = append(ret, fmt.Sprint(item))
		}
		return ret
	}

	if toolCmd.User.ListUser {
		users, err := db.GetAllUsers()
		if err != nil {
			fmt.Println("Query all user failed:", err)
			return
		}
		content := [][]string{}
		for _, item := range users {
			content = append(content, makeLine(item.ID, item.Name, item.Password, item.IsAdmin, item.Banned))
		}
		fmt.Print(RenderTable(userHeader, content))
		return
	}

	if toolCmd.User.SetUser {
		err := service.GetLoginService().PutUser(db.User{
			ID:       toolCmd.User.ID,
			Name:     toolCmd.User.NickName,
			Password: toolCmd.User.Password,
			IsAdmin:  toolCmd.User.IsAdmin,
			Banned:   toolCmd.User.Ban,
		})
		if err != nil {
			fmt.Println("set user failed:", err)
		}
		return
	}
}
