package crawler

// Category specifies which category we need to crawl
type Category = int

const (
	CategorySchweiz int = 6
)

// Tenant part of request
// TODO: Need to understand what does it mean
// Using just tenant=2 like in the request from browser
// Testing different values, available range seems to be tenant=[1...15]
// Content a bit different and language also alternates between DE and FR
type Tenant = int

const (
	TenantTwo = 2
)
