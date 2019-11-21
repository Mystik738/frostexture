package frostexture

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//ConvertToDDSandPNG converts the textures-s3 directory to dds and pngs
func ConvertToDDSandPNG(dir string, overwrite bool, png bool) {
	ddsBegin, _ := hex.DecodeString(strings.ReplaceAll("44 44 53 20 7C 00 00 00 07 10 00 00", " ", ""))
	ddsMid, _ := hex.DecodeString(strings.ReplaceAll("00 00 08 0F 00 10 00 00 0A 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 20 00 00 00 04 00 00 00", " ", ""))
	ddsUncompMid, _ := hex.DecodeString(strings.ReplaceAll("00 00 08 0F 00 10 00 00 0A 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 20 00 00 00 41 00 00 00", " ", ""))
	ddsUncompEnd, _ := hex.DecodeString(strings.ReplaceAll("15 00 00 00 20 00 00 00 FF 00 00 00 00 FF 00 00 00 00 FF 00 00 00 00 FF 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00", " ", ""))

	var files []string
	if _, err := os.Stat(dir); err == nil {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			files = append(files, info.Name())
			return nil
		})

		for _, fileName := range files {
			file, _ := os.Open(filepath.Join(dir, fileName))
			if file != nil {
				headerInfo := make([]byte, 16)
				_, err = file.ReadAt(headerInfo, 0)
				checkError(err)
				isDxt5 := string(headerInfo[8:11]) == "DXT"
				isDDS := (binary.LittleEndian.Uint32(headerInfo[8:12]) == 21 &&
					binary.LittleEndian.Uint32(headerInfo[12:]) == 1)

				if isDxt5 || isDDS {
					extIndex := strings.Index(fileName, ".")
					if extIndex == -1 {
						extIndex = len(fileName)
					}

					//If there doesn't exist a .dds file with this name or we're overwriting files
					if _, err := os.Stat(filepath.Join(dir, fileName[:extIndex]+".dds")); overwrite || err != nil {
						width := headerInfo[:4]
						height := headerInfo[4:8]
						//TODO: two files don't have the right dimensions, this is a stopgap but it's wrong, prevents a 4 GB file though
						if hex.EncodeToString(headerInfo[4:8]) == "ffffffff" {
							height = width
						}
						size := binary.LittleEndian.Uint32(width) * binary.LittleEndian.Uint32(height)
						fi, err := file.Stat()
						checkError(err)
						fileSize := fi.Size()

						newFile, err := os.Create(filepath.Join(dir, fileName[:extIndex]+".dds"))
						checkError(err)

						newFile.Write(ddsBegin)
						newFile.Write(height)
						newFile.Write(width)

						if isDxt5 {
							newFile.Write(ddsMid)
							midHeader := make([]byte, 44)
							file.ReadAt(midHeader, 8)
							newFile.Write(midHeader)
						} else {
							newFile.Write(ddsUncompMid)
							newFile.Write(ddsUncompEnd)
							size *= 4
						}

						data := make([]byte, size)
						file.ReadAt(data, fileSize-int64(size))

						newFile.Write(data)
						err = newFile.Close()
						checkError(err)
					}

					if png {
						if _, err := os.Stat(filepath.Join(dir, fileName[:extIndex]+".png")); overwrite || err != nil {
							if _, err := os.Stat(filepath.Join(dir, fileName[:extIndex]+".dds")); os.IsNotExist(err) {
								fmt.Println(fileName[:extIndex]+".dds", " doesn't exist, skipping...")
							} else {
								//Convert to PNG
								cmd := exec.Command("magick", "mogrify", "-format", "png", filepath.Join(dir, fileName[:extIndex]+".dds"))
								err = cmd.Run()
								if err != nil && err.Error() == "exec: \"magick\": executable file not found in $PATH" {
									fmt.Println("Magick not found, skipping PNG conversion")
									png = false
								}
							}
						}
					}
				}

				//Normally would defer this, but can cause a lot of open files
				file.Close()
			}
		}
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
