package main

import (
	"fmt"

	"github.com/m1a9s9a4/autofill"
)

type User struct {
	ID          int64  `autofill:"seq"`
	Email       string `autofill:"email"`
	Age         int
	WorkspaceID int64  // 全員同じワークスペース
	TeamID      string // 全員同じチーム
	Role        string // 全員同じロール
}

func main() {
	fmt.Println("=== Fixed Values Example ===")
	fmt.Println("全てのユーザーが同じWorkspaceID, TeamID, Roleを持つ")

	users := make([]User, 5)
	err := autofill.FillSlice(&users, autofill.Override{
		"Email":       autofill.Seq("user%d@example.com"), // 連番
		"Age":         autofill.SeqInt(25),                 // 連番
		"WorkspaceID": int64(12345),                        // 固定値
		"TeamID":      "engineering",                       // 固定値
		"Role":        "developer",                         // 固定値
	})

	if err != nil {
		panic(err)
	}

	for i, user := range users {
		fmt.Printf("%d: ID=%d Email=%s Age=%d WorkspaceID=%d TeamID=%s Role=%s\n",
			i+1, user.ID, user.Email, user.Age, user.WorkspaceID, user.TeamID, user.Role)
	}

	fmt.Println("\n=== Mixed: Fixed and Sequential ===")

	type Task struct {
		ID          int64  `autofill:"seq"`
		Title       string
		ProjectID   int64  // 全タスクが同じプロジェクトに属する
		Priority    int    // 優先度は連番
		Status      string // 全タスク同じステータス
	}

	tasks := make([]Task, 3)
	err = autofill.FillSlice(&tasks, autofill.Override{
		"Title":     autofill.Seq("Task #%d"),
		"ProjectID": int64(999),      // 固定: 全てプロジェクト999
		"Priority":  autofill.SeqInt(1), // 連番: 1, 2, 3
		"Status":    "pending",       // 固定: 全て"pending"
	})

	if err != nil {
		panic(err)
	}

	for i, task := range tasks {
		fmt.Printf("%d: ID=%d Title=%s ProjectID=%d Priority=%d Status=%s\n",
			i+1, task.ID, task.Title, task.ProjectID, task.Priority, task.Status)
	}
}
