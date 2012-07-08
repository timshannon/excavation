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
	dataFile = path.Join(wd, "exData.tar.gz")
}

//LoadPipeline loads the default pipeline for the engine
func LoadPipeline() (*Resource, error) {
	pipeline := new(Resource)
	pipeline.H3DRes = horde3d.AddResource(horde3d.ResTypes_Pipeline,
		"pipelines/hdr.pipeline.xml", 0)
	if err := pipeline.Load(); err != nil {
		return nil, err
	}
	return pipeline, nil
}

type Resource struct {
	horde3d.H3DRes
}

func (res *Resource) Type() int {
	return horde3d.GetResType(res.H3DRes)
}

func (res *Resource) Name() string {
	return horde3d.GetResName(res.H3DRes)
}

//Resources will be loaded either from a directory or from a
// tar.gz data file.  If the both the data folder and data tar.gz
// file exist, it will first attempt to load from the data folder.
// if the resource doesn't exist in the data folder, then it will
// try to load it from the data file
func (res *Resource) Load() error {
	data, err := ioutil.ReadFile(path.Join(dataDir, res.Name()))

	if os.IsNotExist(err) {
		//TODO: load from tar.gz data file
	}

	if err != nil {
		return err
	}

	good := horde3d.LoadResource(res.H3DRes, data)
	if !good {
		return errors.New("Horde3D was unable to load the resource.")
	}

	return nil

}

func (res *Resource) Clone(cloneName string) *Resource {
	clone := new(Resource)
	clone.H3DRes = horde3d.CloneResource(res.H3DRes, cloneName)
	return clone
}

func (res *Resource) Remove() int {
	return horde3d.RemoveResource(res.H3DRes)
}

func (res *Resource) IsLoaded() bool {
	return horde3d.IsResLoaded(res.H3DRes)
}

func (res *Resource) Unload() {
	horde3d.UnloadResource(res.H3DRes)
}
