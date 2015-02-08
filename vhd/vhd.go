package vhd

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"
	"unicode/utf16"
	"unicode/utf8"
)

const VHD_COOKIE = "0x636f6e6563746978"
const VHD_DYN_COOKIE = "0x6378737061727365"

func fmtField(name, value string) {
	fmt.Printf("%-25s%s\n", name+":", value)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func hexs(a []byte) string {
	return "0x" + hex.EncodeToString(a[:])
}

func uuid(a []byte) string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%04x",
		a[:4],
		a[4:6],
		a[6:8],
		a[8:10],
		a[10:16])
}

/*
	utf16BytesToString converts UTF-16 encoded bytes, in big or
 	little endian byte order, to a UTF-8 encoded string.
 	http://stackoverflow.com/a/15794113
*/
func utf16BytesToString(b []byte, o binary.ByteOrder) string {
	utf := make([]uint16, (len(b)+(2-1))/2)
	for i := 0; i+(2-1) < len(b); i += 2 {
		utf[i/2] = o.Uint16(b[i:])
	}
	if len(b)/2 < len(utf) {
		utf[len(utf)-1] = utf8.RuneError
	}
	return string(utf16.Decode(utf))
}

/* VHD Dynamic and Differential Header */
/*
	Cookie 8
	Data Offset 8
	Table Offset 8
	Header Version 4
	Max Table Entries 4
	Block Size 4
	Checksum 4
	Parent Unique ID 16
	Parent Time Stamp 4
	Reserved 4
	Parent Unicode Name 512
	Parent Locator Entry 1 24
	Parent Locator Entry 2 24
	Parent Locator Entry 3 24
	Parent Locator Entry 4 24
	Parent Locator Entry 5 24
	Parent Locator Entry 6 24
	Parent Locator Entry 7 24
	Parent Locator Entry 8 24
	Reserved 256
*/
type VHDExtraHeader struct {
	Cookie              [8]byte
	DataOffset          [8]byte
	TableOffset         [8]byte
	HeaderVersion       [4]byte
	MaxTableEntries     [4]byte
	BlockSize           [4]byte
	Checksum            [4]byte
	ParentUUID          [16]byte
	ParentTimestamp     [4]byte
	Reserved            [4]byte
	ParentUnicodeName   [512]byte
	ParentLocatorEntry1 [24]byte
	ParentLocatorEntry2 [24]byte
	ParentLocatorEntry3 [24]byte
	ParentLocatorEntry4 [24]byte
	ParentLocatorEntry5 [24]byte
	ParentLocatorEntry6 [24]byte
	ParentLocatorEntry7 [24]byte
	ParentLocatorEntry8 [24]byte
	Reserved2           [256]byte
}

func (header *VHDExtraHeader) CookieString() string {
	return string(header.Cookie[:])
}

/* VHD Header */
/*
 Cookie 8
 Features 4
 File Format Version 4
 Data Offset 8
 Time Stamp 4
 Creator Application 4
 Creator Version 4
 Creator Host OS 4
 Original Size 8
 Current Size 8
 Disk Geometry 4
 Disk Type 4
 Checksum 4
 Unique Id 16
 Saved State 1
 Reserved 427
*/
type VHDHeader struct {
	Cookie             [8]byte
	Features           [4]byte
	FileFormatVersion  [4]byte
	DataOffset         [8]byte
	Timestamp          [4]byte
	CreatorApplication [4]byte
	CreatorVersion     [4]byte
	CreatorHostOS      [4]byte
	OriginalSize       [8]byte
	CurrentSize        [8]byte
	DiskGeometry       [4]byte
	DiskType           [4]byte
	Checksum           [4]byte
	UniqueId           [16]byte
	SavedState         [1]byte
	Reserved           [427]byte
}

func (h *VHDHeader) DiskTypeStr() (dt string) {
	switch h.DiskType[3] {
	case 0x00:
		dt = "None"
	case 0x01:
		dt = "Deprecated"
	case 0x02:
		dt = "Fixed"
	case 0x03:
		dt = "Dynamic"
	case 0x04:
		dt = "Differential"
	case 0x05:
		dt = "Reserved"
	case 0x06:
		dt = "Reserved"
	default:
		panic("Invalid disk type detected!")
	}

	return
}

