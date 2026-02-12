package bencode

import (
    "errors"
    "strings"
    "strconv"
)

type decoder struct {
    in     string
    cursor int
}

func Decode(bc string) (map[string]any, error) {
    d := &decoder{bc, 0}
    if b, err := d.readByte(); err != nil {
        return make(map[string]any), err
    } else if b != 'd' {
        return make(map[string]any),
                errors.New("failed to read dictionary.")
    }
    return d.readDict()
}

func (d *decoder) readByte() (byte, error) {
    if d.cursor >= len(d.in) {
        return ' ', errors.New("trying to read more bytes than input.")
    }
    b := d.in[d.cursor]
    d.cursor++
    return byte(b), nil
}

func (d *decoder) unreadByte() error {
    if d.cursor < 0 {
        return errors.New("trying to read before start of input.")
    }
    d.cursor--
    return nil
}

func (d *decoder) readDict() (map[string]any, error) {
    dict := make(map[string]any)
    for {
        key, err := d.readString()
        if err != nil {
            return nil, err
        }
        value, err := d.readValue()
        if err != nil {
            return nil, err
        }
        dict[key] = value

        b, err := d.readByte()
        if err != nil {
            return nil, err
        }

        if b == 'e' {
            break
        } else if err := d.unreadByte(); err != nil {
            return nil, err
        }
    }
    return dict, nil
}

func (d *decoder) readString() (string, error) {
    l, err := d.readIntUntil(':')
    if err != nil {
        return "", err
    }
    var sLen int64
    var ok bool
    if sLen, ok = l.(int64); !ok {
        return "", errors.New("string length must not surpass int64 limits")
    }
    if sLen < 0 {
        return "", errors.New("string length must not be negative")
    }

    s := strings.Builder{}
    for i := int64(0); i < sLen; i++ {
        b, err := d.readByte()
        if err != nil {
            return "", err
        }
        s.WriteByte(b)
    }
    return s.String(), nil
}

func (d *decoder) readIntUntil(c byte) (any, error) {
    b, err := d.readByte()
    if err != nil {
        return nil, err
    }
    valueString := strings.Builder{}
    for b != c {
        valueString.WriteByte(b)
        b, err = d.readByte()
        if err != nil {
            return nil, err
        }
    }
    value := valueString.String()
    if v, err := strconv.ParseInt(value, 10, 64); err == nil {
        return v, nil
    } else if v, err := strconv.ParseUint(value, 10, 64); err == nil {
        return v, nil
    }

    return nil, err
}

func (d *decoder) readValue() (v any, err error) {
    typ, err := d.readByte()
    if err != nil {
        return nil, err
    }
    switch typ {
    case 'i':
        v, err = d.readInt()
        break
    case 'l':
        v, err = d.readList()
        break
    case 'd':
        v, err = d.readDict()
        break
    default:
        if err = d.unreadByte(); err != nil {
            return nil, err
        }
        v, err = d.readString()
    }

    return v, err
}

func (d *decoder) readInt() (any, error) {
    return d.readIntUntil('e')
}

func (d *decoder) readList() ([]any, error) {
    var l []any
    for {
        v, err := d.readValue()
        if err != nil {
            return nil, err
        }
        l = append(l, v)

        b, err := d.readByte()
        if err != nil {
            return nil, err
        }

        if b == 'e' {
            break
        } else if err := d.unreadByte(); err != nil {
            return nil, err
        }
    }
    return l, nil
}
