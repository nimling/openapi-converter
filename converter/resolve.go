package converter

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

func (n *OpenAPIConverter) ResolveExternalRefs() error {
	if n.doc.Components != nil {
		if err := n.resolveComponentRefs(n.doc.Components); err != nil {
			return fmt.Errorf("failed to resolve component references: %w", err)
		}
	}

	if n.doc.Paths != nil {
		for key, pathItem := range n.doc.Paths {
			if err := n.resolvePathItemRefs(pathItem); err != nil {
				return fmt.Errorf("failed to resolve refs in path %s: %w", key, err)
			}
		}
	}

	return nil
}

func (n *OpenAPIConverter) resolvePathItemRefs(pathItem *PathItem) error {
	for method, op := range pathItem.Operations() {
		if op == nil {
			continue
		}

		relPath := n.filePath

		if op.Ref != nil && isExternalRef(*op.Ref) {
			filePath := resolveRefPath(n.filePath, *op.Ref)
			resolved, err := loadExternalRef[Operation](filePath)
			if err != nil {
				return fmt.Errorf("failed to load external ref: %w", err)
			}
			op = resolved
			relPath = filePath
		}

		if op.Parameters != nil && len(op.Parameters) > 0 {
			for i, param := range op.Parameters {
				internalRelPath := relPath
				if param.Ref != nil && isExternalRef(*param.Ref) {
					internalRelPath = resolveRefPath(relPath, *param.Ref)
					resolved, err := loadExternalRef[Parameter](internalRelPath)
					if err != nil {
						return fmt.Errorf("failed to load external ref: %w", err)
					}
					op.Parameters[i] = resolved
				}

				if param.Schema != nil {
					if err := param.Schema.resolveExternalRefs(&n.doc.Components.Register, internalRelPath); err != nil {
						return fmt.Errorf("failed to load external ref: %w", err)
					}
				}
			}
		}

		if op.RequestBody != nil {
			if op.RequestBody.Ref != nil && isExternalRef(*op.RequestBody.Ref) {
				filePath := resolveRefPath(relPath, *op.RequestBody.Ref)
				resolved, err := loadExternalRef[RequestBody](filePath)
				if err != nil {
					return fmt.Errorf("failed to load external ref: %w", err)
				}
				op.RequestBody = resolved
			}

			if op.RequestBody.Content != nil {
				for _, content := range op.RequestBody.Content {
					if err := content.Schema.resolveExternalRefs(&n.doc.Components.Register, relPath); err != nil {
						return fmt.Errorf("failed to load external ref: %w", err)
					}
				}
			}
		}

		if op.Responses != nil {
			for code, response := range op.Responses {
				if response.Ref != nil && isExternalRef(*response.Ref) {
					filePath := resolveRefPath(relPath, *response.Ref)
					resolved, err := loadExternalRef[Response](filePath)
					if err != nil {
						return fmt.Errorf("failed to load external ref: %w", err)
					}
					op.Responses[code] = resolved
				}

				if response.Content == nil {
					continue
				}

				for contentType, content := range response.Content {
					if content != nil && content.Schema != nil {
						if content.Schema.Ref != nil && isExternalRef(*content.Schema.Ref) {
							filePath := resolveRefPath(relPath, *content.Schema.Ref)
							if existingRef, ok := n.doc.Components.Register[filePath]; ok {
								content.Schema.Ref = &existingRef
							} else {
								resolved, err := loadExternalRef[ResponseContent](filePath)
								if err != nil {
									return fmt.Errorf("failed to load external ref: %w", err)
								}
								content = resolved
								response.Content[contentType] = content
							}
						}

						if err := content.Schema.resolveExternalRefs(&n.doc.Components.Register, relPath); err != nil {
							return fmt.Errorf("failed to load external ref: %w", err)
						}
					}
				}
			}
		}

		pathItem.SetMethodOperation(method, op)
	}

	return nil
}

func (r *Schema) resolveExternalRefs(register *ReferenceRegister, relPath string) error {
	for _, property := range r.Properties {
		if property.Ref == nil || !isExternalRef(*property.Ref) {
			continue
		}

		filePath := resolveRefPath(relPath, *property.Ref)
		if existingRef, ok := (*register)[filePath]; ok {
			property.Ref = &existingRef
		} else {
			resolved, err := loadExternalRef[Schema](filePath)
			if err != nil {
				return fmt.Errorf("failed to load external ref: %w", err)
			}
			property = resolved
			register.SetComponent("schemas", filePath)
		}

		if property.Properties != nil && len(property.Properties) > 0 {
			if err := property.resolveExternalRefs(register, relPath); err != nil {
				return fmt.Errorf("failed to resolve external refs: %w", err)
			}
		}

		if property.Items != nil {
			if err := property.Items.resolveExternalRefs(register, relPath); err != nil {
				return fmt.Errorf("failed to resolve external refs: %w", err)
			}
		}
	}

	if r.Items == nil || r.Items.Ref == nil || !isExternalRef(*r.Items.Ref) {
		return nil
	}

	filePath := resolveRefPath(relPath, *r.Items.Ref)
	if existingRef, ok := (*register)[filePath]; ok {
		r.Items.Ref = &existingRef
	} else {
		resolved, err := loadExternalRef[Schema](filePath)
		if err != nil {
			return fmt.Errorf("failed to load external ref: %w", err)
		}
		r.Items = resolved
		register.SetComponent("schemas", filePath)
	}

	if r.Items.Properties != nil && len(r.Items.Properties) > 0 {
		if err := r.Items.resolveExternalRefs(register, relPath); err != nil {
			return fmt.Errorf("failed to resolve external refs: %w", err)
		}
	}

	if r.Items.Items != nil {
		if err := r.Items.resolveExternalRefs(register, relPath); err != nil {
			return fmt.Errorf("failed to resolve external refs: %w", err)
		}
	}

	return nil
}

func (n *OpenAPIConverter) resolveComponentRefs(components *Components) error {
	if components == nil {
		return nil
	}

	components.Register = make(map[string]string)
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
	if strings.HasPrefix(filePath, "./") || strings.HasPrefix(filePath, "../") {
		baseDir := filepath.Dir(specPath)
		filePath = filepath.Join(baseDir, filePath)
	}

	return filePath
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
