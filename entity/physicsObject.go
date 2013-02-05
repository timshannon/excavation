package entity

import (
	"excavation/engine"
)

type PhysicsObject struct {
	body *engine.PhysicsBody
}

func (p *PhysicsObject) Add(node *engine.Node, args EntityArgs) {

	//fmt.Println("nodeType: ", node.Type())
	//posData := []float32{0, 0, 0,
		//10, 0, 0,
		//0, 10, 0,
		//10, 10, 0}

	//indexData := []uint32{0, 1, 2, 2, 1, 3}
	//normalData := []int16{0, 0, 1,
		//0, 0, 1,
		//0, 0, 1,
		//0, 0, 1}

	//uvData := []float32{
		//0, 0,
		//1, 0,
		//0, 1,
		//1, 1}

	//res := horde3d.H3DRes(horde3d.GetNodeParamI(node.H3DNode, horde3d.MatRes_MaterialElem))

	//geoRes := horde3d.CreateGeometryRes("geoRes", 4, 6, posData, indexData, normalData, []int16{0}, []int16{0},
		//uvData, []float32{0})
	//model := engine.NewNode(horde3d.AddModelNode(horde3d.RootNode, "DynGeoModelNode", geoRes))
	//horde3d.AddMeshNode(model.H3DNode, "DynGeoMesh", res, 0, 6, 0, 3)

	//isInt16 := horde3d.GetResParamI(geoRes, horde3d.GeoRes_GeometryElem, 0,
		//horde3d.GeoRes_GeoIndices16I)
	//fmt.Println("isInt16", isInt16)

	//indexCount := horde3d.GetResParamI(geoRes, horde3d.GeoRes_GeometryElem, 0,
		//horde3d.GeoRes_GeoIndexCountI)
	//fmt.Println("Index Count: ", indexCount)

	//vertCount := horde3d.GetResParamI(geoRes, horde3d.GeoRes_GeometryElem, 0,
		//horde3d.GeoRes_GeoVertexCountI)
	//fmt.Println("VertCount: ", vertCount)
	//vertices, _ := horde3d.MapFloatResStream(geoRes, horde3d.GeoRes_GeometryElem, 0,
		//horde3d.GeoRes_GeoVertPosStream, true, false, vertCount*3)
	//fmt.Println("Vertices: ", vertices)

	p.body = engine.AddPhysicsBody(node, args.Float("mass"))

}

func (p *PhysicsObject) Trigger(value float32) {
	return
}
