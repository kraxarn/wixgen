package main

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	colorReset	= "\033[0;39;49m"
	colorRed	= "\033[0;31;49m"
)

func PrintErr(message interface{}) {
	msg := fmt.Sprintf("%v%v%v\n", colorRed, message, colorReset)
	if _, err := fmt.Fprint(os.Stderr, msg); err != nil {
		fmt.Print(msg)
	}
}

func main() {
	args := Arguments{}
	for i, arg := range os.Args {
		// Always print help
		if arg == "--help" || arg == "-?" {
			PrintVersion()
			PrintUsage()
			os.Exit(0)
		}
		// We already printed version information
		if arg == "--version" {
			PrintVersion()
			os.Exit(0)
		}
		// Ignore last arg as there's no value
		if i >= len(os.Args) - 1 {
			break
		}
		// Switch argument
		switch arg {
		case "--name":			args.ProductName			= os.Args[i + 1]
		case "--version":		args.ProductVersion			= os.Args[i + 1]
		case "--manufacturer":	args.ProductManufacturer	= os.Args[i + 1]
		case "--comments":		args.PackageComments		= os.Args[i + 1]
		case "--dir":			args.InputDirectory			= os.Args[i + 1]
		case "--out":			args.OutputFile				= os.Args[i + 1]
		case "--exec":			args.ExecName				= os.Args[i + 1]
		case "--icon":			args.Icon					= os.Args[i + 1]
		}
	}
	Validate(&args)

	// Prepare root
	root := NewWixFromArgs(args)
	installDir := root.Product.Directory.Directory[0].Directory[0]

	// Check for all files
	i := 0
	di := 0
	err := filepath.Walk(args.InputDirectory, func(path string, info os.FileInfo, err error) error {
		if path == args.InputDirectory {
			return nil
		}
		if info.IsDir() {
			installDir.Directory = append(
				installDir.Directory,
				NewDirectory(fmt.Sprintf("Dir%v", di), info.Name(), nil))
			di++
			return nil
		}

		filePath := path[len(args.InputDirectory) + 1:]
		full := strings.Split(filePath, "/")
		//ci := 0
		if len(full) > 1 {
			for _, subDir := range installDir.Directory {
				if *subDir.Name == full[0] {
					subDir.Component = append(
						subDir.Component,
						NewComponent(
							fmt.Sprintf("File%v", i),
							NewFile(info.Name(), path)))
					i++
					return nil
				}
			}
		}

		installDir.Component = append(installDir.Component,
			NewComponent(
				fmt.Sprintf("File%v", i),
				NewFile(info.Name(), path)))
		i++
		return nil
	})

	// Create needed references
	root.Product.Feature = NewFeature()
	for j := 0; j < i; j++ {
		root.Product.Feature.ComponentRef = append(
			root.Product.Feature.ComponentRef,
			ComponentRef{
				Id: fmt.Sprintf("File%v", j),
			},
		)
	}

	// Create start menu directory
	menuSub := NewDirectory("ProgramMenuSubfolder", args.ProductName, nil)
	menuComp := NewComponent("ApplicationShortcuts", nil)
	menuComp.Guid = GetGuid([]byte(args.ExecName))
	menuSub.Component = append(menuSub.Component, menuComp)
	menuSub.Component[0].Shortcut = NewShortcut(args.ProductName, args.ExecName)
	menuSub.Component[0].RemoveFolder = NewRemoveFolder()
	root.Product.Directory.Directory = append(
		root.Product.Directory.Directory,
		NewDirectory("ProgramMenuFolder", "", menuSub))

	// Create component ref
	root.Product.Feature.ComponentRef = append(
			root.Product.Feature.ComponentRef,
			ComponentRef{
				Id: "ApplicationShortcuts",
			},
		)

	data, err := xml.MarshalIndent(root, "", "\t")
	if err != nil {
		PrintErr(err)
	}

	// See how we should output
	if args.OutputFile == "" {
		fmt.Printf(string(data))
	} else {
		err := ioutil.WriteFile(args.OutputFile, data, 0644)
		if err != nil {
			PrintErr(err)
		}
	}
}

func PrintUsage() {
	fmt.Printf(
		"%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n",
		"--name\t\tProduct name, required",
		"--version\tProduct version, must be x.y.z, optional, default 1.0.0",
		"--manufacturer\tProduct manufacturer, required",
		"--comments\tPackage comments, optional, default \"[name] installer\"",
		"--dir\t\tDirectory with files to bundle, required",
		"--exec\t\tMain executable in directory, required",
		"--icon\t\tIcon for start menu, optional, default no icon",
		"--out\t\tOutput file name, optional, default stdout")
}

func PrintVersion() {
	fmt.Println("wixgen, wix xml/wxs generator")
}

func GetGuid(data []byte) string {
	// [4]-[2]-[2]-[2]-[6]	=> 16 (32)
	dataHash := md5.Sum(data)
	return fmt.Sprintf("%02d%02d%02d%02d-%x-%x-%02d%02d-%x",
		dataHash[0] % 100, dataHash[1] % 100, dataHash[2] % 100, dataHash[3] % 100,
		string(dataHash[4:6]), string(dataHash[6:8]),
		dataHash[8] % 100, dataHash[9] % 100,
		string(dataHash[10:16]))
}