package common

import (
    "bytes"
    "crypto/rand"
    "crypto/rsa"
    "encoding/hex"
    "fmt"
    "github.com/shopspring/decimal"
    "hash"
    "time"
)

// 使用RSA-OAEP加密
func EncryptOAEP(hash hash.Hash, pub *rsa.PublicKey, msg []byte, label []byte) ([]byte, error) {
    buffer := bytes.NewBuffer(nil)
    chunks := bytesSplit(msg, pub.Size()-2*hash.Size()-2)
    for _, chunk := range chunks {
        ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, chunk, label)
        if err != nil {
            return nil, err
        }
        buffer.Write(ciphertext)
    }
    return buffer.Bytes(), nil
}

// 使用RSA-OAEP解密
func DecryptOAEP(hash hash.Hash, priv *rsa.PrivateKey, ciphertext []byte, label []byte) ([]byte, error) {
    buffer := bytes.NewBuffer(nil)
    chunks := bytesSplit(ciphertext, priv.Size())
    for _, chunk := range chunks {
        plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, chunk, label)
        if err != nil {
            return nil, err
        }
        buffer.Write(plaintext)
    }
    return buffer.Bytes(), nil
}

// bytes切割
func bytesSplit(data []byte, limit int) [][]byte {
    var chunk []byte
    chunks := make([][]byte, 0, len(data)/limit+1)
    for len(data) >= limit {
        chunk, data = data[:limit], data[limit:]
        chunks = append(chunks, chunk)
    }
    if len(data) > 0 {
        chunks = append(chunks, data[:len(data)])
    }
    return chunks
}

// bytes转字符串
func BytesToString(data []byte) string {
    n := bytes.IndexByte(data, 0)
    if n == -1 {
        return string(data)
    }
    return string(data[:n])
}

// 字符串转BCD
func StringToBCD(s string, size ...int) []byte {
    if (len(s) & 1) != 0 {
        s = "0" + s
    }

    data := []byte(s)
    bcd := make([]byte, len(s)/2)
    for i := 0; i < len(bcd); i++ {
        high := data[i*2] - '0'
        low := data[i*2+1] - '0'
        bcd[i] = (high << 4) | low
    }

    if len(size) == 0 {
        return bcd
    }

    ret := make([]byte, size[0])
    if size[0] < len(bcd) {
        copy(ret, bcd)
    } else {
        copy(ret[len(ret)-len(bcd):], bcd)
    }
    return ret
}

// BCD转字符串
func BcdToString(data []byte, ignorePadding ...bool) string {
    for {
        if len(data) == 0 {
            return ""
        }
        if data[0] != 0 {
            break
        }
        data = data[1:]
    }

    buf := make([]byte, 0, len(data)*2)
    for i := 0; i < len(data); i++ {
        buf = append(buf, data[i]&0xf0>>4+'0')
        buf = append(buf, data[i]&0x0f+'0')
    }

    if len(ignorePadding) == 0 || !ignorePadding[0] {
        for idx := range buf {
            if buf[idx] != '0' {
                return string(buf[idx:])
            }
        }
    }
    return string(buf)
}

// 转为BCD时间
func toBCDTime(t time.Time) []byte {
    t = time.Unix(t.Unix(), 0)
    s := t.Format("20060102150405")[2:]
    return StringToBCD(s, 6)
}

// 转为time.Time
func fromBCDTime(bcd []byte) (time.Time, error) {
    if len(bcd) != 6 {
        return time.Time{}, nil
    }
    t, err := time.ParseInLocation(
        "20060102150405", "20"+BcdToString(bcd), time.Local)
    if err != nil {
        return time.Time{}, err
    }
    return t, nil
}

// 获取经纬度
func GetGeoPoint(lat uint32, south bool, lng uint32, west bool) (decimal.Decimal, decimal.Decimal) {
    div := decimal.NewFromFloat(1000000)
    fLat := decimal.NewFromInt(int64(lat)).Div(div)
    fLon := decimal.NewFromInt(int64(lng)).Div(div)
    if south {
        fLat = decimal.Zero.Sub(fLat)
    }
    if west {
        fLon = decimal.Zero.Sub(fLon)
    }
    return fLat.Truncate(6), fLon.Truncate(6)
}

func fromStrTime(str []byte) (time.Time, error) {
    if len(str) != 6 {
        return time.Time{}, nil
    }
    t, err := time.ParseInLocation(
        "20060102150405", timeToString(str), time.Local)
    if err != nil {
        return time.Time{}, err
    }
    return t, nil
}

func timeToString(data []byte) string {
    year := "20" + fmt.Sprintf("%02d", data[0])
    month := fmt.Sprintf("%02d", data[1])
    day := fmt.Sprintf("%02d", data[2])
    hour := fmt.Sprintf("%02d", data[3])
    minute := fmt.Sprintf("%02d", data[4])
    second := fmt.Sprintf("%02d", data[5])

    return year + month + day + hour + minute + second
}

func GetHex(data []byte) []byte {
    dstEncode := make([]byte, hex.EncodedLen(len(data)))
    hex.Encode(dstEncode, data)
    rst, err := hex.DecodeString(string(dstEncode))
    if err != nil {
        return nil
    } else {
        return rst
    }
}

func GetBit(value int, offset int) int {
    return value & (1 << offset) >> offset
}

func Transform808(data []byte) []byte {
    buffer := bytes.NewBuffer(nil)
    buffer.Grow(len(data) + 10)

    buffer.WriteByte(0x7e)
    for _, b := range data[1 : len(data)-1] {
        if b == 0x7e {
            buffer.WriteByte(0x7d)
            buffer.WriteByte(0x02)
        } else if b == 0x7d {
            buffer.WriteByte(0x7d)
            buffer.WriteByte(0x01)
        } else {
            buffer.WriteByte(b)
        }
    }
    buffer.WriteByte(0x7e)

    transformed := buffer.Bytes()

    return transformed
}
