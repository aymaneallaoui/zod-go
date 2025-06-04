---
name: Feature Request
about: Suggest an idea for zod-go
title: '[FEATURE] '
labels: ['enhancement', 'needs-triage']
assignees: ''

---

## Feature Description
A clear and concise description of the feature you'd like to see added.

## Problem Statement
Describe the problem this feature would solve. Is your feature request related to a problem? Please describe.

## Proposed Solution
Describe the solution you'd like. How should this feature work?

## Example Usage
Please provide a code example of how you envision this feature being used:

```go
package main

import (
    "github.com/aymaneallaoui/zod-go/zod/validators"
)

func main() {
    // Example of how the new feature would be used
    schema := validators.NewFeature().
        SomeMethod().
        AnotherMethod()
    
    err := schema.Validate(someData)
    // ...
}
```

## Alternatives Considered
Describe any alternative solutions or features you've considered.

## Additional Context
Add any other context, screenshots, or examples about the feature request here.

## Use Cases
Describe specific use cases where this feature would be valuable:

1. **Use Case 1**: Description of when this would be useful
2. **Use Case 2**: Another scenario where this helps
3. **Use Case 3**: Additional use case

## Implementation Ideas
If you have ideas about how this could be implemented, please share them:

- [ ] API design considerations
- [ ] Performance implications
- [ ] Backward compatibility concerns
- [ ] Testing requirements

## Priority
How important is this feature to you?

- [ ] Nice to have
- [ ] Would be helpful
- [ ] Important for my use case
- [ ] Critical/blocking issue

## Checklist
- [ ] I have checked the [documentation](https://github.com/aymaneallaoui/zod-go#readme)
- [ ] I have searched for similar feature requests
- [ ] I have provided example usage
- [ ] I have described the problem this solves
- [ ] I have considered alternative solutions
