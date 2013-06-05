// Copyright 2012 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package engine

import (
	"bitbucket.org/tshannon/gohorde/horde3d"
	"errors"
	"image"
	"io/ioutil"
	"os"
	"path"
)

var (
	dataDir  string
	dataFile string
)

const (
	virtualPath = "virtual://"
)

func init() {
	wd, _ := os.Getwd()
	dataDir = path.Join(wd, "data")
	dataFile = path.Join(dataDir, "exData.tar.gz")
}

func LoadAllResources() error {
	var res = &Resource{horde3d.H3DRes(0)}
	var err error

	for {
		res.H3DRes = horde3d.NextResource(horde3d.ResTypes_Undefined, res.H3DRes)
		if int(res.H3DRes) != 0 {
			err = res.Load()
			if err != nil {
				return err
			}
		} else {
			break
		}
	}
	return nil
}

func ResourcesNotLoaded() []*Resource {
	notLoaded := make([]*Resource, 0)
	var res = &Resource{horde3d.H3DRes(0)}

	for {
		res = &Resource{horde3d.NextResource(horde3d.ResTypes_Undefined, res.H3DRes)}
		if int(res.H3DRes) != 0 {
			if !res.IsLoaded() {
				notLoaded = append(notLoaded, res)
			}
		} else {
			break
		}
	}
	return notLoaded

}

func ResourceList() []*Resource {
	resList := make([]*Resource, 0)
	res := &Resource{horde3d.H3DRes(0)}

	for {
		res = &Resource{horde3d.NextResource(horde3d.ResTypes_Undefined, res.H3DRes)}
		if int(res.H3DRes) != 0 {
			resList = append(resList, res)
		} else {
			break
		}
	}
	return resList

}

type Resource struct {
	horde3d.H3DRes
}

func (res *Resource) String() string {
	return res.Name()
}

const (
	ResTypeUndefined      = horde3d.ResTypes_Undefined
	ResTypeSceneGraph     = horde3d.ResTypes_SceneGraph
	ResTypeGeometry       = horde3d.ResTypes_Geometry
	ResTypeAnimation      = horde3d.ResTypes_Animation
	ResTypeMaterial       = horde3d.ResTypes_Material
	ResTypeCode           = horde3d.ResTypes_Code
	ResTypeShader         = horde3d.ResTypes_Shader
	ResTypeTexture        = horde3d.ResTypes_Texture
	ResTypeParticleEffect = horde3d.ResTypes_ParticleEffect
	ResTypePipeline       = horde3d.ResTypes_Pipeline
)

func (res *Resource) Type() int { return res.H3DRes.Type() }

func (res *Resource) Name() string { return res.H3DRes.Name() }

//virtualData stores resource data dynamically created
// during the operation of the engine.  Overlays,
// geometry, etc
var virtualData = make(map[string][]byte)

//AddVirtualResource adds the passed in data array
// as resource in memory,
func setVirtualResource(resourceName string, data []byte) {
	virtualData[resourceName] = data
}

func removeVirtualResource(resourceName string) {
	delete(virtualData, resourceName)
}

//NewVirtualResource Creates a new virtual resource from the passed in byte array
func NewVirtualResource(name string, resType int, data []byte) *Resource {
	if resType == ResTypeTexture {
		RaiseError(errors.New("Virtualized Textures must use NewVirtualTexture"))
		return nil
	}
	newRes := &Resource{horde3d.AddResource(resType,
		virtualPath+name, 0)}
	if newRes.H3DRes == 0 {
		err := errors.New("Unable to add resource " + name + " in Horde3D.")
		RaiseError(err)
		return nil
	}

	setVirtualResource(virtualPath+name, data)

	return newRes
}

//NewVirtualTexture adds a new texture in memory.  fmt refers to horde stream format enum
func NewVirtualTexture(name string, width, height, format, flags int) *Texture {
	newRes := &Texture{&Resource{horde3d.CreateTexture(name, width, height, format, flags)}}
	if newRes.H3DRes == 0 {
		err := errors.New("Unable to add virtual texture resource " + name + " in Horde3D.")
		RaiseError(err)
		return nil
	}

	//set empty virtual data, because texture data is set in Texture.SetData()
	setVirtualResource(name, []byte{0})

	return newRes
}

//Resources will be loaded either from a directory or from a
// zip data file.  If the both the data folder and data
// file exist, it will first attempt to load from the data folder.
// if the resource doesn't exist in the data folder, then it will
// try to load it from the data file
func (res *Resource) Load() error {
	if !res.IsLoaded() {
		data, err := loadEngineData(res.FullPath())
		if err != nil {
			return err
		}
		good := res.H3DRes.Load(data)
		if !good {
			err := errors.New("Horde3D was unable to load the resource " + res.FullPath() + ".")
			RaiseError(err)
			return err
		}
	}

	return nil

}

func (res *Resource) IsVirtual() bool {
	_, ok := virtualData[res.Name()]
	return ok
}

