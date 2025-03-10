package converter

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

func (n *OpenAPIConverter) ResolveExternalRefs() error {
	// Initialize component register if not already done
	if n.doc.Components == nil {
		n.doc.Components = &Components{}
	}

	if n.doc.Components.Register == nil {
		n.doc.Components.Register = make(ReferenceRegister)
	}

	// Iteratively resolve references until no external refs remain
	maxIterations := 10 // Prevent infinite loops
	for i := 0; i < maxIterations; i++ {
		externalRefsRemain := false

		// First resolve components
		if err := n.resolveComponentRefs(n.doc.Components); err != nil {
			return fmt.Errorf("failed to resolve component references: %w", err)
		}

		// Then resolve path items
		if n.doc.Paths != nil {
			for key, pathItem := range n.doc.Paths {
				if err := n.resolvePathItemRefs(pathItem); err != nil {
					return fmt.Errorf("failed to resolve refs in path %s: %w", key, err)
				}
			}
		}

		// Check if any external references remain
		externalRefsRemain = n.hasExternalRefs()

		if !externalRefsRemain {
			// No more external refs, we're done
			break
		}

		if i == maxIterations-1 {
			return fmt.Errorf("failed to resolve all external references after %d iterations", maxIterations)
		}
	}

	return nil
}

func (n *OpenAPIConverter) resolvePathItemRefs(pathItem *PathItem) error {
	if pathItem == nil {
		return nil
	}

	// Process path parameters if any
	if pathItem.Parameters != nil {
		for i, param := range pathItem.Parameters {
			if param == nil {
				continue
			}

			if param.Ref != nil && isExternalRef(*param.Ref) {
				filePath := resolveRefPath(n.filePath, *param.Ref)
				resolved, err := loadExternalRef[Parameter](filePath)
				if err != nil {
					return fmt.Errorf("failed to load external parameter ref: %w", err)
				}
				pathItem.Parameters[i] = resolved
			}

			// Process parameter schema
			if param.Schema != nil {
				if err := param.Schema.resolveExternalRefs(&n.doc.Components.Register, n.filePath); err != nil {
					return fmt.Errorf("failed to resolve parameter schema: %w", err)
				}
			}
		}
	}

	for method, op := range pathItem.Operations() {
		if op == nil {
			continue
		}

		// Handle operation reference
		relPath := n.filePath
		if op.Ref != nil && isExternalRef(*op.Ref) {
			filePath := resolveRefPath(n.filePath, *op.Ref)
			resolved, err := loadExternalRef[Operation](filePath)
			if err != nil {
				return fmt.Errorf("failed to load external operation ref: %w", err)
			}
			op = resolved
			relPath = filePath
		}

		// Process operation parameters
		if op.Parameters != nil {
			for i, param := range op.Parameters {
				if param == nil {
					continue
				}

				internalRelPath := relPath
				if param.Ref != nil && isExternalRef(*param.Ref) {
					internalRelPath = resolveRefPath(relPath, *param.Ref)
					resolved, err := loadExternalRef[Parameter](internalRelPath)
					if err != nil {
						return fmt.Errorf("failed to load external parameter ref: %w", err)
					}
					op.Parameters[i] = resolved
					param = resolved
				}

				// Process parameter schema
				if param.Schema != nil {
					if err := param.Schema.resolveExternalRefs(&n.doc.Components.Register, internalRelPath); err != nil {
						return fmt.Errorf("failed to resolve parameter schema: %w", err)
					}
				}
			}
		}

		// Process request body
		if op.RequestBody != nil {
			requestRelPath := relPath
			if op.RequestBody.Ref != nil && isExternalRef(*op.RequestBody.Ref) {
				requestRelPath = resolveRefPath(relPath, *op.RequestBody.Ref)
				resolved, err := loadExternalRef[RequestBody](requestRelPath)
				if err != nil {
					return fmt.Errorf("failed to load external request body ref: %w", err)
				}
				op.RequestBody = resolved
			}

			// Process request body content schemas
			if op.RequestBody.Content != nil {
				for mediaType, content := range op.RequestBody.Content {
					if content == nil || content.Schema == nil {
						continue
					}

					if err := content.Schema.resolveExternalRefs(&n.doc.Components.Register, requestRelPath); err != nil {
						return fmt.Errorf("failed to resolve request body schema for %s: %w", mediaType, err)
					}
				}
			}
		}

		// Process responses and their schemas
		if op.Responses != nil {
			for code, response := range op.Responses {
				if response == nil {
					continue
				}

				responseRelPath := relPath
				if response.Ref != nil && isExternalRef(*response.Ref) {
					responseRelPath = resolveRefPath(relPath, *response.Ref)
					resolved, err := loadExternalRef[Response](responseRelPath)
					if err != nil {
						return fmt.Errorf("failed to load external response ref: %w", err)
					}
					op.Responses[code] = resolved
					response = resolved
				}

				// Process response content schemas
				if response.Content != nil {
					for mediaType, content := range response.Content {
						if content == nil || content.Schema == nil {
							continue
						}

						if err := content.Schema.resolveExternalRefs(&n.doc.Components.Register, responseRelPath); err != nil {
							return fmt.Errorf("failed to resolve response schema for %s: %w", mediaType, err)
						}
					}
				}
			}
		}

		// Update the operation in path item
		pathItem.SetMethodOperation(method, op)
	}

	return nil
}

