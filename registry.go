package yamyams

import (
	myapplication "github.com/kris-nova/yamyams/apps/_example"
	staticsite "github.com/kris-nova/yamyams/apps/staticsite"
	yamyams "github.com/kris-nova/yamyams/pkg"
)

func Load() {

	yamyams.Register(myapplication.New())

	yamyams.Register(staticsite.New("default"))

}
