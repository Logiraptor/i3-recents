package socket

import (
	"encoding/binary"
)

// To add support for a new architecture,
// add a file called byte_order_GOARCH.go
// and set the native byte order appropriately
var nativeOrder = binary.LittleEndian