func readVHDExtraHeader(f *os.File) {
	vhdHeader := make([]byte, 1024)
	_, err := f.Read(vhdHeader)
	check(err)

	var header VHDExtraHeader
	binary.Read(bytes.NewBuffer(vhdHeader[:]), binary.BigEndian, &header)

	fmtField("Cookie", fmt.Sprintf("%s (%s)",
	         hexs(header.Cookie[:]), header.CookieString()))
	fmtField("Data offset", hexs(header.DataOffset[:]))
	fmtField("Table offset", hexs(header.TableOffset[:]))
	fmtField("Header version", hexs(header.HeaderVersion[:]))
	fmtField("Max table entries", hexs(header.MaxTableEntries[:]))
	fmtField("Block size", hexs(header.BlockSize[:]))
	fmtField("Checksum", hexs(header.Checksum[:]))
	fmtField("Parent UUID", uuid(header.ParentUUID[:]))

	// Seconds since January 1, 1970 12:00:00 AM in UTC/GMT.
	// 946684800 = January 1, 2000 12:00:00 AM in UTC/GMT.
	tstamp := binary.BigEndian.Uint32(header.ParentTimestamp[:])
	t := time.Unix(int64(946684800+tstamp), 0)
	fmtField("Parent timestamp", fmt.Sprintf("%s", t))

	fmtField("Reserved", hexs(header.Reserved[:]))
	parentName := utf16BytesToString(header.ParentUnicodeName[:],
		binary.BigEndian)
	fmtField("Parent Name", parentName)
	// Parent locator entries ignored since it's a dynamic disk
	sum := 0
	for _, b := range header.Reserved2 {
		sum += int(b)
	}
	fmtField("Reserved2", strconv.Itoa(sum))
}

func readVHDHeader(vhdHeader []byte) VHDHeader {

	var header VHDHeader
	binary.Read(bytes.NewBuffer(vhdHeader[:]), binary.BigEndian, &header)

	//fmtField("Cookie", string(header.Cookie[:]))
	fmtField("Cookie", fmt.Sprintf("%s (%s)",
	         hexs(header.Cookie[:]), string(header.Cookie[:])))
	fmtField("Features", hexs(header.Features[:]))
	fmtField("File format version", hexs(header.FileFormatVersion[:]))

	dataOffset := binary.BigEndian.Uint64(header.DataOffset[:])
	fmtField("Data offset",
		fmt.Sprintf("%s (%d bytes)", hexs(header.DataOffset[:]), dataOffset))

	//// Seconds since January 1, 1970 12:00:00 AM in UTC/GMT.
	//// 946684800 = January 1, 2000 12:00:00 AM in UTC/GMT.
	t := time.Unix(int64(946684800+binary.BigEndian.Uint32(header.Timestamp[:])), 0)
	fmtField("Timestamp", fmt.Sprintf("%s", t))

	fmtField("Creator application", string(header.CreatorApplication[:]))
	fmtField("Creator version", hexs(header.CreatorVersion[:]))
	fmtField("Creator OS", string(header.CreatorHostOS[:]))

	originalSize := binary.BigEndian.Uint64(header.OriginalSize[:])
	fmtField("Original size",
		fmt.Sprintf("%s ( %d bytes )", hexs(header.OriginalSize[:]), originalSize))

	currentSize := binary.BigEndian.Uint64(header.OriginalSize[:])
	fmtField("Current size",
		fmt.Sprintf("%s ( %d bytes )", hexs(header.CurrentSize[:]), currentSize))

	cilinders := int64(binary.BigEndian.Uint16(header.DiskGeometry[:2]))
	heads := int64(header.DiskGeometry[2])
	sectors := int64(header.DiskGeometry[3])
	dsize := cilinders * heads * sectors * 512
	fmtField("Disk geometry",
		fmt.Sprintf("%s (c: %d, h: %d, s: %d) (%d bytes)",
			hexs(header.DiskGeometry[:]),
			cilinders,
			heads,
			sectors,
			dsize))

	fmtField("Disk type",
		fmt.Sprintf("%s (%s)", hexs(header.DiskType[:]), header.DiskTypeStr()))

	fmtField("Checksum", hexs(header.Checksum[:]))
	fmtField("UUID", uuid(header.UniqueId[:]))
	fmtField("Saved state", fmt.Sprintf("%d", header.SavedState[0]))

	return header
}

// Return the number of blocks in the disk, diskSize in bytes
func getMaxTableEntries(diskSize uint64) uint64 {
	return diskSize * (2 * 1024 * 1024) // block size is 2M
}

//func createSparseVHD(size uint64, name string) {
//	cookie, err := hex.DecodeString(VHD_COOKIE)
//	check(err)
//	fmt.Println(hexs(cookie[:]))
//
//	header := VHDHeader{
//		//Cookie: cookie,
//	}
//
//	fmt.Println(header)
//}

func PrintVHDHeaders(f *os.File) {
	vhdHeader := make([]byte, 512)
	_, err := f.Read(vhdHeader)
	check(err)
	header := readVHDHeader(vhdHeader)

	if header.DiskType[3] == 0x3 || header.DiskType[3] == 0x04 {
		fmt.Println("\nReading dynamic/differential VHD header...")
		readVHDExtraHeader(f)
	}
}
