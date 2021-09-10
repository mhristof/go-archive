# go-archive

Extract files from archives in golang.

## Import

```golang
import "github.com/mhristof/go-archive"
```

and `go get github.com/mhristof/go-archive`

## Example

```golang
data, err := archive.NewURL(
    "https://github.com/terraform-linters/tflint/releases/download/v0.31.0/tflint_darwin_amd64.zip"
).ExtractFile("tflint")
if err != nil {
    panic(err)
}

err = os.WriteFile("tflint", data, 0755)
if err != nil {
    panic(err)
}
```
