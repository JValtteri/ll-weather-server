package main

import (
    "testing"
    "encoding/base64"
)

func TestDecode64(t *testing.T) {
    const str string = "abc123-5"
    var encoded string = base64.StdEncoding.EncodeToString([]byte(str))
    var decoded string = decode64(encoded)
    //decodedExp := base64.StdEncoding.DecodeString(encoded)
    if decoded != str {
        t.Errorf("Expected: %s, got: %s\n", str, decoded)
    }
}

func TestSanitizeAllow(t *testing.T) {
    const input string = "asd 123-5"
    var sanitized string = sanitize(input)
    if sanitized != input {
        t.Errorf("Expected: %s, got: %s\n", input, sanitized)
    }
}

func TestSanitizeBlock(t *testing.T) {
    const input string = "<asd=1,23>"
    const expect string = "asd123"
    var sanitized string = sanitize(input)
    if sanitized != expect {
        t.Errorf("Expected: '%s', got: '%s'\n", expect, sanitized)
    }
}

func TestWhitelistAllow(t *testing.T) {
    const input string = "gfs_graphcast025"
    var allow bool = whitelist(input, modelList)
    if !allow {
        t.Errorf("Expected: %v, got: %v\n", true, allow)
    }
}

func TestWhitelistBlock(t *testing.T) {
    const input string = "gfs_"
    var allow bool = whitelist(input, modelList)
    if allow {
        t.Errorf("Expected: %v, got: %v\n", false, allow)
    }
}
