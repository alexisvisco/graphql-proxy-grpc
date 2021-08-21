package virtual

type Package struct {
	Descriptor *Descriptor `json:"-"`

	ProtoPackage string
	Name         string
	Path         string

	Messages map[string]*Message
	Enums    map[string]*Enum
	Services map[string]*Service
	Dir      string
}

func (pkg Package) String() string {
	return pkg.Path
}

func (pkg *Package) addEnum(e *Enum) {
	pkg.Enums[e.String()] = e
}

func (pkg *Package) addMessage(e *Message) {
	pkg.Messages[e.String()] = e
}

func (pkg Package) addService(svc *Service) {
	pkg.Services[svc.String()] = svc
}

func (pkg Package) getMessage(n Name) *Message {
	return pkg.Messages[n.String()]
}

func (pkg Package) getEnum(n Name) *Enum {
	return pkg.Enums[n.String()]
}

func NewPackage(parent *Descriptor, protoPackage, name, path, dir string) *Package {
	return &Package{
		Descriptor:   parent,
		ProtoPackage: protoPackage,
		Name:         name,
		Path:         path,
		Dir:          dir,
		Messages:     make(map[string]*Message),
		Enums:        make(map[string]*Enum),
		Services:     make(map[string]*Service),
	}
}
