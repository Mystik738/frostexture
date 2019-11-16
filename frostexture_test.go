package frostexture_test

import (
	"encoding/hex"
	"frostexture"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestConvertToDDSandPNG(t *testing.T) {
	compressed_data, _ := hex.DecodeString(strings.ReplaceAll("04 00 00 00 FF FF FF FF 44 58 54 35 01 00 00 00 04 00 04 00 10 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 FF 49 92 24 49 92 24 FF FF FF FF 00 00 00 00", " ", ""))
	uncompressed_data, _ := hex.DecodeString(strings.ReplaceAll("04 00 00 00 04 00 00 00 15 00 00 00 01 00 00 00 04 00 04 00 10 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF", " ", ""))

	uncompressed_dds_expected, _ := hex.DecodeString(strings.ReplaceAll("44 44 53 20 7C 00 00 00 07 10 00 00 04 00 00 00 04 00 00 00 00 00 08 0F 00 10 00 00 0A 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 20 00 00 00 41 00 00 00 15 00 00 00 20 00 00 00 FF 00 00 00 00 FF 00 00 00 00 FF 00 00 00 00 FF 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF 00 00 00 FF", " ", ""))
	compressed_dds_expected, _ := hex.DecodeString(strings.ReplaceAll("44 44 53 20 7C 00 00 00 07 10 00 00 04 00 00 00 04 00 00 00 00 00 08 0F 00 10 00 00 0A 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 20 00 00 00 04 00 00 00 44 58 54 35 01 00 00 00 04 00 04 00 10 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 FF 49 92 24 49 92 24 FF FF FF FF 00 00 00 00", " ", ""))

	dir, _ := os.Getwd()
	compressed, _ := os.Create("white")
	uncompressed, _ := os.Create("black")

	compressed.Write(compressed_data)
	uncompressed.Write(uncompressed_data)

	compressed.Close()
	uncompressed.Close()

	frostexture.ConvertToDDSandPNG(dir, true, true)

	if out, _ := ioutil.ReadFile("white.dds"); !reflect.DeepEqual(out, compressed_dds_expected) {
		t.Errorf("Compressed conversion to DDS failed")
	}
	if out, _ := ioutil.ReadFile("black.dds"); !reflect.DeepEqual(out, uncompressed_dds_expected) {
		t.Errorf("Uncompressed conversion to DDS failed")
	}

	//Since ImageMagick handles the PNG, we don't really need to test content, just that they exist
	if _, err := os.Stat("white.png"); os.IsNotExist(err) {
		t.Errorf("Compressed conversion to PNG failed")
	}
	if _, err := os.Stat("black.png"); os.IsNotExist(err) {
		t.Errorf("Uncompressed conversion to PNG failed")
	}

	os.Remove("white")
	os.Remove("black")
	os.Remove("white.dds")
	os.Remove("black.dds")
	os.Remove("white.png")
	os.Remove("black.png")
}
