goraptor
========

A RFC 5053 implementation written by go.

There is an example:

```Go
func TestRaptorEncodeAndDecode(t *testing.T) {
    symbols := make([][]byte, 128)
    for idx := range symbols {
        symbols[idx] = make([]byte, 1024)
        for offset := 0; offset < len(symbols[idx]); offset++ {
            symbols[idx][offset] = byte(idx + offset)
        }
    }

    k := len(symbols)
    encoder := New(k)
    for esi, data := range symbols {
        encoder.Append(esi, data)
    }
    encoder.Encode()

    decoder := New(k)
    for esi := 0; esi < 2*k; esi++ {
        if esi%3 == 0 {
            continue
        }
        sym, _ := encoder.Symbol(esi)
        decoder.Append(esi, sym)
    }
    decoder.Decode()

    for esi, data := range symbols {
        sym, _ := decoder.Symbol(esi)
        if !bytes.Equal(data, sym) {
            t.Errorf("Expect %v at %v for symbol, but got %v.", data, esi, sym)
            break
        }
    }
}
```
