# Generate .wxs file from a folder of files for use in .msi installer

## Usage
| Argument | Description | Required |
| -------- | ----------- | -------- |
| --name   | Product name | Yes |
| --version | Product version, must be x.y.z | No, defaults to "1.0.0" |
| --manufacturer | Product manufacturer | Yes |
| --comments | Package comments | No, defaults to "[name] installer" |
| --dir | Directory with files to bundle | Yes |
| --out | Outfile file name | No, defaults to stdout |

More proper documentation should come soon, as well as example usages