package pool

import (
    "context"
    "fmt"
    "rhzx3519/unique-key-generator/sequencer"
    "slices"
)

func init() {
    for i := 0; i < 26; i++ {
        lookup = append(lookup, byte(i+'a'))
    }
    for i := 0; i < 10; i++ {
        lookup = append(lookup, byte(i+'0'))
    }
    for i := 0; i < 26; i++ {
        lookup = append(lookup, byte(i+'A'))
    }
    lookup = append(lookup, '.')
    lookup = append(lookup, '-')
}

const SEQ_SPAN = 100

var lookup []byte

type IsExistParam struct {
    key     string
    existed chan bool
}

type Pool struct {
    cache         []string
    currentSeq    int64
    sequencer     sequencer.Sequencer
    keyStream     chan string
    isExistStream chan IsExistParam
}

func NewPool() *Pool {
    seq, _ := sequencer.NewMongoSequencer()
    return &Pool{
        keyStream:     make(chan string),
        isExistStream: make(chan IsExistParam),
        sequencer:     seq,
    }
}

func (p *Pool) close() {
    close(p.keyStream)
    close(p.isExistStream)
    p.sequencer.Reset()
}

func (p *Pool) Run(ctx context.Context) {
    go func() {
        defer fmt.Println("Pool closure exited.")
        defer p.close()
        for {
            select {
            case p.keyStream <- p.generate():
            case param := <-p.isExistStream:
                param.existed <- p.checkExist(param.key)
            case <-ctx.Done():
                return
            }
        }
    }()
}

func (p *Pool) Key() string {
    return <-p.keyStream
}

func (p *Pool) IsExist(key string) bool {
    param := IsExistParam{
        key:     key,
        existed: make(chan bool),
    }
    defer close(param.existed)
    p.isExistStream <- param
    return <-param.existed
}

func (p *Pool) generate() string {
    if len(p.cache) != 0 {
        r := p.cache[0]
        p.cache = p.cache[1:]
        return r
    }
    // generate new keys and add them into cache
    seq, _ := p.sequencer.Current()
    for i := 0; i < SEQ_SPAN; i++ {
        p.cache = append(p.cache, base64Encode(seq+int64(i)))
    }
    err := p.sequencer.Save(seq + SEQ_SPAN)
    if err != nil {
        panic(err)
    }
    p.currentSeq = seq + SEQ_SPAN

    r := p.cache[0]
    p.cache = p.cache[1:]
    return r
}

func (p *Pool) checkExist(key string) bool {
    return base64Decode(key) < p.currentSeq
}

// [a-z0-9A-Z.-]
func base64Encode(i int64) string {
    if i == 0 {
        return "a"
    }
    var bytes []byte
    for ; i != 0; i /= 64 {
        bytes = append(bytes, lookup[i%64])
    }
    slices.Reverse(bytes)
    return string(bytes)
}

func base64Decode(k string) int64 {
    var res int64
    if k == "" {
        return res
    }
    n := len(k)
    for i := 0; i < n; i++ {
        b := k[i]
        var j int
        if b >= 'a' && b <= 'z' {
            j = int(b - 'a')
        } else if b >= '0' && b <= '9' {
            j = int(b-'0') + 26
        } else if b >= 'A' && b <= 'Z' {
            j = int(b-'A') + 36
        } else if b == '.' {
            j = 62
        } else if b == '-' {
            j = 63
        }
        res = res*64 + int64(j)
    }
    return res
}
