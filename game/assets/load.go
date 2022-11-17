package assets

import (
	"fmt"
	"sync"
	"time"
)

func init() {
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(3)

	var assetManagers []Manager = []Manager{
		FontManager,
		ImageManager,
		SoundManager,
	}

	for _, assetManager := range assetManagers {
		go func(assetManager Manager) {
			defer wg.Done()
			assetManager.Load()
		}(assetManager)
	}

	wg.Wait()
	fmt.Printf("Loaded assets in %v ms\n", time.Since(start).Milliseconds())
}
