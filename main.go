/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"vhagar/cmd"
	"vhagar/task"
)

func main() {
	defer task.CloseOutputFile()
	cmd.Execute()
}
