package bencode

import (
    "sort"
    "strings"
    "strconv"
)

type encoder struct {
    builder strings.Builder
}

func Encode(v any) string {
    e := &encoder{}
    e.encodeType(v)
    return e.builder.String()
}

func (e *encoder) encodeType(v any) {
    switch t := v.(type) {
    case string:
        e.encodeString(t)
        break
    case int64:
        e.encodeInt(t)
        break
    case []any:
        e.encodeList(t)
        break
    case map[string]any:
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

func (e *encoder) encodeList(l []any) {
    e.builder.WriteByte('l')
    for _, elem := range l {
        e.encodeType(elem)
    }
    e.builder.WriteByte('e')
}

func (e *encoder) encodeDict(d map[string]any) {
    e.builder.WriteByte('d')
    keys := []string{}
    for k := range d {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    for _, k := range keys {
        e.encodeString(k)
        e.encodeType(d[k])
    }
    e.builder.WriteByte('e')
}
