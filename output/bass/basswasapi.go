// +build windows

package bass

// #cgo CPPFLAGS: -Iinclude
// #cgo CXXFLAGS: -Iinclude
// #include "basswasapi.h"
import "C"
import (
	"unsafe"
)

// Additional error codes returned by BASS_ErrorGetCode
const (
	BASS_ERROR_WASAPI          = C.BASS_ERROR_WASAPI
	BASS_ERROR_WASAPI_BUFFER   = C.BASS_ERROR_WASAPI_BUFFER
	BASS_ERROR_WASAPI_CATEGORY = C.BASS_ERROR_WASAPI_CATEGORY
	BASS_ERROR_WASAPI_DENIED   = C.BASS_ERROR_WASAPI_DENIED
)

// BASS_WASAPI_DEVICEINFO "type"
const (
	BASS_WASAPI_TYPE_NETWORKDEVICE = C.BASS_WASAPI_TYPE_NETWORKDEVICE
	BASS_WASAPI_TYPE_SPEAKERS      = C.BASS_WASAPI_TYPE_SPEAKERS
	BASS_WASAPI_TYPE_LINELEVEL     = C.BASS_WASAPI_TYPE_LINELEVEL
	BASS_WASAPI_TYPE_HEADPHONES    = C.BASS_WASAPI_TYPE_HEADPHONES
	BASS_WASAPI_TYPE_MICROPHONE    = C.BASS_WASAPI_TYPE_MICROPHONE
	BASS_WASAPI_TYPE_HEADSET       = C.BASS_WASAPI_TYPE_HEADSET
	BASS_WASAPI_TYPE_HANDSET       = C.BASS_WASAPI_TYPE_HANDSET
	BASS_WASAPI_TYPE_DIGITAL       = C.BASS_WASAPI_TYPE_DIGITAL
	BASS_WASAPI_TYPE_SPDIF         = C.BASS_WASAPI_TYPE_SPDIF
	BASS_WASAPI_TYPE_HDMI          = C.BASS_WASAPI_TYPE_HDMI
	BASS_WASAPI_TYPE_UNKNOWN       = C.BASS_WASAPI_TYPE_UNKNOWN
)

// BASS_WASAPI_DEVICEINFO flags
const (
	BASS_DEVICE_ENABLED   = C.BASS_DEVICE_ENABLED
	BASS_DEVICE_DEFAULT   = C.BASS_DEVICE_DEFAULT
	BASS_DEVICE_INIT      = C.BASS_DEVICE_INIT
	BASS_DEVICE_LOOPBACK  = C.BASS_DEVICE_LOOPBACK
	BASS_DEVICE_INPUT     = C.BASS_DEVICE_INPUT
	BASS_DEVICE_UNPLUGGED = C.BASS_DEVICE_UNPLUGGED
	BASS_DEVICE_DISABLED  = C.BASS_DEVICE_DISABLED
)

// BASS_WASAPI_Init flags
const (
	BASS_WASAPI_EXCLUSIVE  = C.BASS_WASAPI_EXCLUSIVE
	BASS_WASAPI_AUTOFORMAT = C.BASS_WASAPI_AUTOFORMAT
	BASS_WASAPI_BUFFER     = C.BASS_WASAPI_BUFFER
	BASS_WASAPI_EVENT      = C.BASS_WASAPI_EVENT
	BASS_WASAPI_SAMPLES    = C.BASS_WASAPI_SAMPLES
	BASS_WASAPI_DITHER     = C.BASS_WASAPI_DITHER
	BASS_WASAPI_RAW        = C.BASS_WASAPI_RAW
	BASS_WASAPI_ASYNC      = C.BASS_WASAPI_ASYNC
)

