package bls12377

// #cgo CFLAGS: -I./include/
// #include "msm.h"
import "C"

import (
	"github.com/ingonyama-zk/icicle/wrappers/golang/core"
	cr "github.com/ingonyama-zk/icicle/wrappers/golang/cuda_runtime"
	"unsafe"
)

func GetDefaultMSMConfig() core.MSMConfig {
	return core.GetDefaultMSMConfig()
}

func Msm(scalars core.HostOrDeviceSlice, points core.HostOrDeviceSlice, cfg *core.MSMConfig, results core.HostOrDeviceSlice) cr.CudaError {
	core.MsmCheck(scalars, points, cfg, results)
	var scalarsPointer unsafe.Pointer
	if scalars.IsOnDevice() {
		scalarsDevice := scalars.(core.DeviceSlice)
		scalarsDevice.CheckDevice()
		scalarsPointer = scalarsDevice.AsPointer()
	} else {
		scalarsPointer = unsafe.Pointer(&scalars.(core.HostSlice[ScalarField])[0])
	}
	cScalars := (*C.scalar_t)(scalarsPointer)

	var pointsPointer unsafe.Pointer
	if points.IsOnDevice() {
		pointsDevice := points.(core.DeviceSlice)
		pointsDevice.CheckDevice()
		pointsPointer = pointsDevice.AsPointer()
	} else {
		pointsPointer = unsafe.Pointer(&points.(core.HostSlice[Affine])[0])
	}
	cPoints := (*C.affine_t)(pointsPointer)

	var resultsPointer unsafe.Pointer
	if results.IsOnDevice() {
		resultsDevice := results.(core.DeviceSlice)
		resultsDevice.CheckDevice()
		resultsPointer = resultsDevice.AsPointer()
	} else {
		resultsPointer = unsafe.Pointer(&results.(core.HostSlice[Projective])[0])
	}
	cResults := (*C.projective_t)(resultsPointer)

	cSize := (C.int)(scalars.Len() / results.Len())
	cCfg := (*C.MSMConfig)(unsafe.Pointer(cfg))

	__ret := C.bls12_377MSMCuda(cScalars, cPoints, cSize, cCfg, cResults)
	err := (cr.CudaError)(__ret)
	return err
}
