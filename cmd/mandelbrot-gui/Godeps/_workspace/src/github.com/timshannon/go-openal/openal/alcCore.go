// Forked by Tim Shannon 2012
// Copyright 2009 Peter H. Froehlich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// C-level binding for OpenAL's "alc" API.
//
// Please consider using the Go-level binding instead.
//
// Note that "alc" introduces the exact same types as "al"
// but with different names. Check the documentation of
// openal/al for more information about the mapping to
// Go types.
//
// XXX: Sadly we have to returns pointers for both Device
// and Context to avoid problems with implicit assignments
// in clients. It's sad because it makes the overhead a
// lot higher, each of those calls triggers an allocation.
package openal

//#cgo linux LDFLAGS: -lopenal
//#cgo darwin LDFLAGS: -framework OpenAL
//#include <stdlib.h>
//#include "local.h"
/*
ALCdevice *walcOpenDevice(const char *devicename) {
	return alcOpenDevice(devicename);
}
const ALCchar *alcGetString( ALCdevice *device, ALCenum param );
void walcGetIntegerv(ALCdevice *device, ALCenum param, ALCsizei size, void *data) {
	alcGetIntegerv(device, param, size, data);
}
ALCdevice *walcCaptureOpenDevice(const char *devicename, ALCuint frequency, ALCenum format, ALCsizei buffersize) {
	return alcCaptureOpenDevice(devicename, frequency, format, buffersize);
}
ALCint walcGetInteger(ALCdevice *device, ALCenum param) {
	ALCint result;
	alcGetIntegerv(device, param, 1, &result);
	return result;
}
*/
import "C"
import "unsafe"

// Error codes returned by Device.GetError().
const (
	//NoError        = 0
	InvalidDevice  = 0xA001
	InvalidContext = 0xA002
	//InvalidEnum    = 0xA003
	//InvalidValue   = 0xA004
	OutOfMemory = 0xA005
)

const (
	Frequency     = 0x1007 // int Hz
	Refresh       = 0x1008 // int Hz
	Sync          = 0x1009 // bool
	MonoSources   = 0x1010 // int
	StereoSources = 0x1011 // int
)

// The Specifier string for default device?
const (
	DefaultDeviceSpecifier = 0x1004
	DeviceSpecifier        = 0x1005
	Extensions             = 0x1006
)

// ?
const (
	MajorVersion = 0x1000
	MinorVersion = 0x1001
)

// ?
const (
	AttributesSize = 0x1002
	AllAttributes  = 0x1003
)

// Capture extension
const (
	CaptureDeviceSpecifier        = 0x310
	CaptureDefaultDeviceSpecifier = 0x311
	CaptureSamples                = 0x312
)

type Device struct {
	handle *C.ALCdevice
}

// GetError() returns the most recent error generated
// in the AL state machine.
func (self *Device) GetError() uint32 {
	return uint32(C.alcGetError(self.handle))
}

func OpenDevice(name string) *Device {
	// TODO: turn empty string into nil?
	// TODO: what about an error return?
	p := C.CString(name)
	h := C.walcOpenDevice(p)
	C.free(unsafe.Pointer(p))
	return &Device{h}
}

func (self *Device) CloseDevice() bool {
	//TODO: really a method? or not?
	return C.alcCloseDevice(self.handle) != 0
}

func (self *Device) CreateContext() *Context {
	// TODO: really a method?
	// TODO: attrlist support
	return &Context{C.alcCreateContext(self.handle, nil)}
}

func (self *Device) GetIntegerv(param uint32, size uint32) (result []int32) {
	result = make([]int32, size)
	C.walcGetIntegerv(self.handle, C.ALCenum(param), C.ALCsizei(size), unsafe.Pointer(&result[0]))
	return
}

func (self *Device) GetInteger(param uint32) int32 {
	return int32(C.walcGetInteger(self.handle, C.ALCenum(param)))
}

type CaptureDevice struct {
	Device
	sampleSize uint32
}

func CaptureOpenDevice(name string, freq uint32, format uint32, size uint32) *CaptureDevice {
	// TODO: turn empty string into nil?
	// TODO: what about an error return?
	p := C.CString(name)
	h := C.walcCaptureOpenDevice(p, C.ALCuint(freq), C.ALCenum(format), C.ALCsizei(size))
	C.free(unsafe.Pointer(p))
	s := map[uint32]uint32{FormatMono8: 1, FormatMono16: 2, FormatStereo8: 2, FormatStereo16: 4}[format]
	return &CaptureDevice{Device{h}, s}
}

// XXX: Override Device.CloseDevice to make sure the correct
// C function is called even if someone decides to use this
// behind an interface.
func (self *CaptureDevice) CloseDevice() bool {
	return C.alcCaptureCloseDevice(self.handle) != 0
}

func (self *CaptureDevice) CaptureCloseDevice() bool {
	return self.CloseDevice()
}

func (self *CaptureDevice) CaptureStart() {
	C.alcCaptureStart(self.handle)
}

func (self *CaptureDevice) CaptureStop() {
	C.alcCaptureStop(self.handle)
}

func (self *CaptureDevice) CaptureSamples(size uint32) (data []byte) {
	data = make([]byte, size*self.sampleSize)
	C.alcCaptureSamples(self.handle, unsafe.Pointer(&data[0]), C.ALCsizei(size))
	return
}

///// Context ///////////////////////////////////////////////////////

// Context encapsulates the state of a given instance
// of the OpenAL state machine. Only one context can
// be active in a given process.
type Context struct {
	handle *C.ALCcontext
}

// A context that doesn't exist, useful for certain
// context operations (see OpenAL documentation for
// details).
var NullContext Context

// Renamed, was MakeContextCurrent.
func (self *Context) Activate() bool {
	return C.alcMakeContextCurrent(self.handle) != alFalse
}

// Renamed, was ProcessContext.
func (self *Context) Process() {
	C.alcProcessContext(self.handle)
}

// Renamed, was SuspendContext.
func (self *Context) Suspend() {
	C.alcSuspendContext(self.handle)
}

// Renamed, was DestroyContext.
func (self *Context) Destroy() {
	C.alcDestroyContext(self.handle)
	self.handle = nil
}

// Renamed, was GetContextsDevice.
func (self *Context) GetDevice() *Device {
	return &Device{C.alcGetContextsDevice(self.handle)}
}

// Renamed, was GetCurrentContext.
func CurrentContext() *Context {
	return &Context{C.alcGetCurrentContext()}
}
