package MyPipe

import (
	"errors"
	"sync"
)

type pipe struct {
	// 使用mutex在write时保护数据
	wrMu sync.Mutex
	// 写入数据通道 使用[]byte是为了传递数据给read函数
	wrCh chan []byte
	// pipe.read方法调用会返回int类型 保证write和read长度一样
	rdCh chan int
	//接收通道结束 使用空结构体是因为这个不占内存
	done chan struct{}
}

func (p *pipe) read(buf []byte) (n int, err error) {
	//先做一次预处理，防止通道关闭
	select {
	case <-p.done:
		return 0, errors.New("chan is closed")
	default:
	}

	select {
	case bw := <-p.wrCh:
		//数据拷贝到buf中，buf就是我们存放读取出来数据的地方
		nr := copy(buf, bw)
		//传过去，在write时候，知道接下来从哪里开始，把数据继续发送给read函数
		p.rdCh <- nr
		return nr, nil
	case <-p.done:
		return 0, errors.New("chan is closed")
	}
}

func (p *pipe) write(data []byte) (n int, err error) {
	//如果通道没关闭，就Lock，然后再执行代码
	select {
	case <-p.done:
		return 0, errors.New("chan is closed")
	default:
		p.wrMu.Lock()
		defer p.wrMu.Unlock()
	}
	//tip:amaze 判断+只执行一次，独特的用法
	for once := true; len(data) > 0 || once; once = false {
		select {
		//把已经写入的数据放到wrCh中
		case p.wrCh <- data:
			// 这步之前,p.read已经处理过了，并且把已经读取字符的数量返回存入了rdCh
			// 所有下面这几步就是从继续往下存入wrCh
			// tip:good 这里需要这样处理的原因是因为，虽然两者都是[]byte，但是可能会大小不一样，一次取不完
			nw := <-p.rdCh
			data = data[nw:]
			n += nw
		case <-p.done:
			return 0, errors.New("chan is closed")
		}
	}
	return n, nil
}

// *** 下面两个其实就是pipe 通过下面两种处理，使他们分别实现了不同接口 *** //

// PipeReader 是Pipe返回的可读取的值
type PipeReader struct {
	p *pipe
}

// PipeWriter 是Pipe返回的可写入的值
type PipeWriter struct {
	p *pipe
}

// 实现Read接口
func (r *PipeReader) Read(data []byte) (n int, err error) {
	return r.p.read(data)
}

// 实现Write接口
func (w *PipeWriter) Write(data []byte) (n int, err error) {
	return w.p.write(data)
}

// 分别为这两个高级Pipe实现Close,方便退出

func (r *PipeReader) Close() error {
	// tip:worse 这里只是拙劣的实现，源代码对于error处理更优雅
	close(r.p.done)
	return nil
}

func (w *PipeWriter) Close() error {
	close(w.p.done)
	return nil
}

// 对外暴露Pipe方法

func Pipe() (*PipeReader, *PipeWriter) {
	p := &pipe{
		wrCh: make(chan []byte),
		rdCh: make(chan int),
		done: make(chan struct{}),
	}
	return &PipeReader{p}, &PipeWriter{p}
}
