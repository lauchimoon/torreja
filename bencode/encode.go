package bencode

import (
    "strings"
    "strconv"
)

type encoder struct {
    builder strings.Builder
}

func Encode(v interface{}) string {
    e := &encoder{}
    e.encodeType(v)
    return e.builder.String()
}

func (e *encoder) encodeType(v interface{}) {
    switch t := v.(type) {
    case string:
        e.encodeString(t)
        break
    case int64:
        e.encodeInt(t)
        break
    case []interface{}:
        e.encodeList(t)
        break
    case map[string]interface{}:
        e.encodeDict(t)
        break
    }
}

func (e *encoder) encodeString(s string) {
    sLen := int64(len(s))
    e.builder.WriteString(strconv.FormatInt(sLen, 10))
    e.builder.WriteByte(':')
    e.builder.WriteString(s)
}

func (e *encoder) encodeInt(i int64) {
    e.builder.WriteByte('i')
    e.builder.WriteString(strconv.FormatInt(i, 10))
    e.builder.WriteByte('e')
}

func (e *encoder) encodeList(l []interface{}) {
    e.builder.WriteByte('l')
    for _, elem := range l {
        e.encodeType(elem)
    }
    e.builder.WriteByte('e')
}

func (e *encoder) encodeDict(d map[string]interface{}) {
    e.builder.WriteByte('d')
    for k, v := range d {
        e.encodeType(k)
        e.encodeType(v)
    }
    e.builder.WriteByte('e')
}
