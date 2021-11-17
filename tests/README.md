# YAML Testing

Yes. I am testing your YAML.

Here are all the YAML files we use to exercise naml.
These are real YAML files, with real kubernetes applications.

### Adding a test

Create a new `test_TITLE.yaml` file in this directory.

Create a specific test (with a specific title) in the tests directory.
We prefer to use specific tests for the plumbing so we can easily tell which test failed.

If you created `test_boops.yaml` your test would look like:

```go
func TestBoops(t *testing.T) {
	err := generateCompileRunYAML("test_boops.yaml")
	if err != nil {
		t.Errorf(err.Error())
	}
}

```