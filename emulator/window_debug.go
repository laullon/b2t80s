package emulator

import (
	"C"
	"bufio"
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/disintegration/imaging"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/laullon/b2t80s/debug"
	"github.com/webview/webview"
)

type debugWindow struct {
	web    webview.WebView
	stream *Stream
	tabs   Tabs
}

func NewDebugWindow(name string, machine Machine) Window {
	win := &debugWindow{
		web:    initDebugWindow(name, machine),
		stream: NewStream(),
	}

	win.web.Bind("getStatus", func() string {
		return fmt.Sprintf("time: %s - FPS: %03.2f\n", machine.Clock().Stats(), machine.Monitor().FPS())
	})

	win.web.Bind("getCPU", func() string {
		ui := machine.Control()[win.tabs.Selected()]
		return ui.HTML()
	})

	win.web.Bind("initUI", func() {
		win.tabs.Show()
	})

	http.Handle("/video", win.stream)
	http.Handle("/", http.FileServer(debug.AssetFile()))
	listener, err := net.Listen("tcp", ":")
	if err != nil {
		panic(err)
	}

	go func() {
		panic(http.Serve(listener, nil))
	}()

	fmt.Println("Using port:", listener.Addr())
	url := fmt.Sprintf("http://0:%d/app/?%d", listener.Addr().(*net.TCPAddr).Port, time.Now().Local().Second())
	println("url:", url)
	win.web.Navigate(url)

	machine.Monitor().SetRedraw(func() {
		go func() {
			var b bytes.Buffer
			w := bufio.NewWriter(&b)

			img := imaging.FlipV(machine.Monitor().Screen())
			img = imaging.Resize(img, 600, 0, imaging.NearestNeighbor)

			jpeg.Encode(w, img, &jpeg.Options{Quality: 100})
			win.stream.UpdateJPEG(b.Bytes())

		}()
	})

	win.tabs = NewTabs("tabs", win.web, machine)

	return win
}

func (win *debugWindow) SetOnKey(onKey func(glfw.Key)) {
	// win.main.onKey = onKey
}

func (win *debugWindow) Run() {
	win.web.Run()
	win.web.Destroy()
}

func initDebugWindow(title string, machine Machine) webview.WebView {
	debug := true
	w := webview.New(debug)
	w.SetTitle(title)
	w.SetSize(1200, 600, webview.HintNone)
	// println(w.Window())
	return w
}

type Stream struct {
	m             map[chan []byte]bool
	frame         []byte
	lock          sync.Mutex
	FrameInterval time.Duration
}

func NewStream() *Stream {
	return &Stream{
		m:             make(map[chan []byte]bool),
		frame:         make([]byte, len(headerf)),
		FrameInterval: 60 * time.Millisecond,
	}
}

const boundaryWord = "MJPEGBOUNDARY"
const headerf = "\r\n" +
	"--" + boundaryWord + "\r\n" +
	"Content-Type: image/jpeg\r\n" +
	"Content-Length: %d\r\n" +
	"X-Timestamp: 0.000000\r\n" +
	"\r\n"

func (s *Stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Stream:", r.RemoteAddr, "connected")
	w.Header().Add("Content-Type", "multipart/x-mixed-replace;boundary="+boundaryWord)

	c := make(chan []byte)
	s.lock.Lock()
	s.m[c] = true
	s.lock.Unlock()

	for {
		time.Sleep(s.FrameInterval)
		b := <-c
		_, err := w.Write(b)
		if err != nil {
			break
		}
	}

	s.lock.Lock()
	delete(s.m, c)
	s.lock.Unlock()
	log.Println("Stream:", r.RemoteAddr, "disconnected")
}

func (s *Stream) UpdateJPEG(jpeg []byte) {
	header := fmt.Sprintf(headerf, len(jpeg))
	if len(s.frame) < len(jpeg)+len(header) {
		s.frame = make([]byte, (len(jpeg)+len(header))*2)
	}

	copy(s.frame, header)
	copy(s.frame[len(header):], jpeg)

	s.lock.Lock()
	for c := range s.m {
		select {
		case c <- s.frame:
		default:
		}
	}
	s.lock.Unlock()
}
