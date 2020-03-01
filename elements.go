package main

import "fmt"

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
	Upgrade			*Upgrade
	Condition		*Condition
}

func NewProduct(name, version, manufacturer, packageComments string) Product {
	return Product{
		Id:				"*",
		// Upgrade code always needs to be the same, use product name
		UpgradeCode:	GetGuid([]byte(name)),
		Name:			name,
		Version:		version,
		Manufacturer: 	manufacturer,
		Language: 		"1033",
		Package:		NewPackage(packageComments),
		Media:			NewMedia(),
		Directory:		*NewRootDirectory(name),
		Upgrade:		NewUpgrade(name, version),
		Condition:		NewCondition(),
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

//region Upgrade

type Upgrade struct {
	// Id needs to be the same as Product.UpgradeCode
	Id				string	`xml:",attr"`
	UpgradeVersion	[]UpgradeVersion
}

func NewUpgrade(productName, productVersion string) *Upgrade {
	up := new(Upgrade)
	up.Id 				= GetGuid([]byte(productName))
	up.UpgradeVersion	= GenerateUpgradeVersions(productVersion)
	return up
}

//endregion

//region UpgradeVersion

type UpgradeVersion struct {
	Minimum			string		`xml:",attr"`
	// Optional string arguments
	Maximum			interface{}	`xml:",attr"`
	IncludeMinimum	interface{}	`xml:",attr"`
	IncludeMaximum	interface{}	`xml:",attr"`
	OnlyDetect		interface{}	`xml:",attr"`
	Property		string		`xml:",attr"`
}

func GenerateUpgradeVersions(version string) []UpgradeVersion {
	return []UpgradeVersion{
		{
			Minimum:		version,
			OnlyDetect:		"yes",
			Property:		"NEWERVERSIONDETECTED",
		},
		{
			Minimum:	"0.0.0",
			Maximum:	version,
			IncludeMinimum:	"yes",
			IncludeMaximum:	"no",
			Property:		"OLDVERSIONBEINGUPGRADED",
		},
	}
}

//endregion

//region Condition

type Condition struct {
	Condition	string	`xml:",chardata"`
	Message		string	`xml:",attr"`
}

func NewCondition() *Condition {
	cond := new(Condition)
	cond.Condition	= "NOT NEWERVERSIONDETECTED"
	cond.Message	= "Product is already installed"
	return cond
}

//endregion