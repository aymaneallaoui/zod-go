// Package main demonstrates real-world usage patterns with the enhanced zod-go API
// This file showcases practical examples that developers would encounter in production applications.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aymaneallaoui/zod-go/zod/validators"
)

func main() {
	fmt.Println("🌟 Real-World Examples for Enhanced zod-go")
	fmt.Println("==========================================")

	runUserRegistrationExample()
	runAPIValidationExample()
	runConfigurationValidationExample()
	runECommerceExample()
	runBlogSystemExample()
	runAnalyticsExample()
	runFileUploadExample()
	runSocialMediaExample()
}

// Example 1: User Registration System
func runUserRegistrationExample() {
	fmt.Println("\n📝 Example 1: User Registration System")
	fmt.Println("=====================================")

	// Define comprehensive user registration schema
	userRegistrationSchema := validators.Object(map[string]interface{}{
		"username": validators.String().
			Min(3).WithMinLengthMessage("Username must be at least 3 characters").
			Max(20).WithMaxLengthMessage("Username cannot exceed 20 characters").
			Pattern(`^[a-zA-Z0-9_]+$`).WithPatternMessage("Username can only contain letters, numbers, and underscores").
			Required().WithRequiredMessage("Username is required"),

		"email": validators.Email().
			WithEmailMessage("Please enter a valid email address").
			WithRequiredMessage("Email address is required"),

		"password": validators.String().
			Min(8).WithMinLengthMessage("Password must be at least 8 characters").
			Pattern(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]`).
			WithPatternMessage("Password must contain uppercase, lowercase, digit, and special character").
			Required().WithRequiredMessage("Password is required"),

		"confirmPassword": validators.String().
			Required().WithRequiredMessage("Password confirmation is required"),

		"age": validators.Number().
			Min(13).WithMinMessage("Must be at least 13 years old").
			Max(120).WithMaxMessage("Please enter a valid age").
			Integer().WithIntegerMessage("Age must be a whole number").
			Required().WithRequiredMessage("Age is required"),

		"phoneNumber": validators.String().
			Pattern(`^\+?[1-9]\d{1,14}$`).WithPatternMessage("Please enter a valid phone number").
			Optional(),

		"interests": validators.Array(validators.String().Min(1)).
			MinItems(0).
			MaxItems(10).WithMaxItemsMessage("Cannot select more than 10 interests").
			Optional(),

		"newsletter": validators.Bool().
			Optional().Default(false),

		"termsAccepted": validators.Bool().
			Required().WithRequiredMessage("You must accept the terms and conditions"),
	}).Required().WithRequiredMessage("Registration data is required")

	// Test valid registration
	validRegistration := map[string]interface{}{
		"username":        "john_doe",
		"email":           "john.doe@example.com",
		"password":        "SecurePass123!",
		"confirmPassword": "SecurePass123!",
		"age":             25,
		"phoneNumber":     "+1234567890",
		"interests":       []string{"programming", "music", "travel"},
		"newsletter":      true,
		"termsAccepted":   true,
	}

	fmt.Printf("✅ Valid registration: %s\n", validateAndReport(userRegistrationSchema, validRegistration))

	// Test invalid registration
	invalidRegistration := map[string]interface{}{
		"username":        "jo", // Too short
		"email":           "invalid-email",
		"password":        "weak",
		"confirmPassword": "different",
		"age":             12, // Too young
		"termsAccepted":   false,
	}

	fmt.Printf("❌ Invalid registration: %s\n", validateAndReport(userRegistrationSchema, invalidRegistration))
}

// Example 2: REST API Request/Response Validation
func runAPIValidationExample() {
	fmt.Println("\n🌐 Example 2: REST API Validation")
	fmt.Println("=================================")

	// Blog post creation API
	createPostSchema := validators.Object(map[string]interface{}{
		"title": validators.String().
			Min(1).WithMinLengthMessage("Title cannot be empty").
			Max(200).WithMaxLengthMessage("Title cannot exceed 200 characters").
			Required(),

		"content": validators.String().
			Min(10).WithMinLengthMessage("Content must be at least 10 characters").
			Max(10000).WithMaxLengthMessage("Content cannot exceed 10,000 characters").
			Required(),

		"tags": validators.Array(validators.String().Min(1).Max(50)).
			MinItems(1).WithMinItemsMessage("At least one tag is required").
			MaxItems(5).WithMaxItemsMessage("Cannot have more than 5 tags").
			Required(),

		"publishedAt": validators.String().Optional(), // ISO date string

		"featured": validators.Bool().Optional().Default(false),

		"metadata": validators.Object(map[string]interface{}{
			"seoTitle":       validators.String().Max(60).Optional(),
			"seoDescription": validators.String().Max(160).Optional(),
			"canonicalUrl":   validators.URL().Optional(),
		}).Optional(),
	}).Required()

	// API response schema
	apiResponseSchema := validators.Object(map[string]interface{}{
		"success": validators.Bool().Required(),
		"data": validators.Object(map[string]interface{}{
			"id":        validators.Number().Integer().Required(),
			"title":     validators.String().Required(),
			"slug":      validators.String().Required(),
			"createdAt": validators.String().Required(),
			"updatedAt": validators.String().Required(),
		}).Optional(),
		"error": validators.Object(map[string]interface{}{
			"code":    validators.String().Required(),
			"message": validators.String().Required(),
			"details": validators.Array(validators.String()).Optional(),
		}).Optional(),
		"pagination": validators.Object(map[string]interface{}{
			"page":     validators.Number().Integer().Min(1).Required(),
			"pageSize": validators.Number().Integer().Min(1).Max(100).Required(),
			"total":    validators.Number().Integer().Min(0).Required(),
		}).Optional(),
	}).Required()

	// Test API request
	validRequest := map[string]interface{}{
		"title":   "Getting Started with zod-go",
		"content": "This is a comprehensive guide to using zod-go for validation in Go applications...",
		"tags":    []string{"go", "validation", "tutorial"},
		"featured": true,
		"metadata": map[string]interface{}{
			"seoTitle":       "zod-go Tutorial - Complete Guide",
			"seoDescription": "Learn how to use zod-go for type-safe validation in Go",
			"canonicalUrl":   "https://example.com/zod-go-guide",
		},
	}

	fmt.Printf("✅ Valid API request: %s\n", validateAndReport(createPostSchema, validRequest))

	// Test API response
	successResponse := map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":        123,
			"title":     "Getting Started with zod-go",
			"slug":      "getting-started-with-zod-go",
			"createdAt": "2023-12-06T10:30:00Z",
			"updatedAt": "2023-12-06T10:30:00Z",
		},
	}

	fmt.Printf("✅ Valid API response: %s\n", validateAndReport(apiResponseSchema, successResponse))
}

// Example 3: Configuration File Validation
func runConfigurationValidationExample() {
	fmt.Println("\n⚙️ Example 3: Application Configuration")
	fmt.Println("======================================")

	configSchema := validators.Object(map[string]interface{}{
		"server": validators.Object(map[string]interface{}{
			"host": validators.String().Optional().Default("localhost"),
			"port": validators.Number().
				Min(1).Max(65535).
				Integer().
				Optional().Default(8080.0),
			"ssl": validators.Bool().Optional().Default(false),
		}).Required(),

		"database": validators.Object(map[string]interface{}{
			"driver": validators.String().
				Pattern(`^(postgres|mysql|sqlite)$`).
				Required(),
			"host":     validators.String().Required(),
			"port":     validators.Number().Integer().Required(),
			"database": validators.String().Required(),
			"username": validators.String().Required(),
			"password": validators.String().Required(),
			"maxConnections": validators.Number().
				Integer().Min(1).Max(100).
				Optional().Default(10.0),
			"timeout": validators.Number().
				Min(1).Max(300).
				Optional().Default(30.0),
		}).Required(),

		"redis": validators.Object(map[string]interface{}{
			"host":     validators.String().Optional().Default("localhost"),
			"port":     validators.Number().Integer().Optional().Default(6379.0),
			"password": validators.String().Optional(),
			"db":       validators.Number().Integer().Min(0).Max(15).Optional().Default(0.0),
		}).Optional(),

		"logging": validators.Object(map[string]interface{}{
			"level": validators.String().
				Pattern(`^(debug|info|warn|error)$`).
				Optional().Default("info"),
			"file":   validators.String().Optional(),
			"maxSize": validators.Number().Min(1).Optional().Default(100.0), // MB
		}).Optional(),

		"features": validators.Object(map[string]interface{}{
			"enableMetrics":    validators.Bool().Optional().Default(true),
			"enableProfiling":  validators.Bool().Optional().Default(false),
			"enableRateLimit":  validators.Bool().Optional().Default(true),
			"maxRequestsPerMinute": validators.Number().
				Integer().Min(1).Max(10000).
				Optional().Default(1000.0),
		}).Optional(),
	}).Required()

	validConfig := map[string]interface{}{
		"server": map[string]interface{}{
			"host": "0.0.0.0",
			"port": 3000,
			"ssl":  true,
		},
		"database": map[string]interface{}{
			"driver":        "postgres",
			"host":          "db.example.com",
			"port":          5432,
			"database":      "myapp",
			"username":      "dbuser",
			"password":      "securepassword",
			"maxConnections": 50,
			"timeout":       60,
		},
		"redis": map[string]interface{}{
			"host": "redis.example.com",
			"port": 6379,
			"db":   1,
		},
		"logging": map[string]interface{}{
			"level":   "info",
			"file":    "/var/log/myapp.log",
			"maxSize": 200,
		},
		"features": map[string]interface{}{
			"enableMetrics":        true,
			"enableProfiling":      false,
			"enableRateLimit":      true,
			"maxRequestsPerMinute": 5000,
		},
	}

	fmt.Printf("✅ Valid configuration: %s\n", validateAndReport(configSchema, validConfig))
}

// Example 4: E-Commerce Product Schema
func runECommerceExample() {
	fmt.Println("\n🛒 Example 4: E-Commerce Product Management")
	fmt.Println("===========================================")

	productSchema := validators.Object(map[string]interface{}{
		"name": validators.String().
			Min(1).Max(200).
			Required(),

		"description": validators.String().
			Min(10).Max(5000).
			Required(),

		"price": validators.Number().
			Min(0.01).WithMinMessage("Price must be greater than 0").
			Required(),

		"currency": validators.String().
			Pattern(`^[A-Z]{3}$`).WithPatternMessage("Currency must be a 3-letter ISO code").
			Required(),

		"sku": validators.String().
			Pattern(`^[A-Z0-9-]+$`).WithPatternMessage("SKU must contain only uppercase letters, numbers, and hyphens").
			Required(),

		"category": validators.Object(map[string]interface{}{
			"id":   validators.Number().Integer().Required(),
			"name": validators.String().Required(),
			"path": validators.Array(validators.String()).Required(),
		}).Required(),

		"inventory": validators.Object(map[string]interface{}{
			"quantity":    validators.Number().Integer().Min(0).Required(),
			"lowStock":    validators.Number().Integer().Min(0).Optional().Default(5.0),
			"trackStock":  validators.Bool().Optional().Default(true),
			"backorders":  validators.Bool().Optional().Default(false),
		}).Required(),

		"dimensions": validators.Object(map[string]interface{}{
			"weight": validators.Number().Min(0).Required(),
			"length": validators.Number().Min(0).Required(),
			"width":  validators.Number().Min(0).Required(),
			"height": validators.Number().Min(0).Required(),
			"unit":   validators.String().Pattern(`^(cm|in|mm)$`).Optional().Default("cm"),
		}).Optional(),

		"images": validators.Array(validators.Object(map[string]interface{}{
			"url":     validators.URL().Required(),
			"alt":     validators.String().Required(),
			"primary": validators.Bool().Optional().Default(false),
		})).MinItems(1).MaxItems(10).Required(),

		"variants": validators.Array(validators.Object(map[string]interface{}{
			"id":    validators.Number().Integer().Required(),
			"name":  validators.String().Required(),
			"price": validators.Number().Min(0).Optional(),
			"sku":   validators.String().Optional(),
			"attributes": validators.Object(map[string]interface{}{
				"color": validators.String().Optional(),
				"size":  validators.String().Optional(),
				"material": validators.String().Optional(),
			}).Optional(),
		})).Optional(),

		"tags": validators.Array(validators.String().Min(1)).
			MaxItems(20).
			Optional(),

		"status": validators.String().
			Pattern(`^(draft|active|inactive|discontinued)$`).
			Optional().Default("draft"),
	}).Required()

	productData := map[string]interface{}{
		"name":        "Premium Wireless Headphones",
		"description": "High-quality wireless headphones with noise cancellation and premium sound quality. Perfect for music lovers and professionals.",
		"price":       299.99,
		"currency":    "USD",
		"sku":         "WH-PRE-001",
		"category": map[string]interface{}{
			"id":   1,
			"name": "Electronics",
			"path": []string{"Electronics", "Audio", "Headphones"},
		},
		"inventory": map[string]interface{}{
			"quantity":   50,
			"lowStock":   10,
			"trackStock": true,
			"backorders": false,
		},
		"dimensions": map[string]interface{}{
			"weight": 0.35,
			"length": 20.0,
			"width":  18.0,
			"height": 8.0,
			"unit":   "cm",
		},
		"images": []map[string]interface{}{
			{
				"url":     "https://example.com/images/headphones-1.jpg",
				"alt":     "Premium Wireless Headphones - Front View",
				"primary": true,
			},
			{
				"url":     "https://example.com/images/headphones-2.jpg",
				"alt":     "Premium Wireless Headphones - Side View",
				"primary": false,
			},
		},
		"variants": []map[string]interface{}{
			{
				"id":   1,
				"name": "Black",
				"sku":  "WH-PRE-001-BLK",
				"attributes": map[string]interface{}{
					"color": "Black",
				},
			},
			{
				"id":   2,
				"name": "White",
				"sku":  "WH-PRE-001-WHT",
				"attributes": map[string]interface{}{
					"color": "White",
				},
			},
		},
		"tags":   []string{"wireless", "premium", "noise-cancellation", "bluetooth"},
		"status": "active",
	}

	fmt.Printf("✅ Valid product: %s\n", validateAndReport(productSchema, productData))
}

// Example 5: Blog System with Comments
func runBlogSystemExample() {
	fmt.Println("\n📝 Example 5: Blog System with Comments")
	fmt.Println("======================================")

	commentSchema := validators.Object(map[string]interface{}{
		"id":        validators.Number().Integer().Required(),
		"content":   validators.String().Min(1).Max(1000).Required(),
		"author":    validators.String().Required(),
		"email":     validators.Email().Required(),
		"website":   validators.URL().Optional(),
		"createdAt": validators.String().Required(),
		"approved":  validators.Bool().Optional().Default(false),
		"replies": validators.Array(validators.Object(map[string]interface{}{
			"id":        validators.Number().Integer().Required(),
			"content":   validators.String().Min(1).Max(500).Required(),
			"author":    validators.String().Required(),
			"createdAt": validators.String().Required(),
		})).Optional(),
	})

	blogPostSchema := validators.Object(map[string]interface{}{
		"title": validators.String().
			Min(1).Max(200).
			Required(),

		"slug": validators.String().
			Pattern(`^[a-z0-9-]+$`).WithPatternMessage("Slug must contain only lowercase letters, numbers, and hyphens").
			Required(),

		"excerpt": validators.String().
			Min(10).Max(500).
			Required(),

		"content": validators.String().
			Min(100).
			Required(),

		"author": validators.Object(map[string]interface{}{
			"id":       validators.Number().Integer().Required(),
			"name":     validators.String().Required(),
			"email":    validators.Email().Required(),
			"bio":      validators.String().Max(500).Optional(),
			"avatar":   validators.URL().Optional(),
			"website":  validators.URL().Optional(),
		}).Required(),

		"categories": validators.Array(validators.Object(map[string]interface{}{
			"id":   validators.Number().Integer().Required(),
			"name": validators.String().Required(),
			"slug": validators.String().Required(),
		})).MinItems(1).MaxItems(5).Required(),

		"tags": validators.Array(validators.String().Min(1).Max(50)).
			MaxItems(10).
			Optional(),

		"featuredImage": validators.Object(map[string]interface{}{
			"url":    validators.URL().Required(),
			"alt":    validators.String().Required(),
			"width":  validators.Number().Integer().Min(1).Required(),
			"height": validators.Number().Integer().Min(1).Required(),
		}).Optional(),

		"seo": validators.Object(map[string]interface{}{
			"title":       validators.String().Max(60).Optional(),
			"description": validators.String().Max(160).Optional(),
			"keywords":    validators.Array(validators.String()).MaxItems(10).Optional(),
			"canonical":   validators.URL().Optional(),
		}).Optional(),

		"comments": validators.Array(commentSchema).Optional(),

		"publishedAt": validators.String().Optional(),
		"updatedAt":   validators.String().Required(),
		"status":      validators.String().Pattern(`^(draft|published|archived)$`).Required(),
		"viewCount":   validators.Number().Integer().Min(0).Optional().Default(0.0),
		"featured":    validators.Bool().Optional().Default(false),
	}).Required()

	blogPostData := map[string]interface{}{
		"title":   "Building Type-Safe APIs with Go and zod-go",
		"slug":    "building-type-safe-apis-go-zod",
		"excerpt": "Learn how to create robust, type-safe APIs using Go and the zod-go validation library.",
		"content": "In this comprehensive guide, we'll explore how to build type-safe APIs using Go and zod-go...",
		"author": map[string]interface{}{
			"id":      1,
			"name":    "Jane Developer",
			"email":   "jane@example.com",
			"bio":     "Senior Go developer with 5+ years of experience building scalable web applications.",
			"avatar":  "https://example.com/avatars/jane.jpg",
			"website": "https://jane-dev.com",
		},
		"categories": []map[string]interface{}{
			{
				"id":   1,
				"name": "Go Programming",
				"slug": "go-programming",
			},
			{
				"id":   2,
				"name": "Web Development",
				"slug": "web-development",
			},
		},
		"tags": []string{"go", "api", "validation", "type-safety", "web-development"},
		"featuredImage": map[string]interface{}{
			"url":    "https://example.com/images/go-api-featured.jpg",
			"alt":    "Go API Development",
			"width":  1200,
			"height": 630,
		},
		"seo": map[string]interface{}{
			"title":       "Type-Safe APIs with Go | Complete Guide",
			"description": "Complete guide to building type-safe APIs with Go and zod-go validation library.",
			"keywords":    []string{"go", "api", "type-safe", "validation", "zod-go"},
			"canonical":   "https://example.com/blog/building-type-safe-apis-go-zod",
		},
		"comments": []map[string]interface{}{
			{
				"id":        1,
				"content":   "Great article! This really helped me understand how to implement proper validation.",
				"author":    "Bob Reader",
				"email":     "bob@example.com",
				"website":   "https://bob-codes.com",
				"createdAt": "2023-12-06T10:30:00Z",
				"approved":  true,
				"replies": []map[string]interface{}{
					{
						"id":        1,
						"content":   "Thanks Bob! Glad it was helpful.",
						"author":    "Jane Developer",
						"createdAt": "2023-12-06T11:00:00Z",
					},
				},
			},
		},
		"publishedAt": "2023-12-06T09:00:00Z",
		"updatedAt":   "2023-12-06T09:00:00Z",
		"status":      "published",
		"viewCount":   1250,
		"featured":    true,
	}

	fmt.Printf("✅ Valid blog post: %s\n", validateAndReport(blogPostSchema, blogPostData))
}

// Example 6: Analytics and Metrics Data
func runAnalyticsExample() {
	fmt.Println("\n📊 Example 6: Analytics and Metrics")
	fmt.Println("===================================")

	eventSchema := validators.Object(map[string]interface{}{
		"eventId": validators.String().
			Pattern(`^[a-f0-9-]{36}$`).WithPatternMessage("Event ID must be a valid UUID").
			Required(),

		"userId": validators.String().Optional(),

		"sessionId": validators.String().
			Pattern(`^[a-zA-Z0-9-]+$`).
			Required(),

		"timestamp": validators.Number().
			Integer().Min(0).
			Required(),

		"eventType": validators.String().
			Pattern(`^(page_view|click|form_submit|purchase|signup|login|logout)$`).
			Required(),

		"properties": validators.Object(map[string]interface{}{
			"page":     validators.String().Optional(),
			"referrer": validators.String().Optional(),
			"userAgent": validators.String().Optional(),
			"ip":       validators.String().Optional(),
			"country":  validators.String().Optional(),
			"device":   validators.String().Optional(),
			"browser":  validators.String().Optional(),
		}).Optional(),

		"customData": validators.Object(map[string]interface{}{}).Optional(),
	}).Required()

	metricsSchema := validators.Object(map[string]interface{}{
		"timestamp": validators.Number().Integer().Required(),
		"metrics": validators.Object(map[string]interface{}{
			"activeUsers":    validators.Number().Integer().Min(0).Required(),
			"pageViews":      validators.Number().Integer().Min(0).Required(),
			"uniqueVisitors": validators.Number().Integer().Min(0).Required(),
			"bounceRate":     validators.Number().Min(0).Max(1).Required(),
			"avgSessionDuration": validators.Number().Min(0).Required(),
			"conversionRate": validators.Number().Min(0).Max(1).Required(),
		}).Required(),
		"breakdown": validators.Object(map[string]interface{}{
			"byCountry": validators.Object(map[string]interface{}{}).Optional(),
			"byDevice":  validators.Object(map[string]interface{}{}).Optional(),
			"byPage":    validators.Object(map[string]interface{}{}).Optional(),
		}).Optional(),
	}).Required()

	analyticsEvent := map[string]interface{}{
		"eventId":   "123e4567-e89b-12d3-a456-426614174000",
		"userId":    "user_12345",
		"sessionId": "session_abc123",
		"timestamp": time.Now().Unix(),
		"eventType": "purchase",
		"properties": map[string]interface{}{
			"page":      "/checkout/complete",
			"referrer":  "https://google.com",
			"userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
			"ip":        "192.168.1.1",
			"country":   "US",
			"device":    "desktop",
			"browser":   "chrome",
		},
		"customData": map[string]interface{}{
			"orderId":     "order_789",
			"orderValue":  149.99,
			"currency":    "USD",
			"productIds":  []string{"prod_1", "prod_2"},
		},
	}

	metricsData := map[string]interface{}{
		"timestamp": time.Now().Unix(),
		"metrics": map[string]interface{}{
			"activeUsers":        1250,
			"pageViews":         15000,
			"uniqueVisitors":    8500,
			"bounceRate":        0.35,
			"avgSessionDuration": 185.5,
			"conversionRate":    0.025,
		},
		"breakdown": map[string]interface{}{
			"byCountry": map[string]interface{}{
				"US": 5000,
				"UK": 2000,
				"DE": 1500,
			},
			"byDevice": map[string]interface{}{
				"desktop": 8000,
				"mobile":  6000,
				"tablet":  1000,
			},
		},
	}

	fmt.Printf("✅ Valid analytics event: %s\n", validateAndReport(eventSchema, analyticsEvent))
	fmt.Printf("✅ Valid metrics data: %s\n", validateAndReport(metricsSchema, metricsData))
}

// Example 7: File Upload and Media Management
func runFileUploadExample() {
	fmt.Println("\n📁 Example 7: File Upload and Media Management")
	fmt.Println("==============================================")

	fileUploadSchema := validators.Object(map[string]interface{}{
		"filename": validators.String().
			Min(1).Max(255).
			Pattern(`^[a-zA-Z0-9._-]+$`).WithPatternMessage("Filename contains invalid characters").
			Required(),

		"size": validators.Number().
			Integer().Min(1).Max(50*1024*1024). // 50MB max
			WithMinMessage("File cannot be empty").
			WithMaxMessage("File size cannot exceed 50MB").
			Required(),

		"mimeType": validators.String().
			Pattern(`^[a-zA-Z]+/[a-zA-Z0-9.-]+$`).
			Required(),

		"checksum": validators.String().
			Pattern(`^[a-f0-9]{64}$`).WithPatternMessage("Checksum must be a valid SHA-256 hash").
			Required(),

		"metadata": validators.Object(map[string]interface{}{
			"uploadedBy": validators.String().Required(),
			"uploadedAt": validators.String().Required(),
			"tags":       validators.Array(validators.String()).MaxItems(10).Optional(),
			"description": validators.String().Max(500).Optional(),
			"isPublic":   validators.Bool().Optional().Default(false),
		}).Required(),

		"processing": validators.Object(map[string]interface{}{
			"status": validators.String().
				Pattern(`^(pending|processing|completed|failed)$`).
				Optional().Default("pending"),
			"thumbnails": validators.Array(validators.Object(map[string]interface{}{
				"size":   validators.String().Pattern(`^\d+x\d+$`).Required(),
				"url":    validators.URL().Required(),
				"format": validators.String().Pattern(`^(jpg|png|webp)$`).Required(),
			})).Optional(),
			"optimized": validators.Object(map[string]interface{}{
				"url":         validators.URL().Required(),
				"size":        validators.Number().Integer().Min(1).Required(),
				"compression": validators.Number().Min(0).Max(1).Required(),
			}).Optional(),
		}).Optional(),
	}).Required()

	fileUploadData := map[string]interface{}{
		"filename": "profile-photo.jpg",
		"size":     2048576, // 2MB
		"mimeType": "image/jpeg",
		"checksum": "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3",
		"metadata": map[string]interface{}{
			"uploadedBy":  "user_123",
			"uploadedAt":  "2023-12-06T10:30:00Z",
			"tags":        []string{"profile", "avatar", "photo"},
			"description": "User profile photo",
			"isPublic":    true,
		},
		"processing": map[string]interface{}{
			"status": "completed",
			"thumbnails": []map[string]interface{}{
				{
					"size":   "150x150",
					"url":    "https://cdn.example.com/thumbs/profile-photo-150x150.jpg",
					"format": "jpg",
				},
				{
					"size":   "300x300",
					"url":    "https://cdn.example.com/thumbs/profile-photo-300x300.webp",
					"format": "webp",
				},
			},
			"optimized": map[string]interface{}{
				"url":         "https://cdn.example.com/optimized/profile-photo.webp",
				"size":        1024000,
				"compression": 0.85,
			},
		},
	}

	fmt.Printf("✅ Valid file upload: %s\n", validateAndReport(fileUploadSchema, fileUploadData))
}

// Example 8: Social Media Platform Data
func runSocialMediaExample() {
	fmt.Println("\n📱 Example 8: Social Media Platform")
	fmt.Println("===================================")

	postSchema := validators.Object(map[string]interface{}{
		"id": validators.String().
			Pattern(`^[a-f0-9-]{36}$`).
			Required(),

		"content": validators.String().
			Min(1).Max(2000).
			Required(),

		"author": validators.Object(map[string]interface{}{
			"id":       validators.String().Required(),
			"username": validators.String().Pattern(`^[a-zA-Z0-9_]{3,20}$`).Required(),
			"displayName": validators.String().Max(50).Required(),
			"avatar":   validators.URL().Optional(),
			"verified": validators.Bool().Optional().Default(false),
		}).Required(),

		"media": validators.Array(validators.Object(map[string]interface{}{
			"type": validators.String().Pattern(`^(image|video|gif)$`).Required(),
			"url":  validators.URL().Required(),
			"alt":  validators.String().Optional(),
			"dimensions": validators.Object(map[string]interface{}{
				"width":  validators.Number().Integer().Min(1).Required(),
				"height": validators.Number().Integer().Min(1).Required(),
			}).Optional(),
		})).MaxItems(4).Optional(),

		"hashtags": validators.Array(validators.String().
			Pattern(`^[a-zA-Z0-9_]+$`).Min(1).Max(30)).
			MaxItems(10).Optional(),

		"mentions": validators.Array(validators.String().
			Pattern(`^[a-zA-Z0-9_]{3,20}$`)).
			MaxItems(10).Optional(),

		"engagement": validators.Object(map[string]interface{}{
			"likes":    validators.Number().Integer().Min(0).Required(),
			"reposts":  validators.Number().Integer().Min(0).Required(),
			"replies":  validators.Number().Integer().Min(0).Required(),
			"views":    validators.Number().Integer().Min(0).Optional(),
		}).Required(),

		"visibility": validators.String().
			Pattern(`^(public|followers|private)$`).
			Optional().Default("public"),

		"location": validators.Object(map[string]interface{}{
			"name":      validators.String().Required(),
			"latitude":  validators.Number().Min(-90).Max(90).Optional(),
			"longitude": validators.Number().Min(-180).Max(180).Optional(),
		}).Optional(),

		"createdAt": validators.String().Required(),
		"editedAt":  validators.String().Optional(),

		"replies": validators.Array(validators.Object(map[string]interface{}{
			"id":      validators.String().Required(),
			"content": validators.String().Min(1).Max(500).Required(),
			"author":  validators.String().Required(),
			"createdAt": validators.String().Required(),
			"likes":   validators.Number().Integer().Min(0).Required(),
		})).Optional(),
	}).Required()

	socialMediaPost := map[string]interface{}{
		"id":      "123e4567-e89b-12d3-a456-426614174000",
		"content": "Just shipped a new feature using #zod-go for validation! The type safety is incredible 🚀 #golang #webdev",
		"author": map[string]interface{}{
			"id":          "user_12345",
			"username":    "jane_dev",
			"displayName": "Jane Developer",
			"avatar":      "https://example.com/avatars/jane.jpg",
			"verified":    true,
		},
		"media": []map[string]interface{}{
			{
				"type": "image",
				"url":  "https://example.com/media/feature-screenshot.png",
				"alt":  "Screenshot of the new feature",
				"dimensions": map[string]interface{}{
					"width":  1200,
					"height": 800,
				},
			},
		},
		"hashtags": []string{"zod-go", "golang", "webdev", "validation"},
		"mentions": []string{"go_team", "zod_official"},
		"engagement": map[string]interface{}{
			"likes":   42,
			"reposts": 12,
			"replies": 8,
			"views":   1250,
		},
		"visibility": "public",
		"location": map[string]interface{}{
			"name":      "San Francisco, CA",
			"latitude":  37.7749,
			"longitude": -122.4194,
		},
		"createdAt": "2023-12-06T10:30:00Z",
		"replies": []map[string]interface{}{
			{
				"id":        "reply_1",
				"content":   "Looks amazing! Can't wait to try it out.",
				"author":    "user_67890",
				"createdAt": "2023-12-06T10:35:00Z",
				"likes":     5,
			},
		},
	}

	fmt.Printf("✅ Valid social media post: %s\n", validateAndReport(postSchema, socialMediaPost))
}

// Helper function to validate and report results
func validateAndReport(schema interface{ Validate(interface{}) error }, data interface{}) string {
	if err := schema.Validate(data); err != nil {
		return fmt.Sprintf("Validation failed: %s", err.Error())
	}
	return "Validation passed"
}

// HTTP handler example using the validation schemas
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	userSchema := validators.Object(map[string]interface{}{
		"username": validators.String().Min(3).Max(20).Required(),
		"email":    validators.Email(),
		"age":      validators.Number().Min(13).Integer().Required(),
	}).Required()

	var userData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := userSchema.Validate(userData); err != nil {
		response := map[string]interface{}{
			"error":   "Validation failed",
			"details": err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteStatus(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Process valid user data...
	response := map[string]interface{}{
		"success": true,
		"message": "User created successfully",
		"id":      generateUserID(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateUserID() string {
	return "user_" + strconv.FormatInt(time.Now().Unix(), 10)
}

// Example of middleware using validation
func validationMiddleware(schema interface{ Validate(interface{}) error }) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var data map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
				http.Error(w, "Invalid JSON", http.StatusBadRequest)
				return
			}

			if err := schema.Validate(data); err != nil {
				http.Error(w, fmt.Sprintf("Validation failed: %s", err.Error()), http.StatusBadRequest)
				return
			}

			// Add validated data to request context
			// In a real implementation, you'd use context.WithValue
			next.ServeHTTP(w, r)
		})
	}
}

// Example: Command-line tool configuration validation
func validateCliConfig(configFile string) error {
	configSchema := validators.Object(map[string]interface{}{
		"version": validators.String().Pattern(`^\d+\.\d+\.\d+$`).Required(),
		"commands": validators.Array(validators.Object(map[string]interface{}{
			"name":        validators.String().Required(),
			"description": validators.String().Required(),
			"flags": validators.Array(validators.Object(map[string]interface{}{
				"name":     validators.String().Required(),
				"type":     validators.String().Pattern(`^(string|int|bool)$`).Required(),
				"required": validators.Bool().Optional().Default(false),
				"default":  validators.String().Optional(),
			})).Optional(),
		})).Required(),
	}).Required()

	// In a real implementation, you'd read from the file
	_ = configFile
	
	// Example config data
	configData := map[string]interface{}{
		"version": "1.0.0",
		"commands": []map[string]interface{}{
			{
				"name":        "serve",
				"description": "Start the web server",
				"flags": []map[string]interface{}{
					{
						"name":     "port",
						"type":     "int",
						"required": false,
						"default":  "8080",
					},
					{
						"name":     "host",
						"type":     "string",
						"required": false,
						"default":  "localhost",
					},
				},
			},
		},
	}

	return configSchema.Validate(configData)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
