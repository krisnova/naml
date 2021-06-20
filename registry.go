package yamyams

import (
	mydeployment "github.com/kris-nova/yamyams/apps/_deployment"
	myapplication "github.com/kris-nova/yamyams/apps/_example"
	yamyams "github.com/kris-nova/yamyams/pkg"
)

// Load is where we can set up applications.
//
// This is called whenever the yamyams program starts.
func Load() {

	// We can keep them very simple, and hard code all the logic like this one.
	yamyams.Register(myapplication.New())

	// We can also have several instances of the same application like this.
	yamyams.Register(mydeployment.New("default", "example-1", "beeps", 3))
	yamyams.Register(mydeployment.New("default", "example-2", "boops", 1))
	yamyams.Register(mydeployment.New("default", "example-3", "cyber boops", 7))

}
