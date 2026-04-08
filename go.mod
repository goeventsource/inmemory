module github.com/goeventsource/inmemory

go 1.26

require (
	github.com/goeventsource/goeventsource v0.0.0
	github.com/google/uuid v1.6.0
)

// Local layout under new_org/; remove after tagged releases. CI overrides paths via go mod edit.
replace github.com/goeventsource/goeventsource v0.0.0 => ../goeventsource
