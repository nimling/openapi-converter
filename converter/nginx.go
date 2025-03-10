package converter

import "strings"

func (n *OpenAPIConverter) WriteNginxConfiguration() (string, error) {
	// Validate required fields first
	if err := n.ValidateDocument(); err != nil {
		return "", err
	} // TODO:: Excessive??

	var locations []string

	// Use KeysFromNewest() to iterate over the paths
	for key, pathItem := range n.doc.Paths {
		location, err := n.convertPath(key, pathItem)
		if err != nil {
			return "", err
		}
		locations = append(locations, location)
	}

	return strings.Join(locations, "\n\n"), nil
}
