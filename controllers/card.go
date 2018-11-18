package controllers

import (
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/beego/ms304w-client/basis/conf"
	"github.com/beego/ms304w-client/basis/errors"
	"github.com/tarm/goserial"
)

const (
	// 命令头
	HDR byte = 0x1b
)

var cardCom = &serial.Config{
	Name:        conf.String("com_card_name"),
	Baud:        conf.Int("com_card_baud"),
	ReadTimeout: time.Duration(conf.Int("com_card_read_timeout")) * time.Millisecond,
}

var (
	cardRw io.ReadWriteCloser
	// 全局队列
	cardQueue = make(chan byte, 255)
)

func init() {
	/*
		if cardRw != nil {
			return
		}

		s, err := serial.OpenPort(cardCom)
		if err != nil {
			panic(err)
		}

		cardRw = s
	*/

	/*
		// 写数据到串口
		w := NewWriteCardCom()
		go WriteCardToCom(w.ComData())

		// 从串口读数据到队列
		go ReadCardToQueue()

		// 处理队列数据
		go ReadCardQueue()
	*/
}

// -------------------------------------
// 定义写操作结构体
// var dataArray = []string{"1b", "01", "00", "00", "00", "00", "00"}
// var dataArray = []byte{27, 1, 0, 0, 0, 0, 0}
type WriterCardCom struct {
	Hdr   byte // 1命令头 0x1b
	Cmd   byte // 2命令码
	Seq   byte // 3命令序号，指定天线板序号 0..
	State byte // 4命令状态
	Opt1  byte // 5命令操作码，保留
	Opt2  byte // 6命令操作码，保留
	Opt3  byte // 7命令操作码，保留
}

func NewWriteCardCom() *WriterCardCom {
	// writer
	return &WriterCardCom{
		Hdr:   HDR,
		Cmd:   0x01,
		Seq:   0x00,
		State: 0x00,
		Opt1:  0x00,
		Opt2:  0x00,
		Opt3:  0x00,
	}
}

func (w *WriterCardCom) ComData() []byte {
	var bytes = make([]byte, 0)

	t := reflect.TypeOf(*w)
	v := reflect.ValueOf(*w)

	for i := 0; i < t.NumField(); i++ {
		bytes = append(bytes, v.Field(i).Interface().(byte))
	}

	fmt.Printf("write: %X\n", bytes)

	return bytes
}

// ----------------------------------
// 定义读操作结构体
type ReaderCardCom struct {
	Hdr   byte   // 命令头 0x1b
	Cmd   byte   // 命令码
	Seq   byte   // 命令序号，指定天线板序号 0..
	State byte   // 命令状态
	Opt   byte   // 命令操作码，保留
	Len   byte   // 命令数据长度，不超过255
	Data  []byte // 命令数据内容
	Sign  byte   // 校验码
	// 当前位置
	Flag int // 一次读一个字节，依次处理
	// 异常
	Err error
}

func NewReaderCardCom() *ReaderCardCom {
	return &ReaderCardCom{
		Data: make([]byte, 0),
	}
}

func (r *ReaderCardCom) Card() {
	if len(r.Data) > 0 {
		card := hex.EncodeToString(r.Data)
		card = strings.ToUpper(card)
		log.Info("card: %s", card)

		Server.BroadcastTo("login", "loginByCard", card)
	}
}

// 往串口写数据
func WriteCardToCom(bytes []byte) {
	for {
		fmt.Println("writing.....")
		if _, err := cardRw.Write(bytes); err != nil {
			log.Error("WriteCardToCom: %v", errors.As(err))
		}

		time.Sleep(1 * 1e9)
	}
}

// 从串口读数据到队列
func ReadCardToQueue() {
	for {
		fmt.Println("reading to queue.....")
		buf := make([]byte, 1)
		n, err := cardRw.Read(buf)
		if err != nil {
			panic(err)
		}
		fmt.Printf("read n: %d, buf: %X\n", n, buf)

		if n == 1 {
			cardQueue <- buf[0]
		}
	}
}

func ReadCardQueue() {
	rc := NewReaderCardCom()

	// 处理数据
	for {
		fmt.Println("reading queue..........")
		select {
		case q := <-cardQueue:
			fmt.Printf("queue: %X\n", q)

			fmt.Println("Flag.........", rc.Flag)
			switch rc.Flag {
			case 0: // 命令头
				fmt.Println("00000000000", q == HDR)
				if q == HDR {
					rc.Hdr = q
					rc.Flag = rc.Flag + 1
				}

				break
			case 1: // 命令码
				fmt.Println("11111111111")
				rc.Cmd = q
				rc.Flag = rc.Flag + 1

				break

			case 2: // 命令序号
				fmt.Println("22222222222222")
				rc.Seq = q
				rc.Flag = rc.Flag + 1

				break
			case 3: // 命令状态
				fmt.Println("3333333333333")
				rc.State = q
				rc.Flag = rc.Flag + 1

				break
			case 4: // 命令操作码
				fmt.Println("4444444444444")
				rc.Opt = q
				rc.Flag = rc.Flag + 1

				break
			case 5: // 命令数据长度
				fmt.Println("5555555555")
				rc.Len = q
				if int(rc.Len) == 0 {
					// 没有数据内容
					rc.Flag = rc.Flag + 2
				} else {
					rc.Flag = rc.Flag + 1
				}

				break
			case 6: // 数据内容
				fmt.Println("6666666666666", int(rc.Len))
				fmt.Println(int(rc.Len), len(rc.Data))

				rc.Data = append(rc.Data, q)

				rc.Sign = rc.Sign ^ q

				if int(rc.Len) == len(rc.Data) {
					rc.Flag = rc.Flag + 1
				}

				fmt.Printf("Data........%X\n", rc.Data)

				break
			case 7:
				fmt.Println("77777777777777777")
				if rc.Sign == q {
					rc.Sign = q
					rc.Flag = rc.Flag + 1

					fmt.Println("success----------------")

					// 处理业务
					rc.Card()
				}

				// 初始化
				rc = &ReaderCardCom{
					Data: make([]byte, 0),
				}

				break
			default:
				log.Warn("%v", errors.New("undefined").As(rc.Flag))
			}
		}
	}
}
