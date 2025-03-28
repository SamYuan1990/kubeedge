package writer

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"k8s.io/klog/v2"

	"github.com/kubeedge/beehive/pkg/core/socket/wrapper/packer"
)

// PackageWriter package writer
type PackageWriter struct {
	packer *packer.Packer
	conn   net.Conn
	lock   sync.Mutex
}

// NewPackageWriter new package writer
// Comments below is assisted by Gen AI
// // NewPackageWriter creates and returns a new PackageWriter instance initialized with the provided object.
// The function expects the input object to implement the net.Conn interface. If the object is a valid net.Conn,
// it initializes a new PackageWriter with the connection and a new packer.Packer instance.
// If the object does not implement net.Conn, the function logs an error using klog.Errorf and returns nil.
//
// Parameters:
//   - obj: The object to be used as the connection. It must implement the net.Conn interface.
//
// Returns:
//   - *PackageWriter: A pointer to a new PackageWriter instance if the input object is a valid net.Conn.
//   - nil: If the input object is not a valid net.Conn.
func NewPackageWriter(obj interface{}) *PackageWriter {
	if conn, ok := obj.(net.Conn); ok {
		packer := packer.NewPacker()
		return &PackageWriter{
			conn:   conn,
			packer: packer,
		}
	}
	klog.Errorf("bad conn obj")
	return nil
}

// Write write
func (w *PackageWriter) Write(message []byte) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	w.packer.Message = message
	w.packer.Length = int32(len(message))
	err := w.packer.Write(w.conn)
	if err != nil {
		klog.Errorf("failed to packer with error %+v", err)
		return fmt.Errorf("failed to packer, error:%+v", err)
	}
	return nil
}

// WriteJSON write json
func (w *PackageWriter) WriteJSON(obj interface{}) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	objBytes, err := json.Marshal(obj)
	if err != nil {
		klog.Errorf("failed to marshal obj, error:%+v", err)
		return err
	}
	w.packer.Message = objBytes
	w.packer.Length = int32(len(objBytes))
	err = w.packer.Write(w.conn)
	if err != nil {
		klog.Errorf("failed to packer, error:%+v", err)
		return fmt.Errorf("failed to packer, error:%+v", err)
	}
	return nil
}
