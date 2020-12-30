// +build darwin

package build

import (
	"fmt"
	"github.com/leaanthony/slicer"
	"github.com/wailsapp/wails/v2/internal/fs"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
)

func (d *DesktopBuilder) convertToHexLiteral(bytes []byte) string {
	result := ""
	for _, b := range bytes {
		result += fmt.Sprintf("0x%x, ", b)
	}
	return result
}

// We will compile all tray icons found at <projectdir>/icons/tray/*.png into the application
func (d *DesktopBuilder) processTrayIcons(assetDir string, options *Options) error {

	var err error

	// Get all the tray icon filenames
	trayIconDirectory := filepath.Join(options.ProjectData.IconsDir, "tray")
	var trayIconFilenames []string
	if fs.DirExists(trayIconDirectory) {
		trayIconFilenames, err = filepath.Glob(trayIconDirectory + "/*.png")
		if err != nil {
			log.Fatal(err)
			return err
		}
	}
	// Setup target
	targetFilename := "trayicons"
	targetFile := filepath.Join(assetDir, targetFilename+".c")
	d.addFileToDelete(targetFile)

	var dataBytes []byte

	// Use a strings builder
	var cdata strings.Builder

	// Write header
	header := `// trayicons.c
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL.
// This file was auto-generated. DO NOT MODIFY.

`
	cdata.WriteString(header)

	var variableList slicer.StringSlicer

	// Loop over icons
	for count, filename := range trayIconFilenames {

		// Load the tray icon
		dataBytes, err = ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		iconname := strings.TrimSuffix(filepath.Base(filename), ".png")
		trayIconName := fmt.Sprintf("trayIcon%dName", count)
		variableList.Add(trayIconName)
		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { %s0x00 };\n", trayIconName, d.convertToHexLiteral([]byte(iconname))))

		trayIconLength := fmt.Sprintf("trayIcon%dLength", count)
		variableList.Add(trayIconLength)
		lengthAsString := strconv.Itoa(len(dataBytes))
		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { %s0x00 };\n", trayIconLength, d.convertToHexLiteral([]byte(lengthAsString))))

		trayIconData := fmt.Sprintf("trayIcon%dData", count)
		variableList.Add(trayIconData)
		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { ", trayIconData))

		// Convert each byte to hex
		for _, b := range dataBytes {
			cdata.WriteString(fmt.Sprintf("0x%x, ", b))
		}

		cdata.WriteString("0x00 };\n")
	}

	// Write out main trayIcons data
	cdata.WriteString("const unsigned char *trayIcons[] = { ")
	cdata.WriteString(variableList.Join(", "))
	if len(trayIconFilenames) > 0 {
		cdata.WriteString(", ")
	}
	cdata.WriteString("0x00 };\n")

	err = ioutil.WriteFile(targetFile, []byte(cdata.String()), 0600)
	if err != nil {
		return err
	}
	return nil
}

// We will compile all dialog icons found at <projectdir>/icons/dialog/*.png into the application
func (d *DesktopBuilder) processDialogIcons(assetDir string, options *Options) error {

	var err error

	// Get all the dialog icon filenames
	dialogIconDirectory := filepath.Join(options.ProjectData.IconsDir, "dialog")
	var dialogIconFilenames []string

	// If the user has no custom dialog icons, copy the defaults
	if !fs.DirExists(dialogIconDirectory) {
		defaultDialogIconsDirectory := fs.RelativePath("./internal/packager/icons/dialog")
		err := fs.CopyDir(defaultDialogIconsDirectory, dialogIconDirectory)
		if err != nil {
			return err
		}
	}

	dialogIconFilenames, err = filepath.Glob(dialogIconDirectory + "/*.png")
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Setup target
	targetFilename := "userdialogicons"
	targetFile := filepath.Join(assetDir, targetFilename+".c")
	d.addFileToDelete(targetFile)

	var dataBytes []byte

	// Use a strings builder
	var cdata strings.Builder

	// Write header
	header := `// userdialogicons.c
// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL.
// This file was auto-generated. DO NOT MODIFY.

`
	cdata.WriteString(header)

	var variableList slicer.StringSlicer

	// Loop over icons
	for count, filename := range dialogIconFilenames {

		// Load the tray icon
		dataBytes, err = ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		iconname := strings.TrimSuffix(filepath.Base(filename), ".png")
		dialogIconName := fmt.Sprintf("userDialogIcon%dName", count)
		variableList.Add(dialogIconName)
		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { %s0x00 };\n", dialogIconName, d.convertToHexLiteral([]byte(iconname))))

		dialogIconLength := fmt.Sprintf("userDialogIcon%dLength", count)
		variableList.Add(dialogIconLength)
		lengthAsString := strconv.Itoa(len(dataBytes))
		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { %s0x00 };\n", dialogIconLength, d.convertToHexLiteral([]byte(lengthAsString))))

		dialogIconData := fmt.Sprintf("userDialogIcon%dData", count)
		variableList.Add(dialogIconData)
		cdata.WriteString(fmt.Sprintf("const unsigned char %s[] = { ", dialogIconData))

		// Convert each byte to hex
		for _, b := range dataBytes {
			cdata.WriteString(fmt.Sprintf("0x%x, ", b))
		}

		cdata.WriteString("0x00 };\n")
	}

	// Write out main dialogIcons data
	cdata.WriteString("const unsigned char *userDialogIcons[] = { ")
	cdata.WriteString(variableList.Join(", "))
	if len(dialogIconFilenames) > 0 {
		cdata.WriteString(", ")
	}
	cdata.WriteString("0x00 };\n")

	err = ioutil.WriteFile(targetFile, []byte(cdata.String()), 0600)
	if err != nil {
		return err
	}
	return nil
}
