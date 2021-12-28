package nsp

type nspFile struct {
	Header    pfs0Header
	FileEntry []pfs0FileEntry
	FileName  []string
}

type pfs0Header struct {
	Magic        [4]byte
	FileCnt      uint32
	StrTableSize uint32
	Reserved     uint32
}

type pfs0FileEntry struct {
	FileOffset     uint64
	FileSize       uint64
	FilenameOffset uint32
	Reserved       uint32
}

type rsa2048SHA256Ticket struct {
	SigType       uint32
	Signature     [0x100]uint8
	Padding       [0x3C]uint8
	SigIssuer     [0x40]byte
	TitlekeyBlock [0x100]uint8
	Unk1          uint8
	TitlekeyType  uint8
	Unk2          [0x03]uint8
	MasterKeyRev  uint8
	Unk3          [0x0A]uint8
	TicketID      uint64
	DeviceID      uint64
	RightsID      [0x10]uint8
	AccountID     uint32
	Unk4          [0x0C]uint8
}
