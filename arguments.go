package main

import (
	"fmt"
	"os"
	"path"
	"strings"
)

type Arguments struct {
	ProductName			string
	ProductVersion		string
	ProductManufacturer	string
	PackageComments		string
	InputDirectory		string
	OutputFile			string
	ExecName			string
	Icon				string
}

func (args *Arguments) Missing() []string {
	// List with all invalid/missing arguments
	missing := make([]string, 0)
	// ProductName is required
	if args.ProductName == "" {
		missing = append(missing, "--name")
	}
	// Product version is optional
	if args.ProductVersion == "" {
		args.ProductVersion = "1.0.0"
	}
	// Product manufacturer is required
	if args.ProductManufacturer == "" {
		missing = append(missing, "--manufacturer")
	}
	// Package comments is optional
	if args.PackageComments == "" {
		args.PackageComments = fmt.Sprintf("%v installer", args.ProductName)
	}
	// Input directory is required
	if args.InputDirectory == "" {
		missing = append(missing, "--dir")
	}
	// Executable name is required
	if args.ExecName == "" {
		missing = append(missing, "--exec")
	}
	// Output file is option, but we leave it empty as default
	return missing
}

func (args *Arguments) ExecPath() string {
	return path.Join(args.InputDirectory, args.ExecName)
}

func Validate(args *Arguments) {
	// Validate
	missing := args.Missing()
	if len(missing) > 0 {
		PrintErr(fmt.Sprintf("missing arguments: %v", strings.Join(missing, ", ")))
		PrintUsage()
		os.Exit(1)
	}
	// Check if input directory exists
	stat, err := os.Stat(args.InputDirectory)
	if os.IsNotExist(err) || !stat.IsDir() {
		PrintErr(fmt.Sprintf("\"%v\" does not exist or is not a directory", args.InputDirectory))
		PrintUsage()
		os.Exit(2)
	}
	// Check so version number looks correct
	if strings.Count(args.ProductVersion, ".") < 2 {
		PrintErr("warning: version number should be in format x.y.z")
	}
	// Check so main executable exists
	stat, err = os.Stat(args.ExecPath())
	if os.IsNotExist(err) || stat.IsDir() {
		PrintErr(fmt.Sprintf("\"%v\" does not exist", args.ExecPath()))
		PrintUsage()
		os.Exit(2)
	}
}