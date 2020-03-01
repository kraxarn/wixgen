package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//region Wix

type Wix struct {
	XmlNs	string	`xml:"xmlns,attr"`
	Product	Product
}

func NewWix(productName, productVersion, productManufacturer, packageComments string) Wix {
	return Wix{
		XmlNs:		"http://schemas.microsoft.com/wix/2006/wi",
		Product:	NewProduct(productName, productVersion, productManufacturer, packageComments),
	}
}

func NewWixFromArgs(args Arguments) Wix {
	return NewWix(
		args.ProductName,
		args.ProductVersion,
		args.ProductManufacturer,
		args.PackageComments)
}

//endregion

//region Product

type Product struct {
	Id				string	`xml:",attr"`
	UpgradeCode		string	`xml:",attr"`
	Name			string	`xml:",attr"`
	Version			string	`xml:",attr"`
	Manufacturer	string	`xml:",attr"`
	Language		string	`xml:",attr"`
	Package			Package
	Media			Media
	Directory		Directory
	Feature			Feature
}

func NewProduct(name, version, manufacturer, packageComments string) Product {
	return Product{
		Id:				"*",
		UpgradeCode:	"*",
		Name:			name,
		Version:		version,
		Manufacturer: 	manufacturer,
		Language: 		"1033",
		Package:		NewPackage(packageComments),
		Media:			NewMedia(),
		Directory:		*NewRootDirectory(name),
	}
}

//endregion

//region Package

type Package struct {
	InstallerVersion	string	`xml:",attr"`
	Compressed			string	`xml:",attr"`
	Comments			string	`xml:",attr"`
}

func NewPackage(comments string) Package {
	return Package{
		InstallerVersion:	"200",
		Compressed: 		"yes",
		Comments: 			comments,
	}
}

//endregion

//region Media

type Media struct {
	Id			string	`xml:",attr"`
	Cabinet		string	`xml:",attr"`
	EmbedCab	string	`xml:",attr"`
}

func NewMedia() Media {
	return Media{
		Id:			"1",
		Cabinet:	"product.cab",
		EmbedCab:	"yes",
	}
}

//endregion

//region Directory

type Directory struct {
	Id			string	`xml:",attr"`
	Name		*string	`xml:",attr"`
	Directory	[]*Directory
	Component	[]*Component
}

func NewDirectory(id string, name string, directory *Directory) *Directory {
	dir := new(Directory)
	dir.Id = id
	if name != "" {
		dir.Name = &name
	}
	if directory != nil {
		dir.Directory = []*Directory{directory}
	}
	return dir
}

func NewRootDirectory(productName string) *Directory {
	return NewDirectory("TARGETDIR", "SourceDir",
		NewDirectory("ProgramFilesFolder", "",
			NewDirectory("INSTALLDIR", productName, nil)))
}

//endregion

//region Component

type Component struct {
	Id				string	`xml:",attr"`
	Guid			string	`xml:",attr"`
	File			*File
	Shortcut		*Shortcut
	RemoveFolder	*RemoveFolder
}

func NewComponent(id string, file *File) *Component {
	cmp := new(Component)
	cmp.Id 		= id
	cmp.Guid	= "*"
	cmp.File	= file
	return cmp
}

//endregion

//region File

type File struct {
	Id		string	`xml:",attr"`
	Source	string	`xml:",attr"`
}

func NewFile(id, source string) *File {
	file := new(File)
	file.Id		= id
	file.Source	= source
	return file
}

//endregion

//region ComponentRef

type ComponentRef struct {
	Id	string	`xml:",attr"`
}

//endregion

//region Feature

type Feature struct {
	Id				string	`xml:",attr"`
	Level			string	`xml:",attr"`
	ComponentRef	[]ComponentRef
}

func NewFeature() Feature {
	return Feature{
		Id:		"DefaultFeature",
		Level:	"1",
	}
}

//endregion

//region Shortcut

type Shortcut struct {
	Id					string	`xml:",attr"`
	Name				string	`xml:",attr"`
	Description			string	`xml:",attr"`
	Target				string	`xml:",attr"`
	WorkingDirectory	string	`xml:",attr"`
}

func NewShortcut(name, execName string) *Shortcut {
	cut := new(Shortcut)
	cut.Id					= "ApplicationShortcut"
	cut.Name				= name
	cut.Description			= name
	cut.Target				= fmt.Sprintf("[INSTALLDIR]%v", execName)
	cut.WorkingDirectory	= "INSTALLDIR"
	return cut
}

//endregion

//region RemoveFolder

type RemoveFolder struct {
	Id	string	`xml:",attr"`
	On	string	`xml:",attr"`
}

func NewRemoveFolder() *RemoveFolder {
	rm := new(RemoveFolder)
	rm.Id	= "ProgramMenuSubfolder"
	rm.On	= "uninstall"
	return rm
}

//endregion

type Arguments struct {
	ProductName			string
	ProductVersion		string
	ProductManufacturer	string
	PackageComments		string
	InputDirectory		string
	OutputFile			string
	ExecName			string
}

func (args *Arguments) Validate() []string {
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
	menuSub.Component = append(menuSub.Component, NewComponent("ApplicationShortcuts", nil))
	menuSub.Component[0].Shortcut = NewShortcut(args.ProductName, args.ExecName)
	menuSub.Component[0].RemoveFolder = NewRemoveFolder()
	root.Product.Directory.Directory = append(
		root.Product.Directory.Directory,
		NewDirectory("ProgramFilesFolder", "", menuSub))

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
		"%v\n%v\n%v\n%v\n%v\n%v\n%v\n",
		"--name\t\tProduct name, required",
		"--version\tProduct version, must be x.y.z, optional, default 1.0.0",
		"--manufacturer\tProduct manufacturer, required",
		"--comments\tPackage comments, optional, default \"[name] installer\"",
		"--dir\t\tDirectory with files to bundle, required",
		"--exec\t\tMain executable in directory, required",
		"--out\t\tOutput file name, optional, default stdout")
}

func PrintVersion() {
	fmt.Println("wixgen, wix xml/wxs generator, v1.0")
}

func Validate(args *Arguments) {
	// Validate
	missing := args.Validate()
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