func (r *Schema) resolveExternalRefs(register *ReferenceRegister, relPath string) error {
	// First handle the current schema's ref
	if r.Ref != nil && isExternalRef(*r.Ref) {
		filePath := resolveRefPath(relPath, *r.Ref)
		if existingRef, ok := (*register)[filePath]; ok {
			r.Ref = &existingRef
		} else {
			resolved, err := loadExternalRef[Schema](filePath)
			if err != nil {
				return fmt.Errorf("failed to load external ref: %w", err)
			}

			// Set internal reference and copy properties
			(*register).SetComponent("schemas", filePath)

			// Keep the reference but also copy properties for validation
			oldRef := r.Ref
			*r = *resolved
			r.Ref = oldRef

			// Continue resolving the newly loaded schema's references
			if err := resolved.resolveExternalRefs(register, filePath); err != nil {
				return err
			}
		}
	}

	// Process properties map recursively
	if r.Properties != nil {
		for propName, prop := range r.Properties {
			if prop == nil {
				continue
			}

			// Recursively resolve refs in the property
			if err := prop.resolveExternalRefs(register, relPath); err != nil {
				return fmt.Errorf("failed to resolve property '%s': %w", propName, err)
			}
		}
	}

	// Process items for array types
	if r.Items != nil {
		if err := r.Items.resolveExternalRefs(register, relPath); err != nil {
			return fmt.Errorf("failed to resolve array items: %w", err)
		}
	}

	return nil
}

func (n *OpenAPIConverter) resolveComponentRefs(components *Components) error {
	if components == nil {
		return nil
	}

	if components.SecuritySchemes != nil {
		for key, comp := range components.SecuritySchemes {
			if comp.Ref != nil && isExternalRef(*comp.Ref) {
				filePath := resolveRefPath(n.filePath, *comp.Ref)
				res, err := loadExternalRef[SecurityScheme](filePath)
				if err != nil {
					return fmt.Errorf("failed to load external ref: %s, error: %w", err)
				}

				components.Register.SetComponent("securitySchemes", filePath)
				components.SecuritySchemes[key] = res
			}
		}
	}

	if components.Parameters != nil {
		for key, comp := range components.Parameters {
			if comp.Ref != nil && isExternalRef(*comp.Ref) {
				filePath := resolveRefPath(n.filePath, *comp.Ref)
				res, err := loadExternalRef[Parameter](filePath)
				if err != nil {
					return fmt.Errorf("failed to load external ref: %s, error: %w", err)
				}

				components.Register.SetComponent("parameters", filePath)
				components.Parameters[key] = res
			}
		}
	}

	relPath := n.filePath
	if components.Schemas != nil {
		for key, comp := range components.Schemas {
			if comp.Ref != nil && isExternalRef(*comp.Ref) {
				relPath = resolveRefPath(n.filePath, *comp.Ref)
				res, err := loadExternalRef[Schema](relPath)
				if err != nil {
					return fmt.Errorf("failed to load external ref: %s, error: %w", err)
				}

				components.Register.SetComponent("schemas", relPath)
				comp = res
			} else {
				def := components.Register.GetComponent("schemas", key)
				if def != nil {
					relPath = def.FilePath
				}
			}

			if err := comp.resolveExternalRefs(&n.doc.Components.Register, relPath); err != nil {
				return fmt.Errorf("failed to resolve external refs: %w", err)
			}

			components.Schemas[key] = comp
		}
	}

	return nil
}

