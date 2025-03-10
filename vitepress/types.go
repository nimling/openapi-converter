package vitepress

import "time"

// HeadItem represents a head tag item
type HeadItem struct {
	Tag     string            `yaml:"tag" json:"tag"`
	Options map[string]string `yaml:"options,omitempty" json:"options,omitempty"`
}

// HeroImage represents the hero section image
type HeroImage struct {
	Src    string `yaml:"src" json:"src"`
	Alt    string `yaml:"alt,omitempty" json:"alt,omitempty"`
	Width  *int   `yaml:"width,omitempty" json:"width,omitempty"`
	Height *int   `yaml:"height,omitempty" json:"height,omitempty"`
}

// HeroAction represents a hero section button/action
type HeroAction struct {
	Theme  string  `yaml:"theme" json:"theme"` // brand, alt, sponsor
	Text   string  `yaml:"text" json:"text"`
	Link   string  `yaml:"link" json:"link"`
	Target *string `yaml:"target,omitempty" json:"target,omitempty"` // _blank, _self
}

// Hero represents the hero section configuration
type Hero struct {
	Name    string       `yaml:"name" json:"name"`
	Text    *string      `yaml:"text,omitempty" json:"text,omitempty"`
	Tagline *string      `yaml:"tagline,omitempty" json:"tagline,omitempty"`
	Image   *HeroImage   `yaml:"image,omitempty" json:"image,omitempty"`
	Actions []HeroAction `yaml:"actions,omitempty" json:"actions,omitempty"`
}

// FeatureIcon represents a feature section icon
type FeatureIcon struct {
	Src    string `yaml:"src" json:"src"`
	Width  *int   `yaml:"width,omitempty" json:"width,omitempty"`
	Height *int   `yaml:"height,omitempty" json:"height,omitempty"`
	Alt    string `yaml:"alt,omitempty" json:"alt,omitempty"`
}

// Feature represents a feature item
type Feature struct {
	Icon     interface{} `yaml:"icon,omitempty" json:"icon,omitempty"`
	Title    string      `yaml:"title" json:"title"`
	Details  string      `yaml:"details" json:"details"`
	Link     *string     `yaml:"link,omitempty" json:"link,omitempty"`
	LinkText *string     `yaml:"linkText,omitempty" json:"linkText,omitempty"`
	Rel      *string     `yaml:"rel,omitempty" json:"rel,omitempty"`
}

// Navigation represents prev/next navigation
type Navigation struct {
	Text string `yaml:"text" json:"text"`
	Link string `yaml:"link" json:"link"`
}

// SearchOptions represents search configuration
type SearchOptions struct {
	Provider    string  `yaml:"provider" json:"provider"` // local, algolia
	AppID       *string `yaml:"appId,omitempty" json:"appId,omitempty"`
	APIKey      *string `yaml:"apiKey,omitempty" json:"apiKey,omitempty"`
	IndexName   *string `yaml:"indexName,omitempty" json:"indexName,omitempty"`
	Placeholder *string `yaml:"placeholder,omitempty" json:"placeholder,omitempty"`
	ButtonText  *string `yaml:"buttonText,omitempty" json:"buttonText,omitempty"`
}

// SitemapConfig represents sitemap configuration
type SitemapConfig struct {
	Priority   *float64 `yaml:"priority,omitempty" json:"priority,omitempty"`
	ChangeFreq *string  `yaml:"changefreq,omitempty" json:"changefreq,omitempty"`
}

// MarkdownConfig represents markdown behavior configuration
type MarkdownConfig struct {
	LineNumbers *bool `yaml:"lineNumbers,omitempty" json:"lineNumbers,omitempty"`
	TOC         *bool `yaml:"toc,omitempty" json:"toc,omitempty"`
	Breaks      *bool `yaml:"breaks,omitempty" json:"breaks,omitempty"`
	Linkify     *bool `yaml:"linkify,omitempty" json:"linkify,omitempty"`
	Emoji       *bool `yaml:"emoji,omitempty" json:"emoji,omitempty"`
	Math        *bool `yaml:"math,omitempty" json:"math,omitempty"`
}

