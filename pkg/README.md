# Library code

Notice this package does not log because it should not log.

Notice this package does not terminate the program because it should not terminate the program.

Notice this package does have tests because it should have tests.

Notice this package implements both `Install` as well as `Uninstall` because if it implements one it should implement the other.

### Using this

Yes. I know that you can compile the code, put it in a container, and you will still need YAML to "deploy" your "deployer".

You are missing the point. 

The point is that we now have created a library so we can use it *anywhere*. 

We can use this library to create controllers, to create operators, in other tools, in other packages, this is now trustworthy dependable code you can use wherever your cold dead heart desires.

The point is that you TRUST this code, and you can kick your feet up on the desk and sip a mimosa while your code runs because you know it is going to do what you want it to do.

As the needs of you and your team change and grow so will your deployment libraries. 

On day one you might have something very simple like this to start. 

Fast forward a few years and you have two options

 1. Iterate on YAML templates for 2 years
 2. Iterate on Go for 2 years

I don't know about you but I know which choice I would pick.

Go gives you an entire suite of tools, patterns, examples, and policy to build with. 

The more applications you and your team take on, the more complex your deployment toolchain might be. On day one there is nothing wrong with a CLI tool that vendors a library.
If you write the library well, it should be fairly simple to migrate that to multiple tools, and later a controller, or a microservice, or even a full-fledged public facing API.

Invest in yourself. Invest in your team.

Stop interpolating YAML.

Thank you for coming to my TED talk.