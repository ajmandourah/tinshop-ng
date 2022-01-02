package nsp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"unsafe"
)

// IsTicketValid return if ticket is valid or not
func IsTicketValid(file io.ReadSeeker, titleDBKey string, debugTicket bool) (bool, error) {
	newNSP := &nspFile{}

	// Read Header
	nspHeader := pfs0Header{}
	data := make([]byte, unsafe.Sizeof(nspHeader))
	_, _ = file.Read(data)
	buffer := bytes.NewBuffer(data)
	_ = binary.Read(buffer, binary.LittleEndian, &nspHeader)
	newNSP.Header = nspHeader

	if string(newNSP.Header.Magic[:]) != "PFS0" {
		return false, errors.New("header Magic is not present")
	}

	// Read file entry
	for i := 0; i < int(nspHeader.FileCnt); i++ {
		nspEntry := pfs0FileEntry{}
		data := make([]byte, unsafe.Sizeof(nspEntry))
		_, _ = file.Read(data)
		buffer := bytes.NewBuffer(data)
		_ = binary.Read(buffer, binary.LittleEndian, &nspEntry)
		newNSP.FileEntry = append(newNSP.FileEntry, nspEntry)
	}

	// Read nspStrTable + Display file_name
	nspStrTable := make([]byte, nspHeader.StrTableSize)
	_, _ = file.Read(nspStrTable)

	var tikOffset int
	var tikSize uint64
	var ticketFound bool

	for i := 0; i < int(nspHeader.FileCnt); i++ {
		start := newNSP.FileEntry[i].FilenameOffset
		if i != int(nspHeader.FileCnt)-1 {
			end := newNSP.FileEntry[i+1].FilenameOffset - 1
			newNSP.FileName = append(newNSP.FileName, string(nspStrTable[start:end]))
		} else {
			newNSP.FileName = append(newNSP.FileName, string(nspStrTable[start:]))
		}

		// Compute Ticket information
		if newNSP.FileName[i][len(newNSP.FileName[i])-4:] == ".tik" {
			ticketFound = true
			tikOffset = int(unsafe.Sizeof(nspHeader)) + (int(unsafe.Sizeof(newNSP.FileEntry[i])) * len(newNSP.FileEntry)) + len(nspStrTable) + int(newNSP.FileEntry[i].FileOffset)
			tikSize = newNSP.FileEntry[i].FileSize

			if tikSize != eTicketTIKFileSize {
				msg := "Ticket size mismatch (" + fmt.Sprint(tikSize) + "vs" + fmt.Sprint(eTicketTIKFileSize) + ")"
				return false, errors.New(msg)
			}
		}
	}

	// If no ticket we handle it as valid
	if !ticketFound {
		return true, nil
	}

	// Retrieve Ticket content
	ticket := &rsa2048SHA256Ticket{}
	_, _ = file.Seek(int64(tikOffset), 0)

	data = make([]byte, tikSize)
	_, _ = file.Read(data)
	buffer = bytes.NewBuffer(data)
	err := binary.Read(buffer, binary.LittleEndian, ticket)
	if err != nil {
		return false, err
	}

	var titleKey []byte
	for i := 0; i < 16; i++ {
		titleKey = append(titleKey, ticket.TitlekeyBlock[i])
	}
	var ticketKey = strings.ToUpper(hex.EncodeToString(titleKey))

	if debugTicket {
		if ticketKey == "00000000000000000000000000000000" {
			log.Println("Missing Ticket Key")
		}
	}

	if string(ticket.SigIssuer[:26]) != genericIssuer || ticketKey != titleDBKey || ticket.TicketID != 0 || ticket.DeviceID != 0 || ticket.AccountID != 0 || ticket.TitlekeyType != eTicketTitleKeyCommon {
		return false, nil
	}
	return true, nil
}
