package emulator

import (
	"bufio"
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"sync"
	"time"

	"fyne.io/fyne/v2/app"
	"github.com/disintegration/imaging"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/laullon/b2t80s/debug"
	"github.com/laullon/b2t80s/ui"
	"github.com/webview/webview"
)

type debugWindow struct {
	web    webview.WebView
	stream *Stream
}

func NewDebugWindow(name string, machine Machine) Window {
	// TODO: REMOVE, just to prevent crashed
	ui.App = app.NewWithID("io.fyne.test")

	win := &debugWindow{
		web:    initDebugWindow(name, machine),
		stream: NewStream(),
	}

	win.web.Bind("getStatus", func() string {
		return fmt.Sprintf("time: %s - FPS: %03.2f\n", machine.Clock().Stats(), machine.Monitor().FPS())
	})

	cpuUI := machine.Control()["CPU"]
	win.web.Bind("getCPU", func() string {
		return cpuUI.HTML()
	})

	http.Handle("/video", win.stream)
	http.Handle("/", http.FileServer(debug.AssetFile()))
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

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
	return win
}

func (win *debugWindow) SetOnKey(onKey func(glfw.Key)) {
	// win.main.onKey = onKey
}

func (win *debugWindow) Run() {

	println(1)
	win.web.Run()
	println(2)
	win.web.Destroy()
	println(3)
}

func initDebugWindow(title string, machine Machine) webview.WebView {
	debug := true
	w := webview.New(debug)
	w.SetTitle(title)
	w.SetSize(1200, 600, webview.HintNone)
	w.Navigate(fmt.Sprint("http://localhost:8080/debug/static/?", time.Now().Local().Second()))
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
