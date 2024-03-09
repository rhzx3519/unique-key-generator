package pool

import (
    "context"
    "fmt"
    "rhzx3519/unique-key-generator/sequencer"
)

const SEQ_SPAN = 100

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