// Frontmatter represents the main frontmatter structure
type Frontmatter struct {
	// Required Page Metadata
	Title       string `yaml:"title" json:"title"`
	Description string `yaml:"description" json:"description"`
	Layout      string `yaml:"layout" json:"layout"` // doc, home, page, custom

	// Optional Page Metadata
	Aside        interface{} `yaml:"aside,omitempty" json:"aside,omitempty"`             // bool or string (left, right)
	Outline      interface{} `yaml:"outline,omitempty" json:"outline,omitempty"`         // []int or bool
	LastUpdated  interface{} `yaml:"lastUpdated,omitempty" json:"lastUpdated,omitempty"` // bool or string
	EditLink     *bool       `yaml:"editLink,omitempty" json:"editLink,omitempty"`
	EditLinkText *string     `yaml:"editLinkText,omitempty" json:"editLinkText,omitempty"`
	DocsDir      *string     `yaml:"docsDir,omitempty" json:"docsDir,omitempty"`
	DocFooter    interface{} `yaml:"docFooter,omitempty" json:"docFooter,omitempty"`
	Footer       interface{} `yaml:"footer,omitempty" json:"footer,omitempty"`
	FooterHTML   *bool       `yaml:"footerHtml,omitempty" json:"footerHtml,omitempty"`

	// Page Structure
	Head []HeadItem `yaml:"head,omitempty" json:"head,omitempty"`

	// Navigation
	Nav          interface{} `yaml:"nav,omitempty" json:"nav,omitempty"`
	Navbar       interface{} `yaml:"navbar,omitempty" json:"navbar,omitempty"`
	Sidebar      interface{} `yaml:"sidebar,omitempty" json:"sidebar,omitempty"`
	SidebarDepth *int        `yaml:"sidebarDepth,omitempty" json:"sidebarDepth,omitempty"`
	Prev         *Navigation `yaml:"prev,omitempty" json:"prev,omitempty"`
	Next         *Navigation `yaml:"next,omitempty" json:"next,omitempty"`

	// Home Page Specific
	Hero     *Hero     `yaml:"hero,omitempty" json:"hero,omitempty"`
	Features []Feature `yaml:"features,omitempty" json:"features,omitempty"`

	// Advanced Features
	CustomLayout *string `yaml:"customLayout,omitempty" json:"customLayout,omitempty"`
	PageClass    *string `yaml:"pageClass,omitempty" json:"pageClass,omitempty"`
	ContentClass *string `yaml:"contentClass,omitempty" json:"contentClass,omitempty"`

	// Search Configuration
	Search *SearchOptions `yaml:"search,omitempty" json:"search,omitempty"`

	// Internationalization
	Lang *string `yaml:"lang,omitempty" json:"lang,omitempty"`
	Dir  *string `yaml:"dir,omitempty" json:"dir,omitempty"` // ltr, rtl

	// Content Display
	Date               *time.Time     `yaml:"date,omitempty" json:"date,omitempty"`
	Author             interface{}    `yaml:"author,omitempty" json:"author,omitempty"` // string or []string
	Tags               interface{}    `yaml:"tags,omitempty" json:"tags,omitempty"`     // string or []string
	Categories         interface{}    `yaml:"categories,omitempty" json:"categories,omitempty"`
	ExcludeFromSitemap *bool          `yaml:"excludeFromSitemap,omitempty" json:"excludeFromSitemap,omitempty"`
	Sitemap            *SitemapConfig `yaml:"sitemap,omitempty" json:"sitemap,omitempty"`

	// Security & Performance
	External *bool `yaml:"external,omitempty" json:"external,omitempty"`
	Nofollow *bool `yaml:"nofollow,omitempty" json:"nofollow,omitempty"`
	Cache    *bool `yaml:"cache,omitempty" json:"cache,omitempty"`
	Preload  *bool `yaml:"preload,omitempty" json:"preload,omitempty"`

	// Theme Customization
	Theme       *string             `yaml:"theme,omitempty" json:"theme,omitempty"`
	Appearance  *bool               `yaml:"appearance,omitempty" json:"appearance,omitempty"`
	Logo        interface{}         `yaml:"logo,omitempty" json:"logo,omitempty"`
	SocialLinks []map[string]string `yaml:"socialLinks,omitempty" json:"socialLinks,omitempty"`

	// Custom Data
	FrontmatterData map[string]interface{} `yaml:"frontmatter,omitempty" json:"frontmatter,omitempty"`
	RouteMeta       map[string]interface{} `yaml:"routeMeta,omitempty" json:"routeMeta,omitempty"`

	// Container Elements
	ContainerClass *string `yaml:"containerClass,omitempty" json:"containerClass,omitempty"`

	// Markdown Behavior
	Markdown *MarkdownConfig `yaml:"markdown,omitempty" json:"markdown,omitempty"`

	// Advanced Routing
	Alias     interface{} `yaml:"alias,omitempty" json:"alias,omitempty"` // string or []string
	Permalink *string     `yaml:"permalink,omitempty" json:"permalink,omitempty"`
	Dynamic   *bool       `yaml:"dynamic,omitempty" json:"dynamic,omitempty"`

	// Development
	IsDev *bool `yaml:"isDev,omitempty" json:"isDev,omitempty"`
	Debug *bool `yaml:"debug,omitempty" json:"debug,omitempty"`
}
