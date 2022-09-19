package MyMulti

import (
	"errors"
	"io"
)

// 用于存放多个Writer
type multiWriter struct {
	writers []io.Writer
}

type multiReader struct {
	readers []io.Reader
}

func (t *multiWriter) Write(data []byte) (n int, err error) {
	//逐步writers里面每一项的Write方法
	for _, v := range t.writers {
		n, err = v.Write(data)
		if err != nil {
			return
		}
		if n != len(data) {
			err = errors.New("len is error")
		}
	}
	return len(data), nil
}

func MultiWriter(writers ...io.Writer) *multiWriter {
	//声明一个slice，逐步读取，存放writers
	allWriters := make([]io.Writer, 0, len(writers))
	for _, w := range writers {
		//这一步是为了运行writers中有multiWriter类型的，这样的话就允许嵌套
		if nw, ok := w.(*multiWriter); ok {
			allWriters = append(allWriters, nw.writers...)
		} else {
			allWriters = append(allWriters, w)
		}
	}
	return &multiWriter{allWriters}
}

func (mr *multiReader) Read(data []byte) (n int, err error) {
	//tip:这里的目的也和Write类似把slice里面数据取出来一个一个Read
	//不过有意思的是,这里使用迭代的方式一个一个的shift出来执行
	//这样的目的是对只有一个reader情况进行判断，看看他里面有没有嵌套另一个multiReader

	for len(mr.readers) > 0 {
		if len(mr.readers) == 1 {
			if r, ok := mr.readers[0].(*multiReader); ok {
				mr.readers = r.readers
				continue
			}
		}

		n, err = mr.readers[0].Read(data)
		//到这里就代表这次读取读好了
		if err == io.EOF {
			mr.readers[0] = nil
			mr.readers = mr.readers[1:]
		}

		//tip:good 一定要有这几步，这里一定要err = nil
		// err = nil后，外面会做层判断，检查到nil后会for循环再次Read一次
		//这样就可以继续执行下一个reader了
		if n > 0 || err != io.EOF {
			if err == io.EOF && len(mr.readers) > 0 {
				// Don't return EOF yet. More readers remain.
				err = nil
			}
			return
		}
	}
	return 0, io.EOF
}

func MultiReader(readers ...io.Reader) *multiReader {
	//tip:这里目的与MultiWriter相似，不过用了另一套代码
	rs := make([]io.Reader, len(readers))
	copy(rs, readers)
	return &multiReader{rs}
}
