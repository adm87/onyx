package tiled

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/adm87/onyx/pkg/tiled/data"
	"github.com/klauspost/compress/zstd"
)

const (
	FlippedHorizontallyFlag uint32 = 0x80000000
	FlippedVerticallyFlag   uint32 = 0x40000000
	FlippedDiagonallyFlag   uint32 = 0x20000000
	RotatedHexagonal120Flag uint32 = 0x10000000
	GidMask                 uint32 = 0x1FFFFFFF
)

func decodeContent(format data.Encoding, compression data.Compression, content string) ([]Tile, error) {
	content = strings.TrimSpace(content)

	switch format {
	case data.EncodingCSV:
		tileData, err := decodeCsv([]byte(content))
		if err != nil {
			return nil, err
		}
		return decodeTiles(tileData)

	case data.EncodingBase64:
		tileData, err := decodeBase64Compressed([]byte(content), compression)
		if err != nil {
			return nil, err
		}
		return decodeTiles(tileData)

	default:
		return nil, fmt.Errorf("unsupported encoding format: %s", format)
	}
}

func decodeBase64Compressed(content []byte, compression data.Compression) ([]uint32, error) {
	decoded, err := decodeBase64(content)
	if err != nil {
		return nil, err
	}

	var raw []byte
	switch compression {
	case data.CompressionNone:
		raw = decoded
	case data.CompressionGzip:
		raw, err = decodeGzip(decoded)
	case data.CompressionZlib:
		raw, err = decodeZlib(decoded)
	case data.CompressionZstd:
		raw, err = decodeZstd(decoded)
	default:
		return nil, fmt.Errorf("unsupported compression format: %s", compression)
	}
	if err != nil {
		return nil, err
	}

	return decodeLittleEndian(raw)
}

func decodeCsv(tileData []byte) ([]uint32, error) {
	split := bytes.Split(tileData, []byte{','})
	result := make([]uint32, len(split))
	for i, s := range split {
		value, err := strconv.Atoi(string(bytes.TrimSpace(s)))
		if err != nil {
			return nil, fmt.Errorf("failed to parse CSV value: %w", err)
		}
		result[i] = uint32(value)
	}
	return result, nil
}

func decodeBase64(tileData []byte) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(tileData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 content: %w", err)
	}
	return decoded, nil
}

func decompress(tileData []byte, newReader func(io.Reader) (io.ReadCloser, error)) ([]byte, error) {
	reader, err := newReader(bytes.NewReader(tileData))
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}
	defer reader.Close()

	uncompressed, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %w", err)
	}
	return uncompressed, nil
}

func decodeGzip(tileData []byte) ([]byte, error) {
	return decompress(tileData, func(r io.Reader) (io.ReadCloser, error) {
		return gzip.NewReader(r)
	})
}

func decodeZlib(tileData []byte) ([]byte, error) {
	return decompress(tileData, func(r io.Reader) (io.ReadCloser, error) {
		return zlib.NewReader(r)
	})
}

func decodeZstd(tileData []byte) ([]byte, error) {
	return decompress(tileData, func(r io.Reader) (io.ReadCloser, error) {
		decoder, err := zstd.NewReader(r)
		return decoder.IOReadCloser(), err
	})
}

func decodeLittleEndian(tileData []byte) ([]uint32, error) {
	if len(tileData)%4 != 0 {
		return nil, fmt.Errorf("invalid tile data length: expected a multiple of 4, got %d", len(tileData))
	}
	result := make([]uint32, len(tileData)/4)
	for i := range result {
		result[i] = uint32(tileData[i*4]) |
			uint32(tileData[i*4+1])<<8 |
			uint32(tileData[i*4+2])<<16 |
			uint32(tileData[i*4+3])<<24
	}
	return result, nil
}

func decodeTiles(tileData []uint32) ([]Tile, error) {
	tiles := make([]Tile, len(tileData))
	for i, gid := range tileData {
		tiles[i] = Tile{
			id:    gid & GidMask,
			flags: gid &^ GidMask,
		}
	}
	return tiles, nil
}