const (
	BASS_WASAPI_CATEGORY_MASK                   = C.BASS_WASAPI_CATEGORY_MASK
	BASS_WASAPI_CATEGORY_OTHER                  = C.BASS_WASAPI_CATEGORY_OTHER
	BASS_WASAPI_CATEGORY_FOREGROUNDONLYMEDIA    = C.BASS_WASAPI_CATEGORY_FOREGROUNDONLYMEDIA
	BASS_WASAPI_CATEGORY_BACKGROUNDCAPABLEMEDIA = C.BASS_WASAPI_CATEGORY_BACKGROUNDCAPABLEMEDIA
	BASS_WASAPI_CATEGORY_COMMUNICATIONS         = C.BASS_WASAPI_CATEGORY_COMMUNICATIONS
	BASS_WASAPI_CATEGORY_ALERTS                 = C.BASS_WASAPI_CATEGORY_ALERTS
	BASS_WASAPI_CATEGORY_SOUNDEFFECTS           = C.BASS_WASAPI_CATEGORY_SOUNDEFFECTS
	BASS_WASAPI_CATEGORY_GAMEEFFECTS            = C.BASS_WASAPI_CATEGORY_GAMEEFFECTS
	BASS_WASAPI_CATEGORY_GAMEMEDIA              = C.BASS_WASAPI_CATEGORY_GAMEMEDIA
	BASS_WASAPI_CATEGORY_GAMECHAT               = C.BASS_WASAPI_CATEGORY_GAMECHAT
	BASS_WASAPI_CATEGORY_SPEECH                 = C.BASS_WASAPI_CATEGORY_SPEECH
	BASS_WASAPI_CATEGORY_MOVIE                  = C.BASS_WASAPI_CATEGORY_MOVIE
	BASS_WASAPI_CATEGORY_MEDIA                  = C.BASS_WASAPI_CATEGORY_MEDIA
)

// BASS_WASAPI_INFO "format"
const (
	BASS_WASAPI_FORMAT_FLOAT = C.BASS_WASAPI_FORMAT_FLOAT
	BASS_WASAPI_FORMAT_8BIT  = C.BASS_WASAPI_FORMAT_8BIT
	BASS_WASAPI_FORMAT_16BIT = C.BASS_WASAPI_FORMAT_16BIT
	BASS_WASAPI_FORMAT_24BIT = C.BASS_WASAPI_FORMAT_24BIT
	BASS_WASAPI_FORMAT_32BIT = C.BASS_WASAPI_FORMAT_32BIT
)

// BASS_WASAPI_Set/GetVolume modes
const (
	BASS_WASAPI_CURVE_DB      = C.BASS_WASAPI_CURVE_DB
	BASS_WASAPI_CURVE_LINEAR  = C.BASS_WASAPI_CURVE_LINEAR
	BASS_WASAPI_CURVE_WINDOWS = C.BASS_WASAPI_CURVE_WINDOWS
	BASS_WASAPI_VOL_SESSION   = C.BASS_WASAPI_VOL_SESSION
)

var (
	WASAPIPROC_PUSH = 0
	WASAPIPROC_BASS = -1
)

// Device notifications
const (
	BASS_WASAPI_NOTIFY_ENABLED   = C.BASS_WASAPI_NOTIFY_ENABLED
	BASS_WASAPI_NOTIFY_DISABLED  = C.BASS_WASAPI_NOTIFY_DISABLED
	BASS_WASAPI_NOTIFY_DEFOUTPUT = C.BASS_WASAPI_NOTIFY_DEFOUTPUT
	BASS_WASAPI_NOTIFY_DEFINPUT  = C.BASS_WASAPI_NOTIFY_DEFINPUT
	BASS_WASAPI_NOTIFY_FAIL      = C.BASS_WASAPI_NOTIFY_FAIL
)

func BASS_WASAPI_GetVersion() uint {
	return uint(C.BASS_WASAPI_GetVersion())
}

func BASS_WASAPI_SetNotify(proc *C.WASAPINOTIFYPROC, user unsafe.Pointer) bool {
	return C.BASS_WASAPI_SetNotify(proc, user) != 0
}

type BassWasapiDeviceInfo struct {
	Name      string
	ID        string
	Type      uint
	Flags     uint
	MinPeriod float32
	DefPeriod float32
	MixFreq   uint
	MixChans  uint
}

func BASS_WASAPI_GetDeviceInfo(device uint, info *BassWasapiDeviceInfo) bool {
	i := (*C.BASS_WASAPI_DEVICEINFO)(C.malloc(C.sizeof_BASS_WASAPI_DEVICEINFO))
	defer C.free(unsafe.Pointer(i))
	res := C.BASS_WASAPI_GetDeviceInfo(C.DWORD(device), i) != 0
	if res {
		info.ID = C.GoString(i.id)
		info.Name = C.GoString(i.name)
		info.Type = uint(i.deviceType)
		info.Flags = uint(i.flags)
		info.MinPeriod = float32(i.minperiod)
		info.DefPeriod = float32(i.defperiod)
		info.MixFreq = uint(i.mixfreq)
		info.MixChans = uint(i.mixchans)
	}
	return res
}

