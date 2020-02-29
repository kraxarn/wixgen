package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
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
	Id		string
	Name	string
}

func NewDirectory(id, name string) Directory {
	return Directory{
		Id:		id,
		Name:	name,
	}
}

//endregion

type Arguments struct {
	ProductName			string
	ProductVersion		string
	ProductManufacturer	string
	PackageComments		string
	InputDirectory		string
	OutputFile			string
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
	// Output file is option, but we leave it empty as default
	return missing
}

const (
	colorReset	= "\033[0;39;49m"
	colorRed	= "\033[0;31;49m"
)

func main() {
	// Print header for no reason
	PrintVersion()
	args := Arguments{}
	for i, arg := range os.Args {
		// Always print help
		if arg == "--help" || arg == "-?" {
			PrintUsage()
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
		}
	}
	// Validate
	missing := args.Validate()
	if len(missing) > 0 {
		fmt.Printf(
			"%vmissing arguments: %v%v\n\n",
			colorRed,
			strings.Join(missing, ", "),
			colorReset)
		PrintUsage()
		os.Exit(1)
	}
	// Check if input directory exists
	stat, err := os.Stat(args.InputDirectory)
	if os.IsNotExist(err) || !stat.IsDir() {
		fmt.Printf("%v\"%v\" does not exist or is not a directory%v\n\n",
			colorRed, args.InputDirectory, colorReset)
		PrintUsage()
		os.Exit(2)
	}
	// Check so version number looks correct
	if strings.Count(args.ProductVersion, ".") < 2 {
		fmt.Fprintf(os.Stderr, "warning: version number should be in format x.y.z")
	}

	root := NewWixFromArgs(args)
	data, err := xml.MarshalIndent(root, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile("out.xml", data, 0644)
}

func PrintUsage() {
	fmt.Printf(
		"%v\n%v\n%v\n%v\n%v\n\n%v\n",
		"--name\t\tProduct name, required",
		"--version\tProduct version, must be x.y.z, optional, default 1.0.0",
		"--manufacturer\tProduct manufacturer, required",
		"--comments\tPackage comments, optional, default \"[name] installer\"",
		"--dir\t\tDirectory with files to bundle, required",
		"--out\t\tOutput file name, optional, default stdout")
}

func PrintVersion() {
	fmt.Printf("\nwixgen, wix xml/wxs generator, v1.0\n\n")
}