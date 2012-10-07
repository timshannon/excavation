package entity

type Entity interface {
	Load(args []string)
	Remove()
}

var entities []Entity

func init() {
	entities = make([]Entity)
}

func LoadEntity(attachmentData string) {

}
