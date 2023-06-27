package app

import "strings"

// A Twitter card: https://developer.twitter.com/en/docs/twitter-for-websites/cards
type TwitterCard struct {
	// The card type, which will be one of “summary”, “summary_large_image”,
	// “app”, or “player”.
	Card string

	// Username for the website used in the card footer.
	Site string

	// Username for the content creator / author.
	Creator string

	// A concise title for the related content.
	Title string

	// A description that concisely summarizes the content as appropriate for
	// presentation within a Tweet.
	Description string

	// A URL to a unique image representing the content of the page.
	Image string

	// A text description of the image conveying the essential nature of an
	// image to users who are visually impaired. Maximum 420 characters.
	ImageAlt string
}

func (c TwitterCard) toMap() map[string]string {
	m := make(map[string]string)

	if c.Card != "" {
		m["twitter:card"] = c.Card
	}

	if c.Site != "" {
		if !strings.HasPrefix(c.Site, "@") {
			c.Site = "@" + c.Site
		}
		m["twitter:site"] = c.Site
	}

	if c.Creator != "" {
		if !strings.HasPrefix(c.Creator, "@") {
			c.Creator = "@" + c.Creator
		}
		m["twitter:creator"] = c.Creator
	}

	if c.Title != "" {
		m["twitter:title"] = c.Title
	}

	if c.Description != "" {
		m["twitter:description"] = c.Description
	}

	if c.Image != "" {
		m["twitter:image"] = c.Image
	}

	if c.ImageAlt != "" {
		m["twitter:image:alt"] = c.ImageAlt
	}

	return m
}
