# Manifest Testing

Yes. I am testing your YAML.

Here are all the YAML files we use to exercise naml.
These are real YAML files, with real kubernetes applications.

### Adding a test

Create a new `test_TITLE.yaml` file in the `tests/manifest` directory.

The files are picked up automatically at compile/test time.

`make test` will fail if the program cannot compile. 
`make test` will not actually run/install your program in a cluster.