package engine

import (
	"code.google.com/p/gohorde/horde3d"
	"io/ioutil"
)

//Resources will be loaded either from a directory or from a
// tar.gz data file.  If the both the data folder and data tar.gz
// file exist, it will first attempt to load from the data folder.
// if the resource doesn't exist in the data folder, then it will
// try to load it from the data file

//LoadPipeline loads the default pipeline for the engine
func LoadPipeline() Resource {
	pipeline := Resource{}
	pipeline.H3DRes = horde3d.AddResource(horde3d.ResTypes_Pipeline,
		"pipelines/hdr.pipeline.xml", 0)
}

type Resource struct {
	horde3d.H3DRes
	relPath string
}

func (res *Resource) Type() int {
	return horde3d.GetResType(res.H3DRes)
}

func (res *Resource) Load() bool {

}
