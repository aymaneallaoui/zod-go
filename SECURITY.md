# Security Policy

## Supported Versions

We take security seriously and provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | ✅ Yes            |
| 0.x.x   | ❌ No (EOL)       |

## Reporting a Vulnerability

We appreciate your efforts to responsibly disclose security vulnerabilities. If you discover a security vulnerability in zod-go, please follow these steps:

### 1. Do NOT Create a Public Issue

Please **do not** create a public GitHub issue for security vulnerabilities. This helps protect users until a fix can be developed and released.

### 2. Report Privately

Send your security report to: **security@[maintainer-email].com** (or create a private security advisory on GitHub)

Include the following information in your report:

- **Description**: A clear description of the vulnerability
- **Impact**: What could an attacker accomplish with this vulnerability?
- **Reproduction**: Step-by-step instructions to reproduce the issue
- **Proof of Concept**: If applicable, include a minimal code example
- **Suggested Fix**: If you have ideas for how to fix the issue

### 3. What to Expect

- **Acknowledgment**: We will acknowledge receipt of your report within 48 hours
- **Initial Assessment**: We will provide an initial assessment within 5 business days
- **Regular Updates**: We will keep you informed of our progress
- **Coordinated Disclosure**: We will work with you to determine an appropriate disclosure timeline

### 4. Responsible Disclosure Guidelines

We ask that you:

- Give us reasonable time to investigate and fix the issue before disclosure
- Do not exploit the vulnerability beyond what is necessary to demonstrate it
- Do not access, modify, or delete data belonging to others
- Do not perform testing on our production infrastructure

## Security Measures

### Input Validation

zod-go includes several built-in security measures:

- **Regex Safety**: All user-provided regex patterns are validated to prevent ReDoS attacks
- **Input Bounds**: Numeric operations include bounds checking to prevent overflow
- **Memory Safety**: Protection against memory exhaustion attacks through input size limits
- **Type Safety**: Strong type checking prevents injection-style attacks

### Known Security Considerations

When using zod-go, be aware of these security considerations:

1. **Regex Patterns**: User-provided regex patterns should be from trusted sources
2. **Large Inputs**: Very large arrays or objects may consume significant memory
3. **Custom Validators**: Custom validation functions run with full application privileges
4. **Error Messages**: Error messages may contain sensitive input data

### Security Best Practices

#### 1. Input Sanitization
```go
// Always validate untrusted input
userSchema := validators.Object(map[string]zod.Schema{
    "email": validators.String().Email().Required(),
    "name":  validators.String().Min(1).Max(100).Required(), // Prevent excessive length
})
```

#### 2. Regex Safety
```go
// Use pre-defined patterns for security-sensitive validation
emailSchema := validators.String().Email() // Uses built-in safe regex

// If using custom patterns, ensure they're from trusted sources
safePattern := validators.String().Pattern(`^[a-zA-Z0-9_]+$`) // Simple, safe pattern
```

#### 3. Memory Protection
```go
// Limit array sizes to prevent memory exhaustion
tagsSchema := validators.Array(validators.String().Required()).
    Max(100). // Reasonable upper limit
    Unique()
```

#### 4. Error Handling
```go
// Be careful about exposing validation errors to users
if err := schema.Validate(userInput); err != nil {
    // Log detailed error for debugging
    log.Printf("Validation failed: %v", err)
    
    // Return generic error to user
    return errors.New("invalid input")
}
```

## Security Vulnerability Examples

### Regex Denial of Service (ReDoS)

**Risk**: Malicious regex patterns can cause excessive CPU usage.

**Mitigation**: zod-go validates regex patterns and includes timeouts.

```go
// This is protected against ReDoS
schema := validators.String().Pattern(userProvidedPattern)
```

### Memory Exhaustion

**Risk**: Very large inputs can consume excessive memory.

**Mitigation**: Set reasonable limits on input sizes.

```go
// Protect against large inputs
schema := validators.Array(elementSchema).Max(1000)
objectSchema := validators.Object(fields).Strict() // Reject unknown fields
```

### Information Disclosure

**Risk**: Error messages may expose sensitive information.

**Mitigation**: Sanitize error messages before returning to users.

```go
if err := schema.Validate(sensitiveData); err != nil {
    // Don't expose the actual validation error to users
    log.Printf("Internal validation error: %v", err)
    return errors.New("validation failed")
}
```

## Security Updates

We will:

- Release security fixes as soon as possible
- Clearly mark security releases in the changelog
- Provide upgrade guidance for security issues
- Coordinate with package managers and security databases

## Bug Bounty

Currently, we do not have a formal bug bounty program. However, we greatly appreciate security researchers who help improve the security of zod-go and will:

- Acknowledge your contribution in release notes (if desired)
- Provide early access to security fixes for testing
- Consider feature requests that improve security

## Contact

For security-related questions or concerns:

- **Security Issues**: Use GitHub Security Advisories or email security@[maintainer-email].com
- **General Security Questions**: Create a public issue with the "security" label
- **Security Best Practices**: See our documentation and examples

## Security Checklist for Contributors

If you're contributing to zod-go, please ensure:

- [ ] New features include security considerations
- [ ] Input validation is applied to all user-provided data
- [ ] Error messages don't expose sensitive information
- [ ] Performance implications are considered (DoS protection)
- [ ] Dependencies are kept up to date
- [ ] Security tests are included where appropriate

Thank you for helping keep zod-go secure!
