//+build vfsgen

package template

import "net/http"

var FS http.FileSystem = http.Dir("./_data")
