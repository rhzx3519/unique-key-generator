package pool

import (
	"context"
	"fmt"
	"log"
	"rhzx3519/unique-key-generator/generator"
	"time"
)

const OneTimeGeneratedNumber = 100

type IsExistParam struct {
	key     string
	existed chan bool
}

type GetKeyParam struct {
	key chan string
}

type Pool struct {
	cache         []string
	produced      int64
	currentSeq    int64
	generator     generator.Generator
	keyStream     chan *GetKeyParam
	existedStream chan IsExistParam
}

func NewPool() *Pool {
	g, _ := generator.NewGenerator(false)
	return &Pool{
		keyStream:     make(chan *GetKeyParam),
		existedStream: make(chan IsExistParam),
		generator:     g,
	}
}

func (p *Pool) close() {
	close(p.keyStream)
	close(p.existedStream)
}

func (p *Pool) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 30)
	go func() {
		defer fmt.Println("Pool closure exited.")
		defer ticker.Stop()
		defer p.close()

		for {
			select {
			case <-ctx.Done():
				return
			case param := <-p.existedStream:
				param.existed <- p.checkExist(param.key)
			case <-ticker.C:
				p.doGenerate()
			case param := <-p.keyStream:
				param.key <- p.get()
			}
		}
	}()
}

func (p *Pool) Key() string {
	param := &GetKeyParam{
		key: make(chan string),
	}
	p.keyStream <- param
	return <-param.key
}

func (p *Pool) Existed(key string) bool {
	param := IsExistParam{
		key:     key,
		existed: make(chan bool),
	}
	defer close(param.existed)

	p.existedStream <- param
	return <-param.existed
}

func (p *Pool) get() string {
	p.produced++
	if len(p.cache) == 0 {
		p.doGenerate()
	}
	r := p.cache[0]
	p.cache = p.cache[1:]
	return r
}

func (p *Pool) checkExist(key string) bool {
	return p.generator.Existed(key)
}

func (p *Pool) doGenerate() {
	fmt.Println("doGenerate...")
	for i, leftInCache := 0, len(p.cache); i < OneTimeGeneratedNumber-leftInCache; i++ {
		key, err := p.generator.Generate()
		if err != nil {
			log.Fatalln(err)
		}
		p.cache = append(p.cache, key)
	}
}
