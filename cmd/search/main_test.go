package main

import (
	"testing"
)

func TestSearchAPI(t *testing.T) {
	InitialCmdFlags()

	_ = SearchCmd.Flags().Set("url", "https://gitlab.xxxx.com")
	_ = SearchCmd.Flags().Set("token", "k-xxxx")
	_ = SearchCmd.Flags().Set("keyword", "yyy")
	_ = SearchCmd.Execute()
}
