package events

type ServiceAdded struct {
	ID string
}

type ServiceRemoved struct {
	ID string
}

type ServiceUpdated struct {
	ID string
}

type ServiceTagged struct {
	ID  string
	Tag string
}

type ServiceUntagged struct {
	ID  string
	Tag string
}

type TagCreated struct {
	Name string
}