func loadEngineData(resourcePath string) ([]byte, error) {

	//	Loads virtual resource from memory
	if data, ok := virtualData[resourcePath]; ok {
		return data, nil
	}

	if !path.IsAbs(resourcePath) {
		resourcePath = path.Join(dataDir, resourcePath)
	}

	data, err := ioutil.ReadFile(resourcePath)

	if os.IsNotExist(err) {
		//err = nil
		//TODO: load from zip
		//remove respath root
		RaiseError(err)
		return nil, err
	}

	if err != nil {
		RaiseError(err)
		return nil, err
	}

	return data, nil
}

func (res *Resource) FullPath() string {
	if res.IsVirtual() {
		return res.Name()
	}
	return path.Join(dataDir, res.Name())
}

func (res *Resource) Clone(cloneName string) *Resource {
	clone := new(Resource)
	clone.H3DRes = res.H3DRes.Clone(cloneName)
	return clone
}

//IsValid Returns if the resource this struct points to is still a valid
// resource
func (res *Resource) IsValid() bool {
	if res.Type() == ResTypeUndefined {
		return false
	}
	return true
}

func (res *Resource) Remove() {
	if res.IsVirtual() {
		removeVirtualResource(res.Name())
	}
	res.H3DRes.Remove()
}

func (res *Resource) IsLoaded() bool { return res.H3DRes.IsLoaded() }

func (res *Resource) Unload() {
	res.H3DRes.Unload()
}

type Scene struct{ *Resource }

func NewScene(name string) (*Scene, error) {
	scene := &Scene{new(Resource)}
	scene.H3DRes = horde3d.AddResource(horde3d.ResTypes_SceneGraph,
		name, 0)
	if scene.H3DRes == 0 {
		err := errors.New("Unable to add resource " + name + " in Horde3D.")
		RaiseError(err)
		return nil, err
	}

	return scene, nil
}

type Geometry struct{ *Resource }

func NewGeometry(name string) (*Geometry, error) {
	geo := &Geometry{new(Resource)}

	geo.H3DRes = horde3d.AddResource(horde3d.ResTypes_Geometry,
		name, 0)

	if geo.H3DRes == 0 {
		err := errors.New("Unable to add resource " + name + " in Horde3D.")
		RaiseError(err)
		return nil, err
	}
	return geo, nil
}

type Animation struct{ *Resource }

func NewAnimation(name string) (*Animation, error) {
	anim := &Animation{new(Resource)}
	anim.H3DRes = horde3d.AddResource(horde3d.ResTypes_Animation,
		name, 0)
	if anim.H3DRes == 0 {
		err := errors.New("Unable to add resource " + name + " in Horde3D.")
		RaiseError(err)
		return nil, err
	}

	return anim, nil
}

type Material struct{ *Resource }

func NewMaterial(name string) (*Material, error) {
	material := &Material{&Resource{horde3d.AddResource(horde3d.ResTypes_Material,
		name, 0)}}
	if material.H3DRes == 0 {
		err := errors.New("Unable to add resource " + name + " in Horde3D.")
		RaiseError(err)
		return nil, err
	}

	return material, nil
}

func (m *Material) SetUniform(name string, a, b, c, d float32) bool {
	return horde3d.SetMaterialUniform(m.H3DRes, name, a, b, c, d)
}

type ShaderCode struct{ *Resource }

type Shader struct{ *Resource }
type Texture struct{ *Resource }

func (t *Texture) SetData(data image.Image) {
	width := data.Bounds().Size().X
	height := data.Bounds().Size().Y

	stream, err := t.MapUint8ResStream(horde3d.TexRes_ImageElem, 0, horde3d.TexRes_ImgPixelStream,
		false, true, (width*height)*4)

	if err != nil {
		RaiseError(err)
		return
	}

	//TODO: Fix format
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			r, g, b, a := data.At(x, y).RGBA()
			index := (y * (width * 4)) + (x * 4)

			stream[index] = uint8(r)
			stream[index+1] = uint8(g)
			stream[index+2] = uint8(b)
			stream[index+3] = uint8(a)
		}

	}

	t.UnmapResStream()
}

type ParticleEffect struct{ *Resource }

func NewParticleEffect(name string) (*ParticleEffect, error) {
	part := &ParticleEffect{new(Resource)}

	part.H3DRes = horde3d.AddResource(horde3d.ResTypes_ParticleEffect,
		name, 0)
	if part.H3DRes == 0 {
		err := errors.New("Unable to add resource " + name + " in Horde3D.")
		RaiseError(err)
		return nil, err
	}
	return part, nil
}

type Pipeline struct{ *Resource }

func NewPipeline(name string) (*Pipeline, error) {
	pipeline := &Pipeline{new(Resource)}
	pipeline.H3DRes = horde3d.AddResource(horde3d.ResTypes_Pipeline,
		name, 0)

	if pipeline.H3DRes == 0 {
		err := errors.New("Unable to add resource " + name + " in Horde3D.")
		RaiseError(err)
		return nil, err
	}

	return pipeline, nil

}

//loadPipeline loads the default pipeline for the engine
func loadDefaultPipeline() (*Pipeline, error) {
	pipeline, err := NewPipeline("pipelines/default.pipeline.xml")
	if err != nil {
		return nil, err
	}
	if err = pipeline.Load(); err != nil {
		return nil, err
	}
	return pipeline, nil
}

func (p *Pipeline) ResizeBuffers(width, height int) {
	horde3d.ResizePipelineBuffers(p.H3DRes, width, height)
}
