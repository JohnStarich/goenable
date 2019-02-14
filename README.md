# goenable

A launching platform to write Bash loadables (custom builtins) in Go.

## Quick start

```go
make
enable -f ./out/goenable.so goenable

# Print usage
help goenable

# Load a plugin and run it
goenable load ./out/namespace output  # Load the namespace plugin
eval "$output"  # Prepare functions from the namespace plugin
namespace output ./your_script.sh
eval "$output"
# Now your_script.sh has been loaded in with namespaced functions!
```
