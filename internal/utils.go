package internal

import "fmt"

func PrintBanner() {
	banner := `
 ██████╗  █████╗  ██████╗
██╔═══██╗██╔══██╗██╔════╝
██║   ██║███████║██║     
██║   ██║██╔══██║██║     
╚██████╔╝██║  ██║╚██████╗
 ╚═════╝ ╚═╝  ╚═╝ ╚═════╝
OpenAPI Converter v1.0.0
`
	fmt.Println(banner)
}