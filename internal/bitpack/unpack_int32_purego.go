//go:build purego || !amd64

package bitpack

import (
	"encoding/binary"
	"golang.org/x/sys/cpu"
	"github.com/parquet-go/parquet-go/internal/unsafecast"
)

func unpackInt32(dst []int32, src []byte, bitWidth uint) {
	srcLen := (len(src) / 4)
	bits := make([]uint32, srcLen)
	if cpu.IsBigEndian {
		idx := 0
		for k := 0; k < srcLen; k++ {
			bits[k] = binary.LittleEndian.Uint32((src)[idx:(4 + idx)])
			idx += 4
		}
	} else {
		bits = unsafecast.BytesToUint32(src)
	}

	bitMask := uint32(1<<bitWidth) - 1
	bitOffset := uint(0)

	for n := range dst {
		i := bitOffset / 32
		j := bitOffset % 32
		d := (bits[i] & (bitMask << j)) >> j
		if j+bitWidth > 32 {
			k := 32 - j
			d |= (bits[i+1] & (bitMask >> k)) << k
		}
		dst[n] = int32(d)
		bitOffset += bitWidth
	}
}
