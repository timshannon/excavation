// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"errors"
	"io/ioutil"
	"os"
	"path"
)

var (
	dataDir  string
	dataFile string
)

func init() {
	wd, _ := os.Getwd()
	dataDir = path.Join(wd, "data")
	dataFile = path.Join(dataDir, "exData.tar.gz")
}

func LoadAllResources() error {

	var res = &Resource{horde3d.GetNextResource(horde3d.ResTypes_Undefined, horde3d.H3DRes(0))}
	var err error

	for int(res.H3DRes) != 0 {
		res.H3DRes = horde3d.GetNextResource(horde3d.ResTypes_Undefined, res.H3DRes)
		if int(res.H3DRes) != 0 {
			err = res.Load()
			//fmt.Println("Loading: ", res.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//might use later with audio
//type ResHandler interface {
//Type() int
//Name() string
//Load() error
//FullPath() string
//IsLoaded() bool
//Remove()
//}

type Resource struct {
	horde3d.H3DRes
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

func (res *Resource) Type() int { return horde3d.GetResType(res.H3DRes) }

func (res *Resource) Name() string { return horde3d.GetResName(res.H3DRes) }

//Resources will be loaded either from a directory or from a
// tar.gz data file.  If the both the data folder and data tar.gz
// file exist, it will first attempt to load from the data folder.
// if the resource doesn't exist in the data folder, then it will
// try to load it from the data file
func (res *Resource) Load() error {
	if !res.IsLoaded() {
		data, err := loadEngineData(res.FullPath())
		if err != nil {
			return err
		}
		good := horde3d.LoadResource(res.H3DRes, data)
		if !good {
			err := errors.New("Horde3D was unable to load the resource " + res.FullPath() + ".")
			RaiseError(err)
			return err
		}
	}

	return nil

}

func loadEngineData(resourcePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(resourcePath)

	if os.IsNotExist(err) {
		err = nil
		//TODO: load from tar.gz data file
		//remove respath root
	}

	if err != nil {
		RaiseError(err)
		return nil, err
	}

	return data, nil
}

func (res *Resource) FullPath() string {
	return path.Join(dataDir, res.Name())
}

func (res *Resource) Clone(cloneName string) *Resource {
	clone := new(Resource)
	clone.H3DRes = horde3d.CloneResource(res.H3DRes, cloneName)
	return clone
}

func (res *Resource) Remove() { horde3d.RemoveResource(res.H3DRes) }

func (res *Resource) IsLoaded() bool { return horde3d.IsResLoaded(res.H3DRes) }

func (res *Resource) Unload() { horde3d.UnloadResource(res.H3DRes) }

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

func (m *Material) SetUniform(name string, a, b, c, d float32) bool {
	return horde3d.SetMaterialUniform(m.H3DRes, name, a, b, c, d)
}

type ShaderCode struct{ *Resource }

type Shader struct{ *Resource }
type Texture struct{ *Resource }

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

//LoadPipeline loads the default pipeline for the engine
func LoadPipeline() (*Pipeline, error) {
	pipeline, err := NewPipeline("pipelines/hdr.pipeline.xml")
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
