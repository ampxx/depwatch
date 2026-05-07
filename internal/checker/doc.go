// Package checker provides functionality for querying the Go module proxy
// to retrieve the latest available version of a given Go module.
//
// It uses the public Go module proxy (proxy.golang.org) by default, but
// can be configured to use any compatible proxy via the GOPROXY-style
// base URL.
//
// Example usage:
//
//	client := checker.NewClient("https://proxy.golang.org")
//	version, err := client.LatestVersion(ctx, "github.com/some/module")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(version)
package checker
