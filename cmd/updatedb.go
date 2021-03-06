// Copyright 2017 NDP Systèmes. All Rights Reserved.
// See LICENSE file for full licensing details.

package cmd

import (
	"text/template"

	"github.com/npiganeau/yep/yep/models"
	"github.com/npiganeau/yep/yep/server"
	"github.com/spf13/cobra"
)

const updateDBFileName string = "updatedb.go"

var updateDBCmd = &cobra.Command{
	Use:   "updatedb",
	Short: "Update the database schema",
	Long:  `Synchronize the database schema with the models definitions.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectDir := "."
		if len(args) > 0 {
			projectDir = args[0]
		}
		generateAndRunFile(projectDir, updateDBFileName, updateDBTemplate)
	},
}

// UpdateDB updates the database schema. It is meant to be called from
// a project start file which imports all the project's module.
func UpdateDB(config map[string]interface{}) {
	setupConfig(config)
	connectToDB()
	models.BootStrap()
	models.SyncDatabase()
	server.LoadDataRecords()
	log.Info("Database updated successfully")
}

func initUpdateDB() {
	YEPCmd.AddCommand(updateDBCmd)
}

var updateDBTemplate = template.Must(template.New("").Parse(`
// This file is autogenerated by yep-server
// DO NOT MODIFY THIS FILE - ANY CHANGES WILL BE OVERWRITTEN

package main

import (
	"github.com/npiganeau/yep/cmd"
{{ range .Imports }}	_ "{{ . }}"
{{ end }}
)

func main() {
	cmd.UpdateDB({{ .Config }})
}
`))
