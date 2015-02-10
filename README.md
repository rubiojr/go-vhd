# go-vhd

Go package and CLI to work with VHD images (https://technet.microsoft.com/en-us/virtualization/bb676673.aspx)

**Highly Experimental**

* Support for printing VHD headers
* Create Dynamic (sparse) VHD images

```
govhd create foo.vhd 80GiB
```


```
govhd info foo.vhd
Cookie:                  0x636f6e6563746978 (conectix)
Features:                0x00000002
File format version:     0x00010000
Data offset:             0x0000000000000200 (512 bytes)
Timestamp:               2015-02-10 14:17:25 +0100 CET
Creator application:     go-v
Creator version:         0x00000000
Creator OS:              Wi2k
Original size:           0x0000001400000000 ( 85899345920 bytes )
Current size:            0x0000001400000000 ( 85899345920 bytes )
Disk geometry:           0xa0a010ff (c: 41120, h: 16, s: 255) (85898035200 bytes)
Disk type:               0x00000003 (Dynamic)
Checksum:                0xffffee82
UUID:                    16a1614a-f6f9-1708-a42a-3bf58ada0942
Saved state:             0

Reading dynamic/differential VHD header...
Cookie:                  0x6378737061727365 (cxsparse)
Data offset:             0xffffffffffffffff
Table offset:            0x0000000000000600
Header version:          0x00010000
Max table entries:       0x0000a000
Block size:              0x00200000
Checksum:                0xfffff3d7
Parent UUID:             00000000-0000-0000-0000-000000000000
Parent timestamp:        2000-01-01 01:00:00 +0100 CET
Reserved:                0x00000000
Parent Name:
Reserved2:               0
```