func resolveRefPath(specPath, refPath string) string {
	filePath, _ := splitRefPath(refPath)

	if !strings.HasPrefix(filePath, "./") && !strings.HasPrefix(filePath, "../") {
		return filePath
	}

	baseDir := filepath.Dir(specPath)
	absPath := filepath.Join(baseDir, filePath)

	return filepath.Clean(absPath)
}

func loadExternalRef[T any](filePath string) (*T, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filePath, err)
	}

	var result T
	if err = yaml.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	return &result, nil
}

func splitRefPath(refPath string) (string, string) {
	parts := strings.SplitN(refPath, "#", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return refPath, ""
}

func isExternalRef(refPath string) bool {
	// Check if the reference points to an external file
	// External refs typically start with ./ or ../ or are absolute paths
	return strings.HasPrefix(refPath, "./") ||
		strings.HasPrefix(refPath, "../") ||
		strings.HasPrefix(refPath, "/") ||
		strings.Contains(refPath, "://")
}

// Add helper method to check for remaining external refs
func (n *OpenAPIConverter) hasExternalRefs() bool {
	// Check components for external refs
	if n.doc.Components != nil {
		if hasExternalRefsInComponents(n.doc.Components) {
			return true
		}
	}

	// Check paths for external refs
	if n.doc.Paths != nil {
		for _, pathItem := range n.doc.Paths {
			if hasExternalRefsInPathItem(pathItem) {
				return true
			}
		}
	}

	return false
}

// Helper functions to check for external refs
func hasExternalRefsInComponents(components *Components) bool {
	if components == nil {
		return false
	}

	// Check schemas
	if components.Schemas != nil {
		for _, schema := range components.Schemas {
			if hasExternalRefsInSchema(schema) {
				return true
			}
		}
	}

	// Check parameters
	if components.Parameters != nil {
		for _, param := range components.Parameters {
			if param.Ref != nil && isExternalRef(*param.Ref) {
				return true
			}
			if param.Schema != nil && hasExternalRefsInSchema(param.Schema) {
				return true
			}
		}
	}

	// Check security schemes
	if components.SecuritySchemes != nil {
		for _, scheme := range components.SecuritySchemes {
			if scheme.Ref != nil && isExternalRef(*scheme.Ref) {
				return true
			}
		}
	}

	return false
}

func hasExternalRefsInSchema(schema *Schema) bool {
	if schema == nil {
		return false
	}

	// Check direct ref
	if schema.Ref != nil && isExternalRef(*schema.Ref) {
		return true
	}

	// Check properties
	if schema.Properties != nil {
		for _, prop := range schema.Properties {
			if hasExternalRefsInSchema(prop) {
				return true
			}
		}
	}

	// Check items
	if schema.Items != nil && hasExternalRefsInSchema(schema.Items) {
		return true
	}

	return false
}

func hasExternalRefsInPathItem(pathItem *PathItem) bool {
	if pathItem == nil {
		return false
	}

	// Check parameters
	if pathItem.Parameters != nil {
		for _, param := range pathItem.Parameters {
			if param.Ref != nil && isExternalRef(*param.Ref) {
				return true
			}
			if param.Schema != nil && hasExternalRefsInSchema(param.Schema) {
				return true
			}
		}
	}

	// Check operations
	for _, operation := range pathItem.Operations() {
		if operation == nil {
			continue
		}

		// Check operation ref
		if operation.Ref != nil && isExternalRef(*operation.Ref) {
			return true
		}

		// Check parameters
		if operation.Parameters != nil {
			for _, param := range operation.Parameters {
				if param.Ref != nil && isExternalRef(*param.Ref) {
					return true
				}
				if param.Schema != nil && hasExternalRefsInSchema(param.Schema) {
					return true
				}
			}
		}

		// Check request body
		if operation.RequestBody != nil {
			if operation.RequestBody.Ref != nil && isExternalRef(*operation.RequestBody.Ref) {
				return true
			}
			if operation.RequestBody.Content != nil {
				for _, content := range operation.RequestBody.Content {
					if content != nil && content.Schema != nil && hasExternalRefsInSchema(content.Schema) {
						return true
					}
				}
			}
		}

		// Check responses
		if operation.Responses != nil {
			for _, response := range operation.Responses {
				if response.Ref != nil && isExternalRef(*response.Ref) {
					return true
				}
				if response.Content != nil {
					for _, content := range response.Content {
						if content != nil && content.Schema != nil && hasExternalRefsInSchema(content.Schema) {
							return true
						}
					}
				}
			}
		}
	}

	return false
}
