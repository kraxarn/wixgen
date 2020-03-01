# Generate .wxs file from a folder of files for use in .msi installer

## Install
`go get -v github.com/kraxarn/wixgen`

## Usage
| Argument | Description | Required |
| -------- | ----------- | -------- |
| --name   | Product name | Yes |
| --version | Product version, must be x.y.z | No, defaults to "1.0.0" |
| --manufacturer | Product manufacturer | Yes |
| --comments | Package comments | No, defaults to "[name] installer" |
| --dir | Directory with files to bundle | Yes |
| --exec | Executable path relative to input directory | Yes |
| --icon | Icon for uninstall in control panel | No, defaults to no icon |
| --out | Outfile file name | No, defaults to stdout |

More proper documentation should come soon, as well as example usages