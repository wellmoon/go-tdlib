package tdlib

import "C"
import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"unsafe"

	"github.com/sirupsen/logrus"
	"github.com/wellmoon/go-tdlib/config"
)

// #cgo linux CFLAGS: -I/usr/local/include
// #cgo linux LDFLAGS: -Wl,-rpath=/usr/local/lib -ltdjson
// #include <stdlib.h>
// #include "callbacks.h"
import "C"

//export go_td_log_message_callback_ptr
func go_td_log_message_callback_ptr(verbosityLevel C.int, message *C.char) {
	goMessage := C.GoString(message)
	fmt.Printf("Received log message with verbosity level %d: %s\n", int(verbosityLevel), goMessage)
}

type Manager struct {
	handlers *ManagerHandlers

	clientUpdateChannels sync.Map

	options *ManagerOptions
}

// NewManager creates a new Manager instance
func NewManager(handlers *ManagerHandlers, options *ManagerOptions, ms int, cnum int) *Manager {
	if options != nil {
		C.set_log_message_callback(C.int(options.LogVerbosityLevel))

		if options.LogPath != "" {
			cfgBytes, _ := json.Marshal(map[string]interface{}{
				"@type": "setLogStream",
				"log_stream": map[string]interface{}{
					"@type":         "logStreamFile",
					"path":          options.LogPath,
					"max_file_size": 10485760,
				},
			})

			query := C.CString(string(cfgBytes))
			C.td_execute(query)
			C.free(unsafe.Pointer(query))
		}
	}

	manager := &Manager{
		handlers: handlers,
		options:  options,
	}
	c = make(chan int8, cnum)
	go manager.receiveUpdates(ms)

	return manager
}

// NewClient creates a new client instance
func (m *Manager) NewClient(
	apiID int64,
	apiHash string,
	handlers *Handlers,
	cfg *config.Config,
	logger *logrus.Logger,
) *TDLib {

	clientID := m.newClientID()

	updateChannel := make(chan []byte, 10)

	m.clientUpdateChannels.Store(clientID, updateChannel)

	return newClientV2(
		clientID,
		updateChannel,
		apiID,
		apiHash,
		handlers,
		cfg,
		logger,
	)
}

func (m *Manager) newClientID() int {
	result := C.td_create_client_id()

	return int(result)
}

var c chan int8
var cc chan int8 = make(chan int8, 5)

func (m *Manager) receiveNextUpdate(timeout float64) []byte {
	c <- 1
	defer func() {
		<-c
	}()
	result := C.td_receive(C.double(timeout))
	if result == nil {
		return nil
	}
	res := C.GoString(result)
	if len(c) == 5 {
		fmt.Println(res)
	}

	return []byte(res)
}

func (m *Manager) receiveUpdates(ms int) {

	// text := ""
	fmt.Println("================receiveUpdates===============")
	for {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		updateBytes := m.receiveNextUpdate(10)

		if len(updateBytes) == 0 {
			continue
		}
		// text = string(updateBytes)
		// if !strings.Contains(text, "请点击正确答案") && !strings.Contains(text, "您必须完成人机验证才能继续使用") {
		// 	continue
		// }
		if m.handlers != nil && m.handlers.onRawIncomingEvent != nil {
			go m.handlers.onRawIncomingEvent(updateBytes)
		}

		clientID := m.getClientID(updateBytes)
		if clientID != nil {
			if clientChan, ok := m.getClientChannel(*clientID); ok {
				go m.writeClientEvent(clientChan, updateBytes)
				continue
			} else {
				// TODO: Log Received Update For Unknown Client
				fmt.Printf("Received update for unknown client: %d\n", *clientID)
			}
		}

		// TODO: Log General Received Update, Doesn't Belong To Any Client
		fmt.Printf("Received update: %s\n", string(updateBytes))
	}
}

func (m *Manager) getClientChannel(clientID int) (chan []byte, bool) {
	clientValue, ok := m.clientUpdateChannels.Load(clientID)
	if ok {
		clientUpdateChannel, ok := clientValue.(chan []byte)
		if ok {
			return clientUpdateChannel, true
		}
	}

	return nil, false
}

func (m *Manager) writeClientEvent(updateChan chan<- []byte, update []byte) {
	updateChan <- update
}

func (m *Manager) getClientID(update []byte) *int {
	type ClientID struct {
		ID int `json:"@client_id"`
	}

	var clientID ClientID

	err := json.Unmarshal(update, &clientID)
	if err != nil {
		return nil
	}

	return &clientID.ID
}