func BASS_WASAPI_GetDeviceLevel(device uint, channel int) float32 {
	return float32(C.BASS_WASAPI_GetDeviceLevel(C.DWORD(device), C.int(channel)))
}

func BASS_WASAPI_SetDevice(device uint) bool {
	return C.BASS_WASAPI_SetDevice(C.DWORD(device)) != 0
}

func BASS_WASAPI_GetDevice() uint {
	return uint(C.BASS_WASAPI_GetDevice())
}

func BASS_WASAPI_CheckFormat(device uint, freq uint, chans uint, flags uint) uint {
	return uint(C.BASS_WASAPI_CheckFormat(C.DWORD(device), C.DWORD(freq), C.DWORD(chans), C.DWORD(flags)))
}

func BASS_WASAPI_Init(device int, freq uint, chans uint, flags uint, buffer float32, period float32, proc *C.WASAPIPROC, user unsafe.Pointer) bool {
	return C.BASS_WASAPI_Init(C.int(device), C.DWORD(freq), C.DWORD(chans), C.DWORD(flags), C.float(buffer), C.float(period), proc, user) != 0
}

func BASS_WASAPI_Init_Special(device int, freq uint, chans uint, flags uint, buffer float32, period float32, proc int, user unsafe.Pointer) bool {
	return C.BASS_WASAPI_Init_Special(C.int(device), C.DWORD(freq), C.DWORD(chans), C.DWORD(flags), C.float(buffer), C.float(period), (C.int)(proc), user) != 0
}

func BASS_WASAPI_Free() bool {
	return C.BASS_WASAPI_Free() != 0
}

type BassWasapiInfo struct {
	InitFlags uint
	Freq      uint
	Chans     uint
	Format    uint
	BufLen    uint
	VolMax    float32
	VolMin    float32
	VolStep   float32
}

func BASS_WASAPI_GetInfo(info *BassWasapiInfo) bool {
	i := (*C.BASS_WASAPI_INFO)(C.malloc(C.sizeof_BASS_WASAPI_INFO))
	defer C.free(unsafe.Pointer(i))
	res := C.BASS_WASAPI_GetInfo(i) != 0
	if res {
		info.InitFlags = uint(i.initflags)
		info.Freq = uint(i.freq)
		info.Chans = uint(i.chans)
		info.Format = uint(i.format)
		info.BufLen = uint(i.buflen)
		info.VolMax = float32(i.volmax)
		info.VolMin = float32(i.volmin)
		info.VolStep = float32(i.volstep)
	}
	return res
}

func BASS_WASAPI_GetCPU() float32 {
	return float32(C.BASS_WASAPI_GetCPU())
}

func BASS_WASAPI_Lock(lock bool) bool {
	return C.BASS_WASAPI_Lock(bool2Cint(lock)) != 0
}

func BASS_WASAPI_Start() bool {
	return C.BASS_WASAPI_Start() != 0
}

func BASS_WASAPI_Stop(reset bool) bool {
	return C.BASS_WASAPI_Stop(bool2Cint(reset)) != 0
}

func BASS_WASAPI_IsStarted() bool {
	return C.BASS_WASAPI_IsStarted() != 0
}

func BASS_WASAPI_SetVolume(mode uint, volume float32) bool {
	return C.BASS_WASAPI_SetVolume(C.DWORD(mode), C.float(volume)) != 0
}

func BASS_WASAPI_GetVolume(mode uint) float32 {
	return float32(C.BASS_WASAPI_GetVolume(C.DWORD(mode)))
}

func BASS_WASAPI_SetMute(mode uint, mute bool) bool {
	return C.BASS_WASAPI_SetMute(C.DWORD(mode), bool2Cint(mute)) != 0
}

func BASS_WASAPI_GetMute(mode uint) bool {
	return C.BASS_WASAPI_GetMute(C.DWORD(mode)) != 0
}
