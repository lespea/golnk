package lnk

import (
	"fmt"
	"io"
	"os"
)

// LnkFile represents one lnk file.
type LnkFile struct {
	Header     ShellLinkHeaderSection  // File header.
	IDList     LinkTargetIDListSection // LinkTargetIDList.
	LinkInfo   LinkInfoSection         // LinkInfo.
	StringData StringDataSection       // StringData.
	DataBlocks ExtraDataSection        // ExtraData blocks.
}

// Read parses an io.Reader pointing to the contents of an lnk file.
func Read(r io.Reader) (f LnkFile, err error) {

	f.Header, err = Header(r)
	if err != nil {
		return f, fmt.Errorf("golnk.Read: parse Header - %s", err.Error())
	}

	// If HasLinkTargetIDList is set, header is immediately followed by a LinkTargetIDList.
	if f.Header.LinkFlags["HasLinkTargetIDList"] {
		f.IDList, err = LinkTarget(r)
		if err != nil {
			return f, fmt.Errorf("golnk.Read: parse LinkTarget - %s", err.Error())
		}
	}

	// If HasLinkInfo is set, read LinkInfo section.
	if f.Header.LinkFlags["HasLinkInfo"] {
		f.LinkInfo, err = LinkInfo(r)
		if err != nil {
			return f, fmt.Errorf("golnk.Read: parse LinkInfo - %s", err.Error())
		}
	}

	// Read StringData section.
	f.StringData, err = StringData(r, f.Header.LinkFlags)
	if err != nil {
		return f, fmt.Errorf("golnk.Read: parse StringData - %s", err.Error())
	}

	f.DataBlocks, err = DataBlock(r)
	if err != nil {
		return f, fmt.Errorf("golnk.Read: parse ExtraDataBlock - %s", err.Error())
	}

	return f, err
}

// File parses an lnk File.
func File(filename string) (f LnkFile, err error) {
	fi, err := os.Open(filename)
	if err != nil {
		return f, fmt.Errorf("golnk.File: open file - %s", err.Error())
	}
	defer fi.Close()

	return Read(fi)
}
