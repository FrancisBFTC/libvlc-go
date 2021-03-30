package vlc

// #cgo LDFLAGS: -lvlc
// #include <vlc/vlc.h>
// #include <stdlib.h>
import "C"
import (
	"unsafe"
)

// StereoMode defines stereo modes which can be used by an audio output.
type StereoMode int

// Stereo modes.
const (
	StereoModeError StereoMode = iota - 1
	StereoModeNotSet
	StereoModeNormal
	StereoModeReverse
	StereoModeLeft
	StereoModeRight
	StereoModeDolbySurround
	StereoModeHeadphones
)

// AudioOutput is an abstraction for rendering decoded (or pass-through)
// audio samples.
type AudioOutput struct {
	Name        string
	Description string
}

// AudioOutputList returns the list of available audio outputs.
// In order to change the audio output of a media player instance,
// use the Player.SetAudioOutput method.
func AudioOutputList() ([]*AudioOutput, error) {
	if err := inst.assertInit(); err != nil {
		return nil, err
	}

	cOutputs := C.libvlc_audio_output_list_get(inst.handle)
	if cOutputs == nil {
		return nil, errOrDefault(getError(), ErrAudioOutputListMissing)
	}

	var outputs []*AudioOutput
	for n := cOutputs; n != nil; n = n.p_next {
		outputs = append(outputs, &AudioOutput{
			Name:        C.GoString(n.psz_name),
			Description: C.GoString(n.psz_description),
		})
	}

	C.libvlc_audio_output_list_release(cOutputs)
	return outputs, getError()
}

// AudioOutputDevice contains information regarding an audio output device.
type AudioOutputDevice struct {
	Name        string
	Description string
}

// ListAudioOutputDevices returns the list of available devices for the
// specified audio output. Use the AudioOutputList method in order to obtain
// the list of available audio outputs.
// NOTE: Not all audio outputs support this. An empty list of devices does
// not imply that the specified audio output does not work.
// Some audio output devices in the list might not work in some circumstances.
func ListAudioOutputDevices(output string) ([]*AudioOutputDevice, error) {
	if err := inst.assertInit(); err != nil {
		return nil, err
	}

	cOutput := C.CString(output)
	defer C.free(unsafe.Pointer(cOutput))

	cDevices := C.libvlc_audio_output_device_list_get(inst.handle, cOutput)
	return parseAudioOutputDeviceList(cDevices)
}

func parseAudioOutputDeviceList(cDevices *C.libvlc_audio_output_device_t) ([]*AudioOutputDevice, error) {
	if cDevices == nil {
		return nil, errOrDefault(getError(), ErrAudioOutputDeviceListMissing)
	}

	var devices []*AudioOutputDevice
	for n := cDevices; n != nil; n = n.p_next {
		devices = append(devices, &AudioOutputDevice{
			Name:        C.GoString(n.psz_device),
			Description: C.GoString(n.psz_description),
		})
	}

	C.libvlc_audio_output_device_list_release(cDevices)
	return devices, getError()
}

// ModuleDescription contains information about a libVLC module.
type ModuleDescription struct {
	Name      string
	ShortName string
	LongName  string
	Help      string
}

// ListAudioFilters returns the list of available audio filters.
func ListAudioFilters() ([]*ModuleDescription, error) {
	if err := inst.assertInit(); err != nil {
		return nil, err
	}

	return parseFilterList(C.libvlc_audio_filter_list_get(inst.handle))
}

// ListVideoFilters returns the list of available video filters.
func ListVideoFilters() ([]*ModuleDescription, error) {
	if err := inst.assertInit(); err != nil {
		return nil, err
	}

	return parseFilterList(C.libvlc_video_filter_list_get(inst.handle))
}

func parseFilterList(cFilters *C.libvlc_module_description_t) ([]*ModuleDescription, error) {
	if cFilters == nil {
		return nil, errOrDefault(getError(), ErrFilterListMissing)
	}

	var filters []*ModuleDescription
	for n := cFilters; n != nil; n = n.p_next {
		filters = append(filters, &ModuleDescription{
			Name:      C.GoString(n.psz_name),
			ShortName: C.GoString(n.psz_shortname),
			LongName:  C.GoString(n.psz_longname),
			Help:      C.GoString(n.psz_help),
		})
	}

	C.libvlc_module_description_list_release(cFilters)
	return filters, getError()
}
