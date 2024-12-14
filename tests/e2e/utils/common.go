package tests_e2e

import (
	"net/http"
)

var JSON_HEADER http.Header = http.Header{
	"content-type": {"application/json"},
}